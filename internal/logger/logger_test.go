package logger_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/retryctl/internal/logger"
)

func parseEntry(t *testing.T, buf *bytes.Buffer) map[string]any {
	t.Helper()
	var entry map[string]any
	line := strings.TrimSpace(buf.String())
	if err := json.Unmarshal([]byte(line), &entry); err != nil {
		t.Fatalf("failed to parse log entry %q: %v", line, err)
	}
	return entry
}

func TestInfoLevel(t *testing.T) {
	var buf bytes.Buffer
	l := logger.New(&buf)
	l.Info("hello world", nil)
	entry := parseEntry(t, &buf)
	if entry["level"] != "info" {
		t.Errorf("expected level info, got %v", entry["level"])
	}
	if entry["message"] != "hello world" {
		t.Errorf("expected message 'hello world', got %v", entry["message"])
	}
}

func TestWarnLevel(t *testing.T) {
	var buf bytes.Buffer
	l := logger.New(&buf)
	l.Warn("watch out", nil)
	entry := parseEntry(t, &buf)
	if entry["level"] != "warn" {
		t.Errorf("expected level warn, got %v", entry["level"])
	}
}

func TestErrorLevel(t *testing.T) {
	var buf bytes.Buffer
	l := logger.New(&buf)
	l.Error("something broke", map[string]any{"code": 1})
	entry := parseEntry(t, &buf)
	if entry["level"] != "error" {
		t.Errorf("expected level error, got %v", entry["level"])
	}
	fields, ok := entry["fields"].(map[string]any)
	if !ok {
		t.Fatal("expected fields map")
	}
	if fields["code"] == nil {
		t.Error("expected 'code' field to be present")
	}
}

func TestTimestampPresent(t *testing.T) {
	var buf bytes.Buffer
	l := logger.New(&buf)
	l.Info("ts check", nil)
	entry := parseEntry(t, &buf)
	if entry["timestamp"] == "" {
		t.Error("expected non-empty timestamp")
	}
}
