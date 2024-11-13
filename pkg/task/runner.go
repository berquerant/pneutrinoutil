package task

import (
	"context"
	"strings"
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

func (r Runner) String() string {
	ss := make([]string, len(r.tasks))
	for i, t := range r.tasks {
		ss[i] = t.String()
	}
	return strings.Join(ss, "\n")
}

func (r *Runner) Run(ctx context.Context) error {
	stats := []Stat{}
	for _, t := range r.tasks {
		if err := t.Run(ctx); err != nil {
			return err
		}
		stats = append(stats, NewStat(t.Title(), t.Elapsed()))
	}
	r.stats = stats
	return nil
}
