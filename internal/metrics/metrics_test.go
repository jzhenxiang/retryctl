package metrics_test

import (
	"errors"
	"testing"
	"time"

	"github.com/user/retryctl/internal/metrics"
)

func TestNewSummaryIsZero(t *testing.T) {
	s := metrics.New()
	if s.Attempts != 0 || s.Successes != 0 || s.Failures != 0 {
		t.Fatal("expected zeroed summary")
	}
}

func TestRecordSuccess(t *testing.T) {
	s := metrics.New()
	s.RecordAttempt(10*time.Millisecond, nil)

	if s.Attempts != 1 {
		t.Fatalf("expected 1 attempt, got %d", s.Attempts)
	}
	if s.Successes != 1 {
		t.Fatalf("expected 1 success, got %d", s.Successes)
	}
	if s.Failures != 0 {
		t.Fatalf("expected 0 failures, got %d", s.Failures)
	}
	if s.LastError != nil {
		t.Fatalf("expected nil LastError, got %v", s.LastError)
	}
}

func TestRecordFailure(t *testing.T) {
	s := metrics.New()
	err := errors.New("boom")
	s.RecordAttempt(5*time.Millisecond, err)

	if s.Failures != 1 {
		t.Fatalf("expected 1 failure, got %d", s.Failures)
	}
	if !errors.Is(s.LastError, err) {
		t.Fatalf("expected LastError to be %v, got %v", err, s.LastError)
	}
}

func TestAverageElapsed(t *testing.T) {
	s := metrics.New()
	s.RecordAttempt(20*time.Millisecond, nil)
	s.RecordAttempt(40*time.Millisecond, nil)

	want := 30 * time.Millisecond
	if got := s.AverageElapsed(); got != want {
		t.Fatalf("expected average %v, got %v", want, got)
	}
}

func TestAverageElapsedNoAttempts(t *testing.T) {
	s := metrics.New()
	if got := s.AverageElapsed(); got != 0 {
		t.Fatalf("expected 0 for empty summary, got %v", got)
	}
}

func TestSnapshot(t *testing.T) {
	s := metrics.New()
	s.RecordAttempt(10*time.Millisecond, errors.New("err"))

	snap := s.Snapshot()
	if snap.Attempts != s.Attempts || snap.Failures != s.Failures {
		t.Fatal("snapshot does not match source")
	}
}
