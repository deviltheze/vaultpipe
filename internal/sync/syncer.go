package sync

import (
	"fmt"

	"github.com/vaultpipe/internal/config"
	"github.com/vaultpipe/internal/dotenv"
	"github.com/vaultpipe/internal/vault"
)

// Result holds the outcome of a sync operation.
type Result struct {
	SecretsCount int
	OutputFile   string
}

// Syncer orchestrates reading secrets from Vault and writing them to a .env file.
type Syncer struct {
	client *vault.Client
	writer *dotenv.Writer
	cfg    *config.Config
}

// New creates a Syncer from the provided configuration.
func New(cfg *config.Config) (*Syncer, error) {
	client, err := vault.NewClient(cfg.VaultAddress, cfg.VaultToken)
	if err != nil {
		return nil, fmt.Errorf("sync: create vault client: %w", err)
	}

	writer, err := dotenv.NewWriter(cfg.OutputFile)
	if err != nil {
		return nil, fmt.Errorf("sync: create writer: %w", err)
	}

	return &Syncer{client: client, writer: writer, cfg: cfg}, nil
}

// Run performs the sync: reads secrets from Vault and writes them to the output file.
func (s *Syncer) Run() (*Result, error) {
	secrets, err := s.client.ReadSecrets(s.cfg.SecretPath)
	if err != nil {
		return nil, fmt.Errorf("sync: read secrets: %w", err)
	}

	if err := s.writer.Write(secrets); err != nil {
		return nil, fmt.Errorf("sync: write secrets: %w", err)
	}

	return &Result{
		SecretsCount: len(secrets),
		OutputFile:   s.cfg.OutputFile,
	}, nil
}
