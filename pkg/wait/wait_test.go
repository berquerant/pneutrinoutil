package wait_test

import (
	"errors"
	"fmt"
	"slices"
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
	for _, n := range []int{0, 1, 2, 100} {
		t.Run(fmt.Sprintf("concurrenty%d", n), func(t *testing.T) {

			t.Run("cancel", func(t *testing.T) {
				t.Parallel()
				w := wait.New(n)
				p1 := newCancelProc(1)

				type res struct {
					id int
				}
				want := []res{
					{id: 1},
				}

				go func() {
					w.Add(p1)
					time.Sleep(300 * time.Millisecond)
					w.Cancel()
					w.WaitAndClose()
				}()

				got := []res{}
				for r := range w.DoneC() {
					x := r.Waiter.(*cancelProc)
					got = append(got, res{id: x.id})
				}

				assert.Equal(t, want, got)
			})

			t.Run("jobs", func(t *testing.T) {
				w := wait.New(n)
				p1 := &proc{
					id:    1,
					sleep: 150 * time.Millisecond,
				}
				p2 := &proc{
					id:    2,
					sleep: 50 * time.Millisecond,
				}
				pErr := errors.New("proc error")
				p3 := &proc{
					id:    3,
					sleep: 100 * time.Millisecond,
					err:   pErr,
				}

				type res struct {
					id  int
					err error
				}
				want := []res{
					{id: 1},
					{id: 2},
					{id: 3, err: pErr},
				}

				go func() {
					w.Add(p1)
					w.Add(p2)
					w.Add(p3)
					time.Sleep(500 * time.Millisecond)
					w.WaitAndClose()
				}()

				got := []res{}
				for r := range w.DoneC() {
					x := r.Waiter.(*proc)
					got = append(got, res{id: x.id, err: r.Err})
				}
				slices.SortFunc(got, func(a, b res) int { return a.id - b.id })
				assert.Equal(t, want, got)
			})

		})

	}
}
