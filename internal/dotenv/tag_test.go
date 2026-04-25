package dotenv_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/your-org/vaultpipe/internal/dotenv"
)

func TestWriteAndReadTagFile_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env.tags.json")

	rec := dotenv.TagRecord{
		Environment: "staging",
		VaultPath:   "secret/staging",
		SyncedAt:    time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		Tags:        map[string]string{"team": "platform", "region": "us-east-1"},
	}

	if err := dotenv.WriteTagFile(path, rec); err != nil {
		t.Fatalf("WriteTagFile: %v", err)
	}

	got, err := dotenv.ReadTagFile(path)
	if err != nil {
		t.Fatalf("ReadTagFile: %v", err)
	}

	if got.Environment != rec.Environment {
		t.Errorf("environment: got %q, want %q", got.Environment, rec.Environment)
	}
	if got.VaultPath != rec.VaultPath {
		t.Errorf("vault_path: got %q, want %q", got.VaultPath, rec.VaultPath)
	}
	if got.Tags["team"] != "platform" {
		t.Errorf("tags[team]: got %q, want %q", got.Tags["team"], "platform")
	}
}

func TestWriteTagFile_SetsTimestampWhenZero(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env.tags.json")

	rec := dotenv.TagRecord{Environment: "prod"}
	if err := dotenv.WriteTagFile(path, rec); err != nil {
		t.Fatalf("WriteTagFile: %v", err)
	}

	got, _ := dotenv.ReadTagFile(path)
	if got.SyncedAt.IsZero() {
		t.Error("expected SyncedAt to be set automatically")
	}
}

func TestReadTagFile_MissingFile(t *testing.T) {
	got, err := dotenv.ReadTagFile("/nonexistent/.env.tags.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if got.Environment != "" {
		t.Errorf("expected zero-value record, got: %+v", got)
	}
}

func TestWriteTagFile_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env.tags.json")

	if err := dotenv.WriteTagFile(path, dotenv.TagRecord{Environment: "dev"}); err != nil {
		t.Fatalf("WriteTagFile: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Errorf("file permissions: got %o, want 600", perm)
	}
}

func TestTagPath_Convention(t *testing.T) {
	got := dotenv.TagPath(".env")
	want := ".env.tags.json"
	if got != want {
		t.Errorf("TagPath: got %q, want %q", got, want)
	}
}
