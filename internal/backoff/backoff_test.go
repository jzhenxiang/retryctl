package backoff_test

import (
	"testing"
	"time"

	"github.com/your-org/retryctl/internal/backoff"
)

func TestFixedStrategy(t *testing.T) {
	s := backoff.Fixed{Delay_: 5 * time.Second}
	for _, attempt := range []int{1, 2, 10} {
		if got := s.Delay(attempt); got != 5*time.Second {
			t.Errorf("attempt %d: got %v, want 5s", attempt, got)
		}
	}
}

func TestLinearStrategy(t *testing.T) {
	s := backoff.Linear{Base: 2 * time.Second}
	cases := []struct {
		attempt int
		want    time.Duration
	}{
		{1, 2 * time.Second},
		{2, 4 * time.Second},
		{5, 10 * time.Second},
	}
	for _, tc := range cases {
		if got := s.Delay(tc.attempt); got != tc.want {
			t.Errorf("attempt %d: got %v, want %v", tc.attempt, got, tc.want)
		}
	}
}

func TestExponentialStrategy(t *testing.T) {
	s := backoff.Exponential{Base: time.Second, MaxDelay: 0}
	cases := []struct {
		attempt int
		want    time.Duration
	}{
		{1, 1 * time.Second},
		{2, 2 * time.Second},
		{3, 4 * time.Second},
		{4, 8 * time.Second},
	}
	for _, tc := range cases {
		if got := s.Delay(tc.attempt); got != tc.want {
			t.Errorf("attempt %d: got %v, want %v", tc.attempt, got, tc.want)
		}
	}
}

func TestExponentialStrategyNeverExceedsMaxDelay(t *testing.T) {
	max := 10 * time.Second
	s := backoff.Exponential{Base: time.Second, MaxDelay: max}
	for attempt := 1; attempt <= 20; attempt++ {
		if got := s.Delay(attempt); got > max {
			t.Errorf("attempt %d: got %v, exceeds max %v", attempt, got, max)
		}
	}
}

func TestNewStrategy(t *testing.T) {
	cases := []struct {
		name    string
		wantErr bool
	}{
		{"fixed", false},
		{"linear", false},
		{"exponential", false},
		{"unknown", true},
		{"", true},
	}
	for _, tc := range cases {
		_, err := backoff.NewStrategy(tc.name, time.Second, 30*time.Second)
		if (err != nil) != tc.wantErr {
			t.Errorf("NewStrategy(%q): wantErr=%v, got err=%v", tc.name, tc.wantErr, err)
		}
	}
}

func TestExponentialStrategyAttemptOne(t *testing.T) {
	s := backoff.Exponential{Base: 500 * time.Millisecond, MaxDelay: 5 * time.Second}
	if got := s.Delay(1); got != 500*time.Millisecond {
		t.Errorf("got %v, want 500ms", got)
	}
}
