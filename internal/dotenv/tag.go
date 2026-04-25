package dotenv

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// TagRecord holds metadata tags associated with a sync operation.
type TagRecord struct {
	Environment string            `json:"environment"`
	VaultPath   string            `json:"vault_path"`
	SyncedAt    time.Time         `json:"synced_at"`
	Tags        map[string]string `json:"tags,omitempty"`
}

// WriteTagFile writes a TagRecord as JSON to path.
func WriteTagFile(path string, rec TagRecord) error {
	if rec.SyncedAt.IsZero() {
		rec.SyncedAt = time.Now().UTC()
	}
	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return fmt.Errorf("tag: marshal: %w", err)
	}
	return os.WriteFile(path, append(data, '\n'), 0o600)
}

// ReadTagFile reads a TagRecord from path.
// Returns a zero-value TagRecord and no error if the file does not exist.
func ReadTagFile(path string) (TagRecord, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return TagRecord{}, nil
		}
		return TagRecord{}, fmt.Errorf("tag: read: %w", err)
	}
	var rec TagRecord
	if err := json.Unmarshal(data, &rec); err != nil {
		return TagRecord{}, fmt.Errorf("tag: unmarshal: %w", err)
	}
	return rec, nil
}

// TagPath returns the conventional tag file path for a given env file.
func TagPath(envFile string) string {
	return envFile + ".tags.json"
}
