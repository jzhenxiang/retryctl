package cli

import (
	"testing"

	"github.com/spf13/cobra"
)

func newCmdWithThrottleFlag() *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().Int("throttle", 0, "")
	return cmd
}

func TestBuildThrottleDisabled(t *testing.T) {
	cmd := newCmdWithThrottleFlag()
	th, err := buildThrottle(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if th != nil {
		t.Fatal("expected nil throttle when flag is 0")
	}
}

func TestBuildThrottleEnabled(t *testing.T) {
	cmd := newCmdWithThrottleFlag()
	_ = cmd.Flags().Set("throttle", "4")
	th, err := buildThrottle(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if th == nil {
		t.Fatal("expected non-nil throttle")
	}
	if th.Max() != 4 {
		t.Fatalf("expected max=4, got %d", th.Max())
	}
	if th.Available() != 4 {
		t.Fatalf("expected 4 available slots, got %d", th.Available())
	}
}

func TestBuildThrottleNegativeDisabled(t *testing.T) {
	cmd := newCmdWithThrottleFlag()
	_ = cmd.Flags().Set("throttle", "-1")
	th, err := buildThrottle(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if th != nil {
		t.Fatal("expected nil throttle for negative value")
	}
}

func TestBuildThrottleOne(t *testing.T) {
	cmd := newCmdWithThrottleFlag()
	_ = cmd.Flags().Set("throttle", "1")
	th, err := buildThrottle(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if th == nil {
		t.Fatal("expected non-nil throttle for max=1")
	}
}
