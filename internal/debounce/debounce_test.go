package debounce_test

import (
	"testing"
	"time"

	"github.com/user/retryctl/internal/debounce"
)

func TestNewInvalidWindow(t *testing.T) {
	_, err := debounce.New(0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
	_, err = debounce.New(-time.Second)
	if err == nil {
		t.Fatal("expected error for negative window")
	}
}

func TestFirstAllowAlwaysPasses(t *testing.T) {
	d, err := debounce.New(100 * time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := d.Allow(); err != nil {
		t.Fatalf("expected first Allow to pass, got: %v", err)
	}
}

func TestSecondAllowWithinWindowIsDebounced(t *testing.T) {
	now := time.Now()
	d, _ := debounce.New(500 * time.Millisecond)
	// inject controlled clock
	setClock(d, func() time.Time { return now })

	if err := d.Allow(); err != nil {
		t.Fatalf("first allow failed: %v", err)
	}
	// advance by less than window
	setClock(d, func() time.Time { return now.Add(100 * time.Millisecond) })

	if err := d.Allow(); err != debounce.ErrDebounced {
		t.Fatalf("expected ErrDebounced, got: %v", err)
	}
}

func TestAllowAfterWindowPasses(t *testing.T) {
	now := time.Now()
	d, _ := debounce.New(200 * time.Millisecond)
	setClock(d, func() time.Time { return now })

	d.Allow() //nolint:errcheck
	setClock(d, func() time.Time { return now.Add(300 * time.Millisecond) })

	if err := d.Allow(); err != nil {
		t.Fatalf("expected allow after window, got: %v", err)
	}
}

func TestResetClearsState(t *testing.T) {
	now := time.Now()
	d, _ := debounce.New(500 * time.Millisecond)
	setClock(d, func() time.Time { return now })
	d.Allow() //nolint:errcheck

	d.Reset()
	// still at same time — should pass because reset cleared lastSeen
	if err := d.Allow(); err != nil {
		t.Fatalf("expected allow after reset, got: %v", err)
	}
}

func TestRemainingIsZeroBeforeFirstAllow(t *testing.T) {
	d, _ := debounce.New(100 * time.Millisecond)
	if r := d.Remaining(); r != 0 {
		t.Fatalf("expected 0 remaining before first allow, got %v", r)
	}
}

func TestRemainingPositiveWithinWindow(t *testing.T) {
	now := time.Now()
	d, _ := debounce.New(500 * time.Millisecond)
	setClock(d, func() time.Time { return now })
	d.Allow() //nolint:errcheck
	setClock(d, func() time.Time { return now.Add(100 * time.Millisecond) })

	if r := d.Remaining(); r <= 0 {
		t.Fatalf("expected positive remaining, got %v", r)
	}
}

// setClock uses reflection-free approach via exported test hook.
func setClock(d *debounce.Debouncer, fn func() time.Time) {
	d.SetNow(fn)
}
