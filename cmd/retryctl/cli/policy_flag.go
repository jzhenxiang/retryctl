package cli

import (
	"github.com/retryctl/internal/backoff"
	"github.com/retryctl/internal/jitter"
	"github.com/retryctl/internal/predicates"
	"github.com/retryctl/internal/retry"
	"github.com/spf13/cobra"
)

func init() {
	// max-attempts is already registered by root.go; this file wires the
	// remaining policy-level knobs into a retry.Policy at run time.
}

// buildPolicy assembles a retry.Policy from the flags already bound to cmd.
// It relies on backoff strategy, jitter, and predicate helpers defined in
// their respective flag files.
func buildPolicy(cmd *cobra.Command) (retry.Policy, error) {
	maxAttempts, err := cmd.Flags().GetInt("max-attempts")
	if err != nil {
		return retry.Policy{}, err
	}

	strategy, err := buildBackoffStrategy(cmd)
	if err != nil {
		return retry.Policy{}, err
	}

	j := buildJitter(cmd)

	pred, err := buildPredicate(cmd)
	if err != nil {
		return retry.Policy{}, err
	}

	p := retry.Policy{
		MaxAttempts: maxAttempts,
		Strategy:    strategy,
		Jitter:      j,
		ShouldRetry: func(attempt int, code int, output []byte, execErr error) bool {
			return pred(attempt, code, output, execErr)
		},
	}

	if err := p.Validate(); err != nil {
		return retry.Policy{}, err
	}
	return p, nil
}

// buildBackoffStrategy reads the --backoff and --delay / --max-delay flags
// and returns a configured backoff.Strategy.
func buildBackoffStrategy(cmd *cobra.Command) (backoff.Strategy, error) {
	kind, _ := cmd.Flags().GetString("backoff")
	delay, _ := cmd.Flags().GetDuration("delay")
	maxDelay, _ := cmd.Flags().GetDuration("max-delay")

	var k backoff.Kind
	switch kind {
	case "linear":
		k = backoff.Linear
	case "exponential":
		k = backoff.Exponential
	default:
		k = backoff.Fixed
	}
	return backoff.NewStrategy(k, delay, maxDelay), nil
}

// alwaysRetry is a convenience predicate used when no retry condition flags
// are supplied.
func alwaysRetry(_ int, _ int, _ []byte, _ error) bool {
	return predicates.Always(0, 0, nil, nil)
}
