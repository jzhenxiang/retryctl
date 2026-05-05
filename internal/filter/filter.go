package filter

import (
	"regexp"
)

// ExitCodeMatcher decides whether a given exit code should trigger a retry.
type ExitCodeMatcher struct {
	allowedCodes map[int]struct{}
	retryAll     bool
}

// NewExitCodeMatcher creates a matcher. If codes is empty, all non-zero exit
// codes are considered retryable.
func NewExitCodeMatcher(codes []int) *ExitCodeMatcher {
	if len(codes) == 0 {
		return &ExitCodeMatcher{retryAll: true}
	}
	m := make(map[int]struct{}, len(codes))
	for _, c := range codes {
		m[c] = struct{}{}
	}
	return &ExitCodeMatcher{allowedCodes: m}
}

// ShouldRetry returns true when the exit code indicates the command should be
// retried.
func (e *ExitCodeMatcher) ShouldRetry(code int) bool {
	if code == 0 {
		return false
	}
	if e.retryAll {
		return true
	}
	_, ok := e.allowedCodes[code]
	return ok
}

// OutputFilter decides whether command output matches a pattern that should
// suppress further retries (e.g. a fatal error message).
type OutputFilter struct {
	patterns []*regexp.Regexp
}

// NewOutputFilter compiles the given regular-expression patterns. An error is
// returned if any pattern fails to compile.
func NewOutputFilter(patterns []string) (*OutputFilter, error) {
	regs := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		r, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		regs = append(regs, r)
	}
	return &OutputFilter{patterns: regs}, nil
}

// MatchesAbort returns true when the output contains a pattern that should
// abort retries immediately.
func (f *OutputFilter) MatchesAbort(output string) bool {
	for _, r := range f.patterns {
		if r.MatchString(output) {
			return true
		}
	}
	return false
}
