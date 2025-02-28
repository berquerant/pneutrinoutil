package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/berquerant/pneutrinoutil/pkg/echox"
	"github.com/berquerant/pneutrinoutil/pkg/logx"
	"github.com/berquerant/pneutrinoutil/pkg/wait"
	"github.com/berquerant/pneutrinoutil/server/alog"
	"github.com/berquerant/pneutrinoutil/server/config"
	"github.com/berquerant/pneutrinoutil/server/pworker"
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
// Returns the file path.
func (s Start) GetFormFile(c echo.Context) (string, *StatusError) {
	return ReadFormFile(c, "score", s.uploadDir, uploadMaxSizeBytes)
}

func NewStart(c *config.Config, w *wait.Worker) *Start {
	return &Start{
		pneutrinoutil:        c.Pneutrinoutil,
		neutrinoDir:          c.NeutrinoDir,
		pneutrinoutilWorkDir: c.PneutrinoutilWorkDir(),
		uploadDir:            c.UploadDir(),
		logDir:               c.LogDir(),
		processTimeout:       c.ProcessTimeout(),
		shell:                c.Shell,
		worker:               w,
		env: []string{
			fmt.Sprintf("HOME=%s", os.Getenv("HOME")),
			fmt.Sprintf("PWD=%s", os.Getenv("PWD")),
		},
	}
}

type Start struct {
	pneutrinoutil        string
	neutrinoDir          string
	pneutrinoutilWorkDir string
	uploadDir            string
	logDir               string
	env                  []string
	processTimeout       time.Duration
	shell                string

	worker *wait.Worker
}

func (s *Start) NewProcess(c echo.Context) *StatusError {
	score, err := s.GetFormFile(c)
	if err != nil {
		return err
	}

	rid := echox.RequestID(c)
	logFilePath := filepath.Join(s.logDir, rid)
	logFile, rErr := os.Create(logFilePath)
	if err != nil {
		return &StatusError{
			http.StatusInternalServerError,
			fmt.Errorf("%w: open %s", rErr, logFilePath),
			"failed to create log file",
		}
	}

	args := []string{
		"--desc", rid,
		"--neutrinoDir", s.neutrinoDir,
		"--workDir", s.pneutrinoutilWorkDir,
		"--score", score,
		"--env", "all",
		"--shell", s.shell,
	}
	for k, v := range s.GetFormArgs(c) {
		args = append(args, fmt.Sprintf("--%s", k), v)
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.processTimeout)
	cmd := exec.CommandContext(ctx, s.pneutrinoutil, args...)
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	cmd.Env = s.env
	if err := cmd.Start(); err != nil {
		cancel()
		return &StatusError{
			http.StatusInternalServerError,
			fmt.Errorf("%w: failed to start process", err),
			"failed to start process",
		}
	}

	alog.L().Info("start process",
		slog.String("id", rid),
		slog.Int("pid", cmd.Process.Pid),
		slog.String("bin", s.pneutrinoutil),
		slog.String("args", strings.Join(args, " ")),
		slog.String("log", logFilePath),
	)
	s.worker.Add(&pworker.Process{
		RequestID:  rid,
		Cmd:        cmd,
		CancelFunc: cancel,
		Log:        logFile,
	})

	return nil
}
