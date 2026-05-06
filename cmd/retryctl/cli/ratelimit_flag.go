package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"retryctl/internal/ratelimit"
)

var (
	rateLimitMax    int
	rateLimitWindow time.Duration
)

func init() {
	rootCmd.PersistentFlags().IntVar(
		&rateLimitMax,
		"rate-limit-max",
		0,
		"maximum number of attempts within --rate-limit-window (0 = disabled)",
	)
	rootCmd.PersistentFlags().DurationVar(
		&rateLimitWindow,
		"rate-limit-window",
		60*time.Second,
		"sliding window duration for rate limiting",
	)
}

// buildRateLimiter returns a configured *ratelimit.Limiter when rate limiting
// is enabled (--rate-limit-max > 0), or nil when it is disabled.
func buildRateLimiter(cmd *cobra.Command) (*ratelimit.Limiter, error) {
	max, err := cmd.Flags().GetInt("rate-limit-max")
	if err != nil {
		return nil, fmt.Errorf("rate-limit-max: %w", err)
	}
	if max <= 0 {
		return nil, nil //nolint:nilnil // disabled by design
	}
	win, err := cmd.Flags().GetDuration("rate-limit-window")
	if err != nil {
		return nil, fmt.Errorf("rate-limit-window: %w", err)
	}
	lim, err := ratelimit.New(ratelimit.Config{Max: max, Window: win})
	if err != nil {
		return nil, fmt.Errorf("rate limiter: %w", err)
	}
	return lim, nil
}
