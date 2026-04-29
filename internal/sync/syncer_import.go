package sync

import (
	"fmt"
	"log/slog"

	"github.com/yourusername/vaultpipe/internal/dotenv"
)

// writeImportRecord persists an import record for the current sync operation.
// It records the vault path (source), the resolved secrets, and the timestamp.
// Errors are logged as warnings but do not abort the sync.
func writeImportRecord(logger *slog.Logger, outputFile, vaultPath string, secrets map[string]string) {
	rec := dotenv.ImportRecord{
		Source:  vaultPath,
		Secrets: secrets,
	}
	if err := dotenv.WriteImportFile(outputFile, rec); err != nil {
		logger.Warn("import record: failed to write", "error", err, "output", outputFile)
		return
	}
	logger.Debug("import record written",
		"source", vaultPath,
		"output", outputFile,
		"keys", len(secrets),
	)
}

// verifyImportSource checks whether the import record for the output file
// matches the expected vault path. Returns an error if there is a mismatch.
// A missing import file is not treated as an error (first-time sync).
func verifyImportSource(outputFile, expectedSource string) error {
	rec, err := dotenv.ReadImportFile(outputFile)
	if err != nil {
		return fmt.Errorf("import verify: %w", err)
	}
	if rec.Source == "" {
		// No prior import record — nothing to verify.
		return nil
	}
	if rec.Source != expectedSource {
		return fmt.Errorf(
			"import source mismatch: file was last imported from %q, current source is %q",
			rec.Source, expectedSource,
		)
	}
	return nil
}
