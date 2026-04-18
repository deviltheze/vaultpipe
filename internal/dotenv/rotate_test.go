package dotenv

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRotate_RenamesOldFile(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	// Create an initial env file.
	if err := os.WriteFile(envPath, []byte("OLD=1\n"), 0600); err != nil {
		t.Fatal(err)
	}

	secrets := map[string]string{"NEW": "2"}
	if err := Rotate(envPath, secrets, RotateOptions{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Original path must contain new content.
	data, err := os.ReadFile(envPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "NEW") {
		t.Errorf("expected new secrets in file, got: %s", data)
	}

	// A rotated backup must exist.
	matches, _ := filepath.Glob(filepath.Join(dir, ".*.env"))
	if len(matches) == 0 {
		// try alternate naming pattern
		matches, _ = filepath.Glob(filepath.Join(dir, ".*"))
	}
	files, _ := os.ReadDir(dir)
	backupFound := false
	for _, f := range files {
		if f.Name() != ".env" {
			backupFound = true
		}
	}
	if !backupFound {
		t.Error("expected a rotated backup file")
	}
}

func TestRotate_NoExistingFile(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	secrets := map[string]string{"KEY": "val"}
	if err := Rotate(envPath, secrets, RotateOptions{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(envPath); err != nil {
		t.Errorf("expected env file to be created: %v", err)
	}
}

func TestRotate_PrunesBackups(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	// Simulate two existing backups.
	for _, name := range []string{".env.20230101T000000Z", ".env.20230102T000000Z"} {
		os.WriteFile(filepath.Join(dir, name), []byte("X=1"), 0600)
	}
	os.WriteFile(envPath, []byte("A=1"), 0600)

	if err := Rotate(envPath, map[string]string{"B": "2"}, RotateOptions{MaxBackups: 1}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	files, _ := os.ReadDir(dir)
	backups := 0
	for _, f := range files {
		if f.Name() != ".env" {
			backups++
		}
	}
	if backups > 1 {
		t.Errorf("expected at most 1 backup, got %d", backups)
	}
}
