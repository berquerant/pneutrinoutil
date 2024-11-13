package task

import (
	"context"
	"log/slog"
	"strings"

	"github.com/berquerant/pneutrinoutil/pkg/set"
)

//go:generate go run github.com/berquerant/dataclass -type Stat -field "Title string|Elapsed time.Duration" -output stat_dataclass_generated.go

func NewRunner(tasks []*Task) *Runner {
	return &Runner{
		tasks: tasks,
	}
}

type Runner struct {
	tasks []*Task
	stats []Stat
}

func (r Runner) Stats() []Stat { return r.stats }

func (r Runner) Tasks() []string {
	titles := make([]string, len(r.tasks))
	for i, t := range r.tasks {
		titles[i] = t.Title()
	}
	return titles
}

func (r Runner) String() string {
	ss := make([]string, len(r.tasks))
	for i, t := range r.tasks {
		ss[i] = t.String()
	}
	return strings.Join(ss, "\n")
}

//go:generate go run github.com/berquerant/goconfig -configOption Option -field "Include []string|Exclude []string" -option -output runner_config_generated.go

func (r *Runner) Run(ctx context.Context, opt ...Option) error {
	c := NewConfigBuilder().
		Include([]string{}).
		Exclude([]string{}).
		Build()
	c.Apply(opt...)

	var (
		include = set.New(c.Include.Get())
		exclude = set.New(c.Exclude.Get())
		accept  = func(v string) bool {
			return (include.Len() == 0 || include.In(v)) && !exclude.In(v)
		}
	)

	stats := []Stat{}
	for _, t := range r.tasks {
		if !accept(t.Title()) {
			slog.Info("Skip", slog.String("title", t.Title()))
			continue
		}
		if err := t.Run(ctx); err != nil {
			return err
		}
		stats = append(stats, NewStat(t.Title(), t.Elapsed()))
	}
	r.stats = stats
	return nil
}
