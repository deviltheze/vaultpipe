package dotenv

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// TTLRecord stores the sync timestamp and expiry duration for a secret path.
type TTLRecord struct {
	Path      string        `json:"path"`
	SyncedAt  time.Time     `json:"synced_at"`
	TTL       time.Duration `json:"ttl_ns"`
}

// IsExpired returns true if the TTL has elapsed since SyncedAt.
func (r TTLRecord) IsExpired() bool {
	if r.TTL == 0 {
		return false
	}
	return time.Since(r.SyncedAt) > r.TTL
}

// WriteTTLFile writes a TTLRecord as JSON to path.
func WriteTTLFile(path string, rec TTLRecord) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("ttl: open %s: %w", path, err)
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(rec)
}

// ReadTTLFile reads a TTLRecord from a JSON file at path.
func ReadTTLFile(path string) (TTLRecord, error) {
	f, err := os.Open(path)
	if err != nil {
		return TTLRecord{}, fmt.Errorf("ttl: open %s: %w", path, err)
	}
	defer f.Close()
	var rec TTLRecord
	if err := json.NewDecoder(f).Decode(&rec); err != nil {
		return TTLRecord{}, fmt.Errorf("ttl: decode: %w", err)
	}
	return rec, nil
}
