package dotenv

import (
	"strings"
	"testing"
)

func TestFormat_EmptySecrets(t *testing.T) {
	result := Format(map[string]string{}, FormatOptions{})
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestFormat_SortedKeys(t *testing.T) {
	secrets := map[string]string{
		"ZEBRA": "1",
		"ALPHA": "2",
		"MANGO": "3",
	}
	out := Format(secrets, FormatOptions{SortKeys: true})
	lines := nonEmpty(strings.Split(out, "\n"))
	if lines[0] != "ALPHA=2" {
		t.Errorf("expected ALPHA first, got %q", lines[0])
	}
	if lines[2] != "ZEBRA=1" {
		t.Errorf("expected ZEBRA last, got %q", lines[2])
	}
}

func TestFormat_QuotesSpecialValues(t *testing.T) {
	secrets := map[string]string{
		"KEY": "hello world",
	}
	out := Format(secrets, FormatOptions{})
	if !strings.Contains(out, `KEY="hello world"`) {
		t.Errorf("expected quoted value, got %q", out)
	}
}

func TestFormat_NoQuotesForSimpleValues(t *testing.T) {
	secrets := map[string]string{
		"KEY": "simplevalue",
	}
	out := Format(secrets, FormatOptions{})
	if strings.Contains(out, `"`) {
		t.Errorf("expected no quotes, got %q", out)
	}
}

func TestFormat_GroupByPrefix(t *testing.T) {
	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"AWS_KEY":  "abc",
		"AWS_SECRET": "xyz",
	}
	out := Format(secrets, FormatOptions{SortKeys: true, GroupByPrefix: true})
	if !strings.Contains(out, "# AWS") {
		t.Errorf("expected AWS group header, got:\n%s", out)
	}
	if !strings.Contains(out, "# DB") {
		t.Errorf("expected DB group header, got:\n%s", out)
	}
	awsIdx := strings.Index(out, "# AWS")
	dbIdx := strings.Index(out, "# DB")
	if awsIdx > dbIdx {
		t.Errorf("expected AWS before DB alphabetically")
	}
}

func TestFormat_GroupByPrefix_NoPrefixKey(t *testing.T) {
	secrets := map[string]string{
		"STANDALONE": "value",
	}
	out := Format(secrets, FormatOptions{GroupByPrefix: true})
	if !strings.Contains(out, "# STANDALONE") {
		t.Errorf("expected STANDALONE as its own group, got:\n%s", out)
	}
}

func TestFormat_CustomSeparator(t *testing.T) {
	secrets := map[string]string{
		"DB.HOST": "localhost",
		"DB.PORT": "5432",
	}
	out := Format(secrets, FormatOptions{GroupByPrefix: true, PrefixSeparator: "."})
	if !strings.Contains(out, "# DB") {
		t.Errorf("expected DB group with dot separator, got:\n%s", out)
	}
}

// nonEmpty filters out empty strings from a slice.
func nonEmpty(ss []string) []string {
	var out []string
	for _, s := range ss {
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}
