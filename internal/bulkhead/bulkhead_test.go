package bulkhead_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/user/retryctl/internal/bulkhead"
)

func TestNewEmptyLimitsReturnsError(t *testing.T) {
	_, err := bulkhead.New(map[string]int{})
	if err == nil {
		t.Fatal("expected error for empty limits")
	}
}

func TestNewZeroMaxReturnsError(t *testing.T) {
	_, err := bulkhead.New(map[string]int{"db": 0})
	if err == nil {
		t.Fatal("expected error for zero max")
	}
}

func TestNewNegativeMaxReturnsError(t *testing.T) {
	_, err := bulkhead.New(map[string]int{"db": -1})
	if err == nil {
		t.Fatal("expected error for negative max")
	}
}

func TestAcquireUnknownPartition(t *testing.T) {
	b, _ := bulkhead.New(map[string]int{"db": 2})
	err := b.Acquire("cache")
	if !errors.Is(err, bulkhead.ErrUnknownPartition) {
		t.Fatalf("expected ErrUnknownPartition, got %v", err)
	}
}

func TestAcquireAndRelease(t *testing.T) {
	b, _ := bulkhead.New(map[string]int{"db": 2})
	if err := b.Acquire("db"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := b.Available("db"); got != 1 {
		t.Fatalf("expected 1 available, got %d", got)
	}
	_ = b.Release("db")
	if got := b.Available("db"); got != 2 {
		t.Fatalf("expected 2 available after release, got %d", got)
	}
}

func TestAcquireExceedsLimit(t *testing.T) {
	b, _ := bulkhead.New(map[string]int{"db": 1})
	if err := b.Acquire("db"); err != nil {
		t.Fatalf("first acquire failed: %v", err)
	}
	err := b.Acquire("db")
	if !errors.Is(err, bulkhead.ErrPartitionFull) {
		t.Fatalf("expected ErrPartitionFull, got %v", err)
	}
}

func TestAvailableUnknownReturnsNegativeOne(t *testing.T) {
	b, _ := bulkhead.New(map[string]int{"db": 2})
	if got := b.Available("missing"); got != -1 {
		t.Fatalf("expected -1 for unknown partition, got %d", got)
	}
}

func TestConcurrentAcquireRespectLimit(t *testing.T) {
	const max = 5
	b, _ := bulkhead.New(map[string]int{"svc": max})
	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		allowed int
		denied  int
	)
	for i := 0; i < max*2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := b.Acquire("svc")
			mu.Lock()
			defer mu.Unlock()
			if err == nil {
				allowed++
			} else {
				denied++
			}
		}()
	}
	wg.Wait()
	if allowed > max {
		t.Fatalf("allowed %d exceeds max %d", allowed, max)
	}
	if allowed+denied != max*2 {
		t.Fatalf("expected %d total attempts, got %d", max*2, allowed+denied)
	}
}
