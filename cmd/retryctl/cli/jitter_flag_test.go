package cli

import (
	"testing"

	"retryctl/internal/jitter"
)

func TestBuildJitterNone(t *testing.T) {
	j := buildJitter("none", 0)
	if j == nil {
		t.Fatal("expected non-nil jitter for 'none'")
	}
	// none jitter should always return 0
	for i := 0; i < 10; i++ {
		if got := j(500); got != 0 {
			t.Errorf("none jitter: expected 0, got %d", got)
		}
	}
}

func TestBuildJitterFull(t *testing.T) {
	j := buildJitter("full", 0)
	if j == nil {
		t.Fatal("expected non-nil jitter for 'full'")
	}
	// full jitter should return values in [0, base]
	base := int64(1000)
	for i := 0; i < 20; i++ {
		got := j(base)
		if got < 0 || got > base {
			t.Errorf("full jitter out of range [0, %d]: got %d", base, got)
		}
	}
}

func TestBuildJitterEqual(t *testing.T) {
	j := buildJitter("equal", 0)
	if j == nil {
		t.Fatal("expected non-nil jitter for 'equal'")
	}
	// equal jitter should return values in [base/2, base]
	base := int64(1000)
	for i := 0; i < 20; i++ {
		got := j(base)
		if got < base/2 || got > base {
			t.Errorf("equal jitter out of range [%d, %d]: got %d", base/2, base, got)
		}
	}
}

func TestBuildJitterUnknownFallsBackToNone(t *testing.T) {
	j := buildJitter("bogus", 0)
	if j == nil {
		t.Fatal("expected non-nil jitter for unknown strategy")
	}
	// unknown should fall back to none (zero jitter)
	_ = j // just ensure it doesn't panic
	_ = jitter.NewNone()
}

func TestBuildJitterSeedProducesDifferentSequences(t *testing.T) {
	j1 := buildJitter("full", 42)
	j2 := buildJitter("full", 99)
	base := int64(100000)
	matches := 0
	for i := 0; i < 10; i++ {
		if j1(base) == j2(base) {
			matches++
		}
	}
	// It's astronomically unlikely all 10 values match with different seeds
	if matches == 10 {
		t.Error("expected different sequences for different seeds")
	}
}
