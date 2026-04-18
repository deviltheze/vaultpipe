package sync

import (
	"strings"
	"testing"
)

func TestValidateSecrets_Valid(t *testing.T) {
	secrets := map[string]string{
		"API_KEY":  "abc123",
		"DB_PASS":  "s3cr3t",
	}
	if err := validateSecrets("secret/app", secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateSecrets_BadKey(t *testing.T) {
	secrets := map[string]string{"bad-key": "value"}
	err := validateSecrets("secret/app", secrets)
	if err == nil {
		t.Fatal("expected validation error")
	}
	if !strings.Contains(err.Error(), "secret/app") {
		t.Errorf("error should include path, got: %v", err)
	}
}

func TestValidateSecrets_ValueTooLong(t *testing.T) {
	secrets := map[string]string{"HUGE": strings.Repeat("z", maxSecretValueLen+1)}
	err := validateSecrets("secret/big", secrets)
	if err == nil {
		t.Fatal("expected error for oversized value")
	}
}
