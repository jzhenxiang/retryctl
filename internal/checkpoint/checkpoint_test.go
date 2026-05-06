package checkpoint_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"retryctl/internal/checkpoint"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "checkpoint.json")
}

func TestLoadMissingFileReturnsZero(t *testing.T) {
	store := checkpoint.New(tempPath(t))
	st, err := store.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if st.Attempt != 0 || st.Command != "" {
		t.Fatalf("expected zero state, got %+v", st)
	}
}

func TestSaveAndLoad(t *testing.T) {
	store := checkpoint.New(tempPath(t))
	want := checkpoint.State{
		Command:    "echo hello",
		Attempt:    3,
		LastFailed: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
	}
	if err := store.Save(want); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Command != want.Command {
		t.Errorf("Command: got %q want %q", got.Command, want.Command)
	}
	if got.Attempt != want.Attempt {
		t.Errorf("Attempt: got %d want %d", got.Attempt, want.Attempt)
	}
	if !got.LastFailed.Equal(want.LastFailed) {
		t.Errorf("LastFailed: got %v want %v", got.LastFailed, want.LastFailed)
	}
}

func TestSaveOverwritesPreviousState(t *testing.T) {
	store := checkpoint.New(tempPath(t))
	first := checkpoint.State{Command: "cmd", Attempt: 1}
	second := checkpoint.State{Command: "cmd", Attempt: 5}

	if err := store.Save(first); err != nil {
		t.Fatalf("first Save: %v", err)
	}
	if err := store.Save(second); err != nil {
		t.Fatalf("second Save: %v", err)
	}
	got, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Attempt != 5 {
		t.Errorf("expected attempt 5, got %d", got.Attempt)
	}
}

func TestRemoveDeletesFile(t *testing.T) {
	path := tempPath(t)
	store := checkpoint.New(path)

	if err := store.Save(checkpoint.State{Command: "x", Attempt: 1}); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if err := store.Remove(); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("expected file to be removed")
	}
}

func TestRemoveIsIdempotent(t *testing.T) {
	store := checkpoint.New(tempPath(t))
	if err := store.Remove(); err != nil {
		t.Fatalf("Remove on non-existent file: %v", err)
	}
}
