package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestVersionOutput(t *testing.T) {
	buf := new(bytes.Buffer)
	versionCmd.SetOut(buf)
	versionCmd.SetErr(buf)

	Version = "1.2.3"
	Commit = "abc123"
	BuildDate = "2024-01-01"

	versionCmd.ExecuteContext(context.Background())

	out := buf.String()
	if !strings.Contains(out, "1.2.3") {
		t.Errorf("expected version '1.2.3' in output, got: %s", out)
	}
	if !strings.Contains(out, "abc123") {
		t.Errorf("expected commit 'abc123' in output, got: %s", out)
	}
	if !strings.Contains(out, "2024-01-01") {
		t.Errorf("expected build date '2024-01-01' in output, got: %s", out)
	}
}

func TestVersionDefault(t *testing.T) {
	Version = "dev"
	Commit = "none"
	BuildDate = "unknown"

	buf := new(bytes.Buffer)
	versionCmd.SetOut(buf)
	versionCmd.ExecuteContext(context.Background())

	out := buf.String()
	if !strings.Contains(out, "dev") {
		t.Errorf("expected default version 'dev' in output, got: %s", out)
	}
}
