package wait

import (
	"sync"
	"sync/atomic"

	"github.com/berquerant/pneutrinoutil/pkg/syncx"
	"github.com/berquerant/pneutrinoutil/pkg/uuid"
)

type Waiter interface {
	// Wait blocks until done.
	Wait() error
	// Cancel cancels the job.
	Cancel()
}

type Result struct {
	Waiter Waiter
	Err    error
}

// Worker waits [Waiter] and sends results.
type Worker struct {
	wg     sync.WaitGroup
	doneC  chan *Result
	jobs   *syncx.Map[string, Waiter]
	closed atomic.Bool
}

const (
	workerDoneChannelSize = 100
)

// New returns a new [Worker].
func New() *Worker {
	var w Worker
	w.jobs = syncx.NewMap[string, Waiter]()
	w.doneC = make(chan *Result, workerDoneChannelSize)
	return &w
}

// DoneC returns the channel to retrieve completed jobs.
func (w *Worker) DoneC() <-chan *Result { return w.doneC }

// Cancel cancels all jobs.
func (w *Worker) Cancel() {
	w.closed.Store(true)
	w.jobs.WalkWithLock(func(_ string, value Waiter) (Waiter, bool) {
		value.Cancel()
		return value, false
	})
}

// WaitAndClose waits all jobs until done and close doneC.
func (w *Worker) WaitAndClose() {
	w.closed.Store(true)
	w.wg.Wait()
	close(w.doneC)
}

// Add adds a job.
func (w *Worker) Add(waiter Waiter) {
	if w.closed.Load() {
		return
	}

	id := uuid.New()
	w.wg.Add(1)
	w.jobs.Set(id, waiter)

	go func() {
		defer func() {
			w.jobs.Del(id)
			w.wg.Done()
		}()

		err := waiter.Wait()
		w.doneC <- &Result{
			Waiter: waiter,
			Err:    err,
		}
	}()
}
