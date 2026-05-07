// Package predicates provides composable retry predicates that determine
// whether a failed attempt should be retried based on exit code and output.
package predicates

import (
	"errors"
	"fmt"
)

// Predicate reports whether the given attempt result warrants a retry.
type Predicate func(exitCode int, output string) bool

// Always returns a Predicate that always recommends retrying.
func Always() Predicate {
	return func(_ int, _ string) bool { return true }
}

// Never returns a Predicate that never recommends retrying.
func Never() Predicate {
	return func(_ int, _ string) bool { return false }
}

// OnExitCodes returns a Predicate that retries only when the exit code is one
// of the provided codes. An empty codes list causes an error.
func OnExitCodes(codes ...int) (Predicate, error) {
	if len(codes) == 0 {
		return nil, errors.New("predicates: at least one exit code required")
	}
	set := make(map[int]struct{}, len(codes))
	for _, c := range codes {
		set[c] = struct{}{}
	}
	return func(exitCode int, _ string) bool {
		_, ok := set[exitCode]
		return ok
	}, nil
}

// OnOutputContains returns a Predicate that retries when the command output
// contains the given substring.
func OnOutputContains(sub string) (Predicate, error) {
	if sub == "" {
		return nil, errors.New("predicates: substring must not be empty")
	}
	return func(_ int, output string) bool {
		return contains(output, sub)
	}, nil
}

// Any returns a Predicate that retries if at least one of the provided
// predicates recommends retrying.
func Any(ps ...Predicate) (Predicate, error) {
	if len(ps) == 0 {
		return nil, errors.New("predicates: Any requires at least one predicate")
	}
	return func(exitCode int, output string) bool {
		for _, p := range ps {
			if p(exitCode, output) {
				return true
			}
		}
		return false
	}, nil
}

// All returns a Predicate that retries only when every provided predicate
// recommends retrying.
func All(ps ...Predicate) (Predicate, error) {
	if len(ps) == 0 {
		return nil, fmt.Errorf("predicates: All requires at least one predicate")
	}
	return func(exitCode int, output string) bool {
		for _, p := range ps {
			if !p(exitCode, output) {
				return false
			}
		}
		return true
	}, nil
}

// contains is a simple substring check to avoid importing strings at call sites.
func contains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
