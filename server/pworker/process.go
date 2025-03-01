package pworker

import (
	"io"
	"log/slog"
	"os/exec"

	"github.com/berquerant/pneutrinoutil/pkg/logx"
	"github.com/berquerant/pneutrinoutil/pkg/wait"
	"github.com/berquerant/pneutrinoutil/server/alog"
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
	logRid := slog.String("id", p.RequestID)
	alog.L().Info("starting process", logRid)

	defer p.Log.Close()
	if err := p.Cmd.Start(); err != nil {
		alog.L().Error("failed to start process", logRid, logx.Err(err))
		return err
	}

	alog.L().Info("started process", logRid, slog.Int("pid", p.Cmd.Process.Pid))
	return p.Cmd.Wait()
}
