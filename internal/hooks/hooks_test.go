package hooks_test

import (
	"runtime"
	"testing"

	"github.com/example/retryctl/internal/hooks"
)

func TestHasHooksTrue(t *testing.T) {
	r := hooks.New([]hooks.Hook{
		{Event: hooks.EventAfterSuccess, Command: "echo ok"},
	})
	if !r.HasHooks(hooks.EventAfterSuccess) {
		t.Fatal("expected HasHooks to return true")
	}
}

func TestHasHooksFalse(t *testing.T) {
	r := hooks.New([]hooks.Hook{})
	if r.HasHooks(hooks.EventBeforeAttempt) {
		t.Fatal("expected HasHooks to return false")
	}
}

func TestRunNoHooks(t *testing.T) {
	r := hooks.New(nil)
	if err := r.Run(hooks.EventAfterFinal, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunMatchingHook(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping on windows")
	}
	r := hooks.New([]hooks.Hook{
		{Event: hooks.EventAfterSuccess, Command: "true"},
	})
	if err := r.Run(hooks.EventAfterSuccess, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunNonMatchingHookSkipped(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping on windows")
	}
	// Register a hook that would fail, but for a different event.
	r := hooks.New([]hooks.Hook{
		{Event: hooks.EventAfterFailure, Command: "false"},
	})
	if err := r.Run(hooks.EventAfterSuccess, nil); err != nil {
		t.Fatalf("non-matching hook should be skipped, got error: %v", err)
	}
}

func TestRunFailingHookReturnsError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping on windows")
	}
	r := hooks.New([]hooks.Hook{
		{Event: hooks.EventBeforeAttempt, Command: "false"},
	})
	if err := r.Run(hooks.EventBeforeAttempt, nil); err == nil {
		t.Fatal("expected error from failing hook command")
	}
}

func TestRunEmptyCommandSkipped(t *testing.T) {
	r := hooks.New([]hooks.Hook{
		{Event: hooks.EventAfterFinal, Command: ""},
	})
	if err := r.Run(hooks.EventAfterFinal, nil); err != nil {
		t.Fatalf("empty command should be skipped, got: %v", err)
	}
}
