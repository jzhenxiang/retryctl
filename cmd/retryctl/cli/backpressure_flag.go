package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/user/retryctl/internal/backpressure"
)

func init() {
	rootCmd.PersistentFlags().Int("bp-capacity", 0,
		"token-bucket capacity for back-pressure (0 = disabled)")
	rootCmd.PersistentFlags().Duration("bp-refill", 200*time.Millisecond,
		"how often one token is added back to the bucket")
}

// buildBackpressure returns a *backpressure.Limiter configured from the
// command flags, or nil when back-pressure is disabled (capacity == 0).
func buildBackpressure(cmd *cobra.Command) (*backpressure.Limiter, error) {
	capacity, err := cmd.Flags().GetInt("bp-capacity")
	if err != nil {
		return nil, fmt.Errorf("bp-capacity: %w", err)
	}
	if capacity <= 0 {
		return nil, nil //nolint:nilnil // disabled
	}

	refill, err := cmd.Flags().GetDuration("bp-refill")
	if err != nil {
		return nil, fmt.Errorf("bp-refill: %w", err)
	}

	l, err := backpressure.New(capacity, refill)
	if err != nil {
		return nil, fmt.Errorf("back-pressure: %w", err)
	}
	return l, nil
}
