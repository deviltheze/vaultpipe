// Package sync orchestrates reading secrets from Vault and writing them
// to a local .env file, with optional audit logging.
package sync

import (
	"fmt"

	"github.com/yourusername/vaultpipe/internal/audit"
	"github.com/yourusername/vaultpipe/internal/config"
	"github.com/yourusername/vaultpipe/internal/dotenv"
	"github.com/yourusername/vaultpipe/internal/vault"
)

// VaultReader is the interface for reading secrets.
type VaultReader interface {
	ReadSecrets(path string) (map[string]string, error)
}

// Syncer coordinates the sync operation.
type Syncer struct {
	cfg    *config.Config
	vault  VaultReader
	writer *dotenv.Writer
	audit  *audit.Logger
}

// New constructs a Syncer from the given config.
func New(cfg *config.Config, auditLog *audit.Logger) (*Syncer, error) {
	v, err := vault.NewClient(cfg.VaultAddress, cfg.VaultToken)
	if err != nil {
		return nil, fmt.Errorf("syncer: init vault client: %w", err)
	}
	w, err := dotenv.NewWriter(cfg.OutputFile)
	if err != nil {
		return nil, fmt.Errorf("syncer: init writer: %w", err)
	}
	if auditLog == nil {
		auditLog, _ = audit.NewLogger("")
	}
	return &Syncer{cfg: cfg, vault: v, writer: w, audit: auditLog}, nil
}

// Run performs the secret sync.
func (s *Syncer) Run() error {
	secrets, err := s.vault.ReadSecrets(s.cfg.SecretPath)
	if err != nil {
		_ = s.audit.LogError(s.cfg.SecretPath, err)
		return fmt.Errorf("syncer: read secrets: %w", err)
	}
	if err := s.writer.Write(secrets); err != nil {
		_ = s.audit.LogError(s.cfg.SecretPath, err)
		return fmt.Errorf("syncer: write env: %w", err)
	}
	_ = s.audit.LogSync(s.cfg.SecretPath, s.cfg.OutputFile, len(secrets))
	return nil
}
