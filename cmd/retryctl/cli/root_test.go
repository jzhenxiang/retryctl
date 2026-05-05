package cli

import (
	"bytes"
	"context"
	"testing"
)

func TestExecuteHelp(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"--help"})

	err := rootCmd.ExecuteContext(context.Background())
	if err != nil {
		t.Fatalf("expected no error for --help, got: %v", err)
	}
}

func TestExecuteMissingCommand(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{})

	err := rootCmd.ExecuteContext(context.Background())
	if err == nil {
		t.Fatal("expected error when no command provided")
	}
}

func TestExecuteInvalidBackoff(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"--backoff", "invalid", "--", "echo", "hi"})

	err := rootCmd.ExecuteContext(context.Background())
	if err == nil {
		t.Fatal("expected error for invalid backoff strategy")
	}
}

func TestExecuteSuccess(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"-n", "1", "--", "true"})

	err := rootCmd.ExecuteContext(context.Background())
	if err != nil {
		t.Fatalf("expected no error for 'true' command, got: %v", err)
	}
}
