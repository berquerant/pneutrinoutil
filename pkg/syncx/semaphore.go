package syncx

import (
	"context"

	"golang.org/x/sync/semaphore"
)

type Semaphore interface {
	Wait(ctx context.Context) error
	Signal()
}

var (
	_ Semaphore = &UnlimitedSemaphore{}
	_ Semaphore = &FixedSemaphore{}
)

func NewSemaphore(n int) Semaphore {
	if n < 1 {
		return &UnlimitedSemaphore{}
	}
	return &FixedSemaphore{
		semaphore.NewWeighted(int64(n)),
	}
}

type UnlimitedSemaphore struct{}

func (UnlimitedSemaphore) Wait(_ context.Context) error { return nil }
func (UnlimitedSemaphore) Signal()                      {}

type FixedSemaphore struct {
	*semaphore.Weighted
}

func (s *FixedSemaphore) Wait(ctx context.Context) error { return s.Acquire(ctx, 1) }
func (s *FixedSemaphore) Signal()                        { s.Release(1) }
