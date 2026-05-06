package notify_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/example/retryctl/internal/notify"
)

func makeEvent(success bool) notify.Event {
	var lastErr error
	if !success {
		lastErr = errors.New("exit status 1")
	}
	return notify.Event{
		Command:   "myapp --flag",
		Success:   success,
		Attempts:  3,
		Elapsed:   150 * time.Millisecond,
		LastError: lastErr,
	}
}

func TestLogNotifierSuccess(t *testing.T) {
	var buf bytes.Buffer
	n := notify.NewLogNotifier(&buf)
	if err := n.Notify(makeEvent(true)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "SUCCESS") {
		t.Errorf("expected SUCCESS in output, got: %s", out)
	}
	if !strings.Contains(out, "myapp --flag") {
		t.Errorf("expected command in output, got: %s", out)
	}
}

func TestLogNotifierFailure(t *testing.T) {
	var buf bytes.Buffer
	n := notify.NewLogNotifier(&buf)
	if err := n.Notify(makeEvent(false)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "FAILURE") {
		t.Errorf("expected FAILURE in output, got: %s", out)
	}
	if !strings.Contains(out, "exit status 1") {
		t.Errorf("expected error message in output, got: %s", out)
	}
}

func TestLogNotifierNilWriterDefaultsToStderr(t *testing.T) {
	// Should not panic.
	n := notify.NewLogNotifier(nil)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestMultiNotifier(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	m := notify.NewMulti(
		notify.NewLogNotifier(&buf1),
		notify.NewLogNotifier(&buf2),
	)
	if err := m.Notify(makeEvent(true)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf1.Len() == 0 || buf2.Len() == 0 {
		t.Error("expected both notifiers to have received the event")
	}
}

func TestMultiNotifierEmpty(t *testing.T) {
	m := notify.NewMulti()
	if err := m.Notify(makeEvent(false)); err != nil {
		t.Fatalf("unexpected error on empty multi: %v", err)
	}
}
