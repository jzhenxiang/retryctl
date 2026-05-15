package cli

import (
	"testing"

	"github.com/spf13/cobra"
)

func newCmdWithLabelFlags() *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().StringSlice("label", nil, "")
	return cmd
}

func TestBuildLabelerNoFlags(t *testing.T) {
	cmd := newCmdWithLabelFlags()
	l, err := buildLabeler(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(l.Base()) != 0 {
		t.Errorf("expected empty labels")
	}
}

func TestBuildLabelerValidPairs(t *testing.T) {
	cmd := newCmdWithLabelFlags()
	_ = cmd.Flags().Set("label", "env=prod")
	_ = cmd.Flags().Set("label", "team=sre")
	l, err := buildLabeler(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	base := l.Base()
	if base["env"] != "prod" {
		t.Errorf("env: got %q want prod", base["env"])
	}
	if base["team"] != "sre" {
		t.Errorf("team: got %q want sre", base["team"])
	}
}

func TestBuildLabelerMalformedPair(t *testing.T) {
	cmd := newCmdWithLabelFlags()
	_ = cmd.Flags().Set("label", "nodash")
	_, err := buildLabeler(cmd)
	if err == nil {
		t.Fatal("expected error for malformed label pair, got nil")
	}
}

func TestBuildLabelerWithMerge(t *testing.T) {
	cmd := newCmdWithLabelFlags()
	_ = cmd.Flags().Set("label", "env=prod")
	l, err := buildLabeler(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	merged := l.With(map[string]string{"attempt": "2"})
	if merged["env"] != "prod" {
		t.Errorf("env: got %q want prod", merged["env"])
	}
	if merged["attempt"] != "2" {
		t.Errorf("attempt: got %q want 2", merged["attempt"])
	}
}
