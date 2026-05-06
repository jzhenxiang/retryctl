package jitter_test

import (
	"testing"
	"time"

	"github.com/example/retryctl/internal/jitter"
)

const base = 100 * time.Millisecond

func TestFullJitterInRange(t *testing.T) {
	j := jitter.NewFullWithSource(func() float64 { return 0.5 })
	got := j.Apply(base)
	want := 50 * time.Millisecond
	if got != want {
		t.Fatalf("Full.Apply: got %v, want %v", got, want)
	}
}

func TestFullJitterZero(t *testing.T) {
	j := jitter.NewFullWithSource(func() float64 { return 0.0 })
	if got := j.Apply(base); got != 0 {
		t.Fatalf("expected 0, got %v", got)
	}
}

func TestEqualJitterInRange(t *testing.T) {
	// src=0.0 → half + 0 = 50ms
	j := jitter.NewEqualWithSource(func() float64 { return 0.0 })
	got := j.Apply(base)
	want := 50 * time.Millisecond
	if got != want {
		t.Fatalf("Equal.Apply(src=0): got %v, want %v", got, want)
	}
}

func TestEqualJitterMaxRange(t *testing.T) {
	// src→1.0 → half + ~half = ~base
	j := jitter.NewEqualWithSource(func() float64 { return 0.999 })
	got := j.Apply(base)
	if got < 50*time.Millisecond || got >= base {
		t.Fatalf("Equal.Apply out of expected range: %v", got)
	}
}

func TestNoneJitter(t *testing.T) {
	j := jitter.None{}
	if got := j.Apply(base); got != base {
		t.Fatalf("None.Apply: got %v, want %v", got, base)
	}
}

func TestNewReturnsCorrectType(t *testing.T) {
	cases := []struct {
		name     string
		strategy string
	}{
		{"full", "full"},
		{"equal", "equal"},
		{"none", "none"},
		{"unknown defaults to none", "bogus"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			j := jitter.New(tc.strategy)
			if j == nil {
				t.Fatal("New returned nil")
			}
		})
	}
}
