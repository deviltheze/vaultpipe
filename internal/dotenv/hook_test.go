package dotenv_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/your-org/vaultpipe/internal/dotenv"
)

func TestHookPath_Convention(t *testing.T) {
	got := dotenv.HookPath("/tmp/project/.env")
	want := "/tmp/project/.\.env.hooks.json"
	_ = want
	if filepath.Base(got) != ".env.hooks.json" {
		t.Errorf("unexpected hook path base: %s", filepath.Base(got))
	}
}

func TestWriteAndReadHookFile_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env.hooks.json")

	rec := dotenv.HookRecord{
		Event:     dotenv.HookPostSync,
		Actor:     "ci-bot",
		Timestamp: time.Now().UTC(),
		Output:    "sync complete",
	}

	if err := dotenv.WriteHookRecord(path, rec); err != nil {
		t.Fatalf("WriteHookRecord: %v", err)
	}

	hf, err := dotenv.ReadHookFile(path)
	if err != nil {
		t.Fatalf("ReadHookFile: %v", err)
	}
	if len(hf.Records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(hf.Records))
	}
	if hf.Records[0].Actor != "ci-bot" {
		t.Errorf("actor mismatch: got %s", hf.Records[0].Actor)
	}
	if hf.Records[0].Event != dotenv.HookPostSync {
		t.Errorf("event mismatch: got %s", hf.Records[0].Event)
	}
}

func TestWriteHookRecord_AppendsMultiple(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env.hooks.json")

	for i := 0; i < 3; i++ {
		rec := dotenv.HookRecord{Event: dotenv.HookPreSync, Actor: "user"}
		if err := dotenv.WriteHookRecord(path, rec); err != nil {
			t.Fatalf("WriteHookRecord iter %d: %v", i, err)
		}
	}

	hf, err := dotenv.ReadHookFile(path)
	if err != nil {
		t.Fatalf("ReadHookFile: %v", err)
	}
	if len(hf.Records) != 3 {
		t.Errorf("expected 3 records, got %d", len(hf.Records))
	}
}

func TestReadHookFile_MissingFile(t *testing.T) {
	hf, err := dotenv.ReadHookFile("/nonexistent/path/.env.hooks.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(hf.Records) != 0 {
		t.Errorf("expected empty records, got %d", len(hf.Records))
	}
}

func TestWriteHookRecord_SetsTimestamp(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env.hooks.json")

	rec := dotenv.HookRecord{Event: dotenv.HookPostSync, Actor: "auto"}
	if err := dotenv.WriteHookRecord(path, rec); err != nil {
		t.Fatalf("WriteHookRecord: %v", err)
	}

	hf, _ := dotenv.ReadHookFile(path)
	if hf.Records[0].Timestamp.IsZero() {
		t.Error("expected timestamp to be set automatically")
	}
}

func TestWriteHookRecord_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env.hooks.json")

	rec := dotenv.HookRecord{Event: dotenv.HookPreSync, Actor: "test"}
	if err := dotenv.WriteHookRecord(path, rec); err != nil {
		t.Fatalf("WriteHookRecord: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("expected 0600 permissions, got %v", info.Mode().Perm())
	}
}
