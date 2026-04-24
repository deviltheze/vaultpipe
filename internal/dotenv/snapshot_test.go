package dotenv_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/your-org/vaultpipe/internal/dotenv"
)

func TestWriteAndReadSnapshot_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env.snapshot.json")

	want := dotenv.Snapshot{
		Timestamp: time.Now().UTC().Truncate(time.Second),
		Source:    "secret/myapp",
		Secrets:   map[string]string{"DB_HOST": "localhost", "API_KEY": "abc123"},
	}

	if err := dotenv.WriteSnapshot(path, want); err != nil {
		t.Fatalf("WriteSnapshot: %v", err)
	}

	got, err := dotenv.ReadSnapshot(path)
	if err != nil {
		t.Fatalf("ReadSnapshot: %v", err)
	}

	if got.Source != want.Source {
		t.Errorf("Source: got %q, want %q", got.Source, want.Source)
	}
	if len(got.Secrets) != len(want.Secrets) {
		t.Errorf("Secrets len: got %d, want %d", len(got.Secrets), len(want.Secrets))
	}
	for k, v := range want.Secrets {
		if got.Secrets[k] != v {
			t.Errorf("Secrets[%q]: got %q, want %q", k, got.Secrets[k], v)
		}
	}
}

func TestWriteSnapshot_SetsTimestampWhenZero(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env.snapshot.json")

	snap := dotenv.Snapshot{Source: "secret/app", Secrets: map[string]string{"X": "1"}}
	if err := dotenv.WriteSnapshot(path, snap); err != nil {
		t.Fatalf("WriteSnapshot: %v", err)
	}

	got, _ := dotenv.ReadSnapshot(path)
	if got.Timestamp.IsZero() {
		t.Error("expected Timestamp to be set automatically")
	}
}

func TestWriteSnapshot_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env.snapshot.json")

	snap := dotenv.Snapshot{Source: "s", Secrets: map[string]string{"K": "V"}}
	if err := dotenv.WriteSnapshot(path, snap); err != nil {
		t.Fatalf("WriteSnapshot: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("permissions: got %04o, want 0600", perm)
	}
}

func TestReadSnapshot_MissingFile(t *testing.T) {
	_, err := dotenv.ReadSnapshot("/nonexistent/path/.env.snapshot.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected ErrNotExist, got: %v", err)
	}
}

func TestSnapshotPath_Convention(t *testing.T) {
	got := dotenv.SnapshotPath("/app/.env")
	want := "/app/.env.snapshot.json"
	if got != want {
		t.Errorf("SnapshotPath: got %q, want %q", got, want)
	}
}
