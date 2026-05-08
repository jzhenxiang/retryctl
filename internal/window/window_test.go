package window

import (
	"testing"
	"time"
)

func TestNewInvalidSize(t *testing.T) {
	_, err := New(0, 10)
	if err == nil {
		t.Fatal("expected error for zero size")
	}
}

func TestNewInvalidBuckets(t *testing.T) {
	_, err := New(time.Second, 0)
	if err == nil {
		t.Fatal("expected error for zero buckets")
	}
}

func TestNewValid(t *testing.T) {
	w, err := New(time.Second, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w == nil {
		t.Fatal("expected non-nil window")
	}
}

func TestCountZeroInitially(t *testing.T) {
	w, _ := New(time.Second, 10)
	if got := w.Count(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestRecordIncreasesCount(t *testing.T) {
	w, _ := New(time.Second, 10)
	w.Record()
	w.Record()
	w.Record()
	if got := w.Count(); got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}
}

func TestResetClearsCount(t *testing.T) {
	w, _ := New(time.Second, 10)
	w.Record()
	w.Record()
	w.Reset()
	if got := w.Count(); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestExpiredEventsNotCounted(t *testing.T) {
	now := time.Now()
	calls := 0
	clock := func() time.Time {
		calls++
		if calls <= 2 {
			// first two calls: record events in the past
			return now.Add(-2 * time.Second)
		}
		// subsequent calls: current time (window has passed)
		return now
	}
	w, _ := newWithClock(time.Second, 10, clock)
	w.Record() // recorded in the past
	w.Record() // recorded in the past
	if got := w.Count(); got != 0 {
		t.Fatalf("expected 0 for expired events, got %d", got)
	}
}

func TestConcurrentRecords(t *testing.T) {
	w, _ := New(5*time.Second, 10)
	const goroutines = 20
	done := make(chan struct{}, goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			w.Record()
			done <- struct{}{}
		}()
	}
	for i := 0; i < goroutines; i++ {
		<-done
	}
	if got := w.Count(); got != goroutines {
		t.Fatalf("expected %d, got %d", goroutines, got)
	}
}
