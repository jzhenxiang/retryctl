// Package jitter provides helpers for adding randomness to backoff delays
// to avoid thundering-herd problems when many clients retry simultaneously.
package jitter

import (
	"math/rand"
	"time"
)

// Source is a function that returns a pseudo-random float64 in [0.0, 1.0).
// Swappable for deterministic tests.
type Source func() float64

// defaultSource uses the global math/rand generator.
var defaultSource Source = rand.Float64

// Applier applies jitter to a base duration.
type Applier interface {
	Apply(base time.Duration) time.Duration
}

// Full returns a duration in [0, base).
type Full struct{ src Source }

// NewFull creates a Full jitter applier.
func NewFull() *Full { return &Full{src: defaultSource} }

// NewFullWithSource creates a Full jitter applier with a custom random source.
func NewFullWithSource(src Source) *Full { return &Full{src: src} }

// Apply implements Applier.
func (f *Full) Apply(base time.Duration) time.Duration {
	return time.Duration(f.src() * float64(base))
}

// Equal returns a duration in [base/2, base).
type Equal struct{ src Source }

// NewEqual creates an Equal jitter applier.
func NewEqual() *Equal { return &Equal{src: defaultSource} }

// NewEqualWithSource creates an Equal jitter applier with a custom random source.
func NewEqualWithSource(src Source) *Equal { return &Equal{src: src} }

// Apply implements Applier.
func (e *Equal) Apply(base time.Duration) time.Duration {
	half := base / 2
	return half + time.Duration(e.src()*float64(half))
}

// None is a no-op applier that returns the base duration unchanged.
type None struct{}

// Apply implements Applier.
func (n None) Apply(base time.Duration) time.Duration { return base }

// New returns an Applier for the given strategy name.
// Valid names: "full", "equal", "none" (default).
func New(strategy string) Applier {
	switch strategy {
	case "full":
		return NewFull()
	case "equal":
		return NewEqual()
	default:
		return None{}
	}
}
