package sync

import (
	"os"
	"path/filepath"
	"testing"
)

func TestApplyNamespace_NoNamespace_ReturnsOriginal(t *testing.T) {
	secrets := map[string]string{"FOO": "bar"}
	result, err := applyNamespace(secrets, "", "/tmp/.env")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %v", result)
	}
}

func TestApplyNamespace_PrefixesKeys(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, ".env")
	secrets := map[string]string{"HOST": "localhost", "PORT": "5432"}

	result, err := applyNamespace(secrets, "DB", output)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %v", result)
	}
	if result["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %v", result)
	}
	if _, ok := result["HOST"]; ok {
		t.Error("original key HOST should be removed after namespace apply")
	}
}

func TestApplyNamespace_WritesRecordFile(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, ".env")
	secrets := map[string]string{"KEY": "value"}

	_, err := applyNamespace(secrets, "SVC", output)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	recordPath := filepath.Join(dir, ".vaultpipe.namespace.SVC.json")
	if _, err := os.Stat(recordPath); os.IsNotExist(err) {
		t.Errorf("expected namespace record at %s, not found", recordPath)
	}
}

func TestStripNamespace_NoNamespace_ReturnsOriginal(t *testing.T) {
	secrets := map[string]string{"FOO": "bar"}
	result := stripNamespace(secrets, "")
	if result["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %v", result)
	}
}

func TestStripNamespace_RemovesPrefixedKeys(t *testing.T) {
	secrets := map[string]string{"APP_FOO": "one", "APP_BAR": "two", "UNRELATED": "three"}
	result := stripNamespace(secrets, "APP")
	if result["FOO"] != "one" {
		t.Errorf("expected FOO=one, got %v", result)
	}
	if _, ok := result["UNRELATED"]; ok {
		t.Error("UNRELATED should be excluded after strip")
	}
}
