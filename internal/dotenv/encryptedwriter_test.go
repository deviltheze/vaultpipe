package dotenv

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func testEncKey() []byte { return bytes.Repeat([]byte("e"), 32) }

func TestWriteEncrypted_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env.enc")
	secrets := map[string]string{
		"APP_NAME":     "vaultpipe",
		"DATABASE_PASSWORD": "s3cr3t",
		"API_KEY":      "key-abc",
	}
	if err := WriteEncrypted(path, secrets, testEncKey()); err != nil {
		t.Fatalf("write: %v", err)
	}
	got, err := ReadEncrypted(path, testEncKey())
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	for k, want := range secrets {
		if got[k] != want {
			t.Errorf("key %q: got %q want %q", k, got[k], want)
		}
	}
}

func TestWriteEncrypted_SensitiveValuesObfuscated(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env.enc")
	secrets := map[string]string{"API_KEY": "plain-value"}
	if err := WriteEncrypted(path, secrets, testEncKey()); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(path)
	if bytes.Contains(data, []byte("plain-value")) {
		t.Error("sensitive value should not appear in plain text")
	}
}

func TestWriteEncrypted_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env.enc")
	_ = WriteEncrypted(path, map[string]string{"X": "y"}, testEncKey())
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("expected 0600, got %v", info.Mode().Perm())
	}
}

func TestReadEncrypted_WrongKey(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env.enc")
	_ = WriteEncrypted(path, map[string]string{"API_KEY": "secret"}, testEncKey())
	wrongKey := bytes.Repeat([]byte("x"), 32)
	_, err := ReadEncrypted(path, wrongKey)
	if err == nil {
		t.Error("expected error when decrypting with wrong key")
	}
}
