package handler

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"path/filepath"
	"time"

	"github.com/berquerant/pneutrinoutil/cli/ctl"
	"github.com/berquerant/pneutrinoutil/pkg/alog"
	"github.com/berquerant/pneutrinoutil/pkg/domain"
	"github.com/berquerant/pneutrinoutil/pkg/echox"
	"github.com/berquerant/pneutrinoutil/pkg/logx"
	"github.com/berquerant/pneutrinoutil/pkg/repo"
	"github.com/goccy/go-yaml"
	"github.com/labstack/echo/v4"
)

func NewGet(
	processGetter repo.ProcessGetter,
	detailsGetter repo.ProcessDetailsGetter,
	objectReader repo.ObjectReader,
	objectGetter repo.ObjectGetter,
) *Get {
	return &Get{
		processGetter: processGetter,
		detailsGetter: detailsGetter,
		objectReader:  objectReader,
		objectGetter:  objectGetter,
	}
}

type GetParam struct {
	RequestID string `param:"id" validate:"required"`
}

type Get struct {
	processGetter repo.ProcessGetter
	detailsGetter repo.ProcessDetailsGetter
	objectReader  repo.ObjectReader
	objectGetter  repo.ObjectGetter
}

func (*Get) bind(c echo.Context) (*GetParam, *StatusError) {
	var p GetParam
	if err := c.Bind(&p); err != nil {
		return nil, NewStatusError(http.StatusBadRequest, err, "bad request")
	}
	return &p, nil
}

type result struct {
	requestID      string
	basename       string
	command        *string
	statusID       domain.ProcessStatus
	createdAt      time.Time
	startedAt      *time.Time
	completedAt    *time.Time
	scoreObjectID  int
	logObjectID    *int
	resultObjectID *int
}

func (g *Get) withResult(
	f func(echo.Context, *result) error,
) func(echo.Context) error {
	return func(c echo.Context) error {
		p, bErr := g.bind(c)
		if bErr != nil {
			return bErr.Respond(c)
		}

		proc, err := g.processGetter.GetProcessByRequestId(c.Request().Context(), p.RequestID)
		if err != nil {
			alog.L().Error("missing process", slog.String("id", echox.RequestID(c)), slog.String("param_id", p.RequestID), logx.Err(err))
			return Error(c, http.StatusNotFound, "not found")
		}
		details, err := g.detailsGetter.GetProcessDetails(c.Request().Context(), proc.DetailsID)
		if err != nil {
			alog.L().Error("missing details", slog.String("id", echox.RequestID(c)), slog.String("param_id", p.RequestID), logx.Err(err))
			return Error(c, http.StatusNotFound, "not found")
		}

		r := &result{
			requestID:      proc.RequestID,
			basename:       details.Title,
			command:        details.Command,
			statusID:       proc.Status,
			createdAt:      proc.CreatedAt,
			startedAt:      proc.StartedAt,
			completedAt:    proc.CompletedAt,
			scoreObjectID:  details.ScoreObjectID,
			logObjectID:    details.LogObjectID,
			resultObjectID: details.ResultObjectID,
		}
		return f(c, r)
	}
}

func (g *Get) withStorageObject(objectID int, f func(echo.Context, repo.ReadObjectResponse) error) func(echo.Context) error {
	return func(c echo.Context) error {
		r, err := g.objectReader.ReadObject(c.Request().Context(), objectID)
		if err != nil {
			alog.L().Error("missing object", slog.String("id", echox.RequestID(c)), slog.Int("objectID", objectID), logx.Err(err))
			return Error(c, http.StatusNotFound, "not found")
		}
		return f(c, r)
	}
}

func (g *Get) withStorageObjectFile(objectID int, f func(echo.Context, repo.ReadObjectResponse) error) func(echo.Context) error {
	return g.withStorageObject(objectID, func(c echo.Context, r repo.ReadObjectResponse) error {
		if _, ok := r.Storage(); !ok {
			return Error(c, http.StatusNotFound, "not found")
		}
		return f(c, r)
	})
}

func (g *Get) withStorageObjectFileBlob(objectID int, contentType string) func(echo.Context) error {
	return g.withStorageObjectFile(objectID, func(c echo.Context, r repo.ReadObjectResponse) error {
		storage, _ := r.Storage()
		blob, err := io.ReadAll(storage.Blob)
		if err != nil {
			alog.L().Error("failed to read storage object", slog.String("id", echox.RequestID(c)), slog.Int("objectID", objectID), logx.Err(err))
			return Error(c, http.StatusInternalServerError, "read blob")
		}
		return c.Blob(http.StatusOK, contentType, blob)
	})
}
func (g *Get) withResultObjectFile(objectID int, contentType, path string, f func(int, string) func(echo.Context) error) func(echo.Context) error {
	return g.withStorageObject(objectID, func(c echo.Context, r repo.ReadObjectResponse) error {
		if _, ok := r.Storage(); ok {
			alog.L().Error("result object should be a directory", slog.String("id", echox.RequestID(c)), slog.Int("objectID", objectID))
			return Error(c, http.StatusNotFound, "not found")
		}
		newPath := filepath.Join(r.Object().Path, path)
		alog.L().Debug("find result object", slog.String("id", echox.RequestID(c)), slog.Int("objectID", objectID), slog.String("path", newPath))
		got, err := g.objectGetter.GetObjectByPath(c.Request().Context(), r.Object().Bucket, newPath)
		if err != nil {
			alog.L().Error("result object not found", slog.String("id", echox.RequestID(c)), slog.String("bucket", r.Object().Bucket), slog.String("path", newPath), logx.Err(err))
			return Error(c, http.StatusNotFound, "not found")
		}
		return f(got.ID, contentType)(c)
	})
}

