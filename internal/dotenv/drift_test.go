package dotenv_test

import (
	"strings"
	"testing"

	"github.com/your-org/vaultpipe/internal/dotenv"
)

func TestDetectDrift_Clean(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	vault := map[string]string{"FOO": "bar", "BAZ": "qux"}

	report := dotenv.DetectDrift(env, vault)
	if !report.Clean {
		t.Fatalf("expected clean report, got: %s", report)
	}
	if len(report.Entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(report.Entries))
	}
}

func TestDetectDrift_Modified(t *testing.T) {
	env := map[string]string{"FOO": "old"}
	vault := map[string]string{"FOO": "new"}

	report := dotenv.DetectDrift(env, vault)
	if report.Clean {
		t.Fatal("expected dirty report")
	}
	if len(report.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(report.Entries))
	}
	e := report.Entries[0]
	if e.Kind != "modified" {
		t.Errorf("expected kind=modified, got %q", e.Kind)
	}
	if e.Key != "FOO" {
		t.Errorf("expected key=FOO, got %q", e.Key)
	}
}

func TestDetectDrift_MissingInEnv(t *testing.T) {
	env := map[string]string{}
	vault := map[string]string{"SECRET": "val"}

	report := dotenv.DetectDrift(env, vault)
	if report.Clean {
		t.Fatal("expected dirty report")
	}
	if report.Entries[0].Kind != "missing_in_env" {
		t.Errorf("expected missing_in_env, got %q", report.Entries[0].Kind)
	}
}

func TestDetectDrift_ExtraInEnv(t *testing.T) {
	env := map[string]string{"EXTRA": "val"}
	vault := map[string]string{}

	report := dotenv.DetectDrift(env, vault)
	if report.Clean {
		t.Fatal("expected dirty report")
	}
	if report.Entries[0].Kind != "extra_in_env" {
		t.Errorf("expected extra_in_env, got %q", report.Entries[0].Kind)
	}
}

func TestDetectDrift_StringSummary(t *testing.T) {
	env := map[string]string{"A": "old"}
	vault := map[string]string{"A": "new", "B": "val"}

	report := dotenv.DetectDrift(env, vault)
	summary := report.String()
	if !strings.Contains(summary, "drift detected") {
		t.Errorf("expected 'drift detected' in summary, got: %s", summary)
	}
}

func TestDetectDrift_CleanString(t *testing.T) {
	report := dotenv.DriftReport{Clean: true}
	if report.String() != "no drift detected" {
		t.Errorf("unexpected clean string: %s", report.String())
	}
}
