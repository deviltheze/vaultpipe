package dotenv

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// LineageRecord captures metadata about a single sync event for auditing purposes.
type LineageRecord struct {
	Timestamp time.Time         `json:"timestamp"`
	Source    string            `json:"source"`
	OutputFile string          `json:"output_file"`
	Added     int               `json:"added"`
	Updated   int               `json:"updated"`
	Removed   int               `json:"removed"`
	Checksum  string            `json:"checksum"`
	Meta      map[string]string `json:"meta,omitempty"`
}

// WriteLineage appends a LineageRecord as a JSON line to the given file path.
func WriteLineage(path string, rec LineageRecord) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("lineage: open %s: %w", path, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	if err := enc.Encode(rec); err != nil {
		return fmt.Errorf("lineage: encode record: %w", err)
	}
	return nil
}

// ReadLineage reads all LineageRecords from the given file path.
func ReadLineage(path string) ([]LineageRecord, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("lineage: open %s: %w", path, err)
	}
	defer f.Close()

	var records []LineageRecord
	dec := json.NewDecoder(f)
	for dec.More() {
		var rec LineageRecord
		if err := dec.Decode(&rec); err != nil {
			return nil, fmt.Errorf("lineage: decode record: %w", err)
		}
		records = append(records, rec)
	}
	return records, nil
}
