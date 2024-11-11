package cmd

import (
	"time"

	"github.com/berquerant/pneutrinoutil/pkg/ctl"
	"github.com/berquerant/pneutrinoutil/pkg/task"
)

func newTasks(dir *task.Dir, c *ctl.Config, now time.Time, play bool) []*task.Task {
	g := task.NewGenerator(dir, c, now, play)
	return []*task.Task{
		g.DisplayEnv(),
		g.Init(),
		g.Prepare(),
		g.MusicXMLToLabel(),
		g.NEUTRINO(),
		g.NSF(),
		g.WORLD(),
		g.Cleanup(),
	}
}
