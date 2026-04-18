package audit_test

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpipe/internal/audit"
)

func TestNewLogger_Stderr(t *testing.T) {
	l, err := audit.NewLogger("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestNewLogger_File(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "audit.log")
	l, err := audit.NewLogger(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestNewLogger_BadPath(t *testing.T) {
	_, err := audit.NewLogger("/nonexistent/dir/audit.log")
	if err == nil {
		t.Fatal("expected error for bad path")
	}
}

func TestLogSync_WritesJSON(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "audit.log")
	l, _ := audit.NewLogger(tmp)

	if err := l.LogSync("secret/myapp", ".env", 5); err != nil {
		t.Fatalf("LogSync error: %v", err)
	}

	f, _ := os.Open(tmp)
	defer f.Close()

	var entry audit.Entry
	if err := json.NewDecoder(bufio.NewReader(f)).Decode(&entry); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if entry.Event != "sync_success" {
		t.Errorf("expected sync_success, got %s", entry.Event)
	}
	if entry.KeyCount != 5 {
		t.Errorf("expected key_count 5, got %d", entry.KeyCount)
	}
	if entry.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestLogError_WritesJSON(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "audit.log")
	l, _ := audit.NewLogger(tmp)

	if err := l.LogError("secret/myapp", errors.New("vault unreachable")); err != nil {
		t.Fatalf("LogError error: %v", err)
	}

	f, _ := os.Open(tmp)
	defer f.Close()

	var entry audit.Entry
	json.NewDecoder(f).Decode(&entry)

	if entry.Event != "sync_error" {
		t.Errorf("expected sync_error, got %s", entry.Event)
	}
	if entry.Error != "vault unreachable" {
		t.Errorf("unexpected error string: %s", entry.Error)
	}
}
