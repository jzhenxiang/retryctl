// Package sampling provides probabilistic retry sampling, allowing only a
// fraction of retry attempts to proceed. This is useful for reducing load
// during degraded conditions while still exercising retry paths.
package sampling

import (
	"errors"
	"math/rand"
	"sync"
)

// Sampler decides whether a retry attempt should be allowed based on a
// configured probability in the range (0, 1].
type Sampler struct {
	mu   sync.Mutex
	rng  *rand.Rand
	rate float64
}

// New creates a Sampler that allows attempts with the given probability.
// rate must be in the range (0, 1].
func New(rate float64, src rand.Source) (*Sampler, error) {
	if rate <= 0 || rate > 1 {
		return nil, errors.New("sampling: rate must be in the range (0, 1]")
	}
	if src == nil {
		src = rand.NewSource(rand.Int63())
	}
	return &Sampler{
		rng:  rand.New(src), //nolint:gosec
		rate: rate,
	}, nil
}

// Allow returns true if the current attempt should proceed based on the
// configured sampling rate. It is safe for concurrent use.
func (s *Sampler) Allow() bool {
	s.mu.Lock()
	v := s.rng.Float64()
	s.mu.Unlock()
	return v < s.rate
}

// Rate returns the configured sampling rate.
func (s *Sampler) Rate() float64 {
	return s.rate
}
