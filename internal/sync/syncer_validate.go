package sync

import (
	"fmt"

	"github.com/yourusername/vaultpipe/internal/dotenv"
)

// maxSecretValueLen is the upper bound on a single secret value in bytes.
const maxSecretValueLen = 8192

// validateSecrets runs dotenv.Validate against the fetched secrets and wraps
// any error with context about the Vault path.
func validateSecrets(path string, secrets map[string]string) error {
	if err := dotenv.Validate(secrets, maxSecretValueLen); err != nil {
		return fmt.Errorf("secrets from %q failed validation: %w", path, err)
	}
	return nil
}
