package handler

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"path"
	"path/filepath"
	"time"

	"github.com/berquerant/pneutrinoutil/pkg/alog"
	"github.com/berquerant/pneutrinoutil/pkg/domain"
	"github.com/berquerant/pneutrinoutil/pkg/echox"
	"github.com/berquerant/pneutrinoutil/pkg/logx"
	"github.com/berquerant/pneutrinoutil/pkg/repo"
	"github.com/berquerant/pneutrinoutil/pkg/task"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
)

// Start a process.
//
// @summary start a process
// @description start a pneutrinoutil process with given arguments
// @param score formData file true "musicxml"
// @param enhanceBreathiness formData number false "[0, 100]%, default: 0"
// @param formantShift formData number false "default: 1.0"
// @param inference formData integer false "[2, 3, 4], default: 2"
// @param model formData string false "default: MERROW"
// @param pitchShiftNsf formData number false "default: 0"
// @param pitchShiftWorld formData number false "default: 0"
// @param smoothFormant formData number false "[0, 100]%, default: 0"
// @param smoothPitch formData number false "[0, 100]%, default: 0"
// @param styleShift formData integer false "default: 0"
// @produce json
// @success 202 {object} handler.SuccessResponse[string] "new process started"
// @header 202 {string} string x-request-id "request id, or just id"
// @failure 400 {object} handler.ErrorResponse "bad score"
// @failure 413 {object} handler.ErrorResponse "too big score"
// @failure 500 {object} handler.ErrorResponse
// @router /proc [post]
func (s *Start) Handler(c echo.Context) error {
	err := s.NewProcess(c)
	if err != nil {
		rid := echox.RequestID(c)
		alog.L().Error("failed to start process", slog.String("id", rid), logx.Err(err))
		return err.Respond(c)
	}

	return Success(c, http.StatusAccepted, "accepted")
}

// GetFormArgs extracts arguments for pneutrinoutil from the form values.
func (Start) GetFormArgs(c echo.Context) map[string]string {
	keys := []string{
		"enhanceBreathiness",
		"formantShift",
		"inference",
		"model",
		"pitchShiftNsf",
		"pitchShiftWorld",
		"smoothFormant",
		"smoothPitch",
		"styleShift",
	}
	d := map[string]string{}
	for _, k := range keys {
		if v := c.FormValue(k); v != "" {
			d[k] = v
		}
	}
	return d
}

// GetFormFile reads a musicxml file from the form file `score`.
// Returns the file content.
func (s Start) GetFormFile(c echo.Context) (*ReadFromFileResult, *StatusError) {
	r, err := ReadFormFile(c, "score", uploadMaxSizeBytes)
	if err != nil {
		return nil, err.AppendMessageToErr("failed to read score file from form")
	}
	return r, nil
}

func NewStart(
	client *asynq.Client,
	processTimeout time.Duration,
	bucket string,
	path string,
	objectWriter repo.ObjectWriter,
	detailsCreator repo.ProcessDetailsCreator,
	processCreator repo.ProcessCreator,
) *Start {
	return &Start{
		client:         client,
		processTimeout: processTimeout,
		bucket:         bucket,
		path:           path,
		objectWriter:   objectWriter,
		detailsCreator: detailsCreator,
		processCreator: processCreator,
	}
}

type Start struct {
	client         *asynq.Client
	processTimeout time.Duration
	objectWriter   repo.ObjectWriter
	detailsCreator repo.ProcessDetailsCreator
	processCreator repo.ProcessCreator
	bucket         string
	path           string
}

func (s *Start) NewProcess(c echo.Context) *StatusError {
	rid := echox.RequestID(c)
	score, fErr := s.GetFormFile(c)
	if fErr != nil {
		return fErr
	}

	obj, err := s.objectWriter.WriteObject(c.Request().Context(), &repo.WriteObjectRequest{
		Type:   domain.ObjectTypeFile,
		Bucket: s.bucket,
		Path:   filepath.Join(s.path, rid, score.Name),
		Blob:   bytes.NewBuffer(score.Blob),
	})
	if err != nil {
		return NewStatusError(http.StatusInternalServerError, err, "failed to upload score")
	}

	title := score.Name
	if ext := path.Ext(title); ext != "" {
		title = title[:len(title)-len(ext)]
	}
	details, err := s.detailsCreator.CreateProcessDetails(c.Request().Context(), &repo.CreateProcessDetailsRequest{
		Title:         title,
		ScoreObjectId: obj.Object().ID,
	})
	if err != nil {
		return NewStatusError(http.StatusInternalServerError, err, "failed to create process details")
	}

	proc, err := s.processCreator.CreateProcess(c.Request().Context(), &repo.CreateProcessRequest{
		RequestId: rid,
		DetailsId: details.ID,
		Status:    domain.ProcessStatusPending,
	})
	if err != nil {
		return NewStatusError(http.StatusInternalServerError, err, "failed to create process")
	}

	var args []string
	for k, v := range s.GetFormArgs(c) {
		args = append(args, fmt.Sprintf("--%s", k), v)
	}
	atask, err := task.NewPneutrinoutilStart(task.PneutrinoutilStartPayload{
		RequestID: rid,
		Args:      args,
	})
	if err != nil {
		return NewStatusError(http.StatusInternalServerError, err, "failed to create task")
	}

	info, err := s.client.EnqueueContext(c.Request().Context(), atask, asynq.Timeout(s.processTimeout))
	if err != nil {
		return NewStatusError(http.StatusInternalServerError, err, "failed to enqueue task")
	}

	alog.L().Info("new task enqueued",
		slog.String("id", rid), slog.Int("processID", proc.ID), slog.String("taskID", info.ID),
	)

	return nil
}
