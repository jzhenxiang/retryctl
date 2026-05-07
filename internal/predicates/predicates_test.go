package predicates_test

import (
	"testing"

	"retryctl/internal/predicates"
)

func TestAlwaysReturnsTrue(t *testing.T) {
	p := predicates.Always()
	if !p(0, "") {
		t.Fatal("Always should return true")
	}
	if !p(1, "some output") {
		t.Fatal("Always should return true for non-zero exit")
	}
}

func TestNeverReturnsFalse(t *testing.T) {
	p := predicates.Never()
	if p(1, "output") {
		t.Fatal("Never should return false")
	}
}

func TestOnExitCodesMatch(t *testing.T) {
	p, err := predicates.OnExitCodes(1, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !p(1, "") {
		t.Error("exit code 1 should match")
	}
	if !p(2, "") {
		t.Error("exit code 2 should match")
	}
	if p(3, "") {
		t.Error("exit code 3 should not match")
	}
}

func TestOnExitCodesEmptyReturnsError(t *testing.T) {
	_, err := predicates.OnExitCodes()
	if err == nil {
		t.Fatal("expected error for empty codes")
	}
}

func TestOnOutputContainsMatch(t *testing.T) {
	p, err := predicates.OnOutputContains("timeout")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !p(0, "connection timeout occurred") {
		t.Error("output containing substring should match")
	}
	if p(0, "all good") {
		t.Error("output without substring should not match")
	}
}

func TestOnOutputContainsEmptyReturnsError(t *testing.T) {
	_, err := predicates.OnOutputContains("")
	if err == nil {
		t.Fatal("expected error for empty substring")
	}
}

func TestAnyRetryIfOneMatches(t *testing.T) {
	p, err := predicates.Any(predicates.Never(), predicates.Always())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !p(1, "") {
		t.Error("Any should return true when at least one predicate matches")
	}
}

func TestAnyEmptyReturnsError(t *testing.T) {
	_, err := predicates.Any()
	if err == nil {
		t.Fatal("expected error for empty predicate list")
	}
}

func TestAllRetryOnlyWhenAllMatch(t *testing.T) {
	p, err := predicates.All(predicates.Always(), predicates.Always())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !p(1, "") {
		t.Error("All should return true when every predicate matches")
	}

	p2, _ := predicates.All(predicates.Always(), predicates.Never())
	if p2(1, "") {
		t.Error("All should return false when one predicate does not match")
	}
}

func TestAllEmptyReturnsError(t *testing.T) {
	_, err := predicates.All()
	if err == nil {
		t.Fatal("expected error for empty predicate list")
	}
}
