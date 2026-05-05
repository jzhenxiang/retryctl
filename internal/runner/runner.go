package runner

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/user/retryctl/internal/backoff"
	"github.com/user/retryctl/internal/config"
	"github.com/user/retryctl/internal/logger"
	"github.com/user/retryctl/internal/metrics"
)

// Runner executes a command with retry logic.
type Runner struct {
	cfg      *config.Config
	log      *logger.Logger
	strategy backoff.Strategy
	metrics  *metrics.Summary
}

// New creates a Runner from the provided config and logger.
func New(cfg *config.Config, log *logger.Logger, strategy backoff.Strategy) *Runner {
	return &Runner{
		cfg:      cfg,
		log:      log,
		strategy: strategy,
		metrics:  metrics.New(),
	}
}

// Metrics returns the accumulated run summary.
func (r *Runner) Metrics() *metrics.Summary {
	return r.metrics
}

// Run executes the configured command, retrying on failure.
func (r *Runner) Run(ctx context.Context) error {
	var lastErr error

	for attempt := 1; attempt <= r.cfg.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		r.log.Info("attempt", map[string]any{
			"attempt": attempt,
			"max":     r.cfg.MaxAttempts,
			"command": r.cfg.Command,
		})

		start := time.Now()
		//nolint:gosec
		cmd := exec.CommandContext(ctx, r.cfg.Command, r.cfg.Args...)
		err := cmd.Run()
		elapsed := time.Since(start)
		r.metrics.RecordAttempt(elapsed, err)

		if err == nil {
			r.log.Info("success", map[string]any{"attempt": attempt, "elapsed": elapsed.String()})
			return nil
		}

		lastErr = err
		r.log.Warn("attempt failed", map[string]any{"attempt": attempt, "error": err.Error()})

		if attempt < r.cfg.MaxAttempts {
			delay := r.strategy.Next(attempt)
			r.log.Info("backing off", map[string]any{"delay": delay.String()})
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	return fmt.Errorf("all %d attempts failed: %w", r.cfg.MaxAttempts, lastErr)
}
