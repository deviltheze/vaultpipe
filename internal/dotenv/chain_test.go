package dotenv_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/dotenv"
)

func TestChainPath_Convention(t *testing.T) {
	got := dotenv.ChainPath("/tmp/app.env")
	want := "/tmp/app.env.chain.json"
	if got != want {
		t.Errorf("ChainPath = %q, want %q", got, want)
	}
}

func TestWriteAndReadChainFile_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "app.env.chain.json")

	rec := dotenv.ChainRecord{
		Source: "vault:secret/app",
		Dest:   "app.env",
		Keys:   []string{"DB_URL", "API_KEY"},
		AppliedAt: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
	}

	if err := dotenv.WriteChainRecord(path, rec); err != nil {
		t.Fatalf("WriteChainRecord: %v", err)
	}

	cf, err := dotenv.ReadChainFile(path)
	if err != nil {
		t.Fatalf("ReadChainFile: %v", err)
	}

	if len(cf.Records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(cf.Records))
	}
	if cf.Records[0].Source != rec.Source {
		t.Errorf("Source = %q, want %q", cf.Records[0].Source, rec.Source)
	}
	if len(cf.Records[0].Keys) != 2 {
		t.Errorf("Keys length = %d, want 2", len(cf.Records[0].Keys))
	}
}

func TestWriteChainRecord_AppendsMultiple(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "app.env.chain.json")

	for i := 0; i < 3; i++ {
		rec := dotenv.ChainRecord{
			Source: "vault:secret/app",
			Dest:   "app.env",
			Keys:   []string{"KEY"},
		}
		if err := dotenv.WriteChainRecord(path, rec); err != nil {
			t.Fatalf("WriteChainRecord iteration %d: %v", i, err)
		}
	}

	cf, err := dotenv.ReadChainFile(path)
	if err != nil {
		t.Fatalf("ReadChainFile: %v", err)
	}
	if len(cf.Records) != 3 {
		t.Errorf("expected 3 records, got %d", len(cf.Records))
	}
}

func TestWriteChainRecord_SetsTimestampWhenZero(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "app.env.chain.json")

	rec := dotenv.ChainRecord{
		Source: "vault:secret/app",
		Dest:   "app.env",
		Keys:   []string{"X"},
	}

	if err := dotenv.WriteChainRecord(path, rec); err != nil {
		t.Fatalf("WriteChainRecord: %v", err)
	}

	cf, _ := dotenv.ReadChainFile(path)
	if cf.Records[0].AppliedAt.IsZero() {
		t.Error("expected AppliedAt to be set automatically")
	}
}

func TestReadChainFile_MissingFile(t *testing.T) {
	cf, err := dotenv.ReadChainFile("/nonexistent/path.chain.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(cf.Records) != 0 {
		t.Errorf("expected empty records, got %d", len(cf.Records))
	}
}

func TestWriteChainRecord_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "app.env.chain.json")

	rec := dotenv.ChainRecord{Source: "s", Dest: "d", Keys: []string{"K"}}
	if err := dotenv.WriteChainRecord(path, rec); err != nil {
		t.Fatalf("WriteChainRecord: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("file mode = %o, want 0600", info.Mode().Perm())
	}
}
