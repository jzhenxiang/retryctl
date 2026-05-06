// Package circuit provides a simple circuit breaker for use with retryctl.
//
// The Breaker tracks consecutive failures. Once the failure count reaches
// the configured threshold the breaker transitions to the Open state and
// calls to Allow return ErrOpen, preventing further attempts.
//
// After the reset timeout elapses the breaker moves to HalfOpen, allowing
// a single probe attempt. A successful probe closes the circuit; another
// failure reopens it.
//
// Usage:
//
//	br, err := circuit.New(5, 30*time.Second)
//	if err != nil { ... }
//
//	if err := br.Allow(); err != nil {
//	    // skip attempt — circuit is open
//	}
//	// run command …
//	br.RecordSuccess() // or br.RecordFailure()
package circuit
