package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/retryctl/internal/replay"
)

func init() {
	rootCmd.PersistentFlags().String("replay-out", "", "path to write replay log (NDJSON); disabled if empty")
}

// buildReplayRecorder returns a *replay.Recorder if --replay-out is set,
// or nil when the flag is empty. The caller is responsible for closing
// the returned *os.File (second return value) when non-nil.
func buildReplayRecorder(cmd *cobra.Command) (*replay.Recorder, *os.File, error) {
	path, err := cmd.Flags().GetString("replay-out")
	if err != nil {
		return nil, nil, fmt.Errorf("replay-out flag: %w", err)
	}
	if path == "" {
		return nil, nil, nil
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return nil, nil, fmt.Errorf("replay: open %q: %w", path, err)
	}

	rec, err := replay.New(f)
	if err != nil {
		_ = f.Close()
		return nil, nil, err
	}
	return rec, f, nil
}
