package dotenv

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAcquireLock_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	if err := AcquireLock(path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer ReleaseLock(path)

	if _, err := os.Stat(lockPath(path)); err != nil {
		t.Fatal("lock file not created")
	}
}

func TestAcquireLock_BlocksWhenLocked(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	if err := AcquireLock(path); err != nil {
		t.Fatalf("first acquire failed: %v", err)
	}
	defer ReleaseLock(path)

	err := AcquireLock(path)
	if err == nil {
		t.Fatal("expected ErrLocked, got nil")
	}
}

func TestAcquireLock_StaleLockIsOverwritten(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	lp := lockPath(path)

	// write a stale lock (60 seconds old)
	stale := fmt.Sprintf("9999\n%d\n", time.Now().Unix()-60)
	if err := os.WriteFile(lp, []byte(stale), 0600); err != nil {
		t.Fatal(err)
	}

	if err := AcquireLock(path); err != nil {
		t.Fatalf("expected stale lock to be overwritten, got: %v", err)
	}
	defer ReleaseLock(path)
}

func TestReleaseLock_RemovesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	_ = AcquireLock(path)
	if err := ReleaseLock(path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(lockPath(path)); !os.IsNotExist(err) {
		t.Fatal("lock file should have been removed")
	}
}

func TestReleaseLock_NoOpWhenMissing(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := ReleaseLock(path); err != nil {
		t.Fatalf("unexpected error on missing lock: %v", err)
	}
}
