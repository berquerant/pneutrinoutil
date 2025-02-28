package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Health returns OK.
//
// @summary health check
// @description health check
// @produce json
// @success 200 {object} handler.SuccessResponse[string]
// @router /health [get]
func Health(c echo.Context) error {
	return Success(c, http.StatusOK, "OK")
}
