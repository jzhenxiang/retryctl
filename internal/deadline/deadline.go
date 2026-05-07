// Package deadline provides a wall-clock deadline guard that prevents
// retries from starting once an absolute point in time has been exceeded.
package deadline

import (
	"errors"
	"fmt"
	"time"
)

// ErrDeadlineExceeded is returned when the deadline has passed.
var ErrDeadlineExceeded = errors.New("deadline exceeded")

// Guard enforces an absolute deadline across retry attempts.
type Guard struct {
	deadline time.Time
	now      func() time.Time
}

// New creates a Guard that expires at the given absolute time.
// Returns an error if deadline is in the past relative to now.
func New(deadline time.Time) (*Guard, error) {
	return newWithClock(deadline, time.Now)
}

func newWithClock(deadline time.Time, now func() time.Time) (*Guard, error) {
	if deadline.IsZero() {
		return nil, fmt.Errorf("deadline must not be zero")
	}
	if !deadline.After(now()) {
		return nil, fmt.Errorf("deadline %v is already in the past", deadline)
	}
	return &Guard{deadline: deadline, now: now}, nil
}

// Allow returns nil if the current time is before the deadline,
// or ErrDeadlineExceeded otherwise.
func (g *Guard) Allow() error {
	if g.now().Before(g.deadline) {
		return nil
	}
	return ErrDeadlineExceeded
}

// Remaining returns the duration until the deadline. It may be negative
// if the deadline has already passed.
func (g *Guard) Remaining() time.Duration {
	return g.deadline.Sub(g.now())
}

// Deadline returns the absolute deadline time.
func (g *Guard) Deadline() time.Time {
	return g.deadline
}
