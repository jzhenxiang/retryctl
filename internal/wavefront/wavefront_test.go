package wavefront_test

import (
	"testing"
	"time"

	"github.com/yourorg/retryctl/internal/wavefront"
)

func TestNewInvalidWindow(t *testing.T) {
	_, err := wavefront.New(0, 0.5)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNewInvalidThresholdNegative(t *testing.T) {
	_, err := wavefront.New(time.Second, -0.1)
	if err == nil {
		t.Fatal("expected error for negative threshold")
	}
}

func TestNewInvalidThresholdAboveOne(t *testing.T) {
	_, err := wavefront.New(time.Second, 1.1)
	if err == nil {
		t.Fatal("expected error for threshold > 1")
	}
}

func TestNewValid(t *testing.T) {
	tr, err := wavefront.New(time.Second, 0.5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr == nil {
		t.Fatal("expected non-nil tracker")
	}
}

func TestFailureRateNoEvents(t *testing.T) {
	tr, _ := wavefront.New(time.Second, 0.5)
	if rate := tr.FailureRate(); rate != 0 {
		t.Fatalf("expected 0, got %f", rate)
	}
}

func TestFailureRateAllSuccess(t *testing.T) {
	tr, _ := wavefront.New(time.Minute, 0.5)
	tr.Record(false)
	tr.Record(false)
	if rate := tr.FailureRate(); rate != 0 {
		t.Fatalf("expected 0, got %f", rate)
	}
}

func TestFailureRateAllFailures(t *testing.T) {
	tr, _ := wavefront.New(time.Minute, 0.5)
	tr.Record(true)
	tr.Record(true)
	if rate := tr.FailureRate(); rate != 1.0 {
		t.Fatalf("expected 1.0, got %f", rate)
	}
}

func TestFailureRateMixed(t *testing.T) {
	tr, _ := wavefront.New(time.Minute, 0.5)
	tr.Record(true)
	tr.Record(false)
	tr.Record(true)
	tr.Record(false)
	rate := tr.FailureRate()
	if rate != 0.5 {
		t.Fatalf("expected 0.5, got %f", rate)
	}
}

func TestExceedsThreshold(t *testing.T) {
	tr, _ := wavefront.New(time.Minute, 0.5)
	tr.Record(true)
	tr.Record(true)
	tr.Record(false)
	// rate = 0.667 > 0.5
	if !tr.ExceedsThreshold() {
		t.Fatal("expected threshold to be exceeded")
	}
}

func TestDoesNotExceedThreshold(t *testing.T) {
	tr, _ := wavefront.New(time.Minute, 0.5)
	tr.Record(false)
	tr.Record(false)
	tr.Record(true)
	// rate = 0.333 <= 0.5
	if tr.ExceedsThreshold() {
		t.Fatal("expected threshold not to be exceeded")
	}
}
