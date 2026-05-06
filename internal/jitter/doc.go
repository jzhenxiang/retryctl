// Package jitter provides Applier implementations that add randomness to
// retry backoff delays, reducing the chance of a thundering-herd scenario
// where many retrying clients hit a recovering service at the same instant.
//
// Three strategies are available:
//
//	"full"  – random value in [0, base)
//	"equal" – random value in [base/2, base)
//	"none"  – base is returned unchanged (default)
//
// Use jitter.New(strategy) to obtain an Applier by name, or construct one
// of the concrete types directly when a custom random source is needed.
package jitter
