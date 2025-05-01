package alog

import (
	"io"
	"log/slog"

	"github.com/berquerant/pneutrinoutil/pkg/logx"
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type NopLogger struct{}

func (NopLogger) Debug(_ string, _ ...any) {}
func (NopLogger) Info(_ string, _ ...any)  {}
func (NopLogger) Warn(_ string, _ ...any)  {}
func (NopLogger) Error(_ string, _ ...any) {}

func L() Logger {
	return instance
}

var instance Logger = &NopLogger{}

func Setup(w io.Writer, level slog.Leveler) {
	instance = logx.NewLTSVLogger(w, level)
}
