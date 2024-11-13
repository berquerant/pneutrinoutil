package task

import (
	"context"
	"time"

	"github.com/berquerant/pneutrinoutil/pkg/script"
)

type Task struct {
	s    *script.Script
	opts []script.ConfigOption
}

func (t Task) Title() string          { return t.s.Title() }
func (t Task) Elapsed() time.Duration { return t.s.Elapsed() }
func (t Task) String() string         { return t.s.String() }

func (t *Task) Run(ctx context.Context) error {
	return t.s.Run(ctx, t.opts...)
}
