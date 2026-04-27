package dotenv_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/dotenv"
)

func TestAuditTrailPath_Convention(t *testing.T) {
	p := dotenv.AuditTrailPath("/tmp/mydir/.env")
	want := "/tmp/mydir/..env.audit.json"
	if p != want {
		t.Errorf("got %q, want %q", p, want)
	}
}

func TestWriteAndReadAuditTrail_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, ".env")

	now := time.Now().UTC().Truncate(time.Second)
	trail := dotenv.AuditTrail{
		SyncedAt: now,
		Events: []dotenv.AuditEvent{
			{Timestamp: now, Key: "DB_HOST", Action: "added", Source: "secret/app"},
			{Timestamp: now, Key: "API_KEY", Action: "updated", Source: "secret/app", Actor: "ci"},
		},
	}

	if err := dotenv.WriteAuditTrail(output, trail); err != nil {
		t.Fatalf("WriteAuditTrail: %v", err)
	}

	got, err := dotenv.ReadAuditTrail(output)
	if err != nil {
		t.Fatalf("ReadAuditTrail: %v", err)
	}

	if got.Output != output {
		t.Errorf("Output: got %q, want %q", got.Output, output)
	}
	if len(got.Events) != 2 {
		t.Fatalf("Events len: got %d, want 2", len(got.Events))
	}
	if got.Events[0].Key != "DB_HOST" {
		t.Errorf("Events[0].Key: got %q, want %q", got.Events[0].Key, "DB_HOST")
	}
	if got.Events[1].Actor != "ci" {
		t.Errorf("Events[1].Actor: got %q, want %q", got.Events[1].Actor, "ci")
	}
}

func TestReadAuditTrail_MissingFile(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, ".env")

	trail, err := dotenv.ReadAuditTrail(output)
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(trail.Events) != 0 {
		t.Errorf("expected empty events, got %d", len(trail.Events))
	}
}

func TestWriteAuditTrail_SetsTimestampWhenZero(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, ".env")

	trail := dotenv.AuditTrail{}
	if err := dotenv.WriteAuditTrail(output, trail); err != nil {
		t.Fatalf("WriteAuditTrail: %v", err)
	}

	got, _ := dotenv.ReadAuditTrail(output)
	if got.SyncedAt.IsZero() {
		t.Error("expected SyncedAt to be set automatically")
	}
}

func TestWriteAuditTrail_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, ".env")

	if err := dotenv.WriteAuditTrail(output, dotenv.AuditTrail{}); err != nil {
		t.Fatalf("WriteAuditTrail: %v", err)
	}

	path := dotenv.AuditTrailPath(output)
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("permissions: got %o, want 0600", perm)
	}
}
