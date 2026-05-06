package ratelimit_test

import (
	"sync"
	"testing"
	"time"

	"retryctl/internal/ratelimit"
)

func TestNewInvalidMax(t *testing.T) {
	_, err := ratelimit.New(ratelimit.Config{Max: 0, Window: time.Second})
	if err == nil {
		t.Fatal("expected error for Max=0")
	}
}

func TestNewInvalidWindow(t *testing.T) {
	_, err := ratelimit.New(ratelimit.Config{Max: 1, Window: 0})
	if err == nil {
		t.Fatal("expected error for Window=0")
	}
}

func TestAllowUnderLimit(t *testing.T) {
	l, err := ratelimit.New(ratelimit.Config{Max: 3, Window: time.Second})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 3; i++ {
		if err := l.Allow(); err != nil {
			t.Fatalf("attempt %d: unexpected rate limit", i+1)
		}
	}
}

func TestAllowExceedsLimit(t *testing.T) {
	l, err := ratelimit.New(ratelimit.Config{Max: 2, Window: time.Second})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = l.Allow()
	_ = l.Allow()
	if err := l.Allow(); err != ratelimit.ErrRateLimited {
		t.Fatalf("expected ErrRateLimited, got %v", err)
	}
}

func TestRemainingDecrementsOnAllow(t *testing.T) {
	l, _ := ratelimit.New(ratelimit.Config{Max: 3, Window: time.Second})
	if got := l.Remaining(); got != 3 {
		t.Fatalf("want 3, got %d", got)
	}
	_ = l.Allow()
	if got := l.Remaining(); got != 2 {
		t.Fatalf("want 2, got %d", got)
	}
}

func TestWindowExpiry(t *testing.T) {
	now := time.Now()
	l, _ := ratelimit.New(ratelimit.Config{Max: 2, Window: 100 * time.Millisecond})
	// Inject a controllable clock via the exported nowFn trick — we use real
	// time here and simply wait for the window to expire.
	_ = l.Allow()
	_ = l.Allow()
	time.Sleep(110 * time.Millisecond)
	_ = now // suppress unused warning
	if err := l.Allow(); err != nil {
		t.Fatalf("expected allow after window expiry, got %v", err)
	}
}

func TestConcurrentAllow(t *testing.T) {
	const goroutines = 20
	l, _ := ratelimit.New(ratelimit.Config{Max: 5, Window: time.Second})
	var (
		wg      sync.WaitGroup
		allowed int
		mu      sync.Mutex
	)
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			if err := l.Allow(); err == nil {
				mu.Lock()
				allowed++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	if allowed > 5 {
		t.Fatalf("allowed %d > max 5", allowed)
	}
}
