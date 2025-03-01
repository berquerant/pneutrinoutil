package logx

import (
	"io"
	"log/slog"
)

func NewTextLogger(w io.Writer, level slog.Leveler) *slog.Logger {
	handler := slog.NewTextHandler(w, &slog.HandlerOptions{
		Level: level,
	})
	return slog.New(handler)
}
