package runner

import (
	"context"
	"os/exec"
	"time"

	"github.com/user/retryctl/internal/backoff"
)

// Result holds the outcome of a single command attempt.
type Result struct {
	Attempt  int
	ExitCode int
	Duration time.Duration
	Err      error
}

// Config configures retry behaviour for a command run.
type Config struct {
	MaxAttempts int
	Strategy    backoff.Strategy
}

// Runner executes a command with retry logic.
type Runner struct {
	cfg Config
}

// New creates a Runner with the given Config.
func New(cfg Config) *Runner {
	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = 1
	}
	return &Runner{cfg: cfg}
}

// Run executes the command up to MaxAttempts times, applying backoff between
// failures. It returns all attempt results and the last error (nil on success).
func (r *Runner) Run(ctx context.Context, name string, args ...string) ([]Result, error) {
	var results []Result

	for attempt := 1; attempt <= r.cfg.MaxAttempts; attempt++ {
		start := time.Now()
		cmd := exec.CommandContext(ctx, name, args...)
		err := cmd.Run()
		duration := time.Since(start)

		exitCode := 0
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				exitCode = exitErr.ExitCode()
			} else {
				exitCode = -1
			}
		}

		res := Result{
			Attempt:  attempt,
			ExitCode: exitCode,
			Duration: duration,
			Err:      err,
		}
		results = append(results, res)

		if err == nil {
			return results, nil
		}

		if attempt < r.cfg.MaxAttempts {
			wait := r.cfg.Strategy.Next(attempt)
			select {
			case <-ctx.Done():
				return results, ctx.Err()
			case <-time.After(wait):
			}
		}
	}

	return results, results[len(results)-1].Err
}