func (g *Get) withResultObjectFileBlob(objectID int, contentType, path string) func(echo.Context) error {
	return g.withResultObjectFile(objectID, contentType, path, g.withStorageObjectFileBlob)
}

type GetDetailResponseData struct {
	RequestID   string `json:"rid"`                // request id, or just id
	Basename    string `json:"basename,omitempty"` // original musicxml file name except extension
	Command     string `json:"command,omitempty"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at,omitempty"`
	StartedAt   string `json:"started_at,omitempty"`
	CompletedAt string `json:"completed_at,omitempty"`
}

// Get process info.
//
// @summary get process info
// @description get process info
// @param id path string true "request id"
// @produce json
// @success 200 {object} handler.SuccessResponse[GetDetailResponseData]
// @failure 404 {object} handler.ErrorResponse
// @router /proc/{id}/detail [get]
func (g *Get) Detail(c echo.Context) error {
	return g.withResult(func(c echo.Context, r *result) error {
		v := GetDetailResponseData{
			RequestID: r.requestID,
			Basename:  r.basename,
			Status:    r.statusID.String(),
			CreatedAt: r.createdAt.Format(time.DateTime),
		}
		if x := r.command; x != nil {
			v.Command = *x
		}
		if x := r.startedAt; x != nil {
			v.StartedAt = x.Format(time.DateTime)
		}
		if x := r.completedAt; x != nil {
			v.CompletedAt = x.Format(time.DateTime)
		}
		return Success(c, http.StatusOK, v)
	})(c)
}

// Download process log.
//
// @summary download log
// @description download process log file
// @param id path string true "request id"
// @success 200 {string} file
// @failure 404 {object} handler.ErrorResponse
// @router /proc/{id}/log [get]
func (g *Get) Log(c echo.Context) error {
	return g.withResult(func(c echo.Context, r *result) error {
		if objectID := r.logObjectID; objectID != nil {
			return g.withStorageObjectFileBlob(*objectID, "text/plain")(c)
		}
		return Error(c, http.StatusNotFound, "not found")
	})(c)
}

// Download pneutrinoutil config file.
//
// @summary download config
// @description download pneutrinoutil config as json
// @param id path string true "request id"
// @produce json
// @success 200 {object} handler.SuccessResponse[ctl.Config]
// @failure 404 {object} handler.ErrorResponse
// @router /proc/{id}/config [get]
func (g *Get) Config(c echo.Context) error {
	return g.withResult(func(c echo.Context, r *result) error {
		if objectID := r.resultObjectID; objectID != nil {
			return g.withResultObjectFile(*objectID, "application/json", "config.yml", func(id int, _ string) func(echo.Context) error {
				return g.withStorageObjectFile(id, func(c echo.Context, r repo.ReadObjectResponse) error {
					stor, _ := r.Storage()
					buf, err := io.ReadAll(stor.Blob)
					if err != nil {
						alog.L().Error("failed to read config blob", slog.String("id", echox.RequestID(c)), slog.Int("objectID", id), logx.Err(err))
						return Error(c, http.StatusInternalServerError, "failed to read config blob")
					}
					var cfg ctl.Config
					if err := yaml.Unmarshal(buf, &cfg); err != nil {
						alog.L().Error("failed to unmarshal config", slog.String("id", echox.RequestID(c)), slog.Int("objectID", id), logx.Err(err))
						return Error(c, http.StatusInternalServerError, "failed to unmarshal config")
					}
					return Success(c, http.StatusOK, cfg)
				})
			})(c)
		}
		return Error(c, http.StatusNotFound, "not found")
	})(c)
}

// Download musicxml.
//
// @summary download musicxml
// @description download musicxml file
// @param id path string true "request id"
// @success 200 {string} file
// @failure 404 {object} handler.ErrorResponse
// @router /proc/{id}/musicxml [get]
func (g *Get) MusicXML(c echo.Context) error {
	return g.withResult(func(c echo.Context, r *result) error {
		if objectID := r.resultObjectID; objectID != nil {
			c.Response().Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.musicxml"`, r.basename))
			return g.withResultObjectFileBlob(*objectID, "application/vnd.recordare.musicxml+xml", r.basename+".musicxml")(c)
		}
		return Error(c, http.StatusNotFound, "not found")
	})(c)
}

// Download wav file generated by pneutrinoutil.
//
// @summary download wav
// @description download wav file generated by pneutrinoutil
// @param id path string true "request id"
// @success 200 {string} file
// @failure 404 {object} handler.ErrorResponse
// @router /proc/{id}/wav [get]
func (g *Get) Wav(c echo.Context) error {
	return g.withResult(func(c echo.Context, r *result) error {
		if objectID := r.resultObjectID; objectID != nil {
			c.Response().Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.wav"`, r.basename))
			return g.withResultObjectFileBlob(*objectID, "audio/wav", r.basename+".wav")(c)
		}
		return Error(c, http.StatusNotFound, "not found")
	})(c)
}

// Download world wav file generated by pneutrinoutil.
//
// @summary download world wav
// @description download world wav file generated by pneutrinoutil (before NEUTRINO v3)
// @param id path string true "request id"
// @success 200 {string} file
// @failure 404 {object} handler.ErrorResponse
// @router /proc/{id}/world_wav [get]
func (g *Get) WorldWav(c echo.Context) error {
	return g.withResult(func(c echo.Context, r *result) error {
		if objectID := r.resultObjectID; objectID != nil {
			c.Response().Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s_world.wav"`, r.basename))
			return g.withResultObjectFileBlob(*objectID, "audio/wav", r.basename+"_world.wav")(c)
		}
		return Error(c, http.StatusNotFound, "not found")
	})(c)
}
