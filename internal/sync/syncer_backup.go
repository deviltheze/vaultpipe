package sync

import (
	"fmt"
	"log/slog"

	"github.com/yourusername/vaultpipe/internal/dotenv"
)

const defaultKeepBackups = 5

// backupAndPrune creates a timestamped backup of outputFile (if it exists)
// and prunes old backups, keeping at most defaultKeepBackups.
// Errors are logged but do not abort the sync.
func backupAndPrune(outputFile string, logger *slog.Logger) {
	dst, err := dotenv.Backup(outputFile)
	if err != nil {
		logger.Warn("backup failed", "file", outputFile, "error", err)
		return
	}
	if dst == "" {
		return // file did not exist yet
	}
	logger.Info("backup created", "backup", dst)

	if err := dotenv.PruneBackups(outputFile, defaultKeepBackups); err != nil {
		logger.Warn("prune failed", "file", outputFile, "error", fmt.Sprintf("%v", err))
	}
}
