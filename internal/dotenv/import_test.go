package dotenv

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestImportPath_Convention(t *testing.T) {
	got := ImportPath("/tmp/myapp/.env")
	want := "/tmp/myapp/..env.import.json"
	if got != want {
		t.Errorf("ImportPath = %q, want %q", got, want)
	}
}

func TestWriteAndReadImportFile_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, ".env")

	rec := ImportRecord{
		Source: "secret/myapp",
		Secrets: map[string]string{"DB_HOST": "localhost", "API_KEY": "abc123"},
	}

	if err := WriteImportFile(output, rec); err != nil {
		t.Fatalf("WriteImportFile: %v", err)
	}

	got, err := ReadImportFile(output)
	if err != nil {
		t.Fatalf("ReadImportFile: %v", err)
	}

	if got.Source != rec.Source {
		t.Errorf("Source = %q, want %q", got.Source, rec.Source)
	}
	if len(got.Keys) != 2 {
		t.Errorf("Keys len = %d, want 2", len(got.Keys))
	}
	if got.Secrets["DB_HOST"] != "localhost" {
		t.Errorf("DB_HOST = %q, want localhost", got.Secrets["DB_HOST"])
	}
}

func TestWriteImportFile_SetsTimestampWhenZero(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, ".env")

	before := time.Now().UTC()
	rec := ImportRecord{Source: "secret/app", Secrets: map[string]string{"X": "1"}}
	if err := WriteImportFile(output, rec); err != nil {
		t.Fatalf("WriteImportFile: %v", err)
	}

	got, _ := ReadImportFile(output)
	if got.ImportedAt.Before(before) {
		t.Errorf("ImportedAt not set correctly: %v", got.ImportedAt)
	}
}

func TestReadImportFile_MissingFile(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, ".env")

	rec, err := ReadImportFile(output)
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if rec.Source != "" {
		t.Errorf("expected empty record, got source=%q", rec.Source)
	}
}

func TestWriteImportFile_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, ".env")

	rec := ImportRecord{Source: "vault/path", Secrets: map[string]string{"K": "V"}}
	if err := WriteImportFile(output, rec); err != nil {
		t.Fatalf("WriteImportFile: %v", err)
	}

	info, err := os.Stat(ImportPath(output))
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("permissions = %o, want 0600", info.Mode().Perm())
	}
}
