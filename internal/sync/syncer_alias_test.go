package sync

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpipe/internal/dotenv"
)

var discardLogger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

func TestApplyAliases_NoAliasFile_ReturnsOriginal(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, ".env")
	secrets := map[string]string{"FOO": "bar"}

	out, err := applyAliases(output, secrets, discardLogger)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 || out["FOO"] != "bar" {
		t.Errorf("expected original secrets, got %v", out)
	}
}

func TestApplyAliases_WithAliasFile_InjectsKeys(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, ".env")

	records := []dotenv.AliasRecord{
		{Alias: "DB", Canonical: "DATABASE_URL"},
	}
	if err := dotenv.WriteAliasFile(output, records); err != nil {
		t.Fatalf("setup: %v", err)
	}

	secrets := map[string]string{"DATABASE_URL": "postgres://localhost"}
	out, err := applyAliases(output, secrets, discardLogger)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB"] != "postgres://localhost" {
		t.Errorf("alias DB not applied, got: %v", out)
	}
	if out["DATABASE_URL"] != "postgres://localhost" {
		t.Error("original key should be preserved")
	}
}

func TestWriteAliasFile_WritesAndLogs(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, ".env")

	records := []dotenv.AliasRecord{
		{Alias: "API", Canonical: "API_KEY"},
	}
	if err := writeAliasFile(output, records, discardLogger); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := dotenv.ReadAliasFile(output)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if len(got) != 1 || got[0].Alias != "API" {
		t.Errorf("unexpected records: %v", got)
	}
}

func TestWriteAliasFile_EmptyRecords_NoOp(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, ".env")

	if err := writeAliasFile(output, nil, discardLogger); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(dotenv.AliasPath(output)); !os.IsNotExist(err) {
		t.Error("alias file should not be created for empty records")
	}
}
