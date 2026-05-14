// Package replay provides a mechanism for recording and replaying
// attempt outcomes, useful for dry-run and simulation modes.
package replay

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"
)

// Entry represents a single recorded attempt outcome.
type Entry struct {
	Attempt  int           `json:"attempt"`
	ExitCode int           `json:"exit_code"`
	Elapsed  time.Duration `json:"elapsed_ns"`
	Err      string        `json:"error,omitempty"`
}

// Recorder writes attempt outcomes as newline-delimited JSON.
type Recorder struct {
	w io.Writer
}

// New returns a Recorder that writes entries to w.
func New(w io.Writer) (*Recorder, error) {
	if w == nil {
		return nil, errors.New("replay: writer must not be nil")
	}
	return &Recorder{w: w}, nil
}

// Record encodes and writes a single Entry to the underlying writer.
func (r *Recorder) Record(e Entry) error {
	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("replay: marshal entry: %w", err)
	}
	b = append(b, '\n')
	_, err = r.w.Write(b)
	return err
}

// Load reads all entries from r, returning them in order.
func Load(r io.Reader) ([]Entry, error) {
	if r == nil {
		return nil, errors.New("replay: reader must not be nil")
	}
	var entries []Entry
	dec := json.NewDecoder(r)
	for {
		var e Entry
		if err := dec.Decode(&e); errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, fmt.Errorf("replay: decode entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}
