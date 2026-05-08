package backpressure_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/retryctl/internal/backpressure"
)

func TestNewInvalidCapacity(t *testing.T) {
	_, err := backpressure.New(0, 10*time.Millisecond)
	if err == nil {
		t.Fatal("expected error for capacity 0")
	}
}

func TestNewInvalidRefillRate(t *testing.T) {
	_, err := backpressure.New(1, 0)
	if err == nil {
		t.Fatal("expected error for zero refill interval")
	}
}

func TestNewValid(t *testing.T) {
	l, err := backpressure.New(3, 100*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer l.Stop()

	if got := l.Available(); got != 3 {
		t.Fatalf("expected 3 tokens, got %d", got)
	}
}

func TestAcquireDecrements(t *testing.T) {
	l, _ := backpressure.New(2, 100*time.Millisecond)
	defer l.Stop()

	ctx := context.Background()
	if err := l.Acquire(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := l.Available(); got != 1 {
		t.Fatalf("expected 1 token remaining, got %d", got)
	}
}

func TestAcquireBlocksWhenEmpty(t *testing.T) {
	l, _ := backpressure.New(1, 50*time.Millisecond)
	defer l.Stop()

	ctx := context.Background()
	_ = l.Acquire(ctx) // drain

	ctx2, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	// Should succeed once the refill fires (~50 ms).
	if err := l.Acquire(ctx2); err != nil {
		t.Fatalf("expected token after refill, got error: %v", err)
	}
}

func TestAcquireCancelledContext(t *testing.T) {
	l, _ := backpressure.New(1, time.Hour) // refill never fires in test
	defer l.Stop()

	ctx := context.Background()
	_ = l.Acquire(ctx) // drain

	ctx2, cancel := context.WithCancel(context.Background())
	cancel() // already cancelled

	if err := l.Acquire(ctx2); err == nil {
		t.Fatal("expected context error, got nil")
	}
}

func TestAvailableNeverExceedsCapacity(t *testing.T) {
	cap := 3
	l, _ := backpressure.New(cap, 5*time.Millisecond)
	defer l.Stop()

	time.Sleep(30 * time.Millisecond) // let several refill ticks fire

	if got := l.Available(); got > cap {
		t.Fatalf("available %d exceeds capacity %d", got, cap)
	}
}

// TestAcquireMultipleConsumesAll verifies that acquiring all tokens from a
// limiter with capacity > 1 correctly drains the pool to zero.
func TestAcquireMultipleConsumesAll(t *testing.T) {
	const capacity = 3
	l, _ := backpressure.New(capacity, time.Hour) // refill never fires in test
	defer l.Stop()

	ctx := context.Background()
	for i := 0; i < capacity; i++ {
		if err := l.Acquire(ctx); err != nil {
			t.Fatalf("acquire %d/%d failed: %v", i+1, capacity, err)
		}
	}

	if got := l.Available(); got != 0 {
		t.Fatalf("expected 0 tokens after draining, got %d", got)
	}
}
