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
		var body bytes.Buffer
		w := multipart.NewWriter(&body)
		fw, err := w.CreateFormFile("score", scoreFileName)
		if !assertNil(t, err) {
			return
		}
		if _, err := fmt.Fprintf(fw, scoreContent); !assertNil(t, err) {
			return
		}
		if !assertNil(t, w.Close()) {
			return
		}
		req, err := http.NewRequest(http.MethodPost, newUrl("/proc"), &body)
		if !assertNil(t, err) {
			return
		}
		req.Header.Set("content-type", w.FormDataContentType())
		resp, err := http.DefaultClient.Do(req)
		if !assertNil(t, err) {
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
