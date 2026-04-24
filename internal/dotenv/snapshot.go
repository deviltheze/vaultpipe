package dotenv

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot captures a point-in-time record of secrets written to a .env file.
type Snapshot struct {
	Timestamp time.Time         `json:"timestamp"`
	Source    string            `json:"source"`
	Secrets   map[string]string `json:"secrets"`
}

// WriteSnapshot serialises a Snapshot as JSON to the given path.
// The file is created with 0600 permissions.
func WriteSnapshot(path string, snap Snapshot) error {
	if snap.Timestamp.IsZero() {
		snap.Timestamp = time.Now().UTC()
	}

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("snapshot: write %s: %w", path, err)
	}
	return nil
}

// ReadSnapshot deserialises a Snapshot from the given JSON file.
// Returns an error wrapping os.ErrNotExist when the file is absent.
func ReadSnapshot(path string) (Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: read %s: %w", path, err)
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: unmarshal: %w", err)
	}
	return snap, nil
}

// SnapshotPath returns the conventional path for a snapshot file
// derived from the .env output path.
func SnapshotPath(envPath string) string {
	return envPath + ".snapshot.json"
}
