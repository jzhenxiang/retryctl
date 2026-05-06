// Package timeout provides per-attempt and overall deadline enforcement
// for retryctl command executions.
package timeout

import (
	"context"
	"fmt"
	"time"
)

// Config holds timeout durations for retry execution.
type Config struct {
	// PerAttempt is the maximum duration allowed for a single attempt.
	// Zero means no per-attempt timeout.
	PerAttempt time.Duration

	// Overall is the maximum duration allowed across all attempts combined.
	// Zero means no overall timeout.
	Overall time.Duration
}

// Enforcer wraps a parent context with timeout constraints.
type Enforcer struct {
	cfg Config
}

// New returns an Enforcer configured with cfg.
func New(cfg Config) *Enforcer {
	return &Enforcer{cfg: cfg}
}

// Overall returns a context that is cancelled after the overall deadline.
// The caller must call the returned cancel function when done.
// If Overall duration is zero, the parent context is returned unchanged.
func (e *Enforcer) Overall(parent context.Context) (context.Context, context.CancelFunc) {
	if e.cfg.Overall <= 0 {
		return context.WithCancel(parent)
	}
	return context.WithTimeout(parent, e.cfg.Overall)
}

// Attempt returns a context that is cancelled after the per-attempt deadline.
// The caller must call the returned cancel function when done.
// If PerAttempt duration is zero, the parent context is returned unchanged.
func (e *Enforcer) Attempt(parent context.Context) (context.Context, context.CancelFunc) {
	if e.cfg.PerAttempt <= 0 {
		return context.WithCancel(parent)
	}
	return context.WithTimeout(parent, e.cfg.PerAttempt)
}

// Validate checks that the Config values are sensible.
func (c Config) Validate() error {
	if c.PerAttempt < 0 {
		return fmt.Errorf("timeout: per-attempt duration must be non-negative, got %s", c.PerAttempt)
	}
	if c.Overall < 0 {
		return fmt.Errorf("timeout: overall duration must be non-negative, got %s", c.Overall)
	}
	if c.Overall > 0 && c.PerAttempt > 0 && c.PerAttempt > c.Overall {
		return fmt.Errorf("timeout: per-attempt (%s) must not exceed overall (%s)", c.PerAttempt, c.Overall)
	}
	return nil
}
