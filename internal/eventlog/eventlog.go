// Package eventlog provides a structured in-memory event log that records
// retry lifecycle events (attempt started, attempt succeeded, attempt failed)
// with timestamps and metadata for post-run inspection or export.
package eventlog

import (
	"sync"
	"time"
)

// EventKind identifies the type of lifecycle event.
type EventKind string

const (
	KindAttemptStarted   EventKind = "attempt_started"
	KindAttemptSucceeded EventKind = "attempt_succeeded"
	KindAttemptFailed    EventKind = "attempt_failed"
	KindGaveUp           EventKind = "gave_up"
)

// Event represents a single lifecycle event emitted during a retry run.
type Event struct {
	Kind      EventKind         `json:"kind"`
	Attempt   int               `json:"attempt"`
	Timestamp time.Time         `json:"timestamp"`
	ExitCode  int               `json:"exit_code,omitempty"`
	Message   string            `json:"message,omitempty"`
	Labels    map[string]string `json:"labels,omitempty"`
}

// Log is a concurrency-safe in-memory store of Events.
type Log struct {
	mu     sync.Mutex
	events []Event
}

// New returns an empty, ready-to-use Log.
func New() *Log {
	return &Log{}
}

// Record appends e to the log.
func (l *Log) Record(e Event) {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.events = append(l.events, e)
}

// All returns a snapshot of all recorded events in insertion order.
func (l *Log) All() []Event {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]Event, len(l.events))
	copy(out, l.events)
	return out
}

// Filter returns only the events whose Kind matches kind.
func (l *Log) Filter(kind EventKind) []Event {
	l.mu.Lock()
	defer l.mu.Unlock()
	var out []Event
	for _, e := range l.events {
		if e.Kind == kind {
			out = append(out, e)
		}
	}
	return out
}

// Len returns the total number of recorded events.
func (l *Log) Len() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.events)
}

// Reset clears all recorded events.
func (l *Log) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.events = l.events[:0]
}
