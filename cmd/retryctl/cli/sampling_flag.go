package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"retryctl/internal/sampling"
)

func init() {
	rootCmd.PersistentFlags().Float64("sample-rate", 1.0,
		"fraction of retry attempts to allow (0 < rate ≤ 1); 1.0 disables sampling")
}

// buildSampler reads the --sample-rate flag from cmd and returns a *sampling.Sampler
// when the rate is less than 1.0, or nil when sampling is disabled.
func buildSampler(cmd *cobra.Command) (*sampling.Sampler, error) {
	rate, err := cmd.Flags().GetFloat64("sample-rate")
	if err != nil {
		return nil, fmt.Errorf("sampling: %w", err)
	}
	if rate >= 1.0 {
		return nil, nil // sampling disabled
	}
	s, err := sampling.New(rate, nil)
	if err != nil {
		return nil, fmt.Errorf("sampling: %w", err)
	}
	return s, nil
}
