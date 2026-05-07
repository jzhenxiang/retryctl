package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/user/retryctl/internal/throttle"
)

func init() {
	rootCmd.PersistentFlags().Int(
		"throttle",
		0,
		"maximum number of concurrent retry attempts (0 = unlimited)",
	)
}

// buildThrottle reads the --throttle flag from cmd and returns a configured
// *throttle.Throttle, or nil when throttling is disabled (value == 0).
func buildThrottle(cmd *cobra.Command) (*throttle.Throttle, error) {
	max, err := cmd.Flags().GetInt("throttle")
	if err != nil {
		return nil, fmt.Errorf("throttle flag: %w", err)
	}
	if max <= 0 {
		return nil, nil
	}
	th, err := throttle.New(max)
	if err != nil {
		return nil, fmt.Errorf("throttle flag: %w", err)
	}
	return th, nil
}
