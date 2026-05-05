package config

import (
	"errors"
	"time"
)

// BackoffStrategy represents the type of backoff to use between retries.
type BackoffStrategy string

const (
	StrategyFixed       BackoffStrategy = "fixed"
	StrategyLinear      BackoffStrategy = "linear"
	StrategyExponential BackoffStrategy = "exponential"

	DefaultMaxAttempts  = 3
	DefaultDelay        = 1 * time.Second
	DefaultStrategy     = StrategyFixed
	DefaultLogLevel     = "info"
)

// Config holds all runtime configuration for a retryctl invocation.
type Config struct {
	// MaxAttempts is the total number of times the command will be tried.
	MaxAttempts int

	// Delay is the base duration to wait between attempts.
	Delay time.Duration

	// Strategy controls how the delay grows between attempts.
	Strategy BackoffStrategy

	// LogLevel sets the minimum log level ("info", "warn", "error").
	LogLevel string

	// Command is the executable to run.
	Command string

	// Args are the arguments passed to Command.
	Args []string
}

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	return &Config{
		MaxAttempts: DefaultMaxAttempts,
		Delay:       DefaultDelay,
		Strategy:    DefaultStrategy,
		LogLevel:    DefaultLogLevel,
	}
}

// Validate checks that the Config contains valid values.
func (c *Config) Validate() error {
	if c.Command == "" {
		return errors.New("config: command must not be empty")
	}
	if c.MaxAttempts < 1 {
		return errors.New("config: max_attempts must be at least 1")
	}
	if c.Delay < 0 {
		return errors.New("config: delay must not be negative")
	}
	switch c.Strategy {
	case StrategyFixed, StrategyLinear, StrategyExponential:
		// valid
	default:
		return errors.New("config: unknown backoff strategy: " + string(c.Strategy))
	}
	return nil
}
