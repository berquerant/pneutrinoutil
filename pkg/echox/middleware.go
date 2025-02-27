package echox

import "github.com/labstack/echo/v4"

// RequestID gets the request id from the context.
func RequestID(c echo.Context) string {
	return c.Response().Header().Get(echo.HeaderXRequestID)
}
