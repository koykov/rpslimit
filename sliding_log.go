package rpslimit

import (
	"context"
	"sync"
	"time"
)

type slidingLog struct {
	mux sync.Mutex
	lim uint64
	buf []time.Time
}

func NewSlidingLog(ctx context.Context, limit uint64) Interface {
	_ = ctx
	return &slidingLog{lim: limit}
}

func (l *slidingLog) Allow() bool {
	l.mux.Lock()
	defer l.mux.Unlock()

	obsolete := time.Now().Add(-time.Second)
	for len(l.buf) > 0 && l.buf[0].Before(obsolete) {
		l.buf = l.buf[1:] // bad way, will lead to allocations all the time
	}

	l.buf = append(l.buf, time.Now())
	return len(l.buf) < int(l.lim)
}
