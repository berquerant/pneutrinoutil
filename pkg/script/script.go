package script

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/berquerant/execx"
)

func New(title, content string) *Script {
	return &Script{
		title:   title,
		content: content,
	}
}

const scriptShell = "bash"

type Script struct {
	title   string
	content string
	elapsed time.Duration
}

func (s Script) Title() string          { return s.title }
func (s Script) Elapsed() time.Duration { return s.elapsed }

func (s Script) String() string {
	return fmt.Sprintf(`# %s
set -ex
%s`, s.title, s.content)
}

//go:generate go run github.com/berquerant/goconfig -field "Env execx.Env|Dir string" -option -output script_config_generated.go

func (s *Script) Run(ctx context.Context, opt ...ConfigOption) error {
	startTime := time.Now()
	slog.Info("Run script", slog.String("title", s.title))
	defer func() {
		elapsedTime := time.Since(startTime)
		s.elapsed = elapsedTime
		slog.Info(
			"End run script",
			slog.String("title", s.title),
			slog.String("elapsed", elapsedTime.String()),
			slog.Float64("elapsedSeconds", elapsedTime.Seconds()),
		)
	}()

	c := NewConfigBuilder().
		Env(execx.NewEnv()).
		Dir(".").
		Build()
	c.Apply(opt...)

	x := execx.NewScript("set -ex\n"+s.content, scriptShell)
	defer x.Close()
	x.Env.Merge(execx.EnvFromEnviron())
	x.Env.Merge(c.Env.Get())

	if err := x.Runner(func(cmd *execx.Cmd) error {
		cmd.Dir = c.Dir.Get()
		_, err := cmd.Run(
			ctx,
			execx.WithStdoutConsumer(s.logConsumer(
				slog.String("dir", cmd.Dir),
				slog.String("fd", "stdout"),
			)),
			execx.WithStderrConsumer(s.logConsumer(
				slog.String("dir", cmd.Dir),
				slog.String("fd", "stderr"),
			)),
		)
		return err
	}); err != nil {
		return fmt.Errorf("%w: %s", err, s.title)
	}
	return nil
}

func (s Script) logConsumer(attr ...any) func(execx.Token) {
	attrs := append([]any{slog.String("title", s.title)}, attr...)
	return func(t execx.Token) {
		// CR as newline
		for _, x := range strings.Split(t.String(), "\r") {
			slog.Info(x, attrs...)
		}
	}
}
