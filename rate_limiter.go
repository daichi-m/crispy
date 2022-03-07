package crispy

import (
	"errors"
	"math"
	"time"
)

/* RateLimiter is the primary interface through which users can rate-limit
 * the goroutines. The function to run in the goroutine is either Runner.Run or provided
 * as a RunnerFunc parameter.
 */
type RateLimiter interface {
	// RateLimit runs the Runner.Run method in a goroutine with the rate-limited settings.
	Go(r Runner) error

	// RateLimitFunc runs the RunnerFunc provided in a goroutine with the rate-limited settings.
	GoFunc(rf RunnerFunc) error

	// Cleanup is called to cleanup all resources allocated for the rate-limiter.
	Cleanup() error
}

type channelRateLimiter struct {
	rateLimit int
	doTimeout bool
	timeout   time.Duration
	onTimeout func() error
	rlChannel chan struct{}
}

var _ RateLimiter = (*channelRateLimiter)(nil)

func NewRateLimiter(options ...RateLimiterOption) RateLimiter {
	rl := &channelRateLimiter{
		rateLimit: math.MaxInt,
		doTimeout: false,
		timeout:   8760 * time.Hour, // 1 Year - mimics infinity
		onTimeout: func() error {
			return nil
		},
		rlChannel: nil,
	}
	for _, opts := range options {
		opts(rl)
	}
	if rl.rateLimit == 1 {
		rl.rlChannel = make(chan struct{})
	} else {
		rl.rlChannel = make(chan struct{}, rl.rateLimit)
	}
	return rl
}

type RateLimiterOption func(*channelRateLimiter)

func WithRateLimit(limit int) RateLimiterOption {
	return func(crl *channelRateLimiter) {
		crl.rateLimit = limit
	}
}

func WithTimeout(timeout time.Duration) RateLimiterOption {
	return func(crl *channelRateLimiter) {
		crl.doTimeout = true
		crl.timeout = timeout
	}
}

func WithOnTimeout(f func() error) RateLimiterOption {
	return func(crl *channelRateLimiter) {
		crl.onTimeout = f
	}
}

func (crl *channelRateLimiter) Go(r Runner) error {
	timer := time.NewTimer(crl.timeout)

	select {
	case <-crl.rlChannel:
		// Do call go
	case <-timer.C:
		// Do timeout
		return errors.New("")
	}
	return nil
}

func (crl *channelRateLimiter) GoFunc(rf RunnerFunc) error {
	return nil
}

func (crl *channelRateLimiter) Cleanup() error {
	close(crl.rlChannel)
	return nil
}

type Runner interface {
	Run() error
}

type RunnerFunc func() error
