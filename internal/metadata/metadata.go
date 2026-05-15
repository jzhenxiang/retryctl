// Package metadata attaches arbitrary key-value pairs to a retry run and
// surfaces them in structured log output and audit records.
package metadata

import (
	"errors"
	"fmt"
	"strings"
)

// Metadata holds a set of key-value annotations for a retry run.
type Metadata struct {
	pairs map[string]string
}

// New creates a Metadata instance from a slice of "key=value" strings.
// Every entry must contain exactly one '=' separator and a non-empty key.
func New(pairs []string) (*Metadata, error) {
	if len(pairs) == 0 {
		return &Metadata{pairs: map[string]string{}}, nil
	}
	m := &Metadata{pairs: make(map[string]string, len(pairs))}
	for _, p := range pairs {
		parts := strings.SplitN(p, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("metadata: invalid pair %q: must be key=value", p)
		}
		key := strings.TrimSpace(parts[0])
		if key == "" {
			return nil, errors.New("metadata: key must not be empty")
		}
		m.pairs[key] = parts[1]
	}
	return m, nil
}

// Get returns the value for key and whether it was present.
func (m *Metadata) Get(key string) (string, bool) {
	v, ok := m.pairs[key]
	return v, ok
}

// All returns a copy of all key-value pairs.
func (m *Metadata) All() map[string]string {
	out := make(map[string]string, len(m.pairs))
	for k, v := range m.pairs {
		out[k] = v
	}
	return out
}

// Merge returns a new Metadata whose pairs are the union of m and other.
// Values in other take precedence on key conflicts.
func (m *Metadata) Merge(other *Metadata) *Metadata {
	out := m.All()
	if other != nil {
		for k, v := range other.pairs {
			out[k] = v
		}
	}
	return &Metadata{pairs: out}
}
