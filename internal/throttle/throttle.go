// Package throttle provides a concurrency limiter that caps the number of
// simultaneous retry attempts running at any given time.
package throttle

import (
	"errors"
	"fmt"
)

// Throttle limits the number of concurrent executions.
type Throttle struct {
	sem chan struct{}
	max int
}

// New creates a Throttle that allows at most max concurrent acquisitions.
// max must be >= 1.
func New(max int) (*Throttle, error) {
	if max < 1 {
		return nil, fmt.Errorf("throttle: max must be >= 1, got %d", max)
	}
	return &Throttle{
		sem: make(chan struct{}, max),
		max: max,
	}, nil
}

// Acquire blocks until a slot is available, then acquires it.
// Returns an error only if the throttle is nil.
func (t *Throttle) Acquire() error {
	if t == nil {
		return errors.New("throttle: nil throttle")
	}
	t.sem <- struct{}{}
	return nil
}

// Release frees a previously acquired slot.
func (t *Throttle) Release() {
	if t == nil {
		return
	}
	<-t.sem
}

// Available returns the number of slots currently free.
func (t *Throttle) Available() int {
	if t == nil {
		return 0
	}
	return t.max - len(t.sem)
}

// Max returns the maximum concurrency configured for this throttle.
func (t *Throttle) Max() int {
	if t == nil {
		return 0
	}
	return t.max
}
