package metrics

import (
	"sync"
	"time"
)

// Summary holds aggregated run statistics.
type Summary struct {
	mu           sync.Mutex
	Attempts     int
	Successes    int
	Failures     int
	TotalElapsed time.Duration
	LastError    error
}

// New returns an initialised Summary.
func New() *Summary {
	return &Summary{}
}

// RecordAttempt records a single attempt outcome.
func (s *Summary) RecordAttempt(elapsed time.Duration, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Attempts++
	s.TotalElapsed += elapsed

	if err != nil {
		s.Failures++
		s.LastError = err
	} else {
		s.Successes++
	}
}

// AverageElapsed returns the mean duration per attempt.
// Returns zero if no attempts have been recorded.
func (s *Summary) AverageElapsed() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Attempts == 0 {
		return 0
	}
	return s.TotalElapsed / time.Duration(s.Attempts)
}

// Snapshot returns a copy of the current summary safe for reading
// outside the lock.
func (s *Summary) Snapshot() Summary {
	s.mu.Lock()
	defer s.mu.Unlock()

	return Summary{
		Attempts:     s.Attempts,
		Successes:    s.Successes,
		Failures:     s.Failures,
		TotalElapsed: s.TotalElapsed,
		LastError:    s.LastError,
	}
}
