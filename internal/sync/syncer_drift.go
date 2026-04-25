package sync

import (
	"fmt"
	"os"

	"github.com/your-org/vaultpipe/internal/dotenv"
)

// checkDrift reads the current .env file (if it exists), compares it against
// the freshly fetched vault secrets, and prints a drift report to stderr.
// It returns true when drift is detected so the caller can decide to abort or
// continue based on configuration.
func checkDrift(outputFile string, vaultSecrets map[string]string) (bool, error) {
	_, err := os.Stat(outputFile)
	if os.IsNotExist(err) {
		// No existing file — nothing to compare against.
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("drift check: stat %s: %w", outputFile, err)
	}

	existing, err := dotenv.Read(outputFile)
	if err != nil {
		return false, fmt.Errorf("drift check: read %s: %w", outputFile, err)
	}

	report := dotenv.DetectDrift(existing, vaultSecrets)
	if !report.Clean {
		fmt.Fprintln(os.Stderr, "[vaultpipe] "+report.String())
		return true, nil
	}
	return false, nil
}
