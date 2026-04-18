package sync_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpipe/internal/audit"
	"github.com/yourusername/vaultpipe/internal/config"
	"github.com/yourusername/vaultpipe/internal/sync"
)

type mockVault struct {
	secrets map[string]string
	err     error
}

func (m *mockVault) ReadSecrets(_ string) (map[string]string, error) {
	return m.secrets, m.err
}

func newTestConfig(t *testing.T) *config.Config {
	t.Helper()
	return &config.Config{
		VaultAddress: "http://127.0.0.1:8200",
		VaultToken:   "root",
		SecretPath:   "secret/myapp",
		OutputFile:   filepath.Join(t.TempDir(), ".env"),
	}
}

func TestRun_WritesSecrets(t *testing.T) {
	cfg := newTestConfig(t)
	auditLog, _ := audit.NewLogger(filepath.Join(t.TempDir(), "audit.log"))

	s, err := sync.New(cfg, auditLog)
	if err != nil {
		t.Fatalf("New error: %v", err)
	}
	s.SetVault(&mockVault{secrets: map[string]string{"KEY": "value"}})

	if err := s.Run(); err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if _, err := os.Stat(cfg.OutputFile); err != nil {
		t.Errorf("output file not created: %v", err)
	}
}

func TestRun_VaultError_LogsAndReturns(t *testing.T) {
	cfg := newTestConfig(t)
	auditPath := filepath.Join(t.TempDir(), "audit.log")
	auditLog, _ := audit.NewLogger(auditPath)

	s, err := sync.New(cfg, auditLog)
	if err != nil {
		t.Fatalf("New error: %v", err)
	}
	s.SetVault(&mockVault{err: errors.New("vault down")})

	if err := s.Run(); err == nil {
		t.Fatal("expected error from Run")
	}
	info, _ := os.Stat(auditPath)
	if info == nil || info.Size() == 0 {
		t.Error("expected audit log entry on error")
	}
}

func TestRun_InvalidVaultAddress(t *testing.T) {
	cfg := newTestConfig(t)
	cfg.VaultAddress = "://bad"
	_, err := sync.New(cfg, nil)
	if err == nil {
		t.Fatal("expected error for invalid vault address")
	}
}
