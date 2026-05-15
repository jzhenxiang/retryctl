// Package tagging enriches retry attempt records with arbitrary key-value
// string tags supplied by the operator at invocation time.
//
// Tags flow from CLI flags through the Tagger into every attempt event,
// making it straightforward to slice metrics and audit logs by environment,
// region, service, or any other dimension without modifying retry logic.
//
// Usage:
//
//	tgr, err := tagging.New([]string{"env=prod", "region=eu-west-1"})
//	if err != nil { ... }
//
//	// At attempt time, merge per-attempt context:
//	tags, err := tgr.Tags(fmt.Sprintf("attempt=%d", n))
package tagging
