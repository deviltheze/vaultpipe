package dotenv

import (
	"testing"
)

func TestCompare_Identical(t *testing.T) {
	a := map[string]string{"FOO": "bar", "BAZ": "qux"}
	b := map[string]string{"FOO": "bar", "BAZ": "qux"}

	r := Compare(a, b)

	if len(r.Identical) != 2 {
		t.Fatalf("expected 2 identical, got %d", len(r.Identical))
	}
	if r.HasDifferences() {
		t.Error("expected no differences")
	}
}

func TestCompare_OnlyInA(t *testing.T) {
	a := map[string]string{"FOO": "bar", "NEW": "value"}
	b := map[string]string{"FOO": "bar"}

	r := Compare(a, b)

	if len(r.OnlyInA) != 1 {
		t.Fatalf("expected 1 only-in-A, got %d", len(r.OnlyInA))
	}
	if r.OnlyInA["NEW"] != "value" {
		t.Errorf("unexpected value for NEW: %q", r.OnlyInA["NEW"])
	}
}

func TestCompare_OnlyInB(t *testing.T) {
	a := map[string]string{"FOO": "bar"}
	b := map[string]string{"FOO": "bar", "EXTRA": "old"}

	r := Compare(a, b)

	if len(r.OnlyInB) != 1 {
		t.Fatalf("expected 1 only-in-B, got %d", len(r.OnlyInB))
	}
	if r.OnlyInB["EXTRA"] != "old" {
		t.Errorf("unexpected value for EXTRA: %q", r.OnlyInB["EXTRA"])
	}
}

func TestCompare_Different(t *testing.T) {
	a := map[string]string{"FOO": "new_value"}
	b := map[string]string{"FOO": "old_value"}

	r := Compare(a, b)

	if len(r.Different) != 1 {
		t.Fatalf("expected 1 different, got %d", len(r.Different))
	}
	pair := r.Different["FOO"]
	if pair[0] != "new_value" || pair[1] != "old_value" {
		t.Errorf("unexpected pair: %v", pair)
	}
	if !r.HasDifferences() {
		t.Error("expected HasDifferences to be true")
	}
}

func TestCompare_EmptyMaps(t *testing.T) {
	r := Compare(map[string]string{}, map[string]string{})

	if r.HasDifferences() {
		t.Error("expected no differences for empty maps")
	}
}

func TestCompare_Summary(t *testing.T) {
	a := map[string]string{"A": "1", "B": "changed"}
	b := map[string]string{"B": "original", "C": "3"}

	r := Compare(a, b)
	s := r.Summary()

	if s == "" {
		t.Error("expected non-empty summary")
	}
	// 1 only in A, 1 only in B, 1 different, 0 identical
	expected := "1 only in source, 1 only in target, 1 changed, 0 identical"
	if s != expected {
		t.Errorf("summary mismatch:\n got:  %q\n want: %q", s, expected)
	}
}
