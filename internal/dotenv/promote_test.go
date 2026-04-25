package dotenv_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpipe/internal/dotenv"
)

func writeTempPromoteEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTempPromoteEnv: %v", err)
	}
	return p
}

func TestPromote_AllKeys(t *testing.T) {
	dir := t.TempDir()
	src := writeTempPromoteEnv(t, dir, ".env.staging", "DB_HOST=staging-db\nAPI_KEY=abc123\n")
	dst := filepath.Join(dir, ".env.production")

	res, err := dotenv.Promote(src, dst, dotenv.PromoteOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Promoted) != 2 {
		t.Errorf("expected 2 promoted, got %d", len(res.Promoted))
	}

	got, err := dotenv.Read(dst)
	if err != nil {
		t.Fatalf("read target: %v", err)
	}
	if got["DB_HOST"] != "staging-db" {
		t.Errorf("DB_HOST: got %q", got["DB_HOST"])
	}
	if got["API_KEY"] != "abc123" {
		t.Errorf("API_KEY: got %q", got["API_KEY"])
	}
}

func TestPromote_SpecificKeys(t *testing.T) {
	dir := t.TempDir()
	src := writeTempPromoteEnv(t, dir, ".env.staging", "DB_HOST=staging-db\nAPI_KEY=abc123\nDEBUG=true\n")
	dst := filepath.Join(dir, ".env.production")

	res, err := dotenv.Promote(src, dst, dotenv.PromoteOptions{Keys: []string{"DB_HOST"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Promoted) != 1 || res.Promoted[0] != "DB_HOST" {
		t.Errorf("expected [DB_HOST] promoted, got %v", res.Promoted)
	}
	if len(res.Skipped) != 2 {
		t.Errorf("expected 2 skipped, got %d", len(res.Skipped))
	}

	got, _ := dotenv.Read(dst)
	if _, ok := got["API_KEY"]; ok {
		t.Error("API_KEY should not have been promoted")
	}
}

func TestPromote_MergesWithExistingTarget(t *testing.T) {
	dir := t.TempDir()
	src := writeTempPromoteEnv(t, dir, ".env.staging", "NEW_KEY=hello\n")
	dst := writeTempPromoteEnv(t, dir, ".env.production", "EXISTING=keep\n")

	_, err := dotenv.Promote(src, dst, dotenv.PromoteOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, _ := dotenv.Read(dst)
	if got["EXISTING"] != "keep" {
		t.Errorf("EXISTING should be preserved, got %q", got["EXISTING"])
	}
	if got["NEW_KEY"] != "hello" {
		t.Errorf("NEW_KEY should be promoted, got %q", got["NEW_KEY"])
	}
}

func TestPromote_DryRun_DoesNotWrite(t *testing.T) {
	dir := t.TempDir()
	src := writeTempPromoteEnv(t, dir, ".env.staging", "FOO=bar\n")
	dst := filepath.Join(dir, ".env.production")

	res, err := dotenv.Promote(src, dst, dotenv.PromoteOptions{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.DryRun {
		t.Error("expected DryRun=true in result")
	}
	if len(res.Promoted) != 1 {
		t.Errorf("expected 1 promoted in dry-run, got %d", len(res.Promoted))
	}
	if _, err := os.Stat(dst); err == nil {
		t.Error("target file should not exist after dry-run")
	}
}

func TestPromote_MissingSource_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	_, err := dotenv.Promote(
		filepath.Join(dir, "nonexistent.env"),
		filepath.Join(dir, "out.env"),
		dotenv.PromoteOptions{},
	)
	if err == nil {
		t.Error("expected error for missing source")
	}
}

func TestPromoteResult_String(t *testing.T) {
	r := dotenv.PromoteResult{Promoted: []string{"A", "B"}, Skipped: []string{"C"}, DryRun: true}
	s := r.String()
	if s != "[dry-run] promoted=2 skipped=1" {
		t.Errorf("unexpected String(): %q", s)
	}
}
