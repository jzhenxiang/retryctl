package runner_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/retryctl/internal/backoff"
	"github.com/user/retryctl/internal/config"
	"github.com/user/retryctl/internal/logger"
	"github.com/user/retryctl/internal/runner"
)

func fixedZero(_ int) time.Duration { return 0 }

type fixedStrategy struct{}

func (f fixedStrategy) Next(_ int) time.Duration { return 0 }

func newTestRunner(cmd string, args []string, maxAttempts int) *runner.Runner {
	cfg := &config.Config{Command: cmd, Args: args, MaxAttempts: maxAttempts}
	log := logger.New(nil)
	return runner.New(cfg, log, fixedStrategy{})
}

func TestRunSuccess(t *testing.T) {
	r := newTestRunner("true", nil, 3)
	if err := r.Run(context.Background()); err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	snap := r.Metrics().Snapshot()
	if snap.Successes != 1 {
		t.Fatalf("expected 1 success, got %d", snap.Successes)
	}
}

func TestRunAllFailures(t *testing.T) {
	r := newTestRunner("false", nil, 3)
	err := r.Run(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	snap := r.Metrics().Snapshot()
	if snap.Failures != 3 {
		t.Fatalf("expected 3 failures, got %d", snap.Failures)
	}
}

func TestRunContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	r := newTestRunner("true", nil, 5)
	if err := r.Run(ctx); err == nil {
		t.Fatal("expected context error")
	}
}

func TestRunDefaultMaxAttempts(t *testing.T) {
	cfg := config.Default()
	cfg.Command = "true"
	log := logger.New(nil)
	strat, _ := backoff.NewStrategy("fixed", 0)
	r := runner.New(cfg, log, strat)
	if err := r.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMetricsAccumulate(t *testing.T) {
	r := newTestRunner("false", nil, 2)
	_ = r.Run(context.Background())
	snap := r.Metrics().Snapshot()
	if snap.Attempts != 2 {
		t.Fatalf("expected 2 attempts, got %d", snap.Attempts)
	}
	if snap.LastError == nil {
		t.Fatal("expected LastError to be set")
	}
}
