// Package backpressure provides a simple token-bucket style mechanism that
// slows retry attempts when the system is under load. Callers acquire a token
// before each attempt; if the bucket is empty the call blocks until a token
// becomes available or the context is cancelled.
package backpressure

import (
	"context"
	"errors"
	"time"
)

// ErrInvalidCapacity is returned when capacity is less than 1.
var ErrInvalidCapacity = errors.New("backpressure: capacity must be >= 1")

// ErrInvalidRefillRate is returned when the refill interval is non-positive.
var ErrInvalidRefillRate = errors.New("backpressure: refill interval must be > 0")

// Limiter controls the rate at which retry attempts are allowed to proceed.
type Limiter struct {
	tokens   chan struct{}
	ticker   *time.Ticker
	stop     chan struct{}
}

// New creates a Limiter with the given token capacity and refill interval.
// One token is added to the bucket every refillInterval until it is full.
func New(capacity int, refillInterval time.Duration) (*Limiter, error) {
	if capacity < 1 {
		return nil, ErrInvalidCapacity
	}
	if refillInterval <= 0 {
		return nil, ErrInvalidRefillRate
	}

	l := &Limiter{
		tokens: make(chan struct{}, capacity),
		ticker: time.NewTicker(refillInterval),
		stop:   make(chan struct{}),
	}

	// Pre-fill the bucket.
	for i := 0; i < capacity; i++ {
		l.tokens <- struct{}{}
	}

	go l.refill()
	return l, nil
}

// Acquire blocks until a token is available or ctx is cancelled.
// Returns ctx.Err() if the context is cancelled before a token is obtained.
func (l *Limiter) Acquire(ctx context.Context) error {
	select {
	case <-l.tokens:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Available returns the number of tokens currently in the bucket.
func (l *Limiter) Available() int {
	return len(l.tokens)
}

// Stop shuts down the background refill goroutine.
func (l *Limiter) Stop() {
	l.ticker.Stop()
	close(l.stop)
}

func (l *Limiter) refill() {
	for {
		select {
		case <-l.ticker.C:
			select {
			case l.tokens <- struct{}{}:
			default: // bucket full
			}
		case <-l.stop:
			return
		}
	}
}
