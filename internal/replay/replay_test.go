package replay_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/retryctl/internal/replay"
)

func TestNewNilWriterReturnsError(t *testing.T) {
	_, err := replay.New(nil)
	if err == nil {
		t.Fatal("expected error for nil writer")
	}
}

func TestRecordAndLoad(t *testing.T) {
	var buf bytes.Buffer
	rec, err := replay.New(&buf)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	entries := []replay.Entry{
		{Attempt: 1, ExitCode: 1, Elapsed: 10 * time.Millisecond, Err: "exit status 1"},
		{Attempt: 2, ExitCode: 0, Elapsed: 5 * time.Millisecond},
	}
	for _, e := range entries {
		if err := rec.Record(e); err != nil {
			t.Fatalf("Record: %v", err)
		}
	}

	loaded, err := replay.Load(&buf)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(loaded) != len(entries) {
		t.Fatalf("expected %d entries, got %d", len(entries), len(loaded))
	}
	for i, e := range entries {
		if loaded[i].Attempt != e.Attempt || loaded[i].ExitCode != e.ExitCode {
			t.Errorf("entry %d mismatch: got %+v want %+v", i, loaded[i], e)
		}
	}
}

func TestLoadNilReaderReturnsError(t *testing.T) {
	_, err := replay.Load(nil)
	if err == nil {
		t.Fatal("expected error for nil reader")
	}
}

func TestLoadEmptyReturnsNil(t *testing.T) {
	entries, err := replay.Load(strings.NewReader(""))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestLoadInvalidJSONReturnsError(t *testing.T) {
	_, err := replay.Load(strings.NewReader("{bad json\n"))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestRecordPreservesError(t *testing.T) {
	var buf bytes.Buffer
	rec, _ := replay.New(&buf)
	e := replay.Entry{Attempt: 1, ExitCode: 2, Err: "something failed"}
	_ = rec.Record(e)

	loaded, err := replay.Load(&buf)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded[0].Err != e.Err {
		t.Errorf("expected Err %q, got %q", e.Err, loaded[0].Err)
	}
}
