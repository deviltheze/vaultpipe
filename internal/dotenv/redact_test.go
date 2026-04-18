package dotenv

import (
	"testing"
)

func TestRedact_SensitiveKeysAreHidden(t *testing.T) {
	input := map[string]string{
		"DB_PASSWORD": "supersecret",
		"API_KEY":     "abc123",
		"APP_TOKEN":   "tok_xyz",
		"DB_HOST":     "localhost",
		"PORT":        "5432",
	}

	result := Redact(input)

	if result["DB_PASSWORD"] != RedactedValue {
		t.Errorf("expected DB_PASSWORD to be redacted, got %q", result["DB_PASSWORD"])
	}
	if result["API_KEY"] != RedactedValue {
		t.Errorf("expected API_KEY to be redacted, got %q", result["API_KEY"])
	}
	if result["APP_TOKEN"] != RedactedValue {
		t.Errorf("expected APP_TOKEN to be redacted, got %q", result["APP_TOKEN"])
	}
	if result["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST to be unchanged, got %q", result["DB_HOST"])
	}
	if result["PORT"] != "5432" {
		t.Errorf("expected PORT to be unchanged, got %q", result["PORT"])
	}
}

func TestRedact_EmptyMap(t *testing.T) {
	result := Redact(map[string]string{})
	if len(result) != 0 {
		t.Errorf("expected empty map, got %d entries", len(result))
	}
}

func TestRedact_DoesNotMutateOriginal(t *testing.T) {
	input := map[string]string{"SECRET_KEY": "original"}
	_ = Redact(input)
	if input["SECRET_KEY"] != "original" {
		t.Error("Redact mutated the original map")
	}
}

func TestRedact_CaseInsensitiveKey(t *testing.T) {
	input := map[string]string{"db_password": "hidden", "db_host": "visible"}
	result := Redact(input)
	if result["db_password"] != RedactedValue {
		t.Errorf("expected lowercase db_password to be redacted, got %q", result["db_password"])
	}
	if result["db_host"] != "visible" {
		t.Errorf("expected db_host unchanged, got %q", result["db_host"])
	}
}
