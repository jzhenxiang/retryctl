// Package labeler attaches arbitrary key-value labels to retry attempts so
// that downstream components (loggers, audit recorders, notifiers) can filter
// or annotate events without coupling to the runner internals.
package labeler

import (
	"errors"
	"fmt"
	"strings"
)

// Labels is an immutable snapshot of key-value pairs.
type Labels map[string]string

// Labeler builds a Labels map from a slice of "key=value" strings.
type Labeler struct {
	base Labels
}

// New returns a Labeler pre-populated with the supplied raw pairs.
// Each element must be in "key=value" form; an error is returned on the first
// malformed entry.
func New(pairs []string) (*Labeler, error) {
	base := make(Labels, len(pairs))
	for _, p := range pairs {
		k, v, ok := strings.Cut(p, "=")
		if !ok || k == "" {
			return nil, fmt.Errorf("labeler: malformed pair %q: want key=value", p)
		}
		base[k] = v
	}
	return &Labeler{base: base}, nil
}

// With returns a new Labels map that merges the receiver's base labels with
// the supplied extra pairs.  Extra pairs take precedence on key collision.
func (l *Labeler) With(extra Labels) Labels {
	out := make(Labels, len(l.base)+len(extra))
	for k, v := range l.base {
		out[k] = v
	}
	for k, v := range extra {
		out[k] = v
	}
	return out
}

// Base returns a copy of the base label set.
func (l *Labeler) Base() Labels {
	out := make(Labels, len(l.base))
	for k, v := range l.base {
		out[k] = v
	}
	return out
}

// ErrEmpty is returned when a nil or empty Labeler is used where one is
// required.
var ErrEmpty = errors.New("labeler: no labels defined")
