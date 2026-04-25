package dotenv

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// ExpiryRecord holds expiry metadata for a synced env file.
type ExpiryRecord struct {
	Path      string    `json:"path"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// IsExpired reports whether the record has passed its expiry time.
// A zero ExpiresAt is treated as never-expiring.
func (r ExpiryRecord) IsExpired() bool {
	if r.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().After(r.ExpiresAt)
}

// ExpiryPath returns the conventional path for an expiry file
// adjacent to the given env file.
func ExpiryPath(envPath string) string {
	return envPath + ".expiry"
}

// WriteExpiryFile serialises an ExpiryRecord to disk next to envPath.
func WriteExpiryFile(envPath string, ttl time.Duration) error {
	now := time.Now().UTC()
	rec := ExpiryRecord{
		Path:      envPath,
		CreatedAt: now,
	}
	if ttl > 0 {
		rec.ExpiresAt = now.Add(ttl)
	}
	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return fmt.Errorf("expire: marshal: %w", err)
	}
	return os.WriteFile(ExpiryPath(envPath), data, 0o600)
}

// ReadExpiryFile deserialises the expiry record for envPath.
// Returns a zero-value record (never-expiring) when the file is absent.
func ReadExpiryFile(envPath string) (ExpiryRecord, error) {
	data, err := os.ReadFile(ExpiryPath(envPath))
	if os.IsNotExist(err) {
		return ExpiryRecord{Path: envPath}, nil
	}
	if err != nil {
		return ExpiryRecord{}, fmt.Errorf("expire: read: %w", err)
	}
	var rec ExpiryRecord
	if err := json.Unmarshal(data, &rec); err != nil {
		return ExpiryRecord{}, fmt.Errorf("expire: unmarshal: %w", err)
	}
	return rec, nil
}
