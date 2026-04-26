package dotenv

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ScopeRecord describes a named scope applied to a set of secrets.
type ScopeRecord struct {
	Name      string    `json:"name"`
	Keys      []string  `json:"keys"`
	CreatedAt time.Time `json:"created_at"`
}

// ScopePath returns the conventional path for a scope file.
func ScopePath(outputFile string) string {
	dir := filepath.Dir(outputFile)
	base := filepath.Base(outputFile)
	return filepath.Join(dir, "."+base+".scope.json")
}

// WriteScopeFile persists a ScopeRecord to disk as JSON.
func WriteScopeFile(path string, record ScopeRecord) error {
	if record.CreatedAt.IsZero() {
		record.CreatedAt = time.Now().UTC()
	}
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return fmt.Errorf("scope: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("scope: write %s: %w", path, err)
	}
	return nil
}

// ReadScopeFile loads a ScopeRecord from disk.
// Returns an empty record (no error) when the file does not exist.
func ReadScopeFile(path string) (ScopeRecord, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return ScopeRecord{}, nil
	}
	if err != nil {
		return ScopeRecord{}, fmt.Errorf("scope: read %s: %w", path, err)
	}
	var record ScopeRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return ScopeRecord{}, fmt.Errorf("scope: unmarshal: %w", err)
	}
	return record, nil
}

// ApplyScope filters secrets to only those whose keys are listed in the
// ScopeRecord. If the record has no keys, the original map is returned
// unchanged.
func ApplyScope(secrets map[string]string, record ScopeRecord) map[string]string {
	if len(record.Keys) == 0 {
		return secrets
	}
	allowed := make(map[string]struct{}, len(record.Keys))
	for _, k := range record.Keys {
		allowed[k] = struct{}{}
	}
	out := make(map[string]string)
	for k, v := range secrets {
		if _, ok := allowed[k]; ok {
			out[k] = v
		}
	}
	return out
}
