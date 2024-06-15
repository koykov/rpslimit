package rpslimit

import (
	"context"
	"time"
)

type LeakyBucket struct {
	c chan struct{}
	t *time.Ticker
}

func NewLeakyBucket(ctx context.Context, limit uint64) *LeakyBucket {
	l := &LeakyBucket{
		c: make(chan struct{}, limit+1),
	}
	dur := time.Duration(time.Second.Nanoseconds() / int64(limit))
	go l.init(ctx, dur)
	return l
}

func (l *LeakyBucket) Allow() bool {
	select {
	case l.c <- struct{}{}:
		return true
	default:
		return false
	}
}

func (l *LeakyBucket) init(ctx context.Context, dur time.Duration) {
	l.t = time.NewTicker(dur)
	for {
		select {
		case <-ctx.Done():
			l.t.Stop()
			return
		case <-l.t.C:
			select {
			case <-l.c:
			default:
			}
		}
	}
}
