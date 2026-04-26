// Package dotenv provides utilities for reading, writing, and managing .env files.
package dotenv

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// NamespaceRecord holds metadata about a namespace mapping.
type NamespaceRecord struct {
	Namespace string            `json:"namespace"`
	Prefix    string            `json:"prefix"`
	Keys      map[string]string `json:"keys"` // original key -> namespaced key
	CreatedAt time.Time         `json:"created_at"`
}

// NamespacePath returns the conventional path for a namespace metadata file.
func NamespacePath(envFile string) string {
	dir := filepath.Dir(envFile)
	base := strings.TrimSuffix(filepath.Base(envFile), filepath.Ext(envFile))
	return filepath.Join(dir, "."+base+".namespace.json")
}

// ApplyNamespace prefixes each key in secrets with the given namespace prefix,
// separated by an underscore. Keys that already carry the prefix are left unchanged.
// Returns a new map; the original is not mutated.
func ApplyNamespace(secrets map[string]string, namespace string) map[string]string {
	if namespace == "" {
		result := make(map[string]string, len(secrets))
		for k, v := range secrets {
			result[k] = v
		}
		return result
	}

	prefix := strings.ToUpper(namespace) + "_"
	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if strings.HasPrefix(k, prefix) {
			result[k] = v
		} else {
			result[prefix+k] = v
		}
	}
	return result
}

// StripNamespace removes the namespace prefix from all keys in secrets.
// Keys that do not carry the prefix are left unchanged.
// Returns a new map; the original is not mutated.
func StripNamespace(secrets map[string]string, namespace string) map[string]string {
	if namespace == "" {
		result := make(map[string]string, len(secrets))
		for k, v := range secrets {
			result[k] = v
		}
		return result
	}

	prefix := strings.ToUpper(namespace) + "_"
	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		result[strings.TrimPrefix(k, prefix)] = v
	}
	return result
}

// WriteNamespaceFile persists a NamespaceRecord as JSON to the conventional path.
func WriteNamespaceFile(envFile, namespace string, keys map[string]string) error {
	record := NamespaceRecord{
		Namespace: namespace,
		Prefix:    strings.ToUpper(namespace) + "_",
		Keys:      keys,
		CreatedAt: time.Now().UTC(),
	}

	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return fmt.Errorf("namespace: marshal: %w", err)
	}

	path := NamespacePath(envFile)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("namespace: write %s: %w", path, err)
	}
	return nil
}

// ReadNamespaceFile loads a NamespaceRecord from the conventional path.
// Returns an empty record (no error) when the file does not exist.
func ReadNamespaceFile(envFile string) (NamespaceRecord, error) {
	path := NamespacePath(envFile)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return NamespaceRecord{}, nil
	}
	if err != nil {
		return NamespaceRecord{}, fmt.Errorf("namespace: read %s: %w", path, err)
	}

	var record NamespaceRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return NamespaceRecord{}, fmt.Errorf("namespace: unmarshal: %w", err)
	}
	return record, nil
}
