package backoff

import (
	"math"
	"time"
)

// Strategy defines the interface for backoff strategies.
type Strategy interface {
	Next(attempt int) time.Duration
}

// FixedStrategy waits a constant duration between retries.
type FixedStrategy struct {
	Delay time.Duration
}

func (f *FixedStrategy) Next(_ int) time.Duration {
	return f.Delay
}

// ExponentialStrategy implements exponential backoff with optional jitter.
type ExponentialStrategy struct {
	InitialDelay time.Duration
	Multiplier   float64
	MaxDelay     time.Duration
}

func (e *ExponentialStrategy) Next(attempt int) time.Duration {
	delay := float64(e.InitialDelay) * math.Pow(e.Multiplier, float64(attempt))
	if e.MaxDelay > 0 && time.Duration(delay) > e.MaxDelay {
		return e.MaxDelay
	}
	return time.Duration(delay)
}

// LinearStrategy increases delay linearly with each attempt.
type LinearStrategy struct {
	InitialDelay time.Duration
	Increment    time.Duration
	MaxDelay     time.Duration
}

func (l *LinearStrategy) Next(attempt int) time.Duration {
	delay := l.InitialDelay + time.Duration(attempt)*l.Increment
	if l.MaxDelay > 0 && delay > l.MaxDelay {
		return l.MaxDelay
	}
	return delay
}

// NewStrategy constructs a Strategy by name with the given base delay.
func NewStrategy(name string, initial, max time.Duration) Strategy {
	switch name {
	case "exponential":
		return &ExponentialStrategy{
			InitialDelay: initial,
			Multiplier:   2.0,
			MaxDelay:     max,
		}
	case "linear":
		return &LinearStrategy{
			InitialDelay: initial,
			Increment:    initial,
			MaxDelay:     max,
		}
	default: // "fixed"
		return &FixedStrategy{Delay: initial}
	}
}
