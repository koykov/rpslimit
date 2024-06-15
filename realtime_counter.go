package rpslimit

import (
	"context"

	"github.com/koykov/counter"
)

type RealtimeCounter struct {
	lim uint64
	c   *counter.Counter
}

func NewRealtimeCounter(ctx context.Context, limit uint64) *RealtimeCounter {
	_ = ctx
	l := &RealtimeCounter{lim: limit, c: counter.NewCounter()}
	return l
}

func (l *RealtimeCounter) Allow() bool {
	l.c.Inc()
	return uint64(l.c.Sum()) < l.lim
}
