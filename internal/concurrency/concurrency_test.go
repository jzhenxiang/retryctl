package concurrency_test

import (
	"sync"
	"testing"

	"github.com/yourorg/retryctl/internal/concurrency"
)

func TestNewInvalidMax(t *testing.T) {
	_, err := concurrency.New(0)
	if err == nil {
		t.Fatal("expected error for max=0")
	}
	_, err = concurrency.New(-1)
	if err == nil {
		t.Fatal("expected error for max=-1")
	}
}

func TestNewValidMax(t *testing.T) {
	g, err := concurrency.New(2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g == nil {
		t.Fatal("expected non-nil guard")
	}
}

func TestAcquireAndRelease(t *testing.T) {
	g, _ := concurrency.New(2)

	ok, rel := g.Acquire()
	if !ok {
		t.Fatal("expected first acquire to succeed")
	}
	if g.Active() != 1 {
		t.Fatalf("expected active=1, got %d", g.Active())
	}
	rel()
	if g.Active() != 0 {
		t.Fatalf("expected active=0 after release, got %d", g.Active())
	}
}

func TestAcquireBlocksAtLimit(t *testing.T) {
	g, _ := concurrency.New(1)

	ok1, rel1 := g.Acquire()
	if !ok1 {
		t.Fatal("expected first acquire to succeed")
	}
	defer rel1()

	ok2, rel2 := g.Acquire()
	if ok2 {
		t.Fatal("expected second acquire to fail when limit reached")
	}
	if rel2 != nil {
		t.Fatal("expected nil release when acquire fails")
	}
}

func TestAvailableDecrementsOnAcquire(t *testing.T) {
	g, _ := concurrency.New(3)

	if g.Available() != 3 {
		t.Fatalf("expected available=3, got %d", g.Available())
	}
	_, rel := g.Acquire()
	defer rel()
	if g.Available() != 2 {
		t.Fatalf("expected available=2, got %d", g.Available())
	}
}

func TestConcurrentAcquire(t *testing.T) {
	const max = 5
	g, _ := concurrency.New(max)

	var wg sync.WaitGroup
	var mu sync.Mutex
	successful := 0

	for i := 0; i < max*2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ok, rel := g.Acquire()
			if ok {
				mu.Lock()
				successful++
				mu.Unlock()
				rel()
			}
		}()
	}
	wg.Wait()

	if successful == 0 {
		t.Fatal("expected at least one successful acquire")
	}
	if g.Active() != 0 {
		t.Fatalf("expected active=0 after all goroutines finish, got %d", g.Active())
	}
}
