// Package notify provides pluggable notification channels that fire
// when a retry sequence finishes (success or terminal failure).
package notify

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Event holds the outcome of a full retry run.
type Event struct {
	Command   string
	Success   bool
	Attempts  int
	Elapsed   time.Duration
	LastError error
}

// Notifier sends an Event to some destination.
type Notifier interface {
	Notify(e Event) error
}

// LogNotifier writes a human-readable line to an io.Writer.
type LogNotifier struct {
	w io.Writer
}

// NewLogNotifier returns a LogNotifier that writes to w.
// If w is nil, os.Stderr is used.
func NewLogNotifier(w io.Writer) *LogNotifier {
	if w == nil {
		w = os.Stderr
	}
	return &LogNotifier{w: w}
}

// Notify implements Notifier.
func (l *LogNotifier) Notify(e Event) error {
	status := "SUCCESS"
	if !e.Success {
		status = "FAILURE"
	}
	errMsg := ""
	if e.LastError != nil {
		errMsg = fmt.Sprintf(" error=%q", e.LastError.Error())
	}
	_, err := fmt.Fprintf(
		l.w,
		"[notify] status=%s command=%q attempts=%d elapsed=%s%s\n",
		status, e.Command, e.Attempts, e.Elapsed.Round(time.Millisecond), errMsg,
	)
	return err
}

// Multi fans an Event out to several Notifiers, collecting all errors.
type Multi struct {
	notifiers []Notifier
}

// NewMulti returns a Multi that delegates to each of the supplied Notifiers.
func NewMulti(nn ...Notifier) *Multi {
	return &Multi{notifiers: nn}
}

// Notify implements Notifier.
func (m *Multi) Notify(e Event) error {
	var combined error
	for _, n := range m.notifiers {
		if err := n.Notify(e); err != nil && combined == nil {
			combined = err
		}
	}
	return combined
}
