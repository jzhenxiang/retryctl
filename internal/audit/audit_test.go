package audit_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/retryctl/internal/audit"
)

func TestRecordWritesJSON(t *testing.T) {
	var buf bytes.Buffer
	rec := audit.New(&buf)

	e := audit.NewEntry(1, 0, 50*time.Millisecond, true, nil)
	if err := rec.Record(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	line := strings.TrimSpace(buf.String())
	var got audit.Entry
	if err := json.Unmarshal([]byte(line), &got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if got.Attempt != 1 {
		t.Errorf("attempt: want 1, got %d", got.Attempt)
	}
	if !got.Success {
		t.Error("expected success=true")
	}
	if got.Error != "" {
		t.Errorf("expected empty error, got %q", got.Error)
	}
}

func TestRecordWithError(t *testing.T) {
	var buf bytes.Buffer
	rec := audit.New(&buf)

	e := audit.NewEntry(2, 1, 10*time.Millisecond, false, errors.New("boom"))
	if err := rec.Record(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got audit.Entry
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if got.Error != "boom" {
		t.Errorf("want error 'boom', got %q", got.Error)
	}
	if got.ExitCode != 1 {
		t.Errorf("want exit_code 1, got %d", got.ExitCode)
	}
}

func TestRecordMultipleEntries(t *testing.T) {
	var buf bytes.Buffer
	rec := audit.New(&buf)

	for i := 1; i <= 3; i++ {
		e := audit.NewEntry(i, 0, time.Duration(i)*time.Millisecond, i == 3, nil)
		if err := rec.Record(e); err != nil {
			t.Fatalf("attempt %d: %v", i, err)
		}
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("want 3 lines, got %d", len(lines))
	}
}

func TestTimestampIsUTC(t *testing.T) {
	e := audit.NewEntry(1, 0, 0, true, nil)
	if e.Timestamp.Location() != time.UTC {
		t.Errorf("expected UTC timestamp, got %v", e.Timestamp.Location())
	}
}

func TestRecordWriteError(t *testing.T) {
	rec := audit.New(&failWriter{})
	e := audit.NewEntry(1, 0, 0, true, nil)
	if err := rec.Record(e); err == nil {
		t.Error("expected write error, got nil")
	}
}

type failWriter struct{}

func (f *failWriter) Write(_ []byte) (int, error) {
	return 0, errors.New("disk full")
}
