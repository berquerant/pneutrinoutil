package wait_test

import (
	"errors"
	"testing"
	"time"

	"github.com/berquerant/pneutrinoutil/pkg/wait"
	"github.com/stretchr/testify/assert"
)

type proc struct {
	id    int
	sleep time.Duration
	err   error
}

func (p *proc) Wait() error {
	time.Sleep(p.sleep)
	return p.err
}

func (proc) Cancel() {}

type cancelProc struct {
	id int
	c  chan struct{}
}

func (p *cancelProc) Wait() error {
	<-p.c
	return nil
}

func (p *cancelProc) Cancel() {
	close(p.c)
}

func newCancelProc(id int) *cancelProc {
	return &cancelProc{
		id: id,
		c:  make(chan struct{}),
	}
}

func TestWorker(t *testing.T) {
	t.Run("cancel", func(t *testing.T) {
		w := wait.New()
		p1 := newCancelProc(1)

		type res struct {
			id int
		}
		want := []res{
			{id: 1},
		}

		go func() {
			w.Add(p1)
			w.WaitAndClose()
		}()

		go func() {
			time.Sleep(300 * time.Millisecond)
			w.Cancel()
		}()

		got := []res{}
		for r := range w.DoneC() {
			x := r.Waiter.(*cancelProc)
			got = append(got, res{id: x.id})
		}

		assert.Equal(t, want, got)
	})

	t.Run("jobs", func(t *testing.T) {
		w := wait.New()
		p1 := &proc{
			id:    1,
			sleep: 100 * time.Millisecond,
		}
		p2 := &proc{
			id:    2,
			sleep: 500 * time.Millisecond,
		}
		pErr := errors.New("proc error")
		p3 := &proc{
			id:    3,
			sleep: 300 * time.Millisecond,
			err:   pErr,
		}

		type res struct {
			id  int
			err error
		}
		want := []res{
			{id: 1},
			{id: 3, err: pErr},
			{id: 2},
		}

		go func() {
			w.Add(p1)
			w.Add(p2)
			w.Add(p3)
			w.WaitAndClose()
		}()

		got := []res{}
		for r := range w.DoneC() {
			x := r.Waiter.(*proc)
			got = append(got, res{id: x.id, err: r.Err})
		}

		assert.Equal(t, want, got)
	})
}
