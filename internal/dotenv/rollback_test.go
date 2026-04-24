package dotenv

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRollback_RestoresMostRecentBackup(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	// Write a current env file.
	if err := os.WriteFile(envPath, []byte("KEY=current\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	// Create two backups with different timestamps.
	older := filepath.Join(dir, ".env.backup.20240101120000")
	newer := filepath.Join(dir, ".env.backup.20240202120000")
	if err := os.WriteFile(older, []byte("KEY=older\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(newer, []byte("KEY=newer\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	result, err := Rollback(envPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.RestoredFrom != newer {
		t.Errorf("expected restored from %s, got %s", newer, result.RestoredFrom)
	}
	if result.RestoredTo != envPath {
		t.Errorf("expected restored to %s, got %s", envPath, result.RestoredTo)
	}

	data, err := os.ReadFile(envPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "KEY=newer\n" {
		t.Errorf("expected restored content %q, got %q", "KEY=newer\n", string(data))
	}
}

func TestRollback_NoBackups_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	_, err := Rollback(envPath)
	if err == nil {
		t.Fatal("expected error when no backups exist, got nil")
	}
}

func TestRollback_SetsTimestamp(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	backup := filepath.Join(dir, ".env.backup.20230615093045")
	if err := os.WriteFile(backup, []byte("KEY=val\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	result, err := Rollback(envPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := time.Date(2023, 6, 15, 9, 30, 45, 0, time.UTC)
	if !result.Timestamp.Equal(expected) {
		t.Errorf("expected timestamp %v, got %v", expected, result.Timestamp)
	}
}

func TestParseBackupTimestamp_InvalidReturnsZero(t *testing.T) {
	ts := parseBackupTimestamp("/some/path/.env.backup.notadate")
	if !ts.IsZero() {
		t.Errorf("expected zero time for invalid timestamp, got %v", ts)
	}
}
