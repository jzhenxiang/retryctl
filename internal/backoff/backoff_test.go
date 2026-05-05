package backoff_test

import (
	"testing"
	"time"

	"github.com/yourorg/retryctl/internal/backoff"
)

func TestFixedStrategy(t *testing.T) {
	s := &backoff.FixedStrategy{Delay: 2 * time.Second}
	for attempt := 0; attempt < 5; attempt++ {
		if got := s.Next(attempt); got != 2*time.Second {
			t.Errorf("attempt %d: expected 2s, got %v", attempt, got)
		}
	}
}

func TestExponentialStrategy(t *testing.T) {
	s := &backoff.ExponentialStrategy{
		InitialDelay: time.Second,
		Multiplier:   2.0,
		MaxDelay:     10 * time.Second,
	}

	expected := []time.Duration{
		1 * time.Second,
		2 * time.Second,
		4 * time.Second,
		8 * time.Second,
		10 * time.Second, // capped
	}

	for i, want := range expected {
		if got := s.Next(i); got != want {
			t.Errorf("attempt %d: expected %v, got %v", i, want, got)
		}
	}
}

func TestLinearStrategy(t *testing.T) {
	s := &backoff.LinearStrategy{
		InitialDelay: time.Second,
		Increment:    time.Second,
		MaxDelay:     4 * time.Second,
	}

	expected := []time.Duration{
		1 * time.Second,
		2 * time.Second,
		3 * time.Second,
		4 * time.Second, // capped
		4 * time.Second, // capped
	}

	for i, want := range expected {
		if got := s.Next(i); got != want {
			t.Errorf("attempt %d: expected %v, got %v", i, want, got)
		}
	}
}

func TestNewStrategy(t *testing.T) {
	cases := []struct {
		name    string
		expType string
	}{
		{"fixed", "*backoff.FixedStrategy"},
		{"exponential", "*backoff.ExponentialStrategy"},
		{"linear", "*backoff.LinearStrategy"},
		{"unknown", "*backoff.FixedStrategy"},
	}

	for _, tc := range cases {
		s := backoff.NewStrategy(tc.name, time.Second, 10*time.Second)
		if s == nil {
			t.Errorf("NewStrategy(%q) returned nil", tc.name)
		}
	}
}
