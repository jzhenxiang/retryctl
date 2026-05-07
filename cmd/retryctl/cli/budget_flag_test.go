package cli

import (
	"bytes"
	"testing"
	"time"

	"github.com/spf13/cobra"
)

func newCmdWithBudgetFlags() *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().Int("budget-max", 0, "")
	cmd.Flags().Duration("budget-window", time.Minute, "")
	cmd.SetOut(&bytes.Buffer{})
	return cmd
}

func TestBuildBudgetDisabled(t *testing.T) {
	cmd := newCmdWithBudgetFlags()
	b, err := buildBudget(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b != nil {
		t.Fatal("expected nil budget when max=0")
	}
}

func TestBuildBudgetEnabled(t *testing.T) {
	cmd := newCmdWithBudgetFlags()
	_ = cmd.Flags().Set("budget-max", "5")
	_ = cmd.Flags().Set("budget-window", "30s")

	b, err := buildBudget(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b == nil {
		t.Fatal("expected non-nil budget")
	}
	if b.Remaining() != 5 {
		t.Fatalf("expected remaining=5, got %d", b.Remaining())
	}
}

func TestBuildBudgetInvalidWindow(t *testing.T) {
	cmd := newCmdWithBudgetFlags()
	_ = cmd.Flags().Set("budget-max", "3")
	_ = cmd.Flags().Set("budget-window", "0s")

	_, err := buildBudget(cmd)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}
