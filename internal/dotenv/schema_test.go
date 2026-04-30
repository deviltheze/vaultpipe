package dotenv

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSchemaPath_Convention(t *testing.T) {
	got := SchemaPath("/tmp/project/.env")
	want := "/tmp/project/..env.schema.json"
	if got != want {
		t.Errorf("SchemaPath = %q, want %q", got, want)
	}
}

func TestWriteAndReadSchemaFile_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "schema.json")

	s := Schema{
		Version: "1",
		Fields: []SchemaField{
			{Key: "DB_HOST", Required: true},
			{Key: "DB_PORT", Required: false, Default: "5432"},
		},
	}

	if err := WriteSchemaFile(path, s); err != nil {
		t.Fatalf("WriteSchemaFile: %v", err)
	}

	got, err := ReadSchemaFile(path)
	if err != nil {
		t.Fatalf("ReadSchemaFile: %v", err)
	}
	if got.Version != s.Version {
		t.Errorf("Version = %q, want %q", got.Version, s.Version)
	}
	if len(got.Fields) != 2 {
		t.Errorf("Fields len = %d, want 2", len(got.Fields))
	}
	if got.Fields[1].Default != "5432" {
		t.Errorf("Default = %q, want \"5432\"", got.Fields[1].Default)
	}
}

func TestReadSchemaFile_MissingFile(t *testing.T) {
	dir := t.TempDir()
	s, err := ReadSchemaFile(filepath.Join(dir, "missing.json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Fields) != 0 {
		t.Errorf("expected empty schema, got %+v", s)
	}
}

func TestWriteSchemaFile_SetsTimestamp(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "schema.json")

	before := time.Now().UTC().Add(-time.Second)
	if err := WriteSchemaFile(path, Schema{Version: "1"}); err != nil {
		t.Fatalf("WriteSchemaFile: %v", err)
	}

	data, _ := os.ReadFile(path)
	var s Schema
	_ = json.Unmarshal(data, &s)
	if s.CreatedAt.Before(before) {
		t.Errorf("CreatedAt not set correctly: %v", s.CreatedAt)
	}
}

func TestEnforceSchema_NoViolations(t *testing.T) {
	s := Schema{Fields: []SchemaField{
		{Key: "API_KEY", Required: true},
	}}
	secrets := map[string]string{"API_KEY": "abc123"}
	violations := EnforceSchema(s, secrets)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %v", violations)
	}
}

func TestEnforceSchema_MissingRequired(t *testing.T) {
	s := Schema{Fields: []SchemaField{
		{Key: "API_KEY", Required: true},
		{Key: "OPTIONAL", Required: false},
	}}
	secrets := map[string]string{}
	violations := EnforceSchema(s, secrets)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "API_KEY" {
		t.Errorf("violation key = %q, want \"API_KEY\"", violations[0].Key)
	}
}

func TestEnforceSchema_DefaultSkipsViolation(t *testing.T) {
	s := Schema{Fields: []SchemaField{
		{Key: "DB_PORT", Required: true, Default: "5432"},
	}}
	secrets := map[string]string{}
	violations := EnforceSchema(s, secrets)
	if len(violations) != 0 {
		t.Errorf("expected no violations when default present, got %v", violations)
	}
}
