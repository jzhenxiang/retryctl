package throttle_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/retryctl/internal/throttle"
)

func TestNewInvalidMax(t *testing.T) {
	_, err := throttle.New(0)
	if err == nil {
		t.Fatal("expected error for max=0")
	}
}

func TestNewValidMax(t *testing.T) {
	th, err := throttle.New(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if th.Max() != 3 {
		t.Fatalf("expected max=3, got %d", th.Max())
	}
}

func TestAvailableDecrements(t *testing.T) {
	th, _ := throttle.New(2)
	if th.Available() != 2 {
		t.Fatalf("expected 2 available, got %d", th.Available())
	}
	_ = th.Acquire()
	if th.Available() != 1 {
		t.Fatalf("expected 1 available after acquire, got %d", th.Available())
	}
	th.Release()
	if th.Available() != 2 {
		t.Fatalf("expected 2 available after release, got %d", th.Available())
	}
}

func TestConcurrencyLimit(t *testing.T) {
	const max = 3
	th, _ := throttle.New(max)

	var active atomic.Int32
	var peak atomic.Int32
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = th.Acquire()
			cur := active.Add(1)
			for {
				old := peak.Load()
				if cur <= old || peak.CompareAndSwap(old, cur) {
					break
				}
			}
			time.Sleep(5 * time.Millisecond)
			active.Add(-1)
			th.Release()
		}()
	}
	wg.Wait()

	if p := peak.Load(); p > int32(max) {
		t.Fatalf("peak concurrency %d exceeded max %d", p, max)
	}
}

func TestNilThrottleReleaseSafe(t *testing.T) {
	var th *throttle.Throttle
	th.Release() // must not panic
}

func TestNilThrottleAcquireError(t *testing.T) {
	var th *throttle.Throttle
	if err := th.Acquire(); err == nil {
		t.Fatal("expected error acquiring nil throttle")
	}
}
