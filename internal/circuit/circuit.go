// Package circuit implements a simple circuit breaker for retryctl.
// When consecutive failures exceed a threshold the breaker opens and
// subsequent calls are rejected until a reset timeout elapses.
package circuit

import (
	"errors"
	"sync"
	"time"
)

// ErrOpen is returned when the circuit breaker is in the open state.
var ErrOpen = errors.New("circuit breaker is open")

// State represents the current state of the circuit breaker.
type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

// Breaker is a thread-safe circuit breaker.
type Breaker struct {
	mu           sync.Mutex
	maxFailures  int
	resetTimeout time.Duration
	failures     int
	state        State
	openedAt     time.Time
}

// New creates a Breaker that opens after maxFailures consecutive failures
// and attempts a reset after resetTimeout.
func New(maxFailures int, resetTimeout time.Duration) (*Breaker, error) {
	if maxFailures <= 0 {
		return nil, errors.New("maxFailures must be greater than zero")
	}
	if resetTimeout <= 0 {
		return nil, errors.New("resetTimeout must be greater than zero")
	}
	return &Breaker{
		maxFailures:  maxFailures,
		resetTimeout: resetTimeout,
		state:        StateClosed,
	}, nil
}

// Allow reports whether the call should be allowed through.
// It returns ErrOpen when the breaker is open.
func (b *Breaker) Allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case StateOpen:
		if time.Since(b.openedAt) >= b.resetTimeout {
			b.state = StateHalfOpen
			return nil
		}
		return ErrOpen
	default:
		return nil
	}
}

// RecordSuccess resets the failure counter and closes the circuit.
func (b *Breaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.state = StateClosed
}

// RecordFailure increments the failure counter and opens the circuit
// when the threshold is reached.
func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	if b.failures >= b.maxFailures {
		b.state = StateOpen
		b.openedAt = time.Now()
	}
}

// State returns the current state of the breaker.
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}
