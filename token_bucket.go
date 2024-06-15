package rpslimit

import (
	"context"
	"time"
)

type TokenBucket struct {
	c chan struct{}
	t *time.Timer
}

func NewTokenBucket(ctx context.Context, limit uint64) *TokenBucket {
	l := &TokenBucket{
		c: make(chan struct{}, limit-1),
	}
	for i := 0; i < int(limit-1); i++ {
		l.c <- struct{}{}
	}

	dur := time.Duration(time.Second.Nanoseconds() / int64(limit))
	go l.init(ctx, dur)

	return l
}

func (l *TokenBucket) Allow() bool {
	select {
	case <-l.c:
		return true
	default:
		return false
	}
}

func (l *TokenBucket) init(ctx context.Context, dur time.Duration) {
	l.t = time.NewTimer(dur)
	for {
		select {
		case <-ctx.Done():
			l.t.Stop()
			return
		case <-l.t.C:
			select {
			case l.c <- struct{}{}:
			default:
			}
		}
	}
}
