package dotenv_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/dotenv"
)

func TestWriteAndReadExpiryFile_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	env := filepath.Join(dir, ".env")

	if err := dotenv.WriteExpiryFile(env, 24*time.Hour); err != nil {
		t.Fatalf("WriteExpiryFile: %v", err)
	}

	rec, err := dotenv.ReadExpiryFile(env)
	if err != nil {
		t.Fatalf("ReadExpiryFile: %v", err)
	}
	if rec.Path != env {
		t.Errorf("path: got %q, want %q", rec.Path, env)
	}
	if rec.ExpiresAt.IsZero() {
		t.Error("ExpiresAt should not be zero when TTL > 0")
	}
	if rec.IsExpired() {
		t.Error("record should not be expired immediately after creation")
	}
}

func TestExpiryRecord_ZeroTTL_NeverExpires(t *testing.T) {
	dir := t.TempDir()
	env := filepath.Join(dir, ".env")

	if err := dotenv.WriteExpiryFile(env, 0); err != nil {
		t.Fatalf("WriteExpiryFile: %v", err)
	}

	rec, err := dotenv.ReadExpiryFile(env)
	if err != nil {
		t.Fatalf("ReadExpiryFile: %v", err)
	}
	if rec.IsExpired() {
		t.Error("zero-TTL record should never be expired")
	}
}

func TestExpiryRecord_Expired(t *testing.T) {
	rec := dotenv.ExpiryRecord{
		ExpiresAt: time.Now().Add(-1 * time.Second),
	}
	if !rec.IsExpired() {
		t.Error("record in the past should be expired")
	}
}

func TestReadExpiryFile_MissingFile(t *testing.T) {
	dir := t.TempDir()
	env := filepath.Join(dir, ".env")

	rec, err := dotenv.ReadExpiryFile(env)
	if err != nil {
		t.Fatalf("unexpected error for missing file: %v", err)
	}
	if rec.IsExpired() {
		t.Error("missing expiry file should be treated as never-expiring")
	}
}

func TestWriteExpiryFile_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	env := filepath.Join(dir, ".env")

	if err := dotenv.WriteExpiryFile(env, time.Hour); err != nil {
		t.Fatalf("WriteExpiryFile: %v", err)
	}

	info, err := os.Stat(dotenv.ExpiryPath(env))
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Errorf("permissions: got %o, want 0600", perm)
	}
}

func TestExpiryPath_Convention(t *testing.T) {
	got := dotenv.ExpiryPath("/tmp/.env")
	want := "/tmp/.env.expiry"
	if got != want {
		t.Errorf("ExpiryPath: got %q, want %q", got, want)
	}
}
