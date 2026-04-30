package sync

import (
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpipe/internal/config"
	"github.com/your-org/vaultpipe/internal/dotenv"
	"log/slog"
	"os"
)

func newHookLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, nil))
}

func TestRunHook_NoCommand_RecordsEvent(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{
		OutputFile:   filepath.Join(dir, ".env"),
		Actor:        "tester",
		PreSyncHook:  "",
		PostSyncHook: "",
	}

	if err := runHook(cfg, dotenv.HookPreSync, newHookLogger()); err != nil {
		t.Fatalf("runHook: %v", err)
	}

	hf, err := dotenv.ReadHookFile(dotenv.HookPath(cfg.OutputFile))
	if err != nil {
		t.Fatalf("ReadHookFile: %v", err)
	}
	if len(hf.Records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(hf.Records))
	}
	if hf.Records[0].Event != dotenv.HookPreSync {
		t.Errorf("event mismatch: %s", hf.Records[0].Event)
	}
	if hf.Records[0].Actor != "tester" {
		t.Errorf("actor mismatch: %s", hf.Records[0].Actor)
	}
}

func TestRunHook_PostSync_RecordsEvent(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{
		OutputFile:   filepath.Join(dir, ".env"),
		Actor:        "ci",
		PostSyncHook: "",
	}

	if err := runHook(cfg, dotenv.HookPostSync, newHookLogger()); err != nil {
		t.Fatalf("runHook post: %v", err)
	}

	hf, _ := dotenv.ReadHookFile(dotenv.HookPath(cfg.OutputFile))
	if len(hf.Records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(hf.Records))
	}
	if hf.Records[0].Event != dotenv.HookPostSync {
		t.Errorf("wrong event: %s", hf.Records[0].Event)
	}
}

func TestRunHook_WithEchoCommand_CapturesOutput(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{
		OutputFile:  filepath.Join(dir, ".env"),
		Actor:       "auto",
		PreSyncHook: "echo hello",
	}

	if err := runHook(cfg, dotenv.HookPreSync, newHookLogger()); err != nil {
		t.Fatalf("runHook: %v", err)
	}

	hf, _ := dotenv.ReadHookFile(dotenv.HookPath(cfg.OutputFile))
	if hf.Records[0].Output != "hello" {
		t.Errorf("expected output 'hello', got %q", hf.Records[0].Output)
	}
}

func TestRunHook_FailingCommand_RecordsError(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{
		OutputFile:   filepath.Join(dir, ".env"),
		Actor:        "auto",
		PostSyncHook: "false",
	}

	// Should not return error — hook failures are logged, not fatal.
	if err := runHook(cfg, dotenv.HookPostSync, newHookLogger()); err != nil {
		t.Fatalf("runHook should not fail on hook error: %v", err)
	}

	hf, _ := dotenv.ReadHookFile(dotenv.HookPath(cfg.OutputFile))
	if hf.Records[0].Error == "" {
		t.Error("expected error to be recorded")
	}
}
