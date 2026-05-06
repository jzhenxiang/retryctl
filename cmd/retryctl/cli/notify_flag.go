package cli

import (
	"os"

	"github.com/example/retryctl/internal/notify"
	"github.com/spf13/cobra"
)

// notifyFlag is the value of --notify collected by the root command.
var notifyFlag string

func init() {
	rootCmd.PersistentFlags().StringVar(
		&notifyFlag,
		"notify",
		"",
		`notification channel after run completes; supported: "log" (default when flag is set)`,
	)
}

// buildNotifier returns a Notifier based on the --notify flag value.
// Returns nil when the flag is empty so callers can skip notification cheaply.
func buildNotifier(cmd *cobra.Command) notify.Notifier {
	switch notifyFlag {
	case "log":
		return notify.NewLogNotifier(cmd.OutOrStdout())
	case "stderr":
		return notify.NewLogNotifier(os.Stderr)
	case "":
		return nil
	default:
		// Unknown channel – fall back to log so the user still gets feedback.
		cmd.PrintErrf("warn: unknown --notify channel %q, falling back to log\n", notifyFlag)
		return notify.NewLogNotifier(cmd.OutOrStdout())
	}
}
