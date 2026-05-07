package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/user/retryctl/internal/budget"
)

func init() {
	rootCmd.Flags().Int("budget-max", 0, "max retries allowed within the budget window (0 = disabled)")
	rootCmd.Flags().Duration("budget-window", time.Minute, "sliding window duration for the retry budget")
}

// buildBudget constructs a *budget.Budget from CLI flags, or returns nil when
// the feature is disabled (budget-max == 0).
func buildBudget(cmd *cobra.Command) (*budget.Budget, error) {
	max, err := cmd.Flags().GetInt("budget-max")
	if err != nil {
		return nil, fmt.Errorf("budget-max: %w", err)
	}
	if max <= 0 {
		return nil, nil
	}

	win, err := cmd.Flags().GetDuration("budget-window")
	if err != nil {
		return nil, fmt.Errorf("budget-window: %w", err)
	}

	b, err := budget.New(max, win)
	if err != nil {
		return nil, fmt.Errorf("retry budget: %w", err)
	}
	return b, nil
}
