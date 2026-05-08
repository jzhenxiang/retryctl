// Package window provides a sliding-window counter used to track
// discrete events (e.g. failures) within a rolling time interval.
package window

import (
	"errors"
	"sync"
	"time"
)

// Window is a thread-safe sliding-window counter.
type Window struct {
	mu       sync.Mutex
	size     time.Duration
	buckets  []int
	nBuckets int
	ticks    []time.Time
	clock    func() time.Time
}

// New creates a Window that tracks events over the given duration split into
// the given number of buckets. Returns an error if either argument is invalid.
func New(size time.Duration, nBuckets int) (*Window, error) {
	return newWithClock(size, nBuckets, time.Now)
}

func newWithClock(size time.Duration, nBuckets int, clock func() time.Time) (*Window, error) {
	if size <= 0 {
		return nil, errors.New("window: size must be positive")
	}
	if nBuckets < 1 {
		return nil, errors.New("window: nBuckets must be at least 1")
	}
	return &Window{
		size:     size,
		buckets:  make([]int, nBuckets),
		ticks:    make([]time.Time, nBuckets),
		nBuckets: nBuckets,
		clock:    clock,
	}, nil
}

// Record increments the counter for the current time bucket.
func (w *Window) Record() {
	w.mu.Lock()
	defer w.mu.Unlock()
	now := w.clock()
	idx := w.bucketIndex(now)
	w.maybeReset(idx, now)
	w.buckets[idx]++
	w.ticks[idx] = now
}

// Count returns the total number of events recorded within the current window.
func (w *Window) Count() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	now := w.clock()
	cutoff := now.Add(-w.size)
	total := 0
	for i, t := range w.ticks {
		if !t.IsZero() && t.After(cutoff) {
			total += w.buckets[i]
		}
	}
	return total
}

// Reset clears all buckets.
func (w *Window) Reset() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.buckets = make([]int, w.nBuckets)
	w.ticks = make([]time.Time, w.nBuckets)
}

func (w *Window) bucketIndex(t time.Time) int {
	bucketDur := w.size / time.Duration(w.nBuckets)
	return int(t.UnixNano()/int64(bucketDur)) % w.nBuckets
}

func (w *Window) maybeReset(idx int, now time.Time) {
	if !w.ticks[idx].IsZero() && now.Sub(w.ticks[idx]) >= w.size {
		w.buckets[idx] = 0
		w.ticks[idx] = time.Time{}
	}
}
