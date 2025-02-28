package server_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/berquerant/pneutrinoutil/server/config"
	"github.com/berquerant/pneutrinoutil/server/server"
	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	t.Run("health", func(t *testing.T) {
		s := server.New(&config.Config{
			AccessLogWriter: os.Stdout,
		})
		srv := httptest.NewServer(s.Echo())
		defer srv.Close()

		r, err := http.Get(srv.URL + "/v1/health")
		if !assert.Nil(t, err) {
			return
		}
		defer r.Body.Close()
		io.Copy(io.Discard, r.Body)
		assert.Equal(t, http.StatusOK, r.StatusCode)
	})
}
