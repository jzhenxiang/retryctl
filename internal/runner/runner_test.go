package runner_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/retryctl/internal/backoff"
	"github.com/user/retryctl/internal/runner"
)

func fixedZero() backoff.Strategy {
	s, _ := backoff.NewStrategy("fixed", 0)
	return s
}

func TestRunSuccess(t *testing.T) {
	r := runner.New(runner.Config{MaxAttempts: 3, Strategy: fixedZero()})
	results, err := r.Run(context.Background(), "true")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 attempt, got %d", len(results))
	}
	if results[0].ExitCode != 0 {
		t.Errorf("expected exit code 0, got %d", results[0].ExitCode)
	}
}

func TestRunAllFailures(t *testing.T) {
	r := runner.New(runner.Config{MaxAttempts: 3, Strategy: fixedZero()})
	results, err := r.Run(context.Background(), "false")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 attempts, got %d", len(results))
	}
	for _, res := range results {
		if res.ExitCode == 0 {
			t.Errorf("attempt %d: expected non-zero exit code", res.Attempt)
		}
	}
}

func TestRunContextCancellation(t *testing.T) {
	s, _ := backoff.NewStrategy("fixed", 500*time.Millisecond)
	r := runner.New(runner.Config{MaxAttempts: 5, Strategy: s})

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	_, err := r.Run(ctx, "false")
	if err == nil {
		t.Fatal("expected context error, got nil")
	}
}

func TestRunDefaultMaxAttempts(t *testing.T) {
	r := runner.New(runner.Config{Strategy: fixedZero()})
	results, _ := r.Run(context.Background(), "false")
	if len(results) != 1 {
		t.Errorf("expected default 1 attempt, got %d", len(results))
	}
}
