package config

import (
	"testing"
	"time"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	if cfg.MaxAttempts != DefaultMaxAttempts {
		t.Errorf("expected MaxAttempts %d, got %d", DefaultMaxAttempts, cfg.MaxAttempts)
	}
	if cfg.Delay != DefaultDelay {
		t.Errorf("expected Delay %v, got %v", DefaultDelay, cfg.Delay)
	}
	if cfg.Strategy != DefaultStrategy {
		t.Errorf("expected Strategy %q, got %q", DefaultStrategy, cfg.Strategy)
	}
	if cfg.LogLevel != DefaultLogLevel {
		t.Errorf("expected LogLevel %q, got %q", DefaultLogLevel, cfg.LogLevel)
	}
}

func TestValidateOK(t *testing.T) {
	cfg := Default()
	cfg.Command = "echo"

	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidateMissingCommand(t *testing.T) {
	cfg := Default()
	// Command intentionally left empty

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for empty command, got nil")
	}
}

func TestValidateInvalidMaxAttempts(t *testing.T) {
	cfg := Default()
	cfg.Command = "echo"
	cfg.MaxAttempts = 0

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for max_attempts=0, got nil")
	}
}

func TestValidateNegativeDelay(t *testing.T) {
	cfg := Default()
	cfg.Command = "echo"
	cfg.Delay = -1 * time.Second

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative delay, got nil")
	}
}

func TestValidateUnknownStrategy(t *testing.T) {
	cfg := Default()
	cfg.Command = "echo"
	cfg.Strategy = BackoffStrategy("random")

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for unknown strategy, got nil")
	}
}

func TestValidateAllStrategies(t *testing.T) {
	strategies := []BackoffStrategy{StrategyFixed, StrategyLinear, StrategyExponential}

	for _, s := range strategies {
		cfg := Default()
		cfg.Command = "echo"
		cfg.Strategy = s

		if err := cfg.Validate(); err != nil {
			t.Errorf("strategy %q: unexpected error: %v", s, err)
		}
	}
}
