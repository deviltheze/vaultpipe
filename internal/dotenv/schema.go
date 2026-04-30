package dotenv

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// SchemaField describes the expected shape of a single secret key.
type SchemaField struct {
	Key      string `json:"key"`
	Required bool   `json:"required"`
	Pattern  string `json:"pattern,omitempty"` // optional regex hint (informational)
	Default  string `json:"default,omitempty"`
}

// Schema holds the full set of expected fields for an environment.
type Schema struct {
	Version   string        `json:"version"`
	Fields    []SchemaField `json:"fields"`
	CreatedAt time.Time     `json:"created_at"`
}

// SchemaViolation describes a single schema enforcement failure.
type SchemaViolation struct {
	Key     string
	Message string
}

func (v SchemaViolation) Error() string {
	return fmt.Sprintf("schema violation: key %q — %s", v.Key, v.Message)
}

// SchemaPath returns the conventional path for a schema file.
func SchemaPath(envFile string) string {
	dir := filepath.Dir(envFile)
	base := filepath.Base(envFile)
	return filepath.Join(dir, "."+base+".schema.json")
}

// WriteSchemaFile serialises a Schema to disk.
func WriteSchemaFile(path string, s Schema) error {
	if s.CreatedAt.IsZero() {
		s.CreatedAt = time.Now().UTC()
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("schema: marshal: %w", err)
	}
	return os.WriteFile(path, data, 0o600)
}

// ReadSchemaFile deserialises a Schema from disk.
// Returns an empty Schema (no error) when the file does not exist.
func ReadSchemaFile(path string) (Schema, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return Schema{}, nil
	}
	if err != nil {
		return Schema{}, fmt.Errorf("schema: read: %w", err)
	}
	var s Schema
	if err := json.Unmarshal(data, &s); err != nil {
		return Schema{}, fmt.Errorf("schema: unmarshal: %w", err)
	}
	return s, nil
}

// EnforceSchema checks secrets against the schema and returns all violations.
// Missing required keys and keys that have no default are reported.
func EnforceSchema(s Schema, secrets map[string]string) []SchemaViolation {
	var violations []SchemaViolation
	for _, field := range s.Fields {
		val, exists := secrets[field.Key]
		if !exists || val == "" {
			if field.Required && field.Default == "" {
				violations = append(violations, SchemaViolation{
					Key:     field.Key,
					Message: "required key is missing or empty",
				})
			}
		}
	}
	return violations
}
