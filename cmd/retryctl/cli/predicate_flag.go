package cli

import (
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"retryctl/internal/predicates"
)

var (
	retryOnCodes  []string
	retryOnOutput string
)

func init() {
	rootCmd.PersistentFlags().StringSliceVar(
		&retryOnCodes,
		"retry-on-codes",
		nil,
		"comma-separated list of exit codes that trigger a retry (default: all non-zero)",
	)
	rootCmd.PersistentFlags().StringVar(
		&retryOnOutput,
		"retry-on-output",
		"",
		"retry when command output contains this substring",
	)
}

// buildPredicate constructs a retry Predicate from the CLI flags.
// When no flags are provided it falls back to predicates.Always so that
// every non-zero exit code triggers a retry — preserving existing behaviour.
func buildPredicate(cmd *cobra.Command) predicates.Predicate {
	var parts []predicates.Predicate

	if len(retryOnCodes) > 0 {
		codes, err := parseIntSlice(retryOnCodes)
		if err == nil {
			if p, err := predicates.OnExitCodes(codes...); err == nil {
				parts = append(parts, p)
			}
		}
	}

	if retryOnOutput != "" {
		if p, err := predicates.OnOutputContains(retryOnOutput); err == nil {
			parts = append(parts, p)
		}
	}

	if len(parts) == 0 {
		return predicates.Always()
	}

	if len(parts) == 1 {
		return parts[0]
	}

	p, err := predicates.Any(parts...)
	if err != nil {
		return predicates.Always()
	}
	return p
}

// parseIntSlice converts a slice of string tokens to ints.
func parseIntSlice(tokens []string) ([]int, error) {
	out := make([]int, 0, len(tokens))
	for _, t := range tokens {
		t = strings.TrimSpace(t)
		v, err := strconv.Atoi(t)
		if err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, nil
}
