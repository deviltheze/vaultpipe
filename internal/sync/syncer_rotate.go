package sync

import (
	"fmt"

	"github.com/yourusername/vaultpipe/internal/dotenv"
)

// rotateEnvFile rotates the current env file and writes fresh secrets.
// It reuses the backup pruning limit from the syncer config when available.
func rotateEnvFile(outputPath string, secrets map[string]string, maxBackups int) error {
	opts := dotenv.RotateOptions{MaxBackups: maxBackups}
	if err := dotenv.Rotate(outputPath, secrets, opts); err != nil {
		return fmt.Errorf("rotate env file: %w", err)
	}
	return nil
}
