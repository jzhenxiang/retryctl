package timeout_test

import (
	"context"
	"testing"
	"time"

	"github.com/yourorg/retryctl/internal/timeout"
)

func TestValidateOK(t *testing.T) {
	cfg := timeout.Config{PerAttempt: time.Second, Overall: 10 * time.Second}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateZeroIsOK(t *testing.T) {
	cfg := timeout.Config{}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error for zero config: %v", err)
	}
}

func TestValidateNegativePerAttempt(t *testing.T) {
	cfg := timeout.Config{PerAttempt: -time.Second}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative per-attempt")
	}
}

func TestValidateNegativeOverall(t *testing.T) {
	cfg := timeout.Config{Overall: -time.Second}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative overall")
	}
}

func TestValidatePerAttemptExceedsOverall(t *testing.T) {
	cfg := timeout.Config{PerAttempt: 10 * time.Second, Overall: time.Second}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error when per-attempt exceeds overall")
	}
}

func TestOverallZeroNoTimeout(t *testing.T) {
	e := timeout.New(timeout.Config{})
	ctx, cancel := e.Overall(context.Background())
	defer cancel()
	select {
	case <-ctx.Done():
		t.Fatal("context should not be cancelled immediately")
	default:
	}
}

func TestOverallDeadlineEnforced(t *testing.T) {
	e := timeout.New(timeout.Config{Overall: 20 * time.Millisecond})
	ctx, cancel := e.Overall(context.Background())
	defer cancel()
	select {
	case <-ctx.Done():
		// expected
	case <-time.After(200 * time.Millisecond):
		t.Fatal("overall context should have timed out")
	}
}

func TestAttemptDeadlineEnforced(t *testing.T) {
	e := timeout.New(timeout.Config{PerAttempt: 20 * time.Millisecond})
	ctx, cancel := e.Attempt(context.Background())
	defer cancel()
	select {
	case <-ctx.Done():
		// expected
	case <-time.After(200 * time.Millisecond):
		t.Fatal("attempt context should have timed out")
	}
}

func TestAttemptZeroNoTimeout(t *testing.T) {
	e := timeout.New(timeout.Config{})
	ctx, cancel := e.Attempt(context.Background())
	defer cancel()
	select {
	case <-ctx.Done():
		t.Fatal("context should not be cancelled immediately")
	default:
	}
}
