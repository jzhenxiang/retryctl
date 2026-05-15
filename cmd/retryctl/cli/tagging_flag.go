package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourorg/retryctl/internal/tagging"
)

func init() {
	rootCmd.PersistentFlags().StringArray(
		"tag",
		nil,
		`attach a key=value tag to every attempt record (repeatable, e.g. --tag env=prod --tag region=us-east-1)`,
	)
}

// buildTagger constructs a *tagging.Tagger from the --tag flags on cmd.
// If no flags are provided a zero-tag Tagger is returned so callers never
// need to handle a nil value.
func buildTagger(cmd *cobra.Command) (*tagging.Tagger, error) {
	pairs, err := cmd.Flags().GetStringArray("tag")
	if err != nil {
		return nil, fmt.Errorf("tagging: reading --tag flags: %w", err)
	}
	tgr, err := tagging.New(pairs)
	if err != nil {
		return nil, fmt.Errorf("tagging: %w", err)
	}
	return tgr, nil
}
