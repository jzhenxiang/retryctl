// Package concurrency provides a guard that limits the number of concurrent
// retry attempts running across multiple goroutines.
package concurrency

import (
	"errors"
	"fmt"
	"sync"
)

// Guard enforces a maximum number of concurrent in-flight attempts.
type Guard struct {
	mu      sync.Mutex
	max     int
	active  int
}

// New creates a Guard that allows at most max simultaneous attempts.
// max must be greater than zero.
func New(max int) (*Guard, error) {
	if max <= 0 {
		return nil, fmt.Errorf("concurrency: max must be greater than zero, got %d", max)
	}
	return &Guard{max: max}, nil
}

// Acquire attempts to claim one concurrency slot. It returns true and a
// release function when a slot is available, or false (with a nil release)
// when the limit has been reached.
func (g *Guard) Acquire() (bool, func()) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.active >= g.max {
		return false, nil
	}

	g.active++
	release := func() {
		g.mu.Lock()
		defer g.mu.Unlock()
		if g.active > 0 {
			g.active--
		}
	}
	return true, release
}

// Active returns the current number of in-flight attempts.
func (g *Guard) Active() int {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.active
}

// Available returns the number of remaining slots before the limit is hit.
func (g *Guard) Available() int {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.max - g.active
}

// ErrAtLimit is returned by helper callers that treat a full guard as an error.
var ErrAtLimit = errors.New("concurrency: limit reached")
