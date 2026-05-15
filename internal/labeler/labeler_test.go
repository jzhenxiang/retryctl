package labeler_test

import (
	"testing"

	"github.com/yourorg/retryctl/internal/labeler"
)

func TestNewValidPairs(t *testing.T) {
	l, err := labeler.New([]string{"env=prod", "team=platform"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	base := l.Base()
	if base["env"] != "prod" {
		t.Errorf("env: got %q want %q", base["env"], "prod")
	}
	if base["team"] != "platform" {
		t.Errorf("team: got %q want %q", base["team"], "platform")
	}
}

func TestNewEmptyPairs(t *testing.T) {
	l, err := labeler.New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(l.Base()) != 0 {
		t.Errorf("expected empty base labels")
	}
}

func TestNewMalformedPair(t *testing.T) {
	_, err := labeler.New([]string{"nodash"})
	if err == nil {
		t.Fatal("expected error for malformed pair, got nil")
	}
}

func TestNewEmptyKey(t *testing.T) {
	_, err := labeler.New([]string{"=value"})
	if err == nil {
		t.Fatal("expected error for empty key, got nil")
	}
}

func TestWithMergesLabels(t *testing.T) {
	l, _ := labeler.New([]string{"env=prod"})
	merged := l.With(labeler.Labels{"attempt": "3"})
	if merged["env"] != "prod" {
		t.Errorf("env: got %q want %q", merged["env"], "prod")
	}
	if merged["attempt"] != "3" {
		t.Errorf("attempt: got %q want %q", merged["attempt"], "3")
	}
}

func TestWithExtraOverridesBase(t *testing.T) {
	l, _ := labeler.New([]string{"env=prod"})
	merged := l.With(labeler.Labels{"env": "staging"})
	if merged["env"] != "staging" {
		t.Errorf("env: got %q want %q", merged["env"], "staging")
	}
}

func TestBaseReturnsCopy(t *testing.T) {
	l, _ := labeler.New([]string{"k=v"})
	b1 := l.Base()
	b1["k"] = "mutated"
	b2 := l.Base()
	if b2["k"] != "v" {
		t.Errorf("Base should return an independent copy")
	}
}
