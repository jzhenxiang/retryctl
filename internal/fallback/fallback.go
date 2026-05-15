// Package fallback provides a mechanism to execute an alternative command
// when the primary command exhausts all retry attempts.
package fallback

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
)

// ErrNoFallback is returned when no fallback command is configured.
var ErrNoFallback = errors.New("fallback: no fallback command configured")

// Result holds the outcome of a fallback execution.
type Result struct {
	Output   []byte
	ExitCode int
}

// Handler executes a pre-configured fallback command.
type Handler struct {
	args []string
}

// New creates a new Handler with the given fallback command and arguments.
// Returns an error if args is empty.
func New(args []string) (*Handler, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("fallback: args must not be empty")
	}
	return &Handler{args: args}, nil
}

// HasFallback returns true when a fallback command is configured.
func (h *Handler) HasFallback() bool {
	return h != nil && len(h.args) > 0
}

// Run executes the fallback command using the provided context.
// It returns a Result containing the combined output and exit code.
func (h *Handler) Run(ctx context.Context) (*Result, error) {
	if !h.HasFallback() {
		return nil, ErrNoFallback
	}

	//nolint:gosec // args are provided by the operator
	cmd := exec.CommandContext(ctx, h.args[0], h.args[1:]...)
	out, err := cmd.CombinedOutput()

	result := &Result{Output: out}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		result.ExitCode = exitErr.ExitCode()
		return result, nil
	}
	if err != nil {
		return nil, fmt.Errorf("fallback: exec: %w", err)
	}

	return result, nil
}
