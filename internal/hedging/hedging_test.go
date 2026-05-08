package hedging_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/retryctl/internal/hedging"
)

func TestNewInvalidMaxAttempts(t *testing.T) {
	_, err := hedging.New(0, time.Millisecond)
	if err == nil {
		t.Fatal("expected error for maxAttempts=0")
	}
}

func TestNewNegativeDelay(t *testing.T) {
	_, err := hedging.New(2, -time.Millisecond)
	if err == nil {
		t.Fatal("expected error for negative delay")
	}
}

func TestNewValid(t *testing.T) {
	h, err := hedging.New(3, 10*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h.MaxAttempts != 3 {
		t.Errorf("expected MaxAttempts=3, got %d", h.MaxAttempts)
	}
}

func TestRunFirstAttemptSucceeds(t *testing.T) {
	h, _ := hedging.New(3, 50*time.Millisecond)
	var calls int32
	err := h.Run(context.Background(), func(_ context.Context, _ int) error {
		atomic.AddInt32(&calls, 1)
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	// Only the first goroutine should have run before cancel propagates.
	if atomic.LoadInt32(&calls) < 1 {
		t.Error("expected at least one call")
	}
}

func TestRunAllFail(t *testing.T) {
	h, _ := hedging.New(3, time.Millisecond)
	sentinel := errors.New("boom")
	err := h.Run(context.Background(), func(_ context.Context, _ int) error {
		return sentinel
	})
	if err == nil {
		t.Fatal("expected an error when all attempts fail")
	}
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestRunSecondAttemptSucceeds(t *testing.T) {
	h, _ := hedging.New(3, 5*time.Millisecond)
	var calls int32
	err := h.Run(context.Background(), func(_ context.Context, idx int) error {
		atomic.AddInt32(&calls, 1)
		if idx == 0 {
			return errors.New("first fails")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestRunContextCancellation(t *testing.T) {
	h, _ := hedging.New(5, 10*time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately
	err := h.Run(ctx, func(ctx context.Context, _ int) error {
		<-ctx.Done()
		return ctx.Err()
	})
	// Should complete without hanging.
	_ = err
}

func TestRunSingleAttemptNoDelay(t *testing.T) {
	h, _ := hedging.New(1, 0)
	err := h.Run(context.Background(), func(_ context.Context, _ int) error {
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
