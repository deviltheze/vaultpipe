package sync

import (
	"fmt"
	"log/slog"

	"github.com/your-org/vaultpipe/internal/config"
	"github.com/your-org/vaultpipe/internal/dotenv"
)

// writeTagFile persists a TagRecord alongside the output env file.
// Tags are sourced from the config and supplemented with sync metadata.
func writeTagFile(cfg *config.Config, logger *slog.Logger) error {
	path := dotenv.TagPath(cfg.OutputFile)

	rec := dotenv.TagRecord{
		Environment: cfg.Environment,
		VaultPath:   cfg.VaultPath,
		Tags:        cfg.Tags,
	}

	if err := dotenv.WriteTagFile(path, rec); err != nil {
		return fmt.Errorf("syncer: write tag file: %w", err)
	}

	logger.Info("tag file written", "path", path, "environment", cfg.Environment)
	return nil
}

// logExistingTags reads and logs any previously written tags for the output file.
func logExistingTags(cfg *config.Config, logger *slog.Logger) {
	path := dotenv.TagPath(cfg.OutputFile)
	rec, err := dotenv.ReadTagFile(path)
	if err != nil {
		logger.Warn("could not read tag file", "path", path, "error", err)
		return
	}
	if rec.Environment == "" {
		return // no prior tag record
	}
	logger.Info("existing tag record found",
		"environment", rec.Environment,
		"vault_path", rec.VaultPath,
		"synced_at", rec.SyncedAt,
	)
}
