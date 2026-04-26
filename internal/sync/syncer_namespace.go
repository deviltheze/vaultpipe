package sync

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/your-org/vaultpipe/internal/dotenv"
)

// applyNamespace prefixes all secret keys with the configured namespace and
// writes a namespace record alongside the output file for auditing purposes.
// If no namespace is configured the secrets are returned unchanged.
func applyNamespace(secrets map[string]string, namespace, outputFile string) (map[string]string, error) {
	if namespace == "" {
		return secrets, nil
	}

	prefixed := dotenv.ApplyNamespace(secrets, namespace)

	keys := make([]string, 0, len(prefixed))
	for k := range prefixed {
		keys = append(keys, k)
	}

	rec := dotenv.NamespaceRecord{
		Namespace: namespace,
		AppliedAt: time.Now().UTC(),
		Keys:      keys,
	}

	dir := filepath.Dir(outputFile)
	if err := dotenv.WriteNamespaceFile(dir, namespace, rec); err != nil {
		return nil, fmt.Errorf("namespace: write record: %w", err)
	}

	return prefixed, nil
}

// stripNamespace removes a namespace prefix from all matching secret keys.
// Keys that do not carry the prefix are silently dropped. If namespace is
// empty the original map is returned as-is.
func stripNamespace(secrets map[string]string, namespace string) map[string]string {
	if namespace == "" {
		return secrets
	}
	return dotenv.StripNamespace(secrets, namespace)
}
