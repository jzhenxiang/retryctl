package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourorg/retryctl/internal/labeler"
)

func init() {
	rootCmd.PersistentFlags().StringSlice(
		"label",
		nil,
		`attach key=value labels to every attempt (repeatable, e.g. --label env=prod)`,
	)
}

// buildLabeler constructs a *labeler.Labeler from the --label flags on cmd.
// If no labels are provided an empty Labeler is returned without error.
func buildLabeler(cmd *cobra.Command) (*labeler.Labeler, error) {
	pairs, err := cmd.Flags().GetStringSlice("label")
	if err != nil {
		return nil, fmt.Errorf("label flag: %w", err)
	}
	l, err := labeler.New(pairs)
	if err != nil {
		return nil, fmt.Errorf("label flag: %w", err)
	}
	return l, nil
}
