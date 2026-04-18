package sync

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/dotenv"
)

func TestCheckTTL_NoFile_Proceeds(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, ".env")
	if err := checkTTL(out); err != nil {
		t.Fatalf("expected nil when no TTL file exists, got: %v", err)
	}
}

func TestCheckTTL_Expired_Proceeds(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, ".env")
	rec := dotenv.TTLRecord{
		Path:     "secret/app",
		SyncedAt: time.Now().Add(-2 * time.Hour),
		TTL:      1 * time.Hour,
	}
	_ = dotenv.WriteTTLFile(ttlPath(out), rec)
	if err := checkTTL(out); err != nil {
		t.Fatalf("expected nil for expired TTL, got: %v", err)
	}
}

func TestCheckTTL_Fresh_Blocks(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, ".env")
	rec := dotenv.TTLRecord{
		Path:     "secret/app",
		SyncedAt: time.Now(),
		TTL:      10 * time.Minute,
	}
	_ = dotenv.WriteTTLFile(ttlPath(out), rec)
	if err := checkTTL(out); err == nil {
		t.Fatal("expected error for fresh TTL")
	}
}

func TestWriteTTL_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, ".env")
	if err := writeTTL(out, "secret/app", 5*time.Minute); err != nil {
		t.Fatalf("writeTTL: %v", err)
	}
	p := ttlPath(out)
	if _, err := os.Stat(p); err != nil {
		t.Fatalf("TTL file not created: %v", err)
	}
	rec, err := dotenv.ReadTTLFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if rec.TTL != 5*time.Minute {
		t.Errorf("TTL mismatch: got %v", rec.TTL)
	}
}
