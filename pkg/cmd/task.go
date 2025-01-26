package cmd

import (
	"time"

	"github.com/berquerant/pneutrinoutil/pkg/ctl"
	"github.com/berquerant/pneutrinoutil/pkg/task"
)

func newTaskRunner(dir *task.Dir, c *ctl.Config, now time.Time, play string) *task.Runner {
	g := task.NewGenerator(dir, c, now, play)
	tasks := []*task.Task{
		g.DisplayEnv(),
		g.Init(),
		g.MusicXMLToLabel(),
		g.NEUTRINO(),
		g.NSF(),
		g.WORLD(),
		g.Cleanup(),
	}
	return task.NewRunner(tasks)
}
