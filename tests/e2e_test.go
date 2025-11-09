package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/berquerant/pneutrinoutil/cli/ctl"
	"github.com/berquerant/pneutrinoutil/server/handler"
	"github.com/stretchr/testify/assert"
)

func run(t *testing.T, executable string, args ...string) error {
	t.Helper()
	t.Logf("run[%s %s]", executable, args)
	x := exec.Command(executable, args...)
	x.Stdout = os.Stdout
	x.Stderr = os.Stdin
	return x.Run()
}

func compile(t *testing.T, dst, src string) error {
	return run(t, "go", "build", "-o", dst, src)
}

func TestE2E(t *testing.T) {
	var (
		mockCliBin = filepath.Join(t.TempDir(), "cli")
		workerBin  = filepath.Join(t.TempDir(), "worker")
		serverBin  = filepath.Join(t.TempDir(), "server")
	)
	if !(assert.Nil(t, compile(t, mockCliBin, "../mockcli/main.go"), "compile cli") &&
		assert.Nil(t, compile(t, workerBin, "../worker/main.go"), "compile worker") &&
		assert.Nil(t, compile(t, serverBin, "../server/main.go"), "compile server")) {
		return
	}

	const (
		serverHost    = "0.0.0.0"
		serverPort    = "9101"
		storageBucket = "test"
	)

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

	if !assert.Nil(t, run(t, "../bin/redis.sh", "-n", redisDB, "flushdb"), "flush redis") {
		return
	}

	t.Log("start worker")
	workerProcess := exec.Command(
		workerBin,
		"--redisDSN", redisDSN,
		"--mysqlDSN", mysqlDSN,
		"--storageS3",
		"--storageBucket", storageBucket,
		"--pneutrinoutil", mockCliBin,
		"--workDir", workerWorkDir,
		"--debug",
	)
	workerProcess.Stdout = os.Stdout
	workerProcess.Stderr = os.Stderr

	if !assert.Nil(t, workerProcess.Start()) {
		return
	}

	t.Log("start server")
	serverProcess := exec.Command(
		serverBin,
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
	if !assert.Nil(t, serverProcess.Start()) {
		return
	}

	newUrl := func(path string) string {
		return fmt.Sprintf("http://%s:%s/v1%s", serverHost, serverPort, path)
	}

	t.Run("health", func(t *testing.T) {
		for i := range 10 {
			t.Logf("health check[%d]", i)
			if i > 0 {
				time.Sleep(300 * time.Millisecond)
			}
			r, err := http.Get(newUrl("/health"))
			if err != nil {
				continue
			}
			r.Body.Close()
			if r.StatusCode == http.StatusOK {
				return
			}
		}
		assert.Fail(t, "healthcheck")
	})

	t.Run("version", func(t *testing.T) {
		r, err := http.Get(newUrl("/version"))
		if !assert.Nil(t, err) {
			return
		}
		r.Body.Close()
		assert.Equal(t, http.StatusOK, r.StatusCode)
	})

	t.Run("debug", func(t *testing.T) {
		r, err := http.Get(newUrl("/debug"))
		if !assert.Nil(t, err) {
			return
		}
		r.Body.Close()
		assert.Equal(t, http.StatusOK, r.StatusCode)
	})

	const (
		scoreFileName = "score_sample.musicxml"
		basename      = "score_sample"
		scoreContent  = "score content"
	)

	var newRid string
	t.Run("new process", func(t *testing.T) {
		var body bytes.Buffer
		w := multipart.NewWriter(&body)
		fw, err := w.CreateFormFile("score", scoreFileName)
		if !assert.Nil(t, err) {
			return
		}
		if _, err := fmt.Fprintf(fw, scoreContent); !assert.Nil(t, err) {
			return
		}
		if !assert.Nil(t, w.Close()) {
			return
		}
		req, err := http.NewRequest(http.MethodPost, newUrl("/proc"), &body)
		if !assert.Nil(t, err) {
			return
		}
		req.Header.Set("content-type", w.FormDataContentType())
		resp, err := http.DefaultClient.Do(req)
		if !assert.Nil(t, err) {
			return
		}
		_ = resp.Body.Close()
		if !assert.Equal(t, http.StatusAccepted, resp.StatusCode) {
			return
		}
		rid := resp.Header.Get("x-request-id")
		assert.NotEqual(t, "", rid)
		newRid = rid
	})

	t.Run("wait process completed", func(t *testing.T) {
		for i := range 20 {
			if i > 0 {
				time.Sleep(300 * time.Millisecond)
			}
			got, ok := assertAndGet[handler.GetDetailResponseData](t, newUrl("/proc/"+newRid+"/detail"))
			if !ok {
				continue
			}
			if got.Status != "succeed" {
				continue
			}
			assert.Equal(t, newRid, got.RequestID)
			assert.Equal(t, basename, got.Basename)
			return
		}
	})

	t.Run("process added", func(t *testing.T) {
		data, ok := assertAndGet[handler.ListProcessResponseData](t, newUrl("/proc"))
		if !assert.True(t, ok) {
			return
		}
		if !assert.Equal(t, 1, data.Len()) {
			return
		}
		got := data[0]
		assert.Equal(t, newRid, got.RequestID)
	})

	t.Run("find process by title", func(t *testing.T) {
		for _, tc := range []struct {
			title  string
			prefix string
			hit    int
		}{
			{
				title:  "not found",
				prefix: "a" + basename,
				hit:    0,
			},
			{
				title:  "hit",
				prefix: basename,
				hit:    1,
			},
			{
				title:  "hit by prefix",
				prefix: basename[:4],
				hit:    1,
			},
		} {
			t.Run(tc.title, func(t *testing.T) {
				data, ok := assertAndGet[handler.ListProcessResponseData](t, newUrl("/proc/title?prefix="+tc.prefix))
				if !assert.True(t, ok) {
					return
				}
				if !assert.Equal(t, tc.hit, data.Len()) {
					return
				}
			})
		}
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
		if !assert.Nil(t, err) {
			return
		}
		r.Body.Close()
		assert.Equal(t, http.StatusOK, r.StatusCode)
	})

	t.Run("download musicxml", func(t *testing.T) {
		r, err := http.Get(newUrl("/proc/" + newRid + "/musicxml"))
		if !assert.Nil(t, err) {
			return
		}
		defer r.Body.Close()
		assert.Equal(t, http.StatusOK, r.StatusCode)
		body, err := io.ReadAll(r.Body)
		if !assert.Nil(t, err) {
			return
		}
		assert.Equal(t, scoreContent, string(body))
	})

	t.Run("download wav", func(t *testing.T) {
		r, err := http.Get(newUrl("/proc/" + newRid + "/wav"))
		if !assert.Nil(t, err) {
			return
		}
		r.Body.Close()
		assert.Equal(t, http.StatusOK, r.StatusCode)
	})

	t.Run("download world wav", func(t *testing.T) {
		r, err := http.Get(newUrl("/proc/" + newRid + "/world_wav"))
		if !assert.Nil(t, err) {
			return
		}
		r.Body.Close()
		assert.Equal(t, http.StatusOK, r.StatusCode)
	})

	t.Log("cancel")
	assert.Nil(t, workerProcess.Process.Signal(syscall.SIGTERM))
	assert.Nil(t, serverProcess.Process.Signal(syscall.SIGTERM))

	t.Log("close")
	assert.Nil(t, workerProcess.Wait())
	assert.Nil(t, serverProcess.Wait())
}

func assertAndGet[T any](t *testing.T, url string) (T, bool) {
	t.Helper()
	var d T
	r, err := http.Get(url)
	if !assert.Nil(t, err) {
		return d, false
	}
	defer r.Body.Close()
	if !assert.Equal(t, http.StatusOK, r.StatusCode) {
		return d, false
	}
	b, err := io.ReadAll(r.Body)
	if !assert.Nil(t, err) {
		return d, false
	}
	var body handler.SuccessResponse[T]
	if !assert.Nil(t, json.Unmarshal(b, &body)) {
		return d, false
	}
	if !assert.True(t, body.OK) {
		return d, false
	}
	return body.Data, true
}
