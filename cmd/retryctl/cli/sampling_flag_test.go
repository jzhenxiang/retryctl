package cli

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func newCmdWithSamplingFlag() *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().Float64("sample-rate", 1.0, "")
	return cmd
}

func TestBuildSamplerDisabledAtOne(t *testing.T) {
	cmd := newCmdWithSamplingFlag()
	_ = cmd.Flags().Set("sample-rate", "1.0")
	s, err := buildSampler(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s != nil {
		t.Fatal("expected nil sampler when rate=1.0")
	}
}

func TestBuildSamplerEnabled(t *testing.T) {
	cmd := newCmdWithSamplingFlag()
	_ = cmd.Flags().Set("sample-rate", "0.5")
	s, err := buildSampler(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil sampler when rate=0.5")
	}
	if s.Rate() != 0.5 {
		t.Fatalf("expected rate 0.5, got %v", s.Rate())
	}
}

func TestBuildSamplerInvalidRate(t *testing.T) {
	cmd := newCmdWithSamplingFlag()
	_ = cmd.Flags().Set("sample-rate", "0.0")
	_, err := buildSampler(cmd)
	if err == nil {
		t.Fatal("expected error for rate=0.0")
	}
}

func TestBuildSamplerAboveOne(t *testing.T) {
	cmd := newCmdWithSamplingFlag()
	_ = cmd.Flags().Set("sample-rate", "1.5")
	// rate >= 1.0 is treated as disabled, no error expected
	s, err := buildSampler(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s != nil {
		t.Fatal("expected nil sampler when rate>1.0")
	}
}

// Ensure the flag is registered on the root command (integration smoke test).
func TestSamplingFlagRegistered(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	f := rootCmd.PersistentFlags().Lookup("sample-rate")
	if f == nil {
		t.Fatal("--sample-rate flag not registered on rootCmd")
	}
}
