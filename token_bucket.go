package rpslimit

import (
	"context"
	"sync/atomic"
	"time"
)

type tokenBucket struct {
	s, w uint32
	c    chan struct{}
	t    *time.Timer
}

func NewTokenBucket(ctx context.Context, limit uint64) Interface {
	l := &tokenBucket{
		c: make(chan struct{}, limit-1),
	}
	for i := 0; i < int(limit-1); i++ {
		l.c <- struct{}{}
	}

	dur := time.Duration(time.Second.Nanoseconds() / int64(limit))
	go func(ctx context.Context, dur time.Duration) {
		l.t = time.NewTimer(dur)
		for {
			select {
			case <-ctx.Done():
				l.t.Stop()
				atomic.StoreUint32(&l.s, 1)
				for atomic.LoadUint32(&l.w) > 0 {
				}
				close(l.c)
				return
			case <-l.t.C:
				if atomic.LoadUint32(&l.s) == 0 {
					atomic.AddUint32(&l.w, 1)
					select {
					case l.c <- struct{}{}:
					default:
					}
					atomic.AddUint32(&l.w, ^uint32(0))
				}
			}
		}
	}(ctx, dur)

	return l
}

func (l *tokenBucket) Allow() bool {
	select {
	case <-l.c:
		return true
	default:
		return false
	}
}
