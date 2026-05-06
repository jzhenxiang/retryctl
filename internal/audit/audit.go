// Package audit provides a structured audit trail for retry attempts,
// recording each attempt's outcome, duration, and exit code to a sink.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Entry represents a single recorded retry attempt.
type Entry struct {
	Timestamp time.Time     `json:"timestamp"`
	Attempt   int           `json:"attempt"`
	ExitCode  int           `json:"exit_code"`
	Elapsed   time.Duration `json:"elapsed_ns"`
	Success   bool          `json:"success"`
	Error     string        `json:"error,omitempty"`
}

// Recorder writes audit entries to an io.Writer as newline-delimited JSON.
type Recorder struct {
	w io.Writer
}

// New returns a Recorder that writes to w.
func New(w io.Writer) *Recorder {
	return &Recorder{w: w}
}

// Record serialises e as a JSON line to the underlying writer.
// It returns any write or encoding error.
func (r *Recorder) Record(e Entry) error {
	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal: %w", err)
	}
	b = append(b, '\n')
	if _, err := r.w.Write(b); err != nil {
		return fmt.Errorf("audit: write: %w", err)
	}
	return nil
}

// NewEntry is a convenience constructor for building an Entry.
func NewEntry(attempt, exitCode int, elapsed time.Duration, success bool, err error) Entry {
	e := Entry{
		Timestamp: time.Now().UTC(),
		Attempt:   attempt,
		ExitCode:  exitCode,
		Elapsed:   elapsed,
		Success:   success,
	}
	if err != nil {
		e.Error = err.Error()
	}
	return e
}
