package rpslimit

import (
	"context"
	"sync"
	"time"
)

type SlidingLogV2 struct {
	mux  sync.Mutex
	head *slentry
	tail *slentry
	ln   uint64
	stk  []*slentry
	lim  uint64
}

func NewSlidingLogV2(ctx context.Context, limit uint64) *SlidingLogV2 {
	_ = ctx
	entry := &slentry{}
	return &SlidingLogV2{
		lim:  limit,
		head: entry,
		tail: entry,
		ln:   1,
	}
}

func (l *SlidingLogV2) Allow() bool {
	l.mux.Lock()
	defer l.mux.Unlock()

	outdate := time.Now().Add(-time.Second)
	for l.ln > 0 && l.head != nil && l.head.t.Before(outdate) {
		chead := l.head
		l.head = l.head.n
		chead.n = nil
		l.stk = append(l.stk, chead)
		l.ln--
	}
	if l.ln == 0 {
		l.tail = nil
	}

	var entry *slentry
	if len(l.stk) > 0 {
		entry = l.stk[len(l.stk)-1]
		l.stk = l.stk[:len(l.stk)-1]
	} else {
		entry = &slentry{}
	}
	entry.t = time.Now()

	if l.head == nil {
		l.head = entry
	}
	if l.tail == nil {
		l.tail = entry
	} else {
		l.tail.n = entry
	}
	l.tail = entry
	l.ln++

	return l.ln <= l.lim
}

type slentry struct {
	t time.Time
	n *slentry
}
