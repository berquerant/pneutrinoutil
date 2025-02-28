package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Debug info.
//
// @summary debug info
// @description debug info
// @produce json
// @success 200 {object} handler.SuccessResponse[DebugResponseData]
// @router /debug [get]
func Debug(c echo.Context) error {
	return Success(c, http.StatusOK, DebugResponseData{
		Routes: c.Echo().Routes(),
	})
}

type DebugResponseData struct {
	Routes any `json:"routes"`
}
