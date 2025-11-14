package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
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

	eventuallyWaitForMax = time.Second * 20
	eventuallyTick       = time.Millisecond * 300
)

func assertNil[T any](t assert.TestingT, v T) bool {
	return assert.Nil(t, v, "%v should be nil", v)
}

func assertAndGet[T any](t assert.TestingT, url string) (T, bool) {
	var d T
	r, err := http.Get(url)
	if !assertNil(t, err) {
		return d, false
	}
	defer r.Body.Close()
	if !assert.Equal(t, http.StatusOK, r.StatusCode) {
		return d, false
	}
	b, err := io.ReadAll(r.Body)
	if !assertNil(t, err) {
		return d, false
	}
	var body handler.SuccessResponse[T]
	if !assertNil(t, json.Unmarshal(b, &body)) {
		return d, false
	}
	if !assert.True(t, body.OK) {
		return d, false
	}
	return body.Data, true
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
		mysqlDSN = fmt.Sprintf("test:test@tcp(%s:%s)/test?parseTime=true",
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

	eventually(t, func(c *assert.CollectT) {
		got, ok := assertAndGet[handler.GetDetailResponseData](c, newUrl("/proc/"+newRid+"/detail"))
		assert.True(c, ok)
		assert.Equal(c, "succeed", got.Status)
		assert.Equal(c, newRid, got.RequestID)
		assert.Equal(c, basename, got.Basename)
	}, "wait process completed")

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

	t.Log("cancel")
	assertNil(t, workerProcess.Process.Signal(syscall.SIGTERM))
	assertNil(t, serverProcess.Process.Signal(syscall.SIGTERM))

	t.Log("close")
	assertNil(t, workerProcess.Wait())
	assertNil(t, serverProcess.Wait())
}
