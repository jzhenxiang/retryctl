package cli

import (
	"testing"

	"github.com/spf13/cobra"
)

func newCmdWithTagFlags(tags []string) *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().StringArray("tag", nil, "")
	for _, t := range tags {
		_ = cmd.Flags().Set("tag", t)
	}
	return cmd
}

func TestBuildTaggerNoFlags(t *testing.T) {
	cmd := newCmdWithTagFlags(nil)
	tgr, err := buildTagger(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tgr.Len() != 0 {
		t.Fatalf("expected 0 tags, got %d", tgr.Len())
	}
}

func TestBuildTaggerValidPairs(t *testing.T) {
	cmd := newCmdWithTagFlags([]string{"env=prod", "region=us-east-1"})
	tgr, err := buildTagger(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tgr.Len() != 2 {
		t.Fatalf("expected 2 tags, got %d", tgr.Len())
	}
	if !tgr.Has("env") {
		t.Fatal("expected tag 'env' to be present")
	}
	if !tgr.Has("region") {
		t.Fatal("expected tag 'region' to be present")
	}
}

func TestBuildTaggerMalformedPair(t *testing.T) {
	cmd := newCmdWithTagFlags([]string{"notapair"})
	_, err := buildTagger(cmd)
	if err == nil {
		t.Fatal("expected error for malformed pair")
	}
}

func TestBuildTaggerEmptyKey(t *testing.T) {
	cmd := newCmdWithTagFlags([]string{"=value"})
	_, err := buildTagger(cmd)
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}
