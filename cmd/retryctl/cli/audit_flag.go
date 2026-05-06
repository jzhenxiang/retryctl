package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourorg/retryctl/internal/audit"
)

var auditFile string

func init() {
	rootCmd.PersistentFlags().StringVar(
		&auditFile,
		"audit-log",
		"",
		"path to write newline-delimited JSON audit entries (default: disabled)",
	)
}

// buildAuditRecorder returns an *audit.Recorder that writes to the path
// specified by --audit-log, or nil when the flag is unset.
// The caller is responsible for closing the returned io.Closer (if non-nil).
func buildAuditRecorder(cmd *cobra.Command) (*audit.Recorder, io.Closer, error) {
	path, err := cmd.Flags().GetString("audit-log")
	if err != nil || path == "" {
		return nil, nil, nil //nolint:nilerr // flag absent is not an error
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, nil, fmt.Errorf("audit: open %q: %w", path, err)
	}
	return audit.New(f), f, nil
}
