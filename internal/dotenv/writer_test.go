package dotenv

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWrite_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, ".env")

	w := NewWriter(out)
	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}

	if err := w.Write(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("could not read output file: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "DB_HOST=localhost") {
		t.Errorf("expected DB_HOST=localhost in output, got:\n%s", content)
	}
	if !strings.Contains(content, "DB_PORT=5432") {
		t.Errorf("expected DB_PORT=5432 in output, got:\n%s", content)
	}
}

func TestWrite_SortedKeys(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, ".env")

	w := NewWriter(out)
	secrets := map[string]string{
		"Z_KEY": "z",
		"A_KEY": "a",
		"M_KEY": "m",
	}

	if err := w.Write(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(out)
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if lines[0] != "A_KEY=a" || lines[1] != "M_KEY=m" || lines[2] != "Z_KEY=z" {
		t.Errorf("keys not sorted, got: %v", lines)
	}
}

func TestWrite_QuotesSpecialValues(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, ".env")

	w := NewWriter(out)
	secrets := map[string]string{
		"MSG": "hello world",
	}

	if err := w.Write(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(out)
	if !strings.Contains(string(data), `"hello world"`) {
		t.Errorf("expected quoted value, got: %s", string(data))
	}
}

func TestWrite_EmptySecrets(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, ".env")

	w := NewWriter(out)
	err := w.Write(map[string]string{})
	if err == nil {
		t.Error("expected error for empty secrets, got nil")
	}
}

func TestWrite_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, ".env")

	w := NewWriter(out)
	_ = w.Write(map[string]string{"KEY": "val"})

	info, err := os.Stat(out)
	if err != nil {
		t.Fatalf("stat failed: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected mode 0600, got %v", info.Mode().Perm())
	}
}
