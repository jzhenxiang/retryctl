package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"retryctl/internal/circuit"
)

func init() {
	rootCmd.PersistentFlags().Int("circuit-max-failures", 0,
		"open circuit breaker after this many consecutive failures (0 = disabled)")
	rootCmd.PersistentFlags().Duration("circuit-reset-timeout", 30*time.Second,
		"how long to wait before moving from open to half-open")
}

// buildCircuitBreaker reads circuit-breaker flags from cmd and returns a
// configured *circuit.Breaker, or nil when the feature is disabled.
func buildCircuitBreaker(cmd *cobra.Command) (*circuit.Breaker, error) {
	maxFailures, err := cmd.Flags().GetInt("circuit-max-failures")
	if err != nil {
		return nil, fmt.Errorf("circuit-max-failures: %w", err)
	}
	if maxFailures <= 0 {
		return nil, nil // feature disabled
	}

	resetTimeout, err := cmd.Flags().GetDuration("circuit-reset-timeout")
	if err != nil {
		return nil, fmt.Errorf("circuit-reset-timeout: %w", err)
	}

	br, err := circuit.New(maxFailures, resetTimeout)
	if err != nil {
		return nil, fmt.Errorf("circuit breaker: %w", err)
	}
	return br, nil
}
