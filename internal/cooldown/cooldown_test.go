package cooldown

import (
	"testing"
	"time"
)

func newCooldown(t *testing.T, windows map[int]time.Duration) *Cooldown {
	t.Helper()
	c, err := New(windows)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return c
}

func TestNewEmptyWindowsReturnsError(t *testing.T) {
	_, err := New(map[int]time.Duration{})
	if err == nil {
		t.Fatal("expected error for empty windows map")
	}
}

func TestNewNonPositiveWindowReturnsError(t *testing.T) {
	_, err := New(map[int]time.Duration{1: 0})
	if err == nil {
		t.Fatal("expected error for zero-duration window")
	}
}

func TestAllowUnknownCodeAlwaysPasses(t *testing.T) {
	c := newCooldown(t, map[int]time.Duration{1: time.Second})
	if !c.Allow(2) {
		t.Error("expected unknown exit code to be allowed")
	}
}

func TestFirstAllowAlwaysPasses(t *testing.T) {
	c := newCooldown(t, map[int]time.Duration{1: time.Second})
	if !c.Allow(1) {
		t.Error("expected first allow to pass")
	}
}

func TestSecondAllowWithinWindowIsBlocked(t *testing.T) {
	now := time.Now()
	c := newCooldown(t, map[int]time.Duration{1: time.Second})
	c.clock = func() time.Time { return now }
	c.Allow(1) // record timestamp
	if c.Allow(1) {
		t.Error("expected second allow within window to be blocked")
	}
}

func TestAllowAfterWindowPasses(t *testing.T) {
	now := time.Now()
	c := newCooldown(t, map[int]time.Duration{1: time.Second})
	c.clock = func() time.Time { return now }
	c.Allow(1)
	c.clock = func() time.Time { return now.Add(2 * time.Second) }
	if !c.Allow(1) {
		t.Error("expected allow after window expiry to pass")
	}
}

func TestResetClearsState(t *testing.T) {
	now := time.Now()
	c := newCooldown(t, map[int]time.Duration{1: time.Second})
	c.clock = func() time.Time { return now }
	c.Allow(1)
	c.Reset(1)
	if !c.Allow(1) {
		t.Error("expected allow after reset to pass")
	}
}

func TestRemainingZeroWhenNotSeen(t *testing.T) {
	c := newCooldown(t, map[int]time.Duration{1: time.Second})
	if r := c.Remaining(1); r != 0 {
		t.Errorf("expected 0 remaining, got %v", r)
	}
}

func TestRemainingPositiveWithinWindow(t *testing.T) {
	now := time.Now()
	c := newCooldown(t, map[int]time.Duration{1: time.Second})
	c.clock = func() time.Time { return now }
	c.Allow(1)
	c.clock = func() time.Time { return now.Add(400 * time.Millisecond) }
	if r := c.Remaining(1); r <= 0 {
		t.Errorf("expected positive remaining, got %v", r)
	}
}

func TestRemainingZeroAfterWindowExpires(t *testing.T) {
	now := time.Now()
	c := newCooldown(t, map[int]time.Duration{1: time.Second})
	c.clock = func() time.Time { return now }
	c.Allow(1)
	c.clock = func() time.Time { return now.Add(2 * time.Second) }
	if r := c.Remaining(1); r != 0 {
		t.Errorf("expected 0 remaining after expiry, got %v", r)
	}
}
