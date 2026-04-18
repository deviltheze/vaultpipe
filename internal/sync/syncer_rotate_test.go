package sync

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRotateEnvFile_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	secrets := map[string]string{"FOO": "bar"}
	if err := rotateEnvFile(envPath, secrets, 0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(envPath); err != nil {
		t.Errorf("expected env file to exist: %v", err)
	}
}

func TestRotateEnvFile_RotatesExisting(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	os.WriteFile(envPath, []byte("OLD=yes\n"), 0600)

	if err := rotateEnvFile(envPath, map[string]string{"NEW": "yes"}, 2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	files, _ := os.ReadDir(dir)
	if len(files) < 2 {
		t.Errorf("expected original + rotated backup, got %d file(s)", len(files))
	}
}

func TestRotateEnvFile_BadPath(t *testing.T) {
	err := rotateEnvFile("/nonexistent/path/.env", map[string]string{"K": "v"}, 0)
	if err == nil {
		t.Error("expected error for bad path")
	}
}
