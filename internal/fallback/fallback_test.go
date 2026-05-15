package fallback_test

import (
	"context"
	"testing"

	"github.com/user/retryctl/internal/fallback"
)

func TestNewEmptyArgsReturnsError(t *testing.T) {
	_, err := fallback.New([]string{})
	if err == nil {
		t.Fatal("expected error for empty args, got nil")
	}
}

func TestNewValidArgs(t *testing.T) {
	h, err := fallback.New([]string{"echo", "ok"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestHasFallbackTrue(t *testing.T) {
	h, _ := fallback.New([]string{"echo"})
	if !h.HasFallback() {
		t.Fatal("expected HasFallback to return true")
	}
}

func TestHasFallbackNilHandler(t *testing.T) {
	var h *fallback.Handler
	if h.HasFallback() {
		t.Fatal("expected HasFallback to return false for nil handler")
	}
}

func TestRunSuccessExitZero(t *testing.T) {
	h, _ := fallback.New([]string{"echo", "hello"})
	res, err := h.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ExitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", res.ExitCode)
	}
	if len(res.Output) == 0 {
		t.Fatal("expected non-empty output")
	}
}

func TestRunNonZeroExitCode(t *testing.T) {
	h, _ := fallback.New([]string{"sh", "-c", "exit 3"})
	res, err := h.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ExitCode != 3 {
		t.Fatalf("expected exit code 3, got %d", res.ExitCode)
	}
}

func TestRunInvalidCommandReturnsError(t *testing.T) {
	h, _ := fallback.New([]string{"/nonexistent/binary/retryctl-fallback-test"})
	_, err := h.Run(context.Background())
	if err == nil {
		t.Fatal("expected error for invalid command, got nil")
	}
}

func TestRunCancelledContextReturnsError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	h, _ := fallback.New([]string{"sleep", "10"})
	_, err := h.Run(ctx)
	if err == nil {
		t.Fatal("expected error for cancelled context, got nil")
	}
}
