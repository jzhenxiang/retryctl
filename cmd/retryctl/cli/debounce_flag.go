package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/user/retryctl/internal/debounce"
)

func init() {
	rootCmd.PersistentFlags().Duration(
		"debounce",
		0,
		"minimum quiet period between retry attempts (e.g. 500ms); 0 disables debouncing",
	)
}

// buildDebouncer reads the --debounce flag from cmd and returns a configured
// *debounce.Debouncer, or nil when debouncing is disabled (window == 0).
func buildDebouncer(cmd *cobra.Command) (*debounce.Debouncer, error) {
	win, err := cmd.Flags().GetDuration("debounce")
	if err != nil {
		return nil, fmt.Errorf("debounce flag: %w", err)
	}
	if win <= 0 {
		return nil, nil
	}
	d, err := debounce.New(win)
	if err != nil {
		return nil, fmt.Errorf("debounce: %w", err)
	}
	return d, nil
}

// debounceGuard calls d.Allow() when d is non-nil and translates ErrDebounced
// into a formatted error that callers can log or skip on.
func debounceGuard(d *debounce.Debouncer, attempt int) error {
	if d == nil {
		return nil
	}
	if err := d.Allow(); err != nil {
		remaining := d.Remaining()
		return fmt.Errorf("attempt %d debounced (retry in %s): %w",
			attempt, remaining.Round(time.Millisecond), err)
	}
	return nil
}
