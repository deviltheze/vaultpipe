package sync

import (
	"fmt"
	"log/slog"

	"github.com/yourusername/vaultpipe/internal/dotenv"
)

// applyAliases reads the alias file for the configured output path and injects
// alias keys into the secrets map. It is a no-op when no alias file exists.
func applyAliases(outputPath string, secrets map[string]string, logger *slog.Logger) (map[string]string, error) {
	records, err := dotenv.ReadAliasFile(outputPath)
	if err != nil {
		return nil, fmt.Errorf("alias: read file: %w", err)
	}
	if len(records) == 0 {
		return secrets, nil
	}
	result := dotenv.ApplyAliases(secrets, records)
	added := len(result) - len(secrets)
	if added > 0 && logger != nil {
		logger.Info("aliases applied", "count", added, "output", outputPath)
	}
	return result, nil
}

// writeAliasFile persists alias records for the given output path.
func writeAliasFile(outputPath string, records []dotenv.AliasRecord, logger *slog.Logger) error {
	if len(records) == 0 {
		return nil
	}
	if err := dotenv.WriteAliasFile(outputPath, records); err != nil {
		return fmt.Errorf("alias: write file: %w", err)
	}
	if logger != nil {
		logger.Info("alias file written", "path", dotenv.AliasPath(outputPath), "count", len(records))
	}
	return nil
}
