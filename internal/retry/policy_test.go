package retry_test

import (
	"testing"
	"time"

	"github.com/retryctl/internal/backoff"
	"github.com/retryctl/internal/jitter"
	"github.com/retryctl/internal/retry"
)

func TestDefaultPolicyIsValid(t *testing.T) {
	p := retry.Default()
	if err := p.Validate(); err != nil {
		t.Fatalf("Default() should be valid, got: %v", err)
	}
}

func TestDefaultMaxAttempts(t *testing.T) {
	p := retry.Default()
	if p.MaxAttempts != 3 {
		t.Fatalf("expected MaxAttempts=3, got %d", p.MaxAttempts)
	}
}

func TestValidateZeroMaxAttempts(t *testing.T) {
	p := retry.Default()
	p.MaxAttempts = 0
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for MaxAttempts=0")
	}
}

func TestValidateNilStrategy(t *testing.T) {
	p := retry.Default()
	p.Strategy = nil
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for nil Strategy")
	}
}

func TestValidateNilJitter(t *testing.T) {
	p := retry.Default()
	p.Jitter = nil
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for nil Jitter")
	}
}

func TestValidateNilShouldRetry(t *testing.T) {
	p := retry.Default()
	p.ShouldRetry = nil
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for nil ShouldRetry")
	}
}

func TestNextDelayUsesStrategy(t *testing.T) {
	p := retry.Policy{
		MaxAttempts: 5,
		Strategy:    backoff.NewStrategy(backoff.Fixed, 2*time.Second, 30*time.Second),
		Jitter:      jitter.NewNone(),
		ShouldRetry: func(_ int, _ int, _ []byte, _ error) bool { return true },
	}
	d := p.NextDelay(1)
	if d != 2*time.Second {
		t.Fatalf("expected 2s delay, got %v", d)
	}
}

func TestShouldRetryFunc(t *testing.T) {
	called := false
	p := retry.Default()
	p.ShouldRetry = func(attempt int, code int, _ []byte, _ error) bool {
		called = true
		return attempt < 3 && code != 0
	}
	if !p.ShouldRetry(1, 1, nil, nil) {
		t.Fatal("expected ShouldRetry=true for attempt 1, code 1")
	}
	if !called {
		t.Fatal("ShouldRetry was never called")
	}
	if p.ShouldRetry(3, 1, nil, nil) {
		t.Fatal("expected ShouldRetry=false for attempt 3")
	}
}
