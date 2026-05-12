// Package wavefront provides a sliding-window failure-rate tracker.
// It records attempt outcomes and reports whether the failure rate
// over the most recent window exceeds a configured threshold.
package wavefront

import (
	"errors"
	"sync"
	"time"
)

// Tracker tracks the failure rate over a sliding time window.
type Tracker struct {
	mu        sync.Mutex
	window    time.Duration
	threshold float64 // 0.0–1.0
	events    []event
}

type event struct {
	at      time.Time
	failed bool
}

// New creates a Tracker with the given sliding window duration and
// failure-rate threshold (0.0–1.0). Returns an error for invalid inputs.
func New(window time.Duration, threshold float64) (*Tracker, error) {
	if window <= 0 {
		return nil, errors.New("wavefront: window must be positive")
	}
	if threshold < 0 || threshold > 1 {
		return nil, errors.New("wavefront: threshold must be between 0.0 and 1.0")
	}
	return &Tracker{window: window, threshold: threshold}, nil
}

// Record registers an attempt outcome. failed=true indicates the attempt failed.
func (t *Tracker) Record(failed bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.prune(time.Now())
	t.events = append(t.events, event{at: time.Now(), failed: failed})
}

// FailureRate returns the fraction of failed attempts in the current window.
// Returns 0 if no events have been recorded.
func (t *Tracker) FailureRate() float64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.prune(time.Now())
	if len(t.events) == 0 {
		return 0
	}
	var failures int
	for _, e := range t.events {
		if e.failed {
			failures++
		}
	}
	return float64(failures) / float64(len(t.events))
}

// ExceedsThreshold reports whether the current failure rate is strictly
// greater than the configured threshold.
func (t *Tracker) ExceedsThreshold() bool {
	return t.FailureRate() > t.threshold
}

// prune removes events older than the window. Must be called with t.mu held.
func (t *Tracker) prune(now time.Time) {
	cutoff := now.Add(-t.window)
	i := 0
	for i < len(t.events) && t.events[i].at.Before(cutoff) {
		i++
	}
	t.events = t.events[i:]
}
