package dotenv

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// PinRecord records a pinned version of a secret key with an optional expiry.
type PinRecord struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	PinnedAt  time.Time `json:"pinned_at"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	Reason    string    `json:"reason,omitempty"`
}

// IsExpired reports whether the pin has passed its expiry time.
// A zero ExpiresAt means the pin never expires.
func (p PinRecord) IsExpired() bool {
	if p.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().After(p.ExpiresAt)
}

// PinPath returns the conventional path for a pin file alongside the env file.
func PinPath(envPath string) string {
	dir := filepath.Dir(envPath)
	base := filepath.Base(envPath)
	return filepath.Join(dir, "."+base+".pins.json")
}

// WritePinFile persists a slice of PinRecords to disk as JSON.
func WritePinFile(path string, pins []PinRecord) error {
	data, err := json.MarshalIndent(pins, "", "  ")
	if err != nil {
		return fmt.Errorf("pin: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("pin: write %s: %w", path, err)
	}
	return nil
}

// ReadPinFile loads PinRecords from disk. Returns an empty slice when the
// file does not exist.
func ReadPinFile(path string) ([]PinRecord, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return []PinRecord{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("pin: read %s: %w", path, err)
	}
	var pins []PinRecord
	if err := json.Unmarshal(data, &pins); err != nil {
		return nil, fmt.Errorf("pin: unmarshal: %w", err)
	}
	return pins, nil
}

// ApplyPins overrides values in secrets with any non-expired pinned values.
// It returns the number of pins that were applied.
func ApplyPins(secrets map[string]string, pins []PinRecord) (map[string]string, int) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}
	applied := 0
	for _, p := range pins {
		if p.IsExpired() {
			continue
		}
		out[p.Key] = p.Value
		applied++
	}
	return out, applied
}
