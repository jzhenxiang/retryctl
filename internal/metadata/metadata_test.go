package metadata_test

import (
	"testing"

	"github.com/yourorg/retryctl/internal/metadata"
)

func TestNewEmptyPairs(t *testing.T) {
	m, err := metadata.New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.All()) != 0 {
		t.Errorf("expected empty map, got %v", m.All())
	}
}

func TestNewValidPairs(t *testing.T) {
	m, err := metadata.New([]string{"env=prod", "region=us-east-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, ok := m.Get("env"); !ok || v != "prod" {
		t.Errorf("expected env=prod, got %q ok=%v", v, ok)
	}
	if v, ok := m.Get("region"); !ok || v != "us-east-1" {
		t.Errorf("expected region=us-east-1, got %q ok=%v", v, ok)
	}
}

func TestNewMalformedPair(t *testing.T) {
	_, err := metadata.New([]string{"noequals"})
	if err == nil {
		t.Fatal("expected error for malformed pair")
	}
}

func TestNewEmptyKey(t *testing.T) {
	_, err := metadata.New([]string{"=value"})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestGetMissingKey(t *testing.T) {
	m, _ := metadata.New([]string{"a=1"})
	_, ok := m.Get("missing")
	if ok {
		t.Error("expected ok=false for missing key")
	}
}

func TestAllReturnsCopy(t *testing.T) {
	m, _ := metadata.New([]string{"x=1"})
	all := m.All()
	all["x"] = "mutated"
	if v, _ := m.Get("x"); v != "1" {
		t.Errorf("original should not be mutated, got %q", v)
	}
}

func TestMergeOtherTakesPrecedence(t *testing.T) {
	a, _ := metadata.New([]string{"k=base", "only-a=yes"})
	b, _ := metadata.New([]string{"k=override", "only-b=yes"})
	c := a.Merge(b)
	if v, _ := c.Get("k"); v != "override" {
		t.Errorf("expected override, got %q", v)
	}
	if _, ok := c.Get("only-a"); !ok {
		t.Error("expected only-a to be present after merge")
	}
	if _, ok := c.Get("only-b"); !ok {
		t.Error("expected only-b to be present after merge")
	}
}

func TestMergeNilOther(t *testing.T) {
	a, _ := metadata.New([]string{"k=v"})
	c := a.Merge(nil)
	if v, _ := c.Get("k"); v != "v" {
		t.Errorf("expected k=v after nil merge, got %q", v)
	}
}
