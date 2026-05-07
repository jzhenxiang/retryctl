package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"retryctl/internal/cooldown"
)

func init() {
	// --cooldown accepts repeated "code=duration" pairs, e.g.
	// --cooldown 1=5s --cooldown 2=10s
	rootCmd.Flags().StringArray("cooldown", nil,
		"per-exit-code cooldown window as CODE=DURATION (repeatable)")
}

// buildCooldown parses --cooldown flags and returns a *cooldown.Cooldown.
// Returns nil when no flags are provided (feature disabled).
func buildCooldown(cmd *cobra.Command) (*cooldown.Cooldown, error) {
	pairs, err := cmd.Flags().GetStringArray("cooldown")
	if err != nil || len(pairs) == 0 {
		return nil, nil //nolint:nilerr // flag not set → feature disabled
	}

	windows := make(map[int]time.Duration, len(pairs))
	for _, p := range pairs {
		parts := strings.SplitN(p, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("cooldown: invalid format %q, expected CODE=DURATION", p)
		}
		var code int
		if _, err := fmt.Sscanf(parts[0], "%d", &code); err != nil {
			return nil, fmt.Errorf("cooldown: invalid exit code %q: %w", parts[0], err)
		}
		d, err := time.ParseDuration(parts[1])
		if err != nil {
			return nil, fmt.Errorf("cooldown: invalid duration %q: %w", parts[1], err)
		}
		windows[code] = d
	}

	return cooldown.New(windows)
}
