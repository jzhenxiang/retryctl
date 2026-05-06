package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourorg/retryctl/internal/timeout"
)

func init() {
	rootCmd.Flags().Duration("timeout-attempt", 0,
		"maximum duration for each individual attempt (0 = no limit)")
	rootCmd.Flags().Duration("timeout-overall", 0,
		"maximum total duration across all attempts (0 = no limit)")
}

// buildTimeout reads timeout flags from cmd and returns a validated
// timeout.Enforcer. An error is returned if the configuration is invalid.
func buildTimeout(cmd *cobra.Command) (*timeout.Enforcer, error) {
	perAttempt, err := cmd.Flags().GetDuration("timeout-attempt")
	if err != nil {
		return nil, fmt.Errorf("timeout-attempt flag: %w", err)
	}

	overall, err := cmd.Flags().GetDuration("timeout-overall")
	if err != nil {
		return nil, fmt.Errorf("timeout-overall flag: %w", err)
	}

	cfg := timeout.Config{
		PerAttempt: perAttempt,
		Overall:    overall,
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return timeout.New(cfg), nil
}

// mustParseDuration is a helper used in tests to avoid repetition.
func mustParseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(fmt.Sprintf("mustParseDuration: %v", err))
	}
	return d
}
