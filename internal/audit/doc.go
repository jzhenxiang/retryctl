// Package audit records structured, newline-delimited JSON entries for every
// retry attempt made by retryctl.
//
// Each Entry captures:
//   - the wall-clock timestamp (UTC)
//   - the attempt number (1-based)
//   - the process exit code
//   - elapsed duration in nanoseconds
//   - whether the attempt was considered successful
//   - an optional error message
//
// Typical usage:
//
//	f, _ := os.Create("retryctl-audit.jsonl")
//	rec := audit.New(f)
//	rec.Record(audit.NewEntry(attempt, exitCode, elapsed, ok, err))
package audit
