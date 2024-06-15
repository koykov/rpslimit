package rpslimit

import (
	"context"
	"sync"
	"time"
)

type SlidingWindow struct {
	mux    sync.Mutex
	lim    uint64
	ct     time.Time
	pc, cc uint64
}

func NewSlidingWindow(ctx context.Context, limit uint64) *SlidingWindow {
	_ = ctx
	return &SlidingWindow{
		lim: limit,
		ct:  time.Now(),
	}
}

func (l *SlidingWindow) Allow() bool {
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
