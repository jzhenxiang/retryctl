// Package budget provides a retry-budget mechanism that limits the total
// number of retry attempts allowed within a sliding time window, preventing
// retry storms in high-failure scenarios.
package budget

import (
	"errors"
	"sync"
	"time"
)

// ErrBudgetExhausted is returned when no retry tokens remain in the window.
var ErrBudgetExhausted = errors.New("retry budget exhausted")

// Budget tracks retry attempt tokens within a rolling time window.
type Budget struct {
	mu       sync.Mutex
	max      int
	window   time.Duration
	attempts []time.Time
	now      func() time.Time
}

// New creates a Budget that allows at most max retries within window.
// max must be >= 1 and window must be > 0.
func New(max int, window time.Duration) (*Budget, error) {
	if max < 1 {
		return nil, errors.New("budget: max must be at least 1")
	}
	if window <= 0 {
		return nil, errors.New("budget: window must be positive")
	}
	return &Budget{
		max:    max,
		window: window,
		now:    time.Now,
	}, nil
}

// Allow reports whether a retry attempt is permitted, consuming one token.
// It returns ErrBudgetExhausted when the budget is depleted.
func (b *Budget) Allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.evict()
	if len(b.attempts) >= b.max {
		return ErrBudgetExhausted
	}
	b.attempts = append(b.attempts, b.now())
	return nil
}

// Remaining returns the number of retry tokens still available in the window.
func (b *Budget) Remaining() int {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.evict()
	r := b.max - len(b.attempts)
	if r < 0 {
		return 0
	}
	return r
}

// Reset clears all recorded attempts, restoring the full budget.
func (b *Budget) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.attempts = b.attempts[:0]
}

// evict removes attempts that have fallen outside the current window.
// Caller must hold b.mu.
func (b *Budget) evict() {
	cutoff := b.now().Add(-b.window)
	i := 0
	for i < len(b.attempts) && b.attempts[i].Before(cutoff) {
		i++
	}
	b.attempts = b.attempts[i:]
}
