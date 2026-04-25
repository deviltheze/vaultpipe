package dotenv

import (
	"strings"
	"testing"
)

func TestSanitize_TrimWhitespace(t *testing.T) {
	input := map[string]string{
		"KEY": "  hello world  ",
		"OTHER": "\tno trim\t",
	}
	out := Sanitize(input, SanitizeOptions{TrimWhitespace: true})
	if got := out["KEY"]; got != "hello world" {
		t.Errorf("KEY: got %q, want %q", got, "hello world")
	}
	if got := out["OTHER"]; got != "no trim" {
		t.Errorf("OTHER: got %q, want %q", got, "no trim")
	}
}

func TestSanitize_StripControlChars(t *testing.T) {
	input := map[string]string{
		"KEY": "hello\x00world\x01",
		"TABS": "keep\ttabs",
		"NEWLINE": "keep\nnewlines",
	}
	out := Sanitize(input, SanitizeOptions{StripControlChars: true})
	if got := out["KEY"]; got != "helloworld" {
		t.Errorf("KEY: got %q, want %q", got, "helloworld")
	}
	if got := out["TABS"]; got != "keep\ttabs" {
		t.Errorf("TABS: tabs should be preserved, got %q", got)
	}
	if got := out["NEWLINE"]; got != "keep\nnewlines" {
		t.Errorf("NEWLINE: newlines should be preserved, got %q", got)
	}
}

func TestSanitize_MaxValueLength(t *testing.T) {
	input := map[string]string{
		"LONG": strings.Repeat("a", 200),
		"SHORT": "hi",
	}
	out := Sanitize(input, SanitizeOptions{MaxValueLength: 100})
	if got := len(out["LONG"]); got != 100 {
		t.Errorf("LONG: got length %d, want 100", got)
	}
	if got := out["SHORT"]; got != "hi" {
		t.Errorf("SHORT: got %q, want %q", got, "hi")
	}
}

func TestSanitize_ZeroMaxLength_NoTruncation(t *testing.T) {
	long := strings.Repeat("x", 500)
	input := map[string]string{"KEY": long}
	out := Sanitize(input, SanitizeOptions{MaxValueLength: 0})
	if got := len(out["KEY"]); got != 500 {
		t.Errorf("expected no truncation, got length %d", got)
	}
}

func TestSanitize_DoesNotMutateOriginal(t *testing.T) {
	input := map[string]string{"KEY": "  value  "}
	_ = Sanitize(input, SanitizeOptions{TrimWhitespace: true})
	if input["KEY"] != "  value  " {
		t.Error("original map was mutated")
	}
}

func TestSanitize_NoOptions_PassThrough(t *testing.T) {
	input := map[string]string{"A": "  raw  ", "B": "\x00ctrl"}
	out := Sanitize(input, SanitizeOptions{})
	if out["A"] != "  raw  " {
		t.Errorf("A: expected passthrough, got %q", out["A"])
	}
	if out["B"] != "\x00ctrl" {
		t.Errorf("B: expected passthrough, got %q", out["B"])
	}
}
