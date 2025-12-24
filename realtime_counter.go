package rpslimit

import (
	"context"

	"github.com/koykov/counter"
)

type realtimeCounter struct {
	lim uint64
	c   *counter.Counter
}

func NewRealtimeCounter(ctx context.Context, limit uint64) Interface {
	_ = ctx
	l := &realtimeCounter{lim: limit, c: counter.NewCounter()}
	return l
}

func (l *realtimeCounter) Allow() bool {
	l.c.Inc()
	return uint64(l.c.Sum()) < l.lim
}
