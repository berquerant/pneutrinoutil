package pworker

import (
	"io"
	"os/exec"

	"github.com/berquerant/pneutrinoutil/pkg/wait"
)

var _ wait.Waiter = &Process{}

type Process struct {
	RequestID  string
	Cmd        *exec.Cmd
	CancelFunc func()
	Log        io.Closer
}

func (p *Process) Cancel() { p.CancelFunc() }

func (p *Process) Wait() error {
	defer p.Log.Close()
	return p.Cmd.Wait()
}
