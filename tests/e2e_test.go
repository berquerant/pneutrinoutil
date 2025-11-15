package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/berquerant/pneutrinoutil/cli/ctl"
	"github.com/berquerant/pneutrinoutil/server/handler"
	"github.com/stretchr/testify/assert"
)

const (
	mockcli = "../dist/pneutrinoutil-mockcli"
	gendata = "../dist/pneutrinoutil-gendata"
	worker  = "../dist/pneutrinoutil-worker"
	server  = "../dist/pneutrinoutil-server"

	serverHost    = "0.0.0.0"
	serverPort    = "9101"
	storageBucket = "test"

	eventuallyWaitForMax = time.Second * 5
	eventuallyTick       = time.Millisecond * 300
)

func assertNil[T any](t assert.TestingT, v T) bool {
	return assert.Nil(t, v, "%v should be nil", v)
}

func assertAndGet[T any](t *testing.T, url string) (T, bool) {
	var (
		d    T
		retb bool
	)
	t.Run("assertAndGet:"+url, func(t *testing.T) {
		r, err := http.Get(url)
		if !assertNil(t, err) {
			return
		}
		defer r.Body.Close()
		if !assert.Equal(t, http.StatusOK, r.StatusCode) {
			return
		}
		b, err := io.ReadAll(r.Body)
		if !assertNil(t, err) {
			return
		}
		var body handler.SuccessResponse[T]
		if !assertNil(t, json.Unmarshal(b, &body)) {
			return
		}
		if !assert.True(t, body.OK) {
			return
		}
		d = body.Data
		retb = true
	})
	return d, retb
}

func eventually(t *testing.T, condition func(c *assert.CollectT), msgAndArgs ...any) bool {
	return assert.EventuallyWithT(t, condition, eventuallyWaitForMax, eventuallyTick, msgAndArgs...)
}

