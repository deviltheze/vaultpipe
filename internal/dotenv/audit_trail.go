package dotenv

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// AuditEvent represents a single recorded change to a secret.
type AuditEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Key       string    `json:"key"`
	Action    string    `json:"action"` // added, updated, removed, unchanged
	Source    string    `json:"source"` // vault path or origin
	Actor     string    `json:"actor,omitempty"`
}

// AuditTrail holds a sequence of audit events for a sync operation.
type AuditTrail struct {
	SyncedAt time.Time    `json:"synced_at"`
	Output   string       `json:"output"`
	Events   []AuditEvent `json:"events"`
}

// AuditTrailPath returns the conventional path for an audit trail file.
func AuditTrailPath(outputFile string) string {
	dir := filepath.Dir(outputFile)
	base := filepath.Base(outputFile)
	return filepath.Join(dir, "."+base+".audit.json")
}

// WriteAuditTrail serialises an AuditTrail to disk at the conventional path.
func WriteAuditTrail(outputFile string, trail AuditTrail) error {
	if trail.SyncedAt.IsZero() {
		trail.SyncedAt = time.Now().UTC()
	}
	trail.Output = outputFile

	path := AuditTrailPath(outputFile)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("audit_trail: open %s: %w", path, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(trail); err != nil {
		return fmt.Errorf("audit_trail: encode: %w", err)
	}
	return nil
}

// ReadAuditTrail reads an AuditTrail from the conventional path.
func ReadAuditTrail(outputFile string) (AuditTrail, error) {
	path := AuditTrailPath(outputFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return AuditTrail{}, nil
		}
		return AuditTrail{}, fmt.Errorf("audit_trail: read %s: %w", path, err)
	}
	var trail AuditTrail
	if err := json.Unmarshal(data, &trail); err != nil {
		return AuditTrail{}, fmt.Errorf("audit_trail: decode: %w", err)
	}
	return trail, nil
}
