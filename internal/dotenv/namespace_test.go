package dotenv_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/your-org/vaultpipe/internal/dotenv"
)

func TestNamespacePath_Convention(t *testing.T) {
	path := dotenv.NamespacePath("/tmp/env", "prod")
	if path != "/tmp/env/.vaultpipe.namespace.prod.json" {
		t.Errorf("unexpected path: %s", path)
	}
}

func TestApplyNamespace_PrefixesKeys(t *testing.T) {
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	result := dotenv.ApplyNamespace(secrets, "APP")
	if result["APP_FOO"] != "bar" {
		t.Errorf("expected APP_FOO=bar, got %v", result)
	}
	if result["APP_BAZ"] != "qux" {
		t.Errorf("expected APP_BAZ=qux, got %v", result)
	}
	if _, ok := result["FOO"]; ok {
		t.Error("original key FOO should not exist after namespace apply")
	}
}

func TestApplyNamespace_EmptyNamespace_ReturnsOriginal(t *testing.T) {
	secrets := map[string]string{"FOO": "bar"}
	result := dotenv.ApplyNamespace(secrets, "")
	if result["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %v", result)
	}
}

func TestStripNamespace_RemovesPrefix(t *testing.T) {
	secrets := map[string]string{"APP_FOO": "bar", "APP_BAZ": "qux", "OTHER": "val"}
	result := dotenv.StripNamespace(secrets, "APP")
	if result["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %v", result)
	}
	if result["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux, got %v", result)
	}
	if _, ok := result["OTHER"]; ok {
		t.Error("OTHER should be excluded when stripping namespace")
	}
}

func TestWriteAndReadNamespaceFile_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	rec := dotenv.NamespaceRecord{
		Namespace: "staging",
		AppliedAt: time.Now().UTC().Truncate(time.Second),
		Keys:      []string{"DB_HOST", "DB_PORT"},
	}
	err := dotenv.WriteNamespaceFile(dir, "staging", rec)
	if err != nil {
		t.Fatalf("WriteNamespaceFile: %v", err)
	}
	got, err := dotenv.ReadNamespaceFile(dir, "staging")
	if err != nil {
		t.Fatalf("ReadNamespaceFile: %v", err)
	}
	if got.Namespace != "staging" {
		t.Errorf("expected namespace staging, got %s", got.Namespace)
	}
	if len(got.Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(got.Keys))
	}
}

func TestReadNamespaceFile_MissingFile(t *testing.T) {
	dir := t.TempDir()
	_, err := dotenv.ReadNamespaceFile(dir, "missing")
	if !os.IsNotExist(err) {
		t.Errorf("expected not-exist error, got %v", err)
	}
}

func TestWriteNamespaceFile_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	rec := dotenv.NamespaceRecord{Namespace: "dev"}
	if err := dotenv.WriteNamespaceFile(dir, "dev", rec); err != nil {
		t.Fatalf("WriteNamespaceFile: %v", err)
	}
	path := dotenv.NamespacePath(dir, "dev")
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %o", info.Mode().Perm())
	}
}

func TestNamespacePath_UsesDir(t *testing.T) {
	path := dotenv.NamespacePath("/data/secrets", "test")
	dir := filepath.Dir(path)
	if dir != "/data/secrets" {
		t.Errorf("expected dir /data/secrets, got %s", dir)
	}
}
