package dotenv_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/dotenv"
)

func TestWriteAndReadLineage_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "lineage.jsonl")

	rec := dotenv.LineageRecord{
		Timestamp:  time.Now().UTC().Truncate(time.Second),
		Source:     "secret/myapp",
		OutputFile: ".env",
		Added:      3,
		Updated:    1,
		Removed:    0,
		Checksum:   "abc123",
		Meta:       map[string]string{"env": "production"},
	}

	if err := dotenv.WriteLineage(path, rec); err != nil {
		t.Fatalf("WriteLineage: %v", err)
	}

	records, err := dotenv.ReadLineage(path)
	if err != nil {
		t.Fatalf("ReadLineage: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}

	got := records[0]
	if got.Source != rec.Source {
		t.Errorf("Source: got %q, want %q", got.Source, rec.Source)
	}
	if got.Added != rec.Added {
		t.Errorf("Added: got %d, want %d", got.Added, rec.Added)
	}
	if got.Checksum != rec.Checksum {
		t.Errorf("Checksum: got %q, want %q", got.Checksum, rec.Checksum)
	}
	if got.Meta["env"] != "production" {
		t.Errorf("Meta[env]: got %q, want %q", got.Meta["env"], "production")
	}
}

func TestWriteLineage_AppendsMultiple(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "lineage.jsonl")

	for i := 0; i < 3; i++ {
		rec := dotenv.LineageRecord{
			Timestamp:  time.Now().UTC(),
			Source:     "secret/app",
			OutputFile: ".env",
			Added:      i,
		}
		if err := dotenv.WriteLineage(path, rec); err != nil {
			t.Fatalf("WriteLineage iteration %d: %v", i, err)
		}
	}

	records, err := dotenv.ReadLineage(path)
	if err != nil {
		t.Fatalf("ReadLineage: %v", err)
	}
	if len(records) != 3 {
		t.Errorf("expected 3 records, got %d", len(records))
	}
}

func TestReadLineage_MissingFile(t *testing.T) {
	records, err := dotenv.ReadLineage("/nonexistent/path/lineage.jsonl")
	if err != nil {
		t.Fatalf("expected nil error for missing file, got %v", err)
	}
	if records != nil {
		t.Errorf("expected nil records for missing file, got %v", records)
	}
}

func TestWriteLineage_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "lineage.jsonl")

	rec := dotenv.LineageRecord{Source: "secret/test", OutputFile: ".env"}
	if err := dotenv.WriteLineage(path, rec); err != nil {
		t.Fatalf("WriteLineage: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("expected permissions 0600, got %04o", perm)
	}
}
