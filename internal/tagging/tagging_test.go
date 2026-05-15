package tagging_test

import (
	"testing"

	"github.com/yourorg/retryctl/internal/tagging"
)

func TestNewEmptyPairs(t *testing.T) {
	tgr, err := tagging.New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tgr.Len() != 0 {
		t.Fatalf("expected 0 tags, got %d", tgr.Len())
	}
}

func TestNewValidPairs(t *testing.T) {
	tgr, err := tagging.New([]string{"env=prod", "region=us-east-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tgr.Len() != 2 {
		t.Fatalf("expected 2 tags, got %d", tgr.Len())
	}
}

func TestNewMalformedPair(t *testing.T) {
	_, err := tagging.New([]string{"nodivider"})
	if err == nil {
		t.Fatal("expected error for malformed pair")
	}
}

func TestNewEmptyKey(t *testing.T) {
	_, err := tagging.New([]string{"=value"})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestTagsReturnsCopy(t *testing.T) {
	tgr, _ := tagging.New([]string{"a=1"})
	tags, err := tgr.Tags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tags["a"] = "mutated"
	again, _ := tgr.Tags()
	if again["a"] != "1" {
		t.Fatal("base tags were mutated")
	}
}

func TestTagsMergesExtra(t *testing.T) {
	tgr, _ := tagging.New([]string{"env=prod"})
	tags, err := tgr.Tags("attempt=3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tags["env"] != "prod" {
		t.Fatalf("expected env=prod, got %q", tags["env"])
	}
	if tags["attempt"] != "3" {
		t.Fatalf("expected attempt=3, got %q", tags["attempt"])
	}
}

func TestTagsExtraOverridesBase(t *testing.T) {
	tgr, _ := tagging.New([]string{"env=prod"})
	tags, _ := tgr.Tags("env=staging")
	if tags["env"] != "staging" {
		t.Fatalf("extra should override base, got %q", tags["env"])
	}
}

func TestTagsMalformedExtra(t *testing.T) {
	tgr, _ := tagging.New(nil)
	_, err := tgr.Tags("bad")
	if err == nil {
		t.Fatal("expected error for malformed extra pair")
	}
}

func TestHas(t *testing.T) {
	tgr, _ := tagging.New([]string{"x=1"})
	if !tgr.Has("x") {
		t.Fatal("expected Has(x) == true")
	}
	if tgr.Has("y") {
		t.Fatal("expected Has(y) == false")
	}
}
