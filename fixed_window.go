package rpslimit

import (
	"context"
	"sync/atomic"
	"time"
)

type fixedWindow struct {
	l uint64
	c uint64
}

func NewFixedWindow(ctx context.Context, limit uint64) Interface {
	l := &fixedWindow{
		l: limit,
	}
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
	return l
}

func (l *fixedWindow) Allow() bool {
	return atomic.AddUint64(&l.c, 1) <= l.l
}
