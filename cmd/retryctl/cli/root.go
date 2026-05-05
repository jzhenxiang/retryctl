package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/retryctl/retryctl/internal/config"
	"github.com/retryctl/retryctl/internal/logger"
	"github.com/retryctl/retryctl/internal/runner"
	"github.com/spf13/cobra"
)

var (
	maxAttempts int
	delay       time.Duration
	backoff     string
	logLevel    string
)

var rootCmd = &cobra.Command{
	Use:   "retryctl [flags] -- <command> [args...]",
	Short: "Wrap flaky commands with configurable retry logic",
	Args:  cobra.MinimumNArgs(1),
	RunE:  run,
}

func init() {
	def := config.Default()
	rootCmd.Flags().IntVarP(&maxAttempts, "attempts", "n", def.MaxAttempts, "maximum number of attempts")
	rootCmd.Flags().DurationVarP(&delay, "delay", "d", def.Delay, "initial delay between retries")
	rootCmd.Flags().StringVarP(&backoff, "backoff", "b", def.Backoff, "backoff strategy: fixed, linear, exponential")
	rootCmd.Flags().StringVarP(&logLevel, "log-level", "l", "info", "log level: debug, info, warn, error")
}

func run(cmd *cobra.Command, args []string) error {
	cfg := &config.Config{
		Command:     args[0],
		Args:        args[1:],
		MaxAttempts: maxAttempts,
		Delay:       delay,
		Backoff:     backoff,
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	log := logger.New(os.Stdout, logLevel)
	r := runner.New(cfg, log)

	return r.Run(cmd.Context())
}

func Execute() error {
	return rootCmd.Execute()
}
