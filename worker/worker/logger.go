package worker

import (
	"fmt"

	"github.com/berquerant/pneutrinoutil/pkg/alog"
	"github.com/hibiken/asynq"
)

var _ asynq.Logger = &AsynqLogger{}

func NewAsynqLogger(logger alog.Logger) *AsynqLogger {
	return &AsynqLogger{
		logger,
	}
}

type AsynqLogger struct {
	alog.Logger
}

func (a *AsynqLogger) Debug(args ...any) {
	if len(args) > 0 {
		a.Logger.Debug(fmt.Sprint(args[0]), args[1:]...)
	}
}

func (a *AsynqLogger) Info(args ...any) {
	if len(args) > 0 {
		a.Logger.Info(fmt.Sprint(args[0]), args[1:]...)
	}
}

func (a *AsynqLogger) Warn(args ...any) {
	if len(args) > 0 {
		a.Logger.Warn(fmt.Sprint(args[0]), args[1:]...)
	}
}

func (a *AsynqLogger) Error(args ...any) {
	if len(args) > 0 {
		a.Logger.Error(fmt.Sprint(args[0]), args[1:]...)
	}
}
func (a *AsynqLogger) Fatal(args ...any) {
	a.Error(args...)
}
