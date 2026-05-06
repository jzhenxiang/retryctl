package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func newCmdWithAuditFlag() *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("audit-log", "", "")
	return cmd
}

func TestBuildAuditRecorderDisabled(t *testing.T) {
	cmd := newCmdWithAuditFlag()
	rec, closer, err := buildAuditRecorder(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec != nil {
		t.Error("expected nil recorder when flag is unset")
	}
	if closer != nil {
		t.Error("expected nil closer when flag is unset")
	}
}

func TestBuildAuditRecorderEnabled(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.jsonl")

	cmd := newCmdWithAuditFlag()
	if err := cmd.Flags().Set("audit-log", path); err != nil {
		t.Fatalf("set flag: %v", err)
	}

	rec, closer, err := buildAuditRecorder(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec == nil {
		t.Fatal("expected non-nil recorder")
	}
	if closer == nil {
		t.Fatal("expected non-nil closer")
	}
	closer.Close()

	if _, err := os.Stat(path); err != nil {
		t.Errorf("audit file not created: %v", err)
	}
}

func TestBuildAuditRecorderBadPath(t *testing.T) {
	cmd := newCmdWithAuditFlag()
	if err := cmd.Flags().Set("audit-log", "/no/such/dir/audit.jsonl"); err != nil {
		t.Fatalf("set flag: %v", err)
	}

	_, _, err := buildAuditRecorder(cmd)
	if err == nil {
		t.Error("expected error for unwritable path, got nil")
	}
}
