package rpslimit

import (
	"context"
	"sync/atomic"
	"time"
)

type leakyBucket struct {
	s, w uint32
	c    chan struct{}
	t    *time.Ticker
}

func NewLeakyBucket(ctx context.Context, limit uint64) Interface {
	l := &leakyBucket{
		c: make(chan struct{}, limit+1),
	}
	dur := time.Duration(time.Second.Nanoseconds() / int64(limit))
	go func() {
		l.t = time.NewTicker(dur)
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
				select {
				case <-l.c:
				default:
				}
			}
		}
	}()
	return l
}

func (l *leakyBucket) Allow() bool {
	if atomic.LoadUint32(&l.s) == 1 {
		return false
	}

	atomic.AddUint32(&l.w, 1)
	defer atomic.AddUint32(&l.w, ^uint32(0))

	select {
	case l.c <- struct{}{}:
		return true
	default:
		return false
	}
}
