package dotenv

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// AliasRecord maps a friendly alias name to a canonical secret key.
type AliasRecord struct {
	Alias     string    `json:"alias"`
	Canonical string    `json:"canonical"`
	CreatedAt time.Time `json:"created_at"`
}

// AliasFile holds all alias records for an env file.
type AliasFile struct {
	Aliases []AliasRecord `json:"aliases"`
}

// AliasPath returns the conventional path for an alias file.
func AliasPath(envPath string) string {
	dir := filepath.Dir(envPath)
	base := filepath.Base(envPath)
	return filepath.Join(dir, "."+base+".aliases.json")
}

// WriteAliasFile persists alias records alongside the env file.
func WriteAliasFile(envPath string, aliases []AliasRecord) error {
	for i := range aliases {
		if aliases[i].CreatedAt.IsZero() {
			aliases[i].CreatedAt = time.Now().UTC()
		}
	}
	af := AliasFile{Aliases: aliases}
	data, err := json.MarshalIndent(af, "", "  ")
	if err != nil {
		return fmt.Errorf("alias: marshal: %w", err)
	}
	return os.WriteFile(AliasPath(envPath), data, 0600)
}

// ReadAliasFile loads alias records from disk. Returns empty slice if missing.
func ReadAliasFile(envPath string) ([]AliasRecord, error) {
	data, err := os.ReadFile(AliasPath(envPath))
	if os.IsNotExist(err) {
		return []AliasRecord{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("alias: read: %w", err)
	}
	var af AliasFile
	if err := json.Unmarshal(data, &af); err != nil {
		return nil, fmt.Errorf("alias: unmarshal: %w", err)
	}
	return af.Aliases, nil
}

// ApplyAliases returns a new map that includes alias keys pointing to the
// same values as their canonical counterparts. Existing keys are not removed.
func ApplyAliases(secrets map[string]string, aliases []AliasRecord) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}
	for _, a := range aliases {
		if val, ok := secrets[a.Canonical]; ok {
			out[a.Alias] = val
		}
	}
	return out
}
