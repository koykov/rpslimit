package rpslimit

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestRPSLimiter(t *testing.T) {
	fn := func(t *testing.T, ctx context.Context, l RPSLimiter, w int) (uint64, uint64) {
		var c, f uint64
		var wg sync.WaitGroup
		for i := 0; i < w; i++ {
			wg.Add(1)
			go func(ctx context.Context) {
				defer wg.Done()
				for {
					select {
					case <-ctx.Done():
						return
					default:
						if l.Allow() {
							atomic.AddUint64(&c, 1)
						} else {
							atomic.AddUint64(&f, 1)
						}
						time.Sleep(time.Millisecond * 10)
					}
				}
			}(ctx)
		}
		wg.Wait()
		return atomic.LoadUint64(&c), atomic.LoadUint64(&f)
	}
	t.Run("fixed window", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		l := NewFixedWindow(ctx, 100)
		c, _ := fn(t, ctx, l, 1)
		_ = cancel
		if c > 200 {
			t.Errorf("got %d, want max 200", c)
		}
	})
	t.Run("sliding log", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		l := NewSlidingLog(context.TODO(), 100)
		c, _ := fn(t, ctx, l, 2)
		_ = cancel
		if c > 100 {
			t.Errorf("got %d, want max 100", c)
		}
	})
	t.Run("sliding log v2", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		l := NewSlidingLogV2(ctx, 100)
		c, _ := fn(t, ctx, l, 2)
		_ = cancel
		if c != 100 {
			t.Errorf("got %d, want 100", c)
		}
	})
}
