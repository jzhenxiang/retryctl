// Package backoff provides retry delay strategies.
package backoff

import (
	"fmt"
	"time"
)

// Strategy computes the delay before the next retry attempt.
// attempt is 1-based (first retry is attempt 1).
type Strategy interface {
	Delay(attempt int) time.Duration
}

// Fixed returns the same delay for every attempt.
type Fixed struct {
	Delay_ time.Duration
}

func (f Fixed) Delay(_ int) time.Duration { return f.Delay_ }

// Linear increases the delay linearly: attempt * base.
type Linear struct {
	Base time.Duration
}

func (l Linear) Delay(attempt int) time.Duration {
	return time.Duration(attempt) * l.Base
}

// Exponential doubles the delay each attempt, capped at MaxDelay.
type Exponential struct {
	Base     time.Duration
	MaxDelay time.Duration
}

func (e Exponential) Delay(attempt int) time.Duration {
	d := e.Base
	for i := 1; i < attempt; i++ {
		d *= 2
		if e.MaxDelay > 0 && d > e.MaxDelay {
			return e.MaxDelay
		}
	}
	if e.MaxDelay > 0 && d > e.MaxDelay {
		return e.MaxDelay
	}
	return d
}

// NewStrategy constructs a Strategy from a name string and base delay.
// Supported names: "fixed", "linear", "exponential".
func NewStrategy(name string, base, maxDelay time.Duration) (Strategy, error) {
	switch name {
	case "fixed":
		return Fixed{Delay_: base}, nil
	case "linear":
		return Linear{Base: base}, nil
	case "exponential":
		return Exponential{Base: base, MaxDelay: maxDelay}, nil
	default:
		return nil, fmt.Errorf("backoff: unknown strategy %q", name)
	}
}
