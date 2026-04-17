package sync

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/vaultpipe/internal/config"
)

func newMockVault(t *testing.T, data map[string]string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := map[string]interface{}{
			"data": data,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(payload)
	}))
}

func TestRun_WritesSecrets(t *testing.T) {
	srv := newMockVault(t, map[string]string{"KEY": "value"})
	defer srv.Close()

	tmpDir := t.TempDir()
	output := filepath.Join(tmpDir, ".env")

	cfg := &config.Config{
		VaultAddress: srv.URL,
		VaultToken:   "test-token",
		SecretPath:   "secret/app",
		OutputFile:   output,
	}

	syncer, err := New(cfg)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	result, err := syncer.Run()
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	if result.SecretsCount != 1 {
		t.Errorf("expected 1 secret, got %d", result.SecretsCount)
	}
	if result.OutputFile != output {
		t.Errorf("unexpected output file: %s", result.OutputFile)
	}

	if _, err := os.Stat(output); os.IsNotExist(err) {
		t.Error("output file was not created")
	}
}

func TestRun_InvalidVaultAddress(t *testing.T) {
	cfg := &config.Config{
		VaultAddress: "://bad-url",
		VaultToken:   "token",
		SecretPath:   "secret/app",
		OutputFile:   filepath.Join(t.TempDir(), ".env"),
	}

	_, err := New(cfg)
	if err == nil {
		t.Error("expected error for invalid vault address")
	}
}
