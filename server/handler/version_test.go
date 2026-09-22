package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/berquerant/pneutrinoutil/pkg/version"
	"github.com/berquerant/pneutrinoutil/server/handler"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	for _, tc := range []struct {
		title        string
		target       string
		wantStatus   int
		wantVersion  string
		wantRevision string
	}{
		{
			title:        "GET /version returns valid version response",
			target:       "/version",
			wantStatus:   http.StatusOK,
			wantVersion:  version.Version,
			wantRevision: version.Revision,
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, tc.target, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.Version(c)
			assert.Nil(t, err)
			assert.Equal(t, tc.wantStatus, rec.Code)

			var resp handler.SuccessResponse[handler.VersionResponseData]
			assert.Nil(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			assert.True(t, resp.OK)
			assert.Equal(t, tc.wantVersion, resp.Data.Version)
			assert.Equal(t, tc.wantRevision, resp.Data.Revision)
		})
	}
}
