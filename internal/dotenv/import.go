package dotenv

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ImportRecord captures metadata about an imported secret set.
type ImportRecord struct {
	Source    string            `json:"source"`
	ImportedAt time.Time        `json:"imported_at"`
	Keys      []string          `json:"keys"`
	Secrets   map[string]string `json:"secrets"`
}

// ImportPath returns the conventional path for an import record file.
func ImportPath(outputFile string) string {
	dir := filepath.Dir(outputFile)
	base := filepath.Base(outputFile)
	return filepath.Join(dir, "."+base+".import.json")
}

// WriteImportFile persists an ImportRecord as JSON next to the output file.
func WriteImportFile(outputFile string, record ImportRecord) error {
	if record.ImportedAt.IsZero() {
		record.ImportedAt = time.Now().UTC()
	}
	if record.Keys == nil {
		record.Keys = sortedKeyList(record.Secrets)
	}
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return fmt.Errorf("import: marshal: %w", err)
	}
	path := ImportPath(outputFile)
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("import: write %s: %w", path, err)
	}
	return nil
}

// ReadImportFile loads the most recent ImportRecord for a given output file.
func ReadImportFile(outputFile string) (ImportRecord, error) {
	path := ImportPath(outputFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return ImportRecord{}, nil
		}
		return ImportRecord{}, fmt.Errorf("import: read %s: %w", path, err)
	}
	var record ImportRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return ImportRecord{}, fmt.Errorf("import: unmarshal: %w", err)
	}
	return record, nil
}

// sortedKeyList returns a sorted slice of keys from a map.
func sortedKeyList(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sortStrings(keys)
	return keys
}
