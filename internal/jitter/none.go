package jitter

// NoneFunc is a Func that always returns the base delay unchanged,
// effectively applying no jitter. It is useful as a no-op default
// when jitter is disabled by the user.
//
// Example:
//
//	f := NewNone()
//	delay := f(500) // always 0 — base delay is used as-is by the caller
type NoneFunc = Func

// NewNone returns a Func that always returns 0, meaning the caller
// should use the raw base delay without any random perturbation.
//
// This is the default jitter strategy when the user passes --jitter=none
// or omits the flag entirely.
func NewNone() Func {
	return func(_ int64) int64 {
		return 0
	}
}
