package circuit

import (
	"testing"
	"time"
)

func newBreaker(t *testing.T, max int, reset time.Duration) *Breaker {
	t.Helper()
	b, err := New(max, reset)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return b
}

func TestNewInvalidMaxFailures(t *testing.T) {
	_, err := New(0, time.Second)
	if err == nil {
		t.Fatal("expected error for maxFailures=0")
	}
}

func TestNewInvalidResetTimeout(t *testing.T) {
	_, err := New(1, 0)
	if err == nil {
		t.Fatal("expected error for resetTimeout=0")
	}
}

func TestInitiallyClosed(t *testing.T) {
	b := newBreaker(t, 3, time.Second)
	if b.State() != StateClosed {
		t.Fatalf("expected StateClosed, got %v", b.State())
	}
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestOpensAfterThreshold(t *testing.T) {
	b := newBreaker(t, 2, time.Second)
	b.RecordFailure()
	if b.State() != StateClosed {
		t.Fatal("should still be closed after 1 failure")
	}
	b.RecordFailure()
	if b.State() != StateOpen {
		t.Fatal("should be open after 2 failures")
	}
	if err := b.Allow(); err != ErrOpen {
		t.Fatalf("expected ErrOpen, got %v", err)
	}
}

func TestSuccessResetsClosed(t *testing.T) {
	b := newBreaker(t, 1, time.Second)
	b.RecordFailure()
	b.RecordSuccess()
	if b.State() != StateClosed {
		t.Fatal("expected StateClosed after success")
	}
	if err := b.Allow(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHalfOpenAfterTimeout(t *testing.T) {
	b := newBreaker(t, 1, 10*time.Millisecond)
	b.RecordFailure()
	time.Sleep(20 * time.Millisecond)
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil in half-open, got %v", err)
	}
	if b.State() != StateHalfOpen {
		t.Fatalf("expected StateHalfOpen, got %v", b.State())
	}
}

func TestHalfOpenSuccessCloses(t *testing.T) {
	b := newBreaker(t, 1, 10*time.Millisecond)
	b.RecordFailure()
	time.Sleep(20 * time.Millisecond)
	_ = b.Allow()
	b.RecordSuccess()
	if b.State() != StateClosed {
		t.Fatalf("expected StateClosed after half-open success")
	}
}
