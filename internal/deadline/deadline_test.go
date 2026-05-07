package deadline

import (
	"testing"
	"time"
)

func TestNewZeroDeadlineReturnsError(t *testing.T) {
	_, err := New(time.Time{})
	if err == nil {
		t.Fatal("expected error for zero deadline")
	}
}

func TestNewPastDeadlineReturnsError(t *testing.T) {
	past := time.Now().Add(-1 * time.Second)
	_, err := New(past)
	if err == nil {
		t.Fatal("expected error for past deadline")
	}
}

func TestNewFutureDeadlineOK(t *testing.T) {
	future := time.Now().Add(10 * time.Second)
	g, err := New(future)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g == nil {
		t.Fatal("expected non-nil guard")
	}
}

func TestAllowBeforeDeadline(t *testing.T) {
	now := time.Now()
	g, _ := newWithClock(now.Add(5*time.Second), func() time.Time { return now })
	if err := g.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestAllowAfterDeadline(t *testing.T) {
	now := time.Now()
	g, _ := newWithClock(now.Add(5*time.Second), func() time.Time { return now })
	// advance clock past deadline
	g.now = func() time.Time { return now.Add(10 * time.Second) }
	if err := g.Allow(); err != ErrDeadlineExceeded {
		t.Fatalf("expected ErrDeadlineExceeded, got %v", err)
	}
}

func TestAllowAtExactDeadlineExceeds(t *testing.T) {
	now := time.Now()
	deadlineTime := now.Add(5 * time.Second)
	g, _ := newWithClock(deadlineTime, func() time.Time { return now })
	g.now = func() time.Time { return deadlineTime }
	if err := g.Allow(); err != ErrDeadlineExceeded {
		t.Fatalf("expected ErrDeadlineExceeded at exact deadline, got %v", err)
	}
}

func TestRemainingPositive(t *testing.T) {
	now := time.Now()
	g, _ := newWithClock(now.Add(5*time.Second), func() time.Time { return now })
	if r := g.Remaining(); r <= 0 {
		t.Fatalf("expected positive remaining, got %v", r)
	}
}

func TestRemainingNegativeAfterDeadline(t *testing.T) {
	now := time.Now()
	g, _ := newWithClock(now.Add(5*time.Second), func() time.Time { return now })
	g.now = func() time.Time { return now.Add(10 * time.Second) }
	if r := g.Remaining(); r >= 0 {
		t.Fatalf("expected negative remaining, got %v", r)
	}
}

func TestDeadlineReturnsConfiguredTime(t *testing.T) {
	now := time.Now()
	expected := now.Add(5 * time.Second)
	g, _ := newWithClock(expected, func() time.Time { return now })
	if !g.Deadline().Equal(expected) {
		t.Fatalf("expected %v, got %v", expected, g.Deadline())
	}
}
