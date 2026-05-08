package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/retryctl/internal/snapshot"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "snapshot.json")
}

func TestNewEmptyPathReturnsError(t *testing.T) {
	_, err := snapshot.New("")
	if err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}

func TestLoadMissingFileReturnsZero(t *testing.T) {
	store, err := snapshot.New(tempPath(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	st, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if st.Attempt != 0 {
		t.Errorf("expected zero Attempt, got %d", st.Attempt)
	}
}

func TestSaveAndLoad(t *testing.T) {
	store, _ := snapshot.New(tempPath(t))
	want := snapshot.State{
		Attempt:      3,
		LastFailedAt: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		Command:      []string{"curl", "-f", "https://example.com"},
	}
	if err := store.Save(want); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Attempt != want.Attempt {
		t.Errorf("Attempt: got %d, want %d", got.Attempt, want.Attempt)
	}
	if !got.LastFailedAt.Equal(want.LastFailedAt) {
		t.Errorf("LastFailedAt: got %v, want %v", got.LastFailedAt, want.LastFailedAt)
	}
	if len(got.Command) != len(want.Command) || got.Command[0] != want.Command[0] {
		t.Errorf("Command: got %v, want %v", got.Command, want.Command)
	}
}

func TestSaveOverwritesPreviousState(t *testing.T) {
	store, _ := snapshot.New(tempPath(t))
	_ = store.Save(snapshot.State{Attempt: 1})
	_ = store.Save(snapshot.State{Attempt: 5})
	st, _ := store.Load()
	if st.Attempt != 5 {
		t.Errorf("expected Attempt 5, got %d", st.Attempt)
	}
}

func TestRemoveDeletesFile(t *testing.T) {
	path := tempPath(t)
	store, _ := snapshot.New(path)
	_ = store.Save(snapshot.State{Attempt: 2})
	if err := store.Remove(); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("expected file to be deleted")
	}
}

func TestRemoveIsNoOpWhenMissing(t *testing.T) {
	store, _ := snapshot.New(tempPath(t))
	if err := store.Remove(); err != nil {
		t.Errorf("Remove on missing file: %v", err)
	}
}
