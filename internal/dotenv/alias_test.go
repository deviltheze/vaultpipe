package dotenv_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/dotenv"
)

func TestAliasPath_Convention(t *testing.T) {
	p := dotenv.AliasPath("/tmp/project/.env")
	if filepath.Base(p) != "..env.aliases.json" {
		t.Fatalf("unexpected alias path: %s", p)
	}
}

func TestWriteAndReadAliasFile_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	records := []dotenv.AliasRecord{
		{Alias: "DB_URL", Canonical: "DATABASE_URL"},
		{Alias: "API", Canonical: "API_KEY"},
	}

	if err := dotenv.WriteAliasFile(envPath, records); err != nil {
		t.Fatalf("write: %v", err)
	}

	got, err := dotenv.ReadAliasFile(envPath)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 records, got %d", len(got))
	}
	if got[0].Alias != "DB_URL" || got[0].Canonical != "DATABASE_URL" {
		t.Errorf("record mismatch: %+v", got[0])
	}
}

func TestWriteAliasFile_SetsTimestamp(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	records := []dotenv.AliasRecord{
		{Alias: "FOO", Canonical: "FOO_ORIGINAL"},
	}
	if err := dotenv.WriteAliasFile(envPath, records); err != nil {
		t.Fatalf("write: %v", err)
	}
	got, _ := dotenv.ReadAliasFile(envPath)
	if got[0].CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestReadAliasFile_MissingFile(t *testing.T) {
	dir := t.TempDir()
	records, err := dotenv.ReadAliasFile(filepath.Join(dir, ".env"))
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(records) != 0 {
		t.Errorf("expected empty slice, got %d records", len(records))
	}
}

func TestWriteAliasFile_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	_ = dotenv.WriteAliasFile(envPath, []dotenv.AliasRecord{
		{Alias: "X", Canonical: "Y", CreatedAt: time.Now()},
	})
	info, err := os.Stat(dotenv.AliasPath(envPath))
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %o", info.Mode().Perm())
	}
}

func TestApplyAliases_MapsValues(t *testing.T) {
	secrets := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"API_KEY":      "secret123",
	}
	aliases := []dotenv.AliasRecord{
		{Alias: "DB_URL", Canonical: "DATABASE_URL"},
		{Alias: "API", Canonical: "API_KEY"},
	}
	out := dotenv.ApplyAliases(secrets, aliases)
	if out["DB_URL"] != "postgres://localhost/db" {
		t.Errorf("DB_URL not aliased correctly: %s", out["DB_URL"])
	}
	if out["API"] != "secret123" {
		t.Errorf("API not aliased correctly: %s", out["API"])
	}
	if _, ok := out["DATABASE_URL"]; !ok {
		t.Error("original key DATABASE_URL should be preserved")
	}
}

func TestApplyAliases_SkipsMissingCanonical(t *testing.T) {
	secrets := map[string]string{"PRESENT": "val"}
	aliases := []dotenv.AliasRecord{
		{Alias: "GHOST", Canonical: "MISSING_KEY"},
	}
	out := dotenv.ApplyAliases(secrets, aliases)
	if _, ok := out["GHOST"]; ok {
		t.Error("alias for missing canonical should not be added")
	}
}

func TestApplyAliases_DoesNotMutateOriginal(t *testing.T) {
	secrets := map[string]string{"FOO": "bar"}
	aliases := []dotenv.AliasRecord{{Alias: "BAZ", Canonical: "FOO"}}
	_ = dotenv.ApplyAliases(secrets, aliases)
	if len(secrets) != 1 {
		t.Error("original secrets map was mutated")
	}
}
