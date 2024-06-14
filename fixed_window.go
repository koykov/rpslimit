package rpslimit

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type FixedWindow struct {
	l uint64
	c uint64
	o sync.Once
}

func NewFixedWindow(ctx context.Context, limit uint64) *FixedWindow {
	l := &FixedWindow{
		l: limit,
	}
	go l.init(ctx)
	return l
}

func (l *FixedWindow) Allow() bool {
	return atomic.AddUint64(&l.c, 1) <= l.l
}

func (l *FixedWindow) init(ctx context.Context) {
	t := time.NewTimer(time.Second)
	go func() {
		for {
			select {
			case <-ctx.Done():
				t.Stop()
				return
			case <-t.C:
				atomic.StoreUint64(&l.c, 0)
			}
		}
	}()
}
