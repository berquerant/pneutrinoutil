package handler

import (
	"net/http"

	"github.com/berquerant/pneutrinoutil/pkg/version"
	"github.com/labstack/echo/v4"
)

// Get server version.
//
// @summary get server version
// @description get server version
// @produce json
// @success 200 {object} handler.SuccessResponse[VersionResponseData]
// @router /version [get]
func Version(c echo.Context) error {
	return Success(c, http.StatusOK, VersionResponseData{
		Version:  version.Version,
		Revision: version.Revision,
	})
}

type VersionResponseData struct {
	Version  string // server version
	Revision string // commit hash
}
