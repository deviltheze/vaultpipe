package dotenv

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// HookEvent describes when a hook fires.
type HookEvent string

const (
	HookPreSync  HookEvent = "pre_sync"
	HookPostSync HookEvent = "post_sync"
)

// HookRecord represents a single hook invocation record.
type HookRecord struct {
	Event     HookEvent `json:"event"`
	Timestamp time.Time `json:"timestamp"`
	Actor     string    `json:"actor"`
	Output    string    `json:"output,omitempty"`
	Error     string    `json:"error,omitempty"`
}

// HookFile is the collection of hook records persisted to disk.
type HookFile struct {
	Records []HookRecord `json:"records"`
}

// HookPath returns the canonical path for the hook log file.
func HookPath(envFile string) string {
	dir := filepath.Dir(envFile)
	base := filepath.Base(envFile)
	return filepath.Join(dir, "."+base+".hooks.json")
}

// WriteHookRecord appends a hook record to the hook log file.
func WriteHookRecord(path string, rec HookRecord) error {
	if rec.Timestamp.IsZero() {
		rec.Timestamp = time.Now().UTC()
	}

	hf, err := ReadHookFile(path)
	if err != nil {
		return fmt.Errorf("hook: read existing: %w", err)
	}

	hf.Records = append(hf.Records, rec)

	data, err := json.MarshalIndent(hf, "", "  ")
	if err != nil {
		return fmt.Errorf("hook: marshal: %w", err)
	}

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("hook: write: %w", err)
	}
	return nil
}

// ReadHookFile reads the hook log from disk. Returns an empty HookFile if missing.
func ReadHookFile(path string) (HookFile, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return HookFile{}, nil
	}
	if err != nil {
		return HookFile{}, fmt.Errorf("hook: read: %w", err)
	}
	var hf HookFile
	if err := json.Unmarshal(data, &hf); err != nil {
		return HookFile{}, fmt.Errorf("hook: unmarshal: %w", err)
	}
	return hf, nil
}
