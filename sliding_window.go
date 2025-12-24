package rpslimit

import (
	"context"
	"sync"
	"time"
)

type slidingWindow struct {
	mux    sync.Mutex
	lim    uint64
	ct     time.Time // current time
	pc, cc uint64    // prev count; curr count
}

func NewSlidingWindow(ctx context.Context, limit uint64) Interface {
	_ = ctx
	return &slidingWindow{
		lim: limit,
		ct:  time.Now(),
	}
}

func (l *slidingWindow) Allow() bool {
	l.mux.Lock()
	defer l.mux.Unlock()

	future := l.ct.Add(time.Second)
	now := time.Now()
	if now.After(future) {
		l.ct = now
		l.pc = l.cc
		l.cc = 0
	}
	eta := now.Sub(l.ct).Seconds()
	if uint64((float64(l.pc)*(float64(time.Second)-eta)/float64(time.Second))+float64(l.cc)) >= l.lim-1 {
		return false
	}
	l.cc++
	return true
}
