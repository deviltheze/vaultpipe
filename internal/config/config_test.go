package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "vaultpipe.yaml")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}
	return p
}

func TestLoad_Valid(t *testing.T) {
	content := `
vault:
  address: "http://127.0.0.1:8200"
  token: "root"
  secrets:
    - path: "secret/data/app"
      key: "password"
      env: "APP_PASSWORD"
output:
  file: ".env"
`
	p := writeTemp(t, content)
	cfg, err := Load(p)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if cfg.Vault.Address != "http://127.0.0.1:8200" {
		t.Errorf("unexpected address: %s", cfg.Vault.Address)
	}
	if len(cfg.Vault.Secrets) != 1 {
		t.Errorf("expected 1 secret, got %d", len(cfg.Vault.Secrets))
	}
	if cfg.Output.File != ".env" {
		t.Errorf("unexpected output file: %s", cfg.Output.File)
	}
}

func TestLoad_DefaultOutputFile(t *testing.T) {
	content := `
vault:
  address: "http://127.0.0.1:8200"
  secrets:
    - path: "secret/data/app"
      env: "APP_SECRET"
`
	p := writeTemp(t, content)
	cfg, err := Load(p)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if cfg.Output.File != ".env" {
		t.Errorf("expected default output file '.env', got %s", cfg.Output.File)
	}
}

func TestLoad_MissingAddress(t *testing.T) {
	content := `
vault:
  secrets:
    - path: "secret/data/app"
      env: "APP_SECRET"
`
	p := writeTemp(t, content)
	_, err := Load(p)
	if err == nil {
		t.Fatal("expected error for missing vault address")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path/vaultpipe.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
