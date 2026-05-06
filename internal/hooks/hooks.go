package hooks

import (
	"os/exec"
	"strings"
)

// Event represents a lifecycle event during command execution.
type Event string

const (
	EventBeforeAttempt Event = "before_attempt"
	EventAfterSuccess  Event = "after_success"
	EventAfterFailure  Event = "after_failure"
	EventAfterFinal    Event = "after_final"
)

// Hook holds a shell command to execute on a given lifecycle event.
type Hook struct {
	Event   Event
	Command string
}

// Runner executes registered hooks for lifecycle events.
type Runner struct {
	hooks []Hook
}

// New creates a new hook Runner with the provided hooks.
func New(hooks []Hook) *Runner {
	return &Runner{hooks: hooks}
}

// Run executes all hooks registered for the given event.
// Extra key=value pairs are injected as environment variables.
func (r *Runner) Run(event Event, env map[string]string) error {
	for _, h := range r.hooks {
		if h.Event != event {
			continue
		}
		parts := strings.Fields(h.Command)
		if len(parts) == 0 {
			continue
		}
		cmd := exec.Command(parts[0], parts[1:]...) //nolint:gosec
		for k, v := range env {
			cmd.Env = append(cmd.Env, k+"="+v)
		}
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

// HasHooks returns true if any hook is registered for the given event.
func (r *Runner) HasHooks(event Event) bool {
	for _, h := range r.hooks {
		if h.Event == event {
			return true
		}
	}
	return false
}
