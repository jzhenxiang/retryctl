package cli

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func newCmdWithOutput(buf *bytes.Buffer) *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	return cmd
}

func TestBuildNotifierEmpty(t *testing.T) {
	notifyFlag = ""
	if n := buildNotifier(newCmdWithOutput(&bytes.Buffer{})); n != nil {
		t.Errorf("expected nil notifier when flag is empty, got %T", n)
	}
}

func TestBuildNotifierLog(t *testing.T) {
	notifyFlag = "log"
	var buf bytes.Buffer
	n := buildNotifier(newCmdWithOutput(&buf))
	if n == nil {
		t.Fatal("expected non-nil notifier for 'log'")
	}
}

func TestBuildNotifierStderr(t *testing.T) {
	notifyFlag = "stderr"
	n := buildNotifier(newCmdWithOutput(&bytes.Buffer{}))
	if n == nil {
		t.Fatal("expected non-nil notifier for 'stderr'")
	}
}

func TestBuildNotifierUnknownFallsBack(t *testing.T) {
	notifyFlag = "slack" // not implemented
	var buf bytes.Buffer
	n := buildNotifier(newCmdWithOutput(&buf))
	if n == nil {
		t.Fatal("expected fallback notifier for unknown channel")
	}
	if !bytes.Contains(buf.Bytes(), []byte("warn:")) {
		t.Errorf("expected warning message, got: %s", buf.String())
	}
	// reset
	notifyFlag = ""
}
