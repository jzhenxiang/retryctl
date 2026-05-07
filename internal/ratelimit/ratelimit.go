// Package ratelimit provides a token-bucket style rate limiter that can
// be used to cap how many retry attempts are made within a sliding window.
package ratelimit

import (
	"errors"
	"sync"
	"time"
)

// ErrRateLimited is returned when an attempt is denied by the limiter.
var ErrRateLimited = errors.New("rate limit exceeded")

// Limiter controls how many attempts are allowed within a window.
type Limiter struct {
	mu       sync.Mutex
	max      int
	window   time.Duration
	buckets  []time.Time
	nowFn    func() time.Time
}

// Config holds the configuration for a Limiter.
type Config struct {
	// Max is the maximum number of attempts allowed within Window.
	Max int
	// Window is the duration of the sliding window.
	Window time.Duration
}

// New creates a new Limiter from cfg.
// It returns an error if Max < 1 or Window <= 0.
func New(cfg Config) (*Limiter, error) {
	if cfg.Max < 1 {
		return nil, errors.New("ratelimit: Max must be at least 1")
	}
	if cfg.Window <= 0 {
		return nil, errors.New("ratelimit: Window must be positive")
	}
	return &Limiter{
		max:    cfg.Max,
		window: cfg.Window,
		nowFn:  time.Now,
	}, nil
}

// Allow reports whether an attempt is permitted at the current time.
// If permitted it records the attempt; otherwise it returns ErrRateLimited.
func (l *Limiter) Allow() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.nowFn()
	cutoff := now.Add(-l.window)

	// evict expired buckets
	valid := l.buckets[:0]
	for _, t := range l.buckets {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	l.buckets = valid

	if len(l.buckets) >= l.max {
		return ErrRateLimited
	}
	l.buckets = append(l.buckets, now)
	return nil
}

// Remaining returns how many more attempts are allowed in the current window.
func (l *Limiter) Remaining() int {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.nowFn()
	cutoff := now.Add(-l.window)
	count := 0
	for _, t := range l.buckets {
		if t.After(cutoff) {
			count++
		}
	}
	remaining := l.max - count
	if remaining < 0 {
		return 0
	}
	return remaining
}

// Reset clears all recorded attempts, effectively resetting the limiter
// as if no attempts have been made. This is useful in tests or when a
// logical boundary (e.g. a new job run) should start with a fresh window.
func (l *Limiter) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.buckets = l.buckets[:0]
}
