// Package snapshot captures and persists the state of a retry run so that
// a subsequent invocation can resume from where it left off.
package snapshot

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

// State holds the persisted state of an in-progress retry run.
type State struct {
	// Attempt is the next attempt number to execute (1-based).
	Attempt int `json:"attempt"`
	// LastFailedAt is the wall-clock time of the most recent failure.
	LastFailedAt time.Time `json:"last_failed_at"`
	// Command is the command that was being retried.
	Command []string `json:"command"`
}

// Store persists and retrieves snapshot state from a file.
type Store struct {
	path string
}

// New returns a Store that uses path as its backing file.
func New(path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("snapshot: path must not be empty")
	}
	return &Store{path: path}, nil
}

// Save writes state to the backing file, overwriting any previous content.
func (s *Store) Save(st State) error {
	f, err := os.Create(s.path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(st)
}

// Load reads state from the backing file. If the file does not exist it
// returns a zero-value State and a nil error so callers can treat a missing
// snapshot as a fresh start.
func (s *Store) Load() (State, error) {
	f, err := os.Open(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return State{}, nil
	}
	if err != nil {
		return State{}, err
	}
	defer f.Close()
	var st State
	if err := json.NewDecoder(f).Decode(&st); err != nil {
		return State{}, err
	}
	return st, nil
}

// Remove deletes the backing file. It is a no-op if the file does not exist.
func (s *Store) Remove() error {
	err := os.Remove(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
