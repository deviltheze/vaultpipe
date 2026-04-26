package dotenv_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/vaultpipe/internal/dotenv"
)

func TestScopePath_Convention(t *testing.T) {
	got := dotenv.ScopePath("/tmp/envs/.env")
	want := "/tmp/envs/..env.scope.json"
	if got != want {
		t.Errorf("ScopePath = %q; want %q", got, want)
	}
}

func TestWriteAndReadScopeFile_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env.scope.json")

	record := dotenv.ScopeRecord{
		Name: "production",
		Keys: []string{"DB_HOST", "DB_PORT", "API_KEY"},
	}

	if err := dotenv.WriteScopeFile(path, record); err != nil {
		t.Fatalf("WriteScopeFile: %v", err)
	}

	got, err := dotenv.ReadScopeFile(path)
	if err != nil {
		t.Fatalf("ReadScopeFile: %v", err)
	}

	if got.Name != record.Name {
		t.Errorf("Name = %q; want %q", got.Name, record.Name)
	}
	if len(got.Keys) != len(record.Keys) {
		t.Errorf("Keys len = %d; want %d", len(got.Keys), len(record.Keys))
	}
	if got.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
}

func TestReadScopeFile_MissingFile(t *testing.T) {
	record, err := dotenv.ReadScopeFile("/nonexistent/path/.env.scope.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if record.Name != "" || len(record.Keys) != 0 {
		t.Error("expected empty record for missing file")
	}
}

func TestWriteScopeFile_SetsTimestamp(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".scope.json")

	before := time.Now().UTC()
	if err := dotenv.WriteScopeFile(path, dotenv.ScopeRecord{Name: "staging"}); err != nil {
		t.Fatalf("WriteScopeFile: %v", err)
	}

	got, _ := dotenv.ReadScopeFile(path)
	if got.CreatedAt.Before(before) {
		t.Errorf("CreatedAt %v is before test start %v", got.CreatedAt, before)
	}
}

func TestWriteScopeFile_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".scope.json")

	if err := dotenv.WriteScopeFile(path, dotenv.ScopeRecord{Name: "dev"}); err != nil {
		t.Fatalf("WriteScopeFile: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("file perm = %o; want 0600", perm)
	}
}

func TestApplyScope_NoKeys_ReturnsAll(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2", "C": "3"}
	got := dotenv.ApplyScope(secrets, dotenv.ScopeRecord{})
	if len(got) != 3 {
		t.Errorf("expected 3 keys, got %d", len(got))
	}
}

func TestApplyScope_FiltersByKeys(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "localhost", "API_KEY": "secret", "DEBUG": "true"}
	record := dotenv.ScopeRecord{Keys: []string{"DB_HOST", "API_KEY"}}
	got := dotenv.ApplyScope(secrets, record)
	if len(got) != 2 {
		t.Errorf("expected 2 keys, got %d", len(got))
	}
	if _, ok := got["DEBUG"]; ok {
		t.Error("DEBUG should have been filtered out")
	}
}

func TestApplyScope_MissingKeyIgnored(t *testing.T) {
	secrets := map[string]string{"A": "1"}
	record := dotenv.ScopeRecord{Keys: []string{"A", "MISSING"}}
	got := dotenv.ApplyScope(secrets, record)
	if len(got) != 1 {
		t.Errorf("expected 1 key, got %d", len(got))
	}
}
