package dotenv_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/dotenv"
)

func TestTTLRecord_NotExpired(t *testing.T) {
	rec := dotenv.TTLRecord{
		Path:     "secret/app",
		SyncedAt: time.Now(),
		TTL:      10 * time.Minute,
	}
	if rec.IsExpired() {
		t.Fatal("expected not expired")
	}
}

func TestTTLRecord_Expired(t *testing.T) {
	rec := dotenv.TTLRecord{
		Path:     "secret/app",
		SyncedAt: time.Now().Add(-2 * time.Hour),
		TTL:      1 * time.Hour,
	}
	if !rec.IsExpired() {
		t.Fatal("expected expired")
	}
}

func TestTTLRecord_ZeroTTL_NeverExpires(t *testing.T) {
	rec := dotenv.TTLRecord{
		Path:     "secret/app",
		SyncedAt: time.Now().Add(-100 * time.Hour),
		TTL:      0,
	}
	if rec.IsExpired() {
		t.Fatal("zero TTL should never expire")
	}
}

func TestWriteAndReadTTLFile_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, ".ttl.json")

	orig := dotenv.TTLRecord{
		Path:     "secret/myapp",
		SyncedAt: time.Now().Truncate(time.Second),
		TTL:      30 * time.Minute,
	}
	if err := dotenv.WriteTTLFile(p, orig); err != nil {
		t.Fatalf("write: %v", err)
	}

	got, err := dotenv.ReadTTLFile(p)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if got.Path != orig.Path {
		t.Errorf("path: got %q want %q", got.Path, orig.Path)
	}
	if got.TTL != orig.TTL {
		t.Errorf("ttl: got %v want %v", got.TTL, orig.TTL)
	}
}

func TestReadTTLFile_Missing(t *testing.T) {
	_, err := dotenv.ReadTTLFile("/nonexistent/.ttl.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestWriteTTLFile_Permissions(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, ".ttl.json")
	if err := dotenv.WriteTTLFile(p, dotenv.TTLRecord{Path: "x", SyncedAt: time.Now(), TTL: time.Minute}); err != nil {
		t.Fatal(err)
	}
	info, _ := os.Stat(p)
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %v", info.Mode().Perm())
	}
}