func generateData(content string, basename ...string) (map[string]string, error) { // basename to rid
	d := map[string]string{}
	if len(basename) == 0 {
		return d, nil
	}

	stdin := bytes.NewBufferString(strings.Join(basename, "\n"))
	var stdout bytes.Buffer
	cmd := exec.Command(gendata, "--content", content)
	cmd.Stdin = stdin
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(&stdout)
	for scanner.Scan() {
		xs := strings.SplitN(scanner.Text(), " ", 2)
		d[xs[1]] = xs[0]
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(d) != len(basename) {
		return nil, fmt.Errorf("want %d data but got %d", len(basename), len(d))
	}
	return d, nil
}

func TestE2E(t *testing.T) {
	var (
		redisDB  = os.Getenv("REDIS_TEST_DB")
		redisDSN = fmt.Sprintf("redis://%s:%s/%s",
			os.Getenv("REDIS_HOST"),
			os.Getenv("REDIS_PORT"),
			redisDB,
		)
		mysqlDSN = fmt.Sprintf("test:test@tcp(%s:%s)/test?parseTime=true&loc=Asia%%2FTokyo",
			os.Getenv("MYSQL_HOST"),
			os.Getenv("MYSQL_PORT"),
		)

		workerWorkDir = filepath.Join(t.TempDir(), "worker_workspace")
	)

	t.Log("start worker")
	workerProcess := exec.Command(
		worker,
		"--redisDSN", redisDSN,
		"--mysqlDSN", mysqlDSN,
		"--storageS3",
		"--storageBucket", storageBucket,
		"--pneutrinoutil", mockcli,
		"--workDir", workerWorkDir,
		"--debug",
	)
	workerProcess.Stdout = os.Stdout
	workerProcess.Stderr = os.Stderr

	if !assertNil(t, workerProcess.Start()) {
		return
	}

	t.Log("start server")
	serverProcess := exec.Command(
		server,
		"--redisDSN", redisDSN,
		"--mysqlDSN", mysqlDSN,
		"--storageS3",
		"--storageBucket", storageBucket,
		"--host", serverHost,
		"--port", serverPort,
		"--debug",
	)
	serverProcess.Stdout = os.Stdout
	serverProcess.Stderr = os.Stderr
	if !assertNil(t, serverProcess.Start()) {
		return
	}

	newUrl := func(path string) string {
		return fmt.Sprintf("http://%s:%s/v1%s", serverHost, serverPort, path)
	}

	eventually(t, func(c *assert.CollectT) {
		r, err := http.Get(newUrl("/health"))
		if !assertNil(c, err) {
			return
		}
		defer r.Body.Close()
		assert.Equal(t, http.StatusOK, r.StatusCode)
	}, "healthcheck")

	t.Run("version", func(t *testing.T) {
		r, err := http.Get(newUrl("/version"))
		if !assertNil(t, err) {
			return
		}
		defer r.Body.Close()
		assert.Equal(t, http.StatusOK, r.StatusCode)
	})

	t.Run("debug", func(t *testing.T) {
		r, err := http.Get(newUrl("/debug"))
		if !assertNil(t, err) {
			return
		}
		defer r.Body.Close()
		assert.Equal(t, http.StatusOK, r.StatusCode)
	})

	const (
		scoreFileName = "score_sample.musicxml"
		basename      = "score_sample"
		scoreContent  = "score content"
	)

	var newRid string
	t.Run("new process", func(t *testing.T) {
		d, err := generateData(scoreContent, basename)
		if !assert.Nil(t, err) {
			return
		}
		newRid = d[basename]
	})

	t.Run("details", func(t *testing.T) {
		got, ok := assertAndGet[handler.GetDetailResponseData](t, newUrl("/proc/"+newRid+"/detail"))
		assert.True(t, ok)
		assert.Equal(t, "succeed", got.Status)
		assert.Equal(t, newRid, got.RequestID)
		assert.Equal(t, basename, got.Basename)
	})

	t.Run("download config", func(t *testing.T) {
		got, ok := assertAndGet[ctl.Config](t, newUrl("/proc/"+newRid+"/config"))
		if !assert.True(t, ok) {
			return
		}
		assert.Equal(t, newRid, got.Description)
		assert.Equal(t, scoreFileName, filepath.Base(got.Score))
		assert.Equal(t, basename, got.Basename())
	})

	t.Run("download log", func(t *testing.T) {
		r, err := http.Get(newUrl("/proc/" + newRid + "/log"))
		if !assertNil(t, err) {
			return
		}
		r.Body.Close()
		assert.Equal(t, http.StatusOK, r.StatusCode)
	})

	t.Run("download musicxml", func(t *testing.T) {
		r, err := http.Get(newUrl("/proc/" + newRid + "/musicxml"))
		if !assertNil(t, err) {
			return
		}
		defer r.Body.Close()
		assert.Equal(t, http.StatusOK, r.StatusCode)
		body, err := io.ReadAll(r.Body)
		if !assertNil(t, err) {
			return
		}
		assert.Equal(t, scoreContent, string(body))
	})

	t.Run("download wav", func(t *testing.T) {
		r, err := http.Get(newUrl("/proc/" + newRid + "/wav"))
		if !assertNil(t, err) {
			return
		}
		r.Body.Close()
		assert.Equal(t, http.StatusOK, r.StatusCode)
	})

	t.Run("download world wav", func(t *testing.T) {
		r, err := http.Get(newUrl("/proc/" + newRid + "/world_wav"))
		if !assertNil(t, err) {
			return
		}
		r.Body.Close()
		assert.Equal(t, http.StatusOK, r.StatusCode)
	})

	t.Run("search", func(t *testing.T) {
		r, ok := assertAndGet[handler.SearchProcessResponseData](t, newUrl("/proc/search"))
		if !ok {
			return
		}
		if !assert.Len(t, r, 1) {
			return
		}
		x := r[0]
		assert.Equal(t, newRid, x.RequestID)
		assert.Equal(t, basename, x.Title)
	})

	t.Run("search complex", func(t *testing.T) {
		joinMap := func(d1, d2 map[string]string) map[string]string {
			for k, v := range d2 {
				d1[k] = v
			}
			return d1
		}

		//
		// data to be prepared
		//
		// basename
		// t=start1
		// a1, a2
		// t=start2
		// b1, b2, a3
		// t=start3
		// c1
		// d1

		d := map[string]string{
			basename: newRid,
		}
		time.Sleep(time.Second)
		// start1 := time.Now()
		{
			x, err := generateData("content1", "a1", "a2")
			if !assert.Nil(t, err) {
				return
			}
			d = joinMap(d, x)
		}
		time.Sleep(time.Second)
		start2 := time.Now()
		time.Sleep(time.Second)
		{
			x, err := generateData("content2", "b1", "b2", "a3")
			if !assert.Nil(t, err) {
				return
			}
			d = joinMap(d, x)
		}
		time.Sleep(time.Second)
		start3 := time.Now()
		time.Sleep(time.Second)
		{
			x, err := generateData("content3", "c1")
			if !assert.Nil(t, err) {
				return
			}
			d = joinMap(d, x)
		}
		time.Sleep(time.Second)
		{
			x, err := generateData("content4", "d1")
			if !assert.Nil(t, err) {
				return
			}
			d = joinMap(d, x)
		}

		t.Run("order by craeted_at desc", func(t *testing.T) {
			r, ok := assertAndGet[handler.SearchProcessResponseData](t,
				newUrl("/proc/search?limit=2"),
			)
			if !ok {
				return
			}
			if !assert.Len(t, r, 2) {
				return
			}
			assert.Equal(t, "d1", r[0].Title)
			assert.Equal(t, "c1", r[1].Title)
		})

		t.Run("show all", func(t *testing.T) {
			r, ok := assertAndGet[handler.SearchProcessResponseData](t,
				newUrl("/proc/search?limit=100"),
			)
			if !ok {
				return
			}
			t.Logf("%v", d)
			for _, x := range r {
				t.Logf("%s %d %s %s", x.CreatedAt, x.CreatedAt.Unix(), x.RequestID, x.Title)
			}
		})

		for _, tc := range []struct {
			title     string
			q         string
			basenames []string
		}{
			{
				title:     "all",
				q:         "limit=100",
				basenames: slices.Collect(maps.Keys(d)),
			},
			{
				title: "no running",
				q:     "status=running",
			},
			{
				title:     "prefix",
				q:         "prefix=a",
				basenames: []string{"a1", "a2", "a3"},
			},
			{
				title:     "title and start",
				q:         fmt.Sprintf("prefix=a&start=%d", start2.Unix()),
				basenames: []string{"a3"},
			},
			{
				title:     "time range",
				q:         fmt.Sprintf("start=%d&end=%d", start2.Unix(), start3.Unix()),
				basenames: []string{"b1", "b2", "a3"},
			},
		} {
			t.Run(tc.title, func(t *testing.T) {
				u := newUrl("/proc/search")
				if tc.q != "" {
					u += "?" + tc.q
				}
				r, ok := assertAndGet[handler.SearchProcessResponseData](t, u)
				if !ok {
					return
				}
				got := map[string]string{}
				for _, x := range r {
					got[x.Title] = x.RequestID
				}
				want := map[string]string{}
				for _, x := range tc.basenames {
					want[x] = d[x]
				}
				assert.Equal(t, want, got)
			})
		}
	})

	t.Log("cancel")
	assertNil(t, workerProcess.Process.Signal(syscall.SIGTERM))
	assertNil(t, serverProcess.Process.Signal(syscall.SIGTERM))

	t.Log("close")
	assertNil(t, workerProcess.Wait())
	assertNil(t, serverProcess.Wait())
}
