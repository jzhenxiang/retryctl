package cli

import (
	"fmt"

	"github.com/example/retryctl/internal/jitter"
	"github.com/spf13/cobra"
)

var jitterStrategy string

func init() {
	rootCmd.PersistentFlags().StringVar(
		&jitterStrategy,
		"jitter",
		"none",
		`jitter strategy applied to each backoff delay ("full", "equal", "none")`,
	)
}

// buildJitter returns a jitter.Applier for the given strategy name.
// It validates that the name is one of the accepted values.
func buildJitter(cmd *cobra.Command) (jitter.Applier, error) {
	strategy, err := cmd.Flags().GetString("jitter")
	if err != nil {
		strategy = "none"
	}
	switch strategy {
	case "full", "equal", "none":
		return jitter.New(strategy), nil
	default:
		return nil, fmt.Errorf("unknown jitter strategy %q: must be full, equal, or none", strategy)
	}
}
