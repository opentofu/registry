package github

import (
	"context"
	"time"

	"golang.org/x/time/rate"
)

type Throttle struct {
	ctx      context.Context
	limiter  *rate.Limiter
	checkout chan int
}

type ThrottleToken struct {
	manager Throttle
	id      int
}

func NewThrottle(ctx context.Context, every time.Duration, concurrent int) Throttle {
	r := Throttle{
		ctx:      ctx,
		limiter:  rate.NewLimiter(rate.Every(every), 1),
		checkout: make(chan int, concurrent),
	}
	for i := 0; i < concurrent; i++ {
		r.checkout <- i
	}
	return r
}

func (m Throttle) Wait() ThrottleToken {
	id := <-m.checkout
	// TODO We might want to handle the returned error from Wait here
	m.limiter.Wait(m.ctx)
	return ThrottleToken{m, id}
}

func (t ThrottleToken) Done() {
	t.manager.checkout <- t.id
}
