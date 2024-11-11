package script

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

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
}

func (s Script) String() string {
	return fmt.Sprintf(`# %s
set -ex
%s`, s.title, s.content)
}

//go:generate go run github.com/berquerant/goconfig -field "Env execx.Env|Dir string" -option -output script_config_generated.go

func (s Script) Run(ctx context.Context, opt ...ConfigOption) error {
	c := NewConfigBuilder().
		Env(execx.NewEnv()).
		Dir(".").
		Build()
	c.Apply(opt...)

	x := execx.NewScript("set -ex\n"+s.content, scriptShell)
	defer x.Close()
	x.Env.Merge(execx.EnvFromEnviron())
	x.Env.Merge(c.Env.Get())

	return x.Runner(func(cmd *execx.Cmd) error {
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
	})
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