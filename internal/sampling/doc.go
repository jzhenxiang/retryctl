// Package sampling implements probabilistic sampling for retry attempts.
//
// A Sampler is constructed with a rate in (0, 1] and consulted before each
// retry. When Allow returns false the caller should skip the attempt, shedding
// load proportionally to (1 - rate).
//
// Example:
//
//	s, _ := sampling.New(0.25, nil) // allow 25 % of retries
//	if s.Allow() {
//	    // proceed with retry
//	}
package sampling
