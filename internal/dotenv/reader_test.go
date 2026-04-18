package dotenv

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestRead_BasicPairs(t *testing.T) {
	p := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	m, err := Read(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["FOO"] != "bar" || m["BAZ"] != "qux" {
		t.Errorf("unexpected map: %v", m)
	}
}

func TestRead_IgnoresComments(t *testing.T) {
	p := writeTempEnv(t, "# comment\nKEY=value\n")
	m, err := Read(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m) != 1 || m["KEY"] != "value" {
		t.Errorf("unexpected map: %v", m)
	}
}

func TestRead_UnquotesDoubleQuotes(t *testing.T) {
	p := writeTempEnv(t, `SECRET="hello world"` + "\n")
	m, err := Read(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["SECRET"] != "hello world" {
		t.Errorf("got %q", m["SECRET"])
	}
}

func TestRead_UnquotesSingleQuotes(t *testing.T) {
	p := writeTempEnv(t, "TOKEN='abc123'\n")
	m, err := Read(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["TOKEN"] != "abc123" {
		t.Errorf("got %q", m["TOKEN"])
	}
}

func TestRead_InvalidSyntax(t *testing.T) {
	p := writeTempEnv(t, "NODEQUALS\n")
	_, err := Read(p)
	if err == nil {
		t.Fatal("expected error for invalid syntax")
	}
}

func TestRead_FileNotFound(t *testing.T) {
	_, err := Read("/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRead_EmptyFile(t *testing.T) {
	p := writeTempEnv(t, "")
	m, err := Read(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m) != 0 {
		t.Errorf("expected empty map, got %v", m)
	}
}
