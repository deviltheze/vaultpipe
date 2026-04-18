// Package sync orchestrates reading secrets from Vault and writing them
// to a local .env file, optionally merging with existing values.
package sync

import (
	"fmt"
	"log"

	"github.com/your-org/vaultpipe/internal/audit"
	"github.com/your-org/vaultpipe/internal/config"
	"github.com/your-org/vaultpipe/internal/dotenv"
	"github.com/your-org/vaultpipe/internal/vault"
)

// Syncer holds dependencies for a sync run.
type Syncer struct {
	cfg    *config.Config
	logger *audit.Logger
}

// New creates a Syncer with the given config and audit logger.
func New(cfg *config.Config, logger *audit.Logger) *Syncer {
	return &Syncer{cfg: cfg, logger: logger}
}

// Run performs the full sync: fetch secrets, diff, merge, write.
func (s *Syncer) Run() error {
	client, err := vault.NewClient(s.cfg.VaultAddress, s.cfg.VaultToken)
	if err != nil {
		s.logger.LogError("vault_client", err)
		return fmt.Errorf("vault client: %w", err)
	}

	incoming, err := client.ReadSecrets(s.cfg.SecretPath)
	if err != nil {
		s.logger.LogError("read_secrets", err)
		return fmt.Errorf("read secrets: %w", err)
	}

	existing := map[string]string{}
	if s.cfg.OutputFile != "" {
		if ex, readErr := dotenv.Read(s.cfg.OutputFile); readErr == nil {
			existing = ex
		}
	}

	entries := dotenv.Diff(existing, incoming)
	summary := dotenv.Summary(entries)
	log.Printf("diff summary: added=%d updated=%d removed=%d unchanged=%d",
		summary[dotenv.ChangeAdded],
		summary[dotenv.ChangeUpdated],
		summary[dotenv.ChangeRemoved],
		summary[dotenv.ChangeUnchanged],
	)

	merged := dotenv.Merge(existing, incoming, s.cfg.OverwriteExisting)

	w := dotenv.NewWriter(s.cfg.OutputFile)
	if err := w.Write(merged); err != nil {
		s.logger.LogError("write_env", err)
		return fmt.Errorf("write env: %w", err)
	}

	s.logger.LogSync(s.cfg.SecretPath, s.cfg.OutputFile, len(merged))
	return nil
}
