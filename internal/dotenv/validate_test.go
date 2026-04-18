package dotenv

import (
	"strings"
	"testing"
)

func TestValidate_ValidKeys(t *testing.T) {
	secrets := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"_PRIVATE":     "value",
		"key123":       "val",
	}
	if err := Validate(secrets, 0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_InvalidKey_StartsWithDigit(t *testing.T) {
	secrets := map[string]string{"1BAD": "value"}
	err := Validate(secrets, 0)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "1BAD") {
		t.Errorf("error should mention bad key, got: %v", err)
	}
}

func TestValidate_InvalidKey_ContainsHyphen(t *testing.T) {
	secrets := map[string]string{"BAD-KEY": "value"}
	err := Validate(secrets, 0)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestValidate_ValueTooLong(t *testing.T) {
	secrets := map[string]string{"MY_SECRET": strings.Repeat("x", 201)}
	err := Validate(secrets, 200)
	if err == nil {
		t.Fatal("expected error for oversized value")
	}
	if !strings.Contains(err.Error(), "MY_SECRET") {
		t.Errorf("error should mention key name, got: %v", err)
	}
}

func TestValidate_ValueLengthSkippedWhenZero(t *testing.T) {
	secrets := map[string]string{"KEY": strings.Repeat("x", 10000)}
	if err := Validate(secrets, 0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_MultipleIssues(t *testing.T) {
	secrets := map[string]string{
		"bad key": "ok",
		"GOOD":    strings.Repeat("y", 50),
	}
	err := Validate(secrets, 10)
	if err == nil {
		t.Fatal("expected error")
	}
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.Issues) < 2 {
		t.Errorf("expected at least 2 issues, got %d", len(ve.Issues))
	}
}
