package eventlog_test

import (
	"testing"
	"time"

	"github.com/user/retryctl/internal/eventlog"
)

func TestNewLogIsEmpty(t *testing.T) {
	l := eventlog.New()
	if l.Len() != 0 {
		t.Fatalf("expected 0 events, got %d", l.Len())
	}
}

func TestRecordIncreasesLen(t *testing.T) {
	l := eventlog.New()
	l.Record(eventlog.Event{Kind: eventlog.KindAttemptStarted, Attempt: 1})
	l.Record(eventlog.Event{Kind: eventlog.KindAttemptFailed, Attempt: 1, ExitCode: 1})
	if l.Len() != 2 {
		t.Fatalf("expected 2 events, got %d", l.Len())
	}
}

func TestRecordSetsTimestampIfZero(t *testing.T) {
	l := eventlog.New()
	before := time.Now().UTC()
	l.Record(eventlog.Event{Kind: eventlog.KindAttemptStarted, Attempt: 1})
	after := time.Now().UTC()

	events := l.All()
	ts := events[0].Timestamp
	if ts.Before(before) || ts.After(after) {
		t.Fatalf("timestamp %v not in expected range [%v, %v]", ts, before, after)
	}
}

func TestRecordPreservesExplicitTimestamp(t *testing.T) {
	l := eventlog.New()
	fixed := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	l.Record(eventlog.Event{Kind: eventlog.KindGaveUp, Timestamp: fixed})

	if got := l.All()[0].Timestamp; !got.Equal(fixed) {
		t.Fatalf("expected %v, got %v", fixed, got)
	}
}

func TestAllReturnsCopy(t *testing.T) {
	l := eventlog.New()
	l.Record(eventlog.Event{Kind: eventlog.KindAttemptSucceeded, Attempt: 1})

	snap := l.All()
	snap[0].Attempt = 99

	if l.All()[0].Attempt == 99 {
		t.Fatal("All() should return an independent copy")
	}
}

func TestFilterReturnsMatchingKinds(t *testing.T) {
	l := eventlog.New()
	l.Record(eventlog.Event{Kind: eventlog.KindAttemptStarted, Attempt: 1})
	l.Record(eventlog.Event{Kind: eventlog.KindAttemptFailed, Attempt: 1})
	l.Record(eventlog.Event{Kind: eventlog.KindAttemptStarted, Attempt: 2})
	l.Record(eventlog.Event{Kind: eventlog.KindAttemptSucceeded, Attempt: 2})

	started := l.Filter(eventlog.KindAttemptStarted)
	if len(started) != 2 {
		t.Fatalf("expected 2 started events, got %d", len(started))
	}
	for _, e := range started {
		if e.Kind != eventlog.KindAttemptStarted {
			t.Fatalf("unexpected kind %q", e.Kind)
		}
	}
}

func TestFilterReturnsNilWhenNoMatch(t *testing.T) {
	l := eventlog.New()
	l.Record(eventlog.Event{Kind: eventlog.KindAttemptStarted, Attempt: 1})

	result := l.Filter(eventlog.KindGaveUp)
	if len(result) != 0 {
		t.Fatalf("expected empty slice, got %d events", len(result))
	}
}

func TestResetClearsAllEvents(t *testing.T) {
	l := eventlog.New()
	l.Record(eventlog.Event{Kind: eventlog.KindAttemptStarted, Attempt: 1})
	l.Record(eventlog.Event{Kind: eventlog.KindAttemptFailed, Attempt: 1})
	l.Reset()

	if l.Len() != 0 {
		t.Fatalf("expected 0 after reset, got %d", l.Len())
	}
}
