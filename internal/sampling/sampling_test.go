package sampling_test

import (
	"math/rand"
	"testing"

	"retryctl/internal/sampling"
)

func TestNewInvalidRateZero(t *testing.T) {
	_, err := sampling.New(0, nil)
	if err == nil {
		t.Fatal("expected error for rate=0")
	}
}

func TestNewInvalidRateAboveOne(t *testing.T) {
	_, err := sampling.New(1.1, nil)
	if err == nil {
		t.Fatal("expected error for rate=1.1")
	}
}

func TestNewValidRate(t *testing.T) {
	s, err := sampling.New(0.5, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Rate() != 0.5 {
		t.Fatalf("expected rate 0.5, got %v", s.Rate())
	}
}

func TestAllowAlwaysWhenRateOne(t *testing.T) {
	s, err := sampling.New(1.0, rand.NewSource(42))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 100; i++ {
		if !s.Allow() {
			t.Fatalf("expected Allow()=true at iteration %d with rate=1.0", i)
		}
	}
}

func TestAllowNeverWhenRateNearZero(t *testing.T) {
	// Use a deterministic source; with rate=1e-9 virtually all draws are >= rate.
	s, err := sampling.New(1e-9, rand.NewSource(99))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	allowed := 0
	for i := 0; i < 1000; i++ {
		if s.Allow() {
			allowed++
		}
	}
	if allowed > 5 {
		t.Fatalf("expected near-zero allows, got %d", allowed)
	}
}

func TestAllowApproximatesRate(t *testing.T) {
	const rate = 0.3
	const iterations = 10_000
	s, err := sampling.New(rate, rand.NewSource(7))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	allowed := 0
	for i := 0; i < iterations; i++ {
		if s.Allow() {
			allowed++
		}
	}
	actual := float64(allowed) / iterations
	if actual < 0.25 || actual > 0.35 {
		t.Fatalf("sampling rate out of expected range: got %.3f", actual)
	}
}
