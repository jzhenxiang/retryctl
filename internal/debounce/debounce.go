// Package debounce provides a mechanism to suppress rapid consecutive retries
// by enforcing a minimum quiet period between attempts.
package debounce

import (
	"errors"
	"sync"
	"time"
)

// ErrDebounced is returned when an attempt is suppressed by the debounce window.
var ErrDebounced = errors.New("attempt suppressed by debounce window")

// Debouncer tracks the last attempt time and rejects calls that arrive
// within the configured quiet window.
type Debouncer struct {
	mu       sync.Mutex
	window   time.Duration
	lastSeen time.Time
	now      func() time.Time
}

// New returns a Debouncer that suppresses attempts within window duration
// of the previous attempt. window must be positive.
func New(window time.Duration) (*Debouncer, error) {
	if window <= 0 {
		return nil, errors.New("debounce window must be positive")
	}
	return &Debouncer{
		window: window,
		now:    time.Now,
	}, nil
}

// Allow returns nil if the attempt should proceed, or ErrDebounced if it
// falls within the quiet window since the last allowed attempt.
func (d *Debouncer) Allow() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.now()
	if !d.lastSeen.IsZero() && now.Sub(d.lastSeen) < d.window {
		return ErrDebounced
	}
	d.lastSeen = now
	return nil
}

// Reset clears the recorded last-seen time, allowing the next call to Allow
// to proceed unconditionally.
func (d *Debouncer) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.lastSeen = time.Time{}
}

// Remaining returns how long until the next attempt will be allowed.
// Returns 0 if the window has already elapsed or no attempt has been recorded.
func (d *Debouncer) Remaining() time.Duration {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.lastSeen.IsZero() {
		return 0
	}
	elapsed := d.now().Sub(d.lastSeen)
	if elapsed >= d.window {
		return 0
	}
	return d.window - elapsed
}
