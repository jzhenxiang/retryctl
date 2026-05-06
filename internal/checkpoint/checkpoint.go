// Package checkpoint provides attempt-progress persistence so that
// retryctl can resume from the last known attempt index after a
// process restart rather than starting from attempt 1 every time.
package checkpoint

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

// State holds the persisted retry progress for a single invocation.
type State struct {
	Command    string    `json:"command"`
	Attempt    int       `json:"attempt"`
	LastFailed time.Time `json:"last_failed"`
}

// Store persists and retrieves checkpoint state via a file on disk.
type Store struct {
	path string
}

// New returns a Store that uses path as the backing file.
// The file is created on first Save; its parent directory must exist.
func New(path string) *Store {
	return &Store{path: path}
}

// Save writes state to the checkpoint file, overwriting any previous
// content. It is safe to call concurrently only when callers serialise
// access externally.
func (s *Store) Save(st State) error {
	f, err := os.CreateTemp("", "checkpoint-*")
	if err != nil {
		return err
	}
	tmpName := f.Name()

	if err := json.NewEncoder(f).Encode(st); err != nil {
		f.Close()
		os.Remove(tmpName)
		return err
	}
	if err := f.Close(); err != nil {
		os.Remove(tmpName)
		return err
	}
	return os.Rename(tmpName, s.path)
}

// Load reads and returns the previously saved State.
// If the checkpoint file does not exist, Load returns a zero State and
// a nil error so callers can treat a missing file as "start fresh".
func (s *Store) Load() (State, error) {
	f, err := os.Open(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return State{}, nil
		}
		return State{}, err
	}
	defer f.Close()

	var st State
	if err := json.NewDecoder(f).Decode(&st); err != nil {
		return State{}, err
	}
	return st, nil
}

// Remove deletes the checkpoint file. It is a no-op when the file does
// not exist.
func (s *Store) Remove() error {
	err := os.Remove(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
