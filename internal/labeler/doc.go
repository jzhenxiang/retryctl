// Package labeler provides a lightweight key-value label system for annotating
// retry attempts.
//
// Labels are defined once at startup via "key=value" CLI flags and are merged
// with per-attempt metadata (attempt number, exit code, etc.) before being
// forwarded to loggers, audit recorders, and notifiers.
//
// Example usage:
//
//	l, err := labeler.New([]string{"env=prod", "service=payments"})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	annotated := l.With(labeler.Labels{"attempt": "1"})
package labeler
