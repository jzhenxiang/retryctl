package budget_test

import (
	"testing"
	"time"

	"github.com/user/retryctl/internal/budget"
)

func TestNewInvalidMax(t *testing.T) {
	_, err := budget.New(0, time.Second)
	if err == nil {
		t.Fatal("expected error for max=0")
	}
}

func TestNewInvalidWindow(t *testing.T) {
	_, err := budget.New(3, 0)
	if err == nil {
		t.Fatal("expected error for window=0")
	}
}

func TestAllowUnderBudget(t *testing.T) {
	b, err := budget.New(3, time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 3; i++ {
		if err := b.Allow(); err != nil {
			t.Fatalf("attempt %d: unexpected error: %v", i, err)
		}
	}
}

func TestAllowExceedsBudget(t *testing.T) {
	b, _ := budget.New(2, time.Minute)
	_ = b.Allow()
	_ = b.Allow()
	if err := b.Allow(); err != budget.ErrBudgetExhausted {
		t.Fatalf("expected ErrBudgetExhausted, got %v", err)
	}
}

func TestRemainingDecrementsOnAllow(t *testing.T) {
	b, _ := budget.New(5, time.Minute)
	if b.Remaining() != 5 {
		t.Fatalf("expected 5, got %d", b.Remaining())
	}
	_ = b.Allow()
	if b.Remaining() != 4 {
		t.Fatalf("expected 4, got %d", b.Remaining())
	}
}

func TestResetRestoresBudget(t *testing.T) {
	b, _ := budget.New(2, time.Minute)
	_ = b.Allow()
	_ = b.Allow()
	b.Reset()
	if b.Remaining() != 2 {
		t.Fatalf("expected 2 after reset, got %d", b.Remaining())
	}
}

func TestWindowEviction(t *testing.T) {
	now := time.Now()
	clock := &now

	b, _ := budget.New(2, 100*time.Millisecond)
	// Inject clock via unexported field is not possible; use real sleep instead.
	_ = b.Allow()
	time.Sleep(150 * time.Millisecond)
	// After the window elapses the old token should be evicted.
	if b.Remaining() != 2 {
		t.Fatalf("expected budget to recover after window; got %d (clock=%v)", b.Remaining(), *clock)
	}
}
