// Package retry provides a composable retry policy that wires together
// backoff, jitter, and predicate logic into a single decision-making unit
// consumed by the runner.
package retry

import (
	"fmt"
	"time"

	"github.com/retryctl/internal/backoff"
	"github.com/retryctl/internal/jitter"
)

// ShouldRetryFunc decides whether a failed attempt should be retried.
// attempt is 1-based. exitCode is the process exit code; err is non-nil
// when the process could not be started at all.
type ShouldRetryFunc func(attempt int, exitCode int, output []byte, err error) bool

// Policy holds the retry decision parameters.
type Policy struct {
	MaxAttempts int
	Strategy    backoff.Strategy
	Jitter      jitter.Jitter
	ShouldRetry ShouldRetryFunc
}

// Default returns a Policy with sensible defaults: 3 attempts, fixed
// 1-second backoff, no jitter, retry on any failure.
func Default() Policy {
	return Policy{
		MaxAttempts: 3,
		Strategy:    backoff.NewStrategy(backoff.Fixed, 1*time.Second, 30*time.Second),
		Jitter:      jitter.NewNone(),
		ShouldRetry: func(_ int, _ int, _ []byte, err error) bool { return true },
	}
}

// NextDelay returns the delay before the next attempt (attempt is 1-based,
// representing the attempt that just finished).
func (p Policy) NextDelay(attempt int) time.Duration {
	base := p.Strategy.Delay(attempt)
	return p.Jitter.Apply(base)
}

// Validate checks that the policy is self-consistent.
func (p Policy) Validate() error {
	if p.MaxAttempts < 1 {
		return fmt.Errorf("retry: MaxAttempts must be >= 1, got %d", p.MaxAttempts)
	}
	if p.Strategy == nil {
		return fmt.Errorf("retry: Strategy must not be nil")
	}
	if p.Jitter == nil {
		return fmt.Errorf("retry: Jitter must not be nil")
	}
	if p.ShouldRetry == nil {
		return fmt.Errorf("retry: ShouldRetry must not be nil")
	}
	return nil
}
