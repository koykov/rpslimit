package rpslimit

import (
	"sync"
	"time"
)

type SlidingLogV2 struct {
	mux sync.Mutex
	lim uint64
	dur time.Duration
	buf []time.Time
}

func NewSlidingLogV2(limit uint64, interval time.Duration) *SlidingLogV2 {
	return &SlidingLogV2{
		lim: limit,
		dur: interval,
	}
}

func (l *SlidingLogV2) Allow() bool {
	l.mux.Lock()
	defer l.mux.Unlock()

	obsolete := time.Now().Add(-l.dur)
	for len(l.buf) > 0 && l.buf[0].Before(obsolete) {
		l.buf = l.buf[1:] // bad way, will lead to allocations all the time
	}

	l.buf = append(l.buf, time.Now())
	return len(l.buf) < int(l.lim)
}
