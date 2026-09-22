package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/berquerant/pneutrinoutil/server/handler"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	for _, tc := range []struct {
		title      string
		target     string
		wantStatus int
		wantData   string
	}{
		{
			title:      "GET /health returns OK",
			target:     "/health",
			wantStatus: http.StatusOK,
			wantData:   "OK",
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, tc.target, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.Health(c)
			assert.Nil(t, err)
			assert.Equal(t, tc.wantStatus, rec.Code)

			var resp handler.SuccessResponse[string]
			assert.Nil(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			assert.True(t, resp.OK)
			assert.Equal(t, tc.wantData, resp.Data)
		})
	}
}
