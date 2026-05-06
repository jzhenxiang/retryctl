package cli

import (
	"testing"
	"time"

	"github.com/spf13/cobra"
)

func newCmdWithCircuitFlags() *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().Int("circuit-max-failures", 0, "")
	cmd.Flags().Duration("circuit-reset-timeout", 30*time.Second, "")
	return cmd
}

func TestBuildCircuitBreakerDisabled(t *testing.T) {
	cmd := newCmdWithCircuitFlags()
	br, err := buildCircuitBreaker(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if br != nil {
		t.Fatal("expected nil breaker when max-failures=0")
	}
}

func TestBuildCircuitBreakerEnabled(t *testing.T) {
	cmd := newCmdWithCircuitFlags()
	_ = cmd.Flags().Set("circuit-max-failures", "3")
	_ = cmd.Flags().Set("circuit-reset-timeout", "10s")

	br, err := buildCircuitBreaker(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if br == nil {
		t.Fatal("expected non-nil breaker")
	}
}

func TestBuildCircuitBreakerDefaultTimeout(t *testing.T) {
	cmd := newCmdWithCircuitFlags()
	_ = cmd.Flags().Set("circuit-max-failures", "5")

	br, err := buildCircuitBreaker(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if br == nil {
		t.Fatal("expected non-nil breaker with default timeout")
	}
}

func TestBuildCircuitBreakerInvalidTimeout(t *testing.T) {
	cmd := newCmdWithCircuitFlags()
	// Override reset-timeout with a zero value by re-registering (simulate bad config)
	cmd2 := &cobra.Command{Use: "test"}
	cmd2.Flags().Int("circuit-max-failures", 5, "")
	cmd2.Flags().Duration("circuit-reset-timeout", 0, "")

	_, err := buildCircuitBreaker(cmd2)
	if err == nil {
		t.Fatal("expected error for zero reset-timeout with max-failures>0")
	}
}
