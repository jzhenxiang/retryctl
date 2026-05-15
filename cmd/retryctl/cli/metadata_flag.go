package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourorg/retryctl/internal/metadata"
)

func init() {
	rootCmd.PersistentFlags().StringArray(
		"meta",
		nil,
		`attach key=value metadata to the run (repeatable, e.g. --meta env=prod)`,
	)
}

// buildMetadata reads --meta flags from cmd and returns a populated Metadata.
// If no flags are provided an empty Metadata is returned without error.
func buildMetadata(cmd *cobra.Command) (*metadata.Metadata, error) {
	pairs, err := cmd.Flags().GetStringArray("meta")
	if err != nil {
		return nil, fmt.Errorf("metadata: %w", err)
	}
	m, err := metadata.New(pairs)
	if err != nil {
		return nil, fmt.Errorf("metadata: %w", err)
	}
	return m, nil
}
