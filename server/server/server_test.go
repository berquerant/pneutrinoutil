package server_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/berquerant/pneutrinoutil/cli/ctl"
	"github.com/berquerant/pneutrinoutil/server/alog"
	"github.com/berquerant/pneutrinoutil/server/config"
	"github.com/berquerant/pneutrinoutil/server/handler"
	"github.com/berquerant/pneutrinoutil/server/server"
	"github.com/berquerant/structconfig"
	"github.com/stretchr/testify/assert"
)

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

func newExecutableFile(name, content string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := os.Chmod(f.Name(), 0755); err != nil {
		return err
	}
	if _, err := fmt.Fprint(f, content); err != nil {
		return err
	}
	return nil
}

func TestServer(t *testing.T) {
	alog.Setup(os.Stdout, slog.LevelDebug)

	const (
		mockCLISource                      = "../../mockcli/"
		defaultScoreContent                = "some musicxml content"
		defaultScoreFileNameBase           = "defaultscore.musicxml"
		defaultScoreFileNameBaseWithoutExt = "defaultscore"
	)
	var (
		rootTempDir          = t.TempDir()
		mockWorkDir          = rootTempDir
		mockCLIBin           = filepath.Join(rootTempDir, "mockclibin")
		mockCLI              = filepath.Join(rootTempDir, "mockcli")
		mockCLILog           = filepath.Join(rootTempDir, "mockcli.log")
		defaultScoreFileName = filepath.Join(rootTempDir, defaultScoreFileNameBase)
		// notify               = filepath.Join(rootTempDir, "notify")
		// notifyJSON           = filepath.Join(rootTempDir, "notify.json")

		mockCLIScript = fmt.Sprintf(`#!/bin/bash
cat <<EOS >> "%[2]s"
---------- mockCLIScript ----------
ARGS=$*
----------
EOS
%[1]s "$@" >> "%[2]s"`,
			mockCLIBin, mockCLILog)

		// 		notifyScript = fmt.Sprintf(`#!/bin/bash
		// cat <<EOS
		// ---------- notifyScript ----------
		// ARGS=$*
		// ----------
		// EOS
		// cat <<EOS > "%s"
		// {
		//   "rid": "${REQUEST_ID}",
		//   "status": ${STATUS}
		// }
		// EOS`, notifyJSON)

		removeMockCLILog = func() {
			os.Remove(mockCLILog)
		}
		readMockCLILog = func() {
			f, err := os.Open(mockCLILog)
			if err != nil {
				return
			}
			defer f.Close()
			t.Log("---------- mockcli log start ----------")
			io.Copy(os.Stdout, f)
			t.Log("---------- mockcli log end ----------")
		}
	)

	// prepare mockcli
	if err := exec.Command("go", "build", "-o", mockCLIBin, mockCLISource).Run(); err != nil {
		t.Fatal(err)
	}
	if err := newExecutableFile(mockCLI, mockCLIScript); err != nil {
		t.Fatal(err)
	}
	// create score file
	{
		f, err := os.Create(defaultScoreFileName)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := f.WriteString(defaultScoreContent); err != nil {
			t.Fatal(err)
		}
		f.Close()
	}

	// prepare notify script
	// if err := newExecutableFile(notify, notifyScript); err != nil {
	// 	t.Fatal(err)
	// }
	// type NotifyResult struct {
	// 	RequestID string `json:"rid"`
	// 	Status    int    `json:"status"`
	// }
	// readNotifyResult := func() (*NotifyResult, error) {
	// 	f, err := os.Open(notifyJSON)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	b, err := io.ReadAll(f)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	var x NotifyResult
	// 	if err := json.Unmarshal(b, &x); err != nil {
	// 		return nil, err
	// 	}
	// 	return &x, nil
	// }

	newDefaultConfig := func() *config.Config {
		var c config.Config
		if err := structconfig.New[config.Config]().FromDefault(&c); err != nil {
			t.Fatal(err)
		}
		return &c
	}

	// create server config
	c := newDefaultConfig()
	c.AccessLogFile = config.AccessLogStderr
	c.WorkDir = mockWorkDir
	c.Pneutrinoutil = mockCLI
	// c.NotificationCommand = notify
	// c.NotificationLogFile = config.NotificationLogStderr
	{
		err := c.Init()
		if err != nil {
			t.Fatal(err)
		}
		defer c.Close()
	}

	// launch server
	srv := server.New(c)
	ctx, cancel := context.WithCancel(context.Background())
	go srv.StartLoop(ctx)
	testSrv := httptest.NewServer(srv.Echo())

	// cleanup
	defer func() {
		cancel()
		testSrv.Close()
		srv.Wait()
	}()

	newURL := func(p string) string {
		return testSrv.URL + p
	}

	t.Run("health", func(t *testing.T) {
		r, err := http.Get(newURL("/v1/health"))
		if !assert.Nil(t, err) {
			return
		}
		r.Body.Close()
		assert.Equal(t, http.StatusOK, r.StatusCode)
	})

	t.Run("version", func(t *testing.T) {
		r, err := http.Get(newURL("/v1/version"))
		if !assert.Nil(t, err) {
			return
		}
		r.Body.Close()
		assert.Equal(t, http.StatusOK, r.StatusCode)
	})

	t.Run("debug", func(t *testing.T) {
		r, err := http.Get(newURL("/v1/debug"))
		if !assert.Nil(t, err) {
			return
		}
		r.Body.Close()
		assert.Equal(t, http.StatusOK, r.StatusCode)
	})

	assertListProcess := func(t *testing.T, want handler.ListResponseData) {
		t.Helper()
		data, ok := assertAndGet[handler.ListResponseData](t, newURL("/v1/proc"))
		if !ok {
			return
		}
		// assert only request ids
		wantRid := make([]string, len(want))
		for i, x := range want {
			wantRid[i] = x.RequestID
		}
		sort.Strings(wantRid)
		gotRid := make([]string, len(data))
		for i, x := range data {
			gotRid[i] = x.RequestID
		}
		sort.Strings(gotRid)
		assert.Equal(t, wantRid, gotRid)
	}

	t.Run("list proc empty", func(t *testing.T) {
		assertListProcess(t, handler.ListResponseData([]*handler.ListResponseDataElement{}))
	})

	t.Run("new process", func(t *testing.T) {
		removeMockCLILog()

		// build request body
		var body bytes.Buffer
		w := multipart.NewWriter(&body)
		score, err := os.Open(defaultScoreFileName)
		if !assert.Nil(t, err) {
			return
		}
		defer score.Close()
		fw, err := w.CreateFormFile("score", score.Name())
		if !assert.Nil(t, err) {
			return
		}
		if _, err = io.Copy(fw, score); !assert.Nil(t, err) {
			return
		}
		if !assert.Nil(t, w.Close()) {
			return
		}
		// build request
		req, err := http.NewRequest(http.MethodPost, newURL("/v1/proc"), &body)
		if !assert.Nil(t, w.Close()) {
			return
		}
		req.Header.Set("Content-Type", w.FormDataContentType())
		// send request
		resp, err := http.DefaultClient.Do(req)
		if !assert.Nil(t, err) {
			return
		}
		// assert response
		_ = resp.Body.Close()
		if !assert.Equal(t, http.StatusAccepted, resp.StatusCode) {
			return
		}
		rid := resp.Header.Get("x-request-id")
		if !assert.NotEqual(t, "", rid) {
			return
		}
		// wait until process completed
		time.Sleep(500 * time.Millisecond)
		readMockCLILog()
		// assert notify result
		// nr, err := readNotifyResult()
		// if !assert.Nil(t, err) {
		// 	return
		// }
		// assert.Equal(t, rid, nr.RequestID)
		// assert.Equal(t, 0, nr.Status, "want success")

		t.Run("process added", func(t *testing.T) {
			assertListProcess(t, handler.ListResponseData([]*handler.ListResponseDataElement{
				{
					RequestID: rid,
				},
			}))
		})

		t.Run("download config", func(t *testing.T) {
			want, err := ctl.NewDefaultConfig()
			if !assert.Nil(t, err) {
				return
			}
			got, ok := assertAndGet[ctl.Config](t, newURL("/v1/proc/"+rid+"/config"))
			if !ok {
				return
			}
			want.Description = rid
			if !assert.Equal(t, filepath.Base(got.Score), filepath.Base(defaultScoreFileName)) {
				return
			}
			want.Score = got.Score // enough if only basename matches
			assert.Equal(t, want, &got)
		})

		t.Run("get process detail", func(t *testing.T) {
			got, ok := assertAndGet[handler.GetDetailResponseData](t, newURL("/v1/proc/"+rid+"/detail"))
			if !ok {
				return
			}
			assert.Equal(t, rid, got.RequestID)
			assert.Equal(t, defaultScoreFileNameBaseWithoutExt, got.Basename)
			assert.Equal(t, "", got.Err)
		})

		t.Run("download process log", func(t *testing.T) {
			r, err := http.Get(newURL("/v1/proc/" + rid + "/log"))
			if !assert.Nil(t, err) {
				return
			}
			r.Body.Close()
			assert.Equal(t, http.StatusOK, r.StatusCode)
		})

		t.Run("download musicxml", func(t *testing.T) {
			r, err := http.Get(newURL("/v1/proc/" + rid + "/musicxml"))
			if !assert.Nil(t, err) {
				return
			}
			defer r.Body.Close()
			assert.Equal(t, http.StatusOK, r.StatusCode)
			body, err := io.ReadAll(r.Body)
			if !assert.Nil(t, err) {
				return
			}
			assert.Equal(t, defaultScoreContent, string(body))
		})

		t.Run("download wav", func(t *testing.T) {
			r, err := http.Get(newURL("/v1/proc/" + rid + "/wav"))
			if !assert.Nil(t, err) {
				return
			}
			r.Body.Close()
			assert.Equal(t, http.StatusOK, r.StatusCode)
		})

		t.Run("download world wav", func(t *testing.T) {
			r, err := http.Get(newURL("/v1/proc/" + rid + "/world_wav"))
			if !assert.Nil(t, err) {
				return
			}
			r.Body.Close()
			assert.Equal(t, http.StatusOK, r.StatusCode)
		})
	})
}
