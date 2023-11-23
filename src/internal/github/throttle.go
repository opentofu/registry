package github

import (
	"context"
	"time"

	"golang.org/x/time/rate"
)

type Throttle func() ThrottleToken

type ThrottleToken func()

func NewThrottle(ctx context.Context, every time.Duration, concurrent int) Throttle {
	limiter := rate.NewLimiter(rate.Every(every), 1)
	checkout := make(chan int, concurrent)

	for i := 0; i < concurrent; i++ {
		checkout <- i
	}
	return func() ThrottleToken {
		id := <-checkout
		// TODO We might want to handle the returned error from Wait here
		limiter.Wait(ctx)
		return func() { checkout <- id }
	}
}
