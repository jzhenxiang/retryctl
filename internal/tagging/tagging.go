// Package tagging provides attempt-level tag enrichment for retry events.
// Tags are key-value string pairs attached to each attempt record, allowing
// downstream consumers (loggers, auditors, metrics) to correlate attempts
// with environment, region, or any other contextual dimension.
package tagging

import (
	"errors"
	"fmt"
	"strings"
)

// Tagger holds a resolved set of tags and can merge additional ones at
// call time.
type Tagger struct {
	base map[string]string
}

// New creates a Tagger from a slice of "key=value" pairs.
// An error is returned if any pair is malformed or has an empty key.
func New(pairs []string) (*Tagger, error) {
	if len(pairs) == 0 {
		return &Tagger{base: map[string]string{}}, nil
	}
	tags := make(map[string]string, len(pairs))
	for _, p := range pairs {
		k, v, ok := strings.Cut(p, "=")
		if !ok {
			return nil, fmt.Errorf("tagging: malformed pair %q (want key=value)", p)
		}
		if k == "" {
			return nil, errors.New("tagging: key must not be empty")
		}
		tags[k] = v
	}
	return &Tagger{base: tags}, nil
}

// Tags returns a copy of the base tags merged with any extra pairs supplied
// at call time. Extra pairs follow the same "key=value" format and override
// base values on collision.
func (t *Tagger) Tags(extra ...string) (map[string]string, error) {
	out := make(map[string]string, len(t.base)+len(extra))
	for k, v := range t.base {
		out[k] = v
	}
	for _, p := range extra {
		k, v, ok := strings.Cut(p, "=")
		if !ok {
			return nil, fmt.Errorf("tagging: malformed extra pair %q", p)
		}
		if k == "" {
			return nil, errors.New("tagging: extra key must not be empty")
		}
		out[k] = v
	}
	return out, nil
}

// Has reports whether the given key exists in the base tags.
func (t *Tagger) Has(key string) bool {
	_, ok := t.base[key]
	return ok
}

// Len returns the number of base tags.
func (t *Tagger) Len() int { return len(t.base) }
