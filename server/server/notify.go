package server

import (
	"context"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"github.com/berquerant/pneutrinoutil/pkg/logx"
	"github.com/berquerant/pneutrinoutil/server/alog"
)

func (s *Server) notify(ctx context.Context, requestID string, success bool) {
	c := s.c.NotificationCommand
	if c == "" {
		return
	}

	iCtx, cancel := context.WithTimeout(ctx, s.c.NotificationTimeout())
	defer cancel()

	status := func() string {
		if success {
			return "0"
		}
		return "1"
	}()

	cmd := exec.CommandContext(iCtx, c, requestID, status)
	cmd.Stdout = s.c.NotificationLogWriter
	cmd.Stderr = s.c.NotificationLogWriter
	cmd.Env = []string{
		"HOME=" + os.Getenv("HOME"),
		"PWD=" + os.Getenv("PWD"),
		"REQUEST_ID=" + requestID,
		"STATUS=" + status,
	}

	logRid := slog.String("id", requestID)
	logCmd := slog.String("command", strings.Join(cmd.Args, " "))
	alog.L().Info("starting notification", logRid, logCmd)
	if err := cmd.Start(); err != nil {
		alog.L().Error("failed to start notification", logRid, logx.Err(err))
		return
	}
	alog.L().Info("started notification", logRid, slog.Int("pid", cmd.Process.Pid))

	if err := cmd.Wait(); err != nil {
		alog.L().Error("failed notification", logRid, logx.Err(err))
		return
	}
	alog.L().Info("end notification", logRid)
}
