package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"retryctl/internal/deadline"
)

func init() {
	rootCmd.PersistentFlags().String(
		"deadline", "",
		"absolute deadline as RFC3339 timestamp (e.g. 2025-12-31T23:59:59Z); "+
			"retries will not start after this time",
	)
}

// buildDeadlineGuard parses the --deadline flag from cmd and returns a
// *deadline.Guard if the flag was set, or nil if it was omitted.
// An error is returned when the value is present but invalid.
func buildDeadlineGuard(cmd *cobra.Command) (*deadline.Guard, error) {
	raw, err := cmd.Flags().GetString("deadline")
	if err != nil || raw == "" {
		return nil, nil //nolint:nilerr // flag not set is intentional
	}

	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		// also try date-only for convenience
		t, err = time.Parse("2006-01-02", raw)
		if err != nil {
			return nil, fmt.Errorf("--deadline: invalid RFC3339 timestamp %q", raw)
		}
		// treat date-only as end of that UTC day
		t = time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, time.UTC)
	}

	g, err := deadline.New(t)
	if err != nil {
		return nil, fmt.Errorf("--deadline: %w", err)
	}
	return g, nil
}
