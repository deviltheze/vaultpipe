package dotenv

import (
	"os"
	"path/filepath"
	"testing"
)

func TestChecksum_Deterministic(t *testing.T) {
	secrets := map[string]string{"B": "2", "A": "1", "C": "3"}
	a := Checksum(secrets)
	b := Checksum(secrets)
	if a != b {
		t.Fatalf("expected same checksum, got %q and %q", a, b)
	}
}

func TestChecksum_OrderIndependent(t *testing.T) {
	s1 := map[string]string{"A": "1", "B": "2"}
	s2 := map[string]string{"B": "2", "A": "1"}
	if Checksum(s1) != Checksum(s2) {
		t.Fatal("checksum should be order-independent")
	}
}

func TestChecksum_DifferentSecrets(t *testing.T) {
	a := Checksum(map[string]string{"KEY": "val1"})
	b := Checksum(map[string]string{"KEY": "val2"})
	if a == b {
		t.Fatal("different secrets should produce different checksums")
	}
}

func TestWriteAndVerifyChecksumFile_Match(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env.sha256")
	secrets := map[string]string{"TOKEN": "abc", "HOST": "localhost"}

	if err := WriteChecksumFile(path, secrets); err != nil {
		t.Fatalf("WriteChecksumFile: %v", err)
	}

	ok, err := VerifyChecksumFile(path, secrets)
	if err != nil {
		t.Fatalf("VerifyChecksumFile: %v", err)
	}
	if !ok {
		t.Fatal("expected checksum to match")
	}
}

func TestVerifyChecksumFile_Mismatch(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env.sha256")
	secrets := map[string]string{"KEY": "original"}

	if err := WriteChecksumFile(path, secrets); err != nil {
		t.Fatalf("WriteChecksumFile: %v", err)
	}

	ok, err := VerifyChecksumFile(path, map[string]string{"KEY": "changed"})
	if err != nil {
		t.Fatalf("VerifyChecksumFile: %v", err)
	}
	if ok {
		t.Fatal("expected checksum mismatch")
	}
}

func TestVerifyChecksumFile_Missing(t *testing.T) {
	ok, err := VerifyChecksumFile("/nonexistent/.env.sha256", map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected false for missing file")
	}
}

func TestWriteChecksumFile_Permissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env.sha256")
	if err := WriteChecksumFile(path, map[string]string{"X": "y"}); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Fatalf("expected 0600, got %v", info.Mode().Perm())
	}
}
