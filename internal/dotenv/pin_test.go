package dotenv_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/your-org/vaultpipe/internal/dotenv"
)

func TestWriteAndReadPinFile_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pins.json")

	pins := []dotenv.PinRecord{
		{Key: "DB_HOST", Value: "localhost", PinnedAt: time.Now().UTC(), Reason: "local dev"},
		{Key: "API_KEY", Value: "abc123", PinnedAt: time.Now().UTC()},
	}

	if err := dotenv.WritePinFile(path, pins); err != nil {
		t.Fatalf("WritePinFile: %v", err)
	}

	got, err := dotenv.ReadPinFile(path)
	if err != nil {
		t.Fatalf("ReadPinFile: %v", err)
	}
	if len(got) != len(pins) {
		t.Fatalf("expected %d pins, got %d", len(pins), len(got))
	}
	if got[0].Key != "DB_HOST" || got[0].Value != "localhost" {
		t.Errorf("unexpected pin[0]: %+v", got[0])
	}
}

func TestReadPinFile_MissingFile(t *testing.T) {
	pins, err := dotenv.ReadPinFile("/nonexistent/pins.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(pins) != 0 {
		t.Errorf("expected empty slice, got %d pins", len(pins))
	}
}

func TestPinRecord_IsExpired(t *testing.T) {
	expired := dotenv.PinRecord{
		Key:       "OLD",
		Value:     "v",
		PinnedAt:  time.Now().Add(-2 * time.Hour),
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	if !expired.IsExpired() {
		t.Error("expected pin to be expired")
	}

	active := dotenv.PinRecord{
		Key:       "NEW",
		Value:     "v",
		PinnedAt:  time.Now(),
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}
	if active.IsExpired() {
		t.Error("expected pin to be active")
	}

	neverExpires := dotenv.PinRecord{Key: "K", Value: "v", PinnedAt: time.Now()}
	if neverExpires.IsExpired() {
		t.Error("expected zero-expiry pin to never expire")
	}
}

func TestApplyPins_OverridesSecrets(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "prod-db", "PORT": "5432"}
	pins := []dotenv.PinRecord{
		{Key: "DB_HOST", Value: "localhost", PinnedAt: time.Now()},
	}

	out, applied := dotenv.ApplyPins(secrets, pins)
	if applied != 1 {
		t.Errorf("expected 1 applied, got %d", applied)
	}
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected pinned value 'localhost', got %q", out["DB_HOST"])
	}
	if out["PORT"] != "5432" {
		t.Errorf("expected PORT unchanged, got %q", out["PORT"])
	}
	if secrets["DB_HOST"] != "prod-db" {
		t.Error("ApplyPins must not mutate original secrets map")
	}
}

func TestApplyPins_SkipsExpired(t *testing.T) {
	secrets := map[string]string{"KEY": "original"}
	pins := []dotenv.PinRecord{
		{Key: "KEY", Value: "pinned", PinnedAt: time.Now().Add(-2 * time.Hour), ExpiresAt: time.Now().Add(-1 * time.Hour)},
	}

	out, applied := dotenv.ApplyPins(secrets, pins)
	if applied != 0 {
		t.Errorf("expected 0 applied, got %d", applied)
	}
	if out["KEY"] != "original" {
		t.Errorf("expected original value, got %q", out["KEY"])
	}
}

func TestPinPath_Convention(t *testing.T) {
	got := dotenv.PinPath("/app/.env")
	want := "/app/..env.pins.json"
	if got != want {
		t.Errorf("PinPath = %q, want %q", got, want)
	}
}

func TestWritePinFile_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pins.json")
	if err := dotenv.WritePinFile(path, []dotenv.PinRecord{}); err != nil {
		t.Fatalf("WritePinFile: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("expected permissions 0600, got %o", perm)
	}
}
