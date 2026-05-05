package filter_test

import (
	"testing"

	"github.com/user/retryctl/internal/filter"
)

// --- ExitCodeMatcher ---

func TestExitCodeMatcherRetryAll(t *testing.T) {
	m := filter.NewExitCodeMatcher(nil)
	if m.ShouldRetry(0) {
		t.Error("exit 0 must not trigger retry")
	}
	for _, code := range []int{1, 2, 127, 255} {
		if !m.ShouldRetry(code) {
			t.Errorf("expected exit %d to be retryable", code)
		}
	}
}

func TestExitCodeMatcherSpecificCodes(t *testing.T) {
	m := filter.NewExitCodeMatcher([]int{1, 42})
	if m.ShouldRetry(0) {
		t.Error("exit 0 must not trigger retry")
	}
	if !m.ShouldRetry(1) {
		t.Error("exit 1 should be retryable")
	}
	if !m.ShouldRetry(42) {
		t.Error("exit 42 should be retryable")
	}
	if m.ShouldRetry(2) {
		t.Error("exit 2 should not be retryable")
	}
	if m.ShouldRetry(127) {
		t.Error("exit 127 should not be retryable")
	}
}

// --- OutputFilter ---

func TestNewOutputFilterInvalidPattern(t *testing.T) {
	_, err := filter.NewOutputFilter([]string{"["})
	if err == nil {
		t.Error("expected error for invalid regex")
	}
}

func TestOutputFilterNoPatterns(t *testing.T) {
	f, err := filter.NewOutputFilter(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.MatchesAbort("fatal: something went wrong") {
		t.Error("no patterns should never match")
	}
}

func TestOutputFilterMatch(t *testing.T) {
	f, err := filter.NewOutputFilter([]string{`fatal:`, `permission denied`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.MatchesAbort("fatal: disk full") {
		t.Error("expected match on 'fatal:'")
	}
	if !f.MatchesAbort("bash: permission denied") {
		t.Error("expected match on 'permission denied'")
	}
	if f.MatchesAbort("temporary network error") {
		t.Error("should not match unrelated output")
	}
}
