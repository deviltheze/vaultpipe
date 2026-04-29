package sync

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpipe/internal/dotenv"
)

func TestWriteImportRecord_WritesFile(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, ".env")
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	secrets := map[string]string{"DB_HOST": "localhost", "API_KEY": "xyz"}
	writeImportRecord(logger, output, "secret/myapp", secrets)

	rec, err := dotenv.ReadImportFile(output)
	if err != nil {
		t.Fatalf("ReadImportFile: %v", err)
	}
	if rec.Source != "secret/myapp" {
		t.Errorf("Source = %q, want secret/myapp", rec.Source)
	}
	if len(rec.Secrets) != 2 {
		t.Errorf("Secrets len = %d, want 2", len(rec.Secrets))
	}
}

func TestVerifyImportSource_NoFile_Passes(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, ".env")

	if err := verifyImportSource(output, "secret/myapp"); err != nil {
		t.Errorf("expected no error for missing import file, got: %v", err)
	}
}

func TestVerifyImportSource_MatchingSource_Passes(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, ".env")
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	writeImportRecord(logger, output, "secret/myapp", map[string]string{"X": "1"})

	if err := verifyImportSource(output, "secret/myapp"); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestVerifyImportSource_Mismatch_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, ".env")
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	writeImportRecord(logger, output, "secret/original", map[string]string{"X": "1"})

	err := verifyImportSource(output, "secret/different")
	if err == nil {
		t.Fatal("expected error for source mismatch, got nil")
	}
}
