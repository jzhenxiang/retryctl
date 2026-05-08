// Package hedging implements a hedged retry strategy that fires speculative
// duplicate attempts after a configurable delay, returning the first success.
package hedging

import (
	"context"
	"errors"
	"sync"
	"time"
)

// ErrAllHedgesFailed is returned when every hedged attempt fails.
var ErrAllHedgesFailed = errors.New("hedging: all attempts failed")

// Attempt represents a single hedged execution result.
type Attempt struct {
	Index int
	Err   error
}

// Hedger fires up to MaxAttempts executions of fn, each separated by Delay.
// It returns the error from the first attempt that succeeds (nil error), or
// ErrAllHedgesFailed if every attempt fails.
type Hedger struct {
	MaxAttempts int
	Delay       time.Duration
}

// New creates a Hedger with the given parameters.
// maxAttempts must be >= 1; delay must be >= 0.
func New(maxAttempts int, delay time.Duration) (*Hedger, error) {
	if maxAttempts < 1 {
		return nil, errors.New("hedging: maxAttempts must be at least 1")
	}
	if delay < 0 {
		return nil, errors.New("hedging: delay must be non-negative")
	}
	return &Hedger{MaxAttempts: maxAttempts, Delay: delay}, nil
}

// Run executes fn up to h.MaxAttempts times with h.Delay between launches.
// The context is forwarded to every attempt; cancelling it aborts pending ones.
func (h *Hedger) Run(ctx context.Context, fn func(ctx context.Context, index int) error) error {
	type result struct {
		index int
		err   error
	}

	results := make(chan result, h.MaxAttempts)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup

	for i := 0; i < h.MaxAttempts; i++ {
		if i > 0 {
			select {
			case <-time.After(h.Delay):
			case <-ctx.Done():
				break
			}
		}

		// Check if a winner was already found before launching more.
		select {
		case <-ctx.Done():
			goto drain
		default:
		}

		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			err := fn(ctx, idx)
			results <- result{index: idx, err: err}
		}(i)
	}

drain:
	go func() {
		wg.Wait()
		close(results)
	}()

	var lastErr error
	seen := 0
	for r := range results {
		seen++
		if r.err == nil {
			cancel()
			return nil
		}
		lastErr = r.err
	}
	if lastErr != nil {
		return lastErr
	}
	return ErrAllHedgesFailed
}
