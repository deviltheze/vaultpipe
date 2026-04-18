package dotenv

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateTemplate_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, ".env.template")

	secrets := map[string]string{
		"APP_NAME": "myapp",
		"PORT":     "8080",
		"DEBUG":    "true",
		"EMPTY":    "",
	}

	if err := GenerateTemplate(out, secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("could not read file: %v", err)
	}
	content := string(data)

	if !strings.Contains(content, "APP_NAME=<string>") {
		t.Errorf("expected APP_NAME=<string>, got:\n%s", content)
	}
	if !strings.Contains(content, "PORT=<number>") {
		t.Errorf("expected PORT=<number>, got:\n%s", content)
	}
	if !strings.Contains(content, "DEBUG=<bool>") {
		t.Errorf("expected DEBUG=<bool>, got:\n%s", content)
	}
	if !strings.Contains(content, "EMPTY=\n") {
		t.Errorf("expected EMPTY= (empty), got:\n%s", content)
	}
}

func TestGenerateTemplate_SortedKeys(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, ".env.template")

	secrets := map[string]string{"Z_KEY": "z", "A_KEY": "a", "M_KEY": "m"}
	if err := GenerateTemplate(out, secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(out)
	lines := strings.Split(string(data), "\n")
	var keys []string
	for _, l := range lines {
		if strings.Contains(l, "=") && !strings.HasPrefix(l, "#") {
			keys = append(keys, strings.SplitN(l, "=", 2)[0])
		}
	}
	expected := []string{"A_KEY", "M_KEY", "Z_KEY"}
	for i, k := range expected {
		if keys[i] != k {
			t.Errorf("expected key[%d]=%s, got %s", i, k, keys[i])
		}
	}
}

func TestGenerateTemplate_HasHeader(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, ".env.template")
	if err := GenerateTemplate(out, map[string]string{"KEY": "val"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(out)
	if !strings.HasPrefix(string(data), "# Auto-generated") {
		t.Errorf("expected header comment at top")
	}
}

func TestValueHint(t *testing.T) {
	cases := []struct{ in, want string }{
		{"", ""},
		{"123", "<number>"},
		{"true", "<bool>"},
		{"FALSE", "<bool>"},
		{"hello", "<string>"},
	}
	for _, c := range cases {
		got := valueHint(c.in)
		if got != c.want {
			t.Errorf("valueHint(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}
