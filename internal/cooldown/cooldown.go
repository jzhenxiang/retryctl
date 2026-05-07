// Package cooldown provides a per-exit-code cooldown mechanism that enforces
// a minimum wait period between retries based on the last observed exit code.
package cooldown

import (
	"errors"
	"sync"
	"time"
)

// Cooldown tracks per-exit-code cooldown windows.
type Cooldown struct {
	mu       sync.Mutex
	windows  map[int]time.Duration
	lastSeen map[int]time.Time
	clock    func() time.Time
}

// New creates a Cooldown with the given per-exit-code windows.
// windows maps an exit code to the minimum duration that must elapse before
// that exit code is allowed to trigger another retry.
func New(windows map[int]time.Duration) (*Cooldown, error) {
	if len(windows) == 0 {
		return nil, errors.New("cooldown: windows map must not be empty")
	}
	for code, d := range windows {
		if d <= 0 {
			return nil, fmt.Errorf("cooldown: window for exit code %d must be positive", code)
		}
	}
	return &Cooldown{
		windows:  windows,
		lastSeen: make(map[int]time.Time),
		clock:    time.Now,
	}, nil
}

// Allow reports whether a retry is permitted for the given exit code.
// If the exit code has no configured window it is always allowed.
// When allowed, the last-seen timestamp for the code is updated.
func (c *Cooldown) Allow(exitCode int) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	window, ok := c.windows[exitCode]
	if !ok {
		return true
	}
	now := c.clock()
	if last, seen := c.lastSeen[exitCode]; seen && now.Sub(last) < window {
		return false
	}
	c.lastSeen[exitCode] = now
	return true
}

// Reset clears the last-seen timestamp for the given exit code.
func (c *Cooldown) Reset(exitCode int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.lastSeen, exitCode)
}

// Remaining returns how much cooldown time is left for the given exit code.
// Returns 0 if the code is not in cooldown or has no configured window.
func (c *Cooldown) Remaining(exitCode int) time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()

	window, ok := c.windows[exitCode]
	if !ok {
		return 0
	}
	last, seen := c.lastSeen[exitCode]
	if !seen {
		return 0
	}
	elapsed := c.clock().Sub(last)
	if elapsed >= window {
		return 0
	}
	return window - elapsed
}
