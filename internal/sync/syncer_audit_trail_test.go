package sync

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpipe/internal/dotenv"
)

func TestBuildAuditTrail_WritesEvents(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, ".env")
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	diffs := []dotenv.DiffEntry{
		{Key: "DB_HOST", Change: dotenv.ChangeAdded},
		{Key: "API_KEY", Change: dotenv.ChangeUpdated},
		{Key: "OLD_VAR", Change: dotenv.ChangeRemoved},
	}

	if err := buildAuditTrail(output, "secret/app", "ci-bot", diffs, logger); err != nil {
		t.Fatalf("buildAuditTrail: %v", err)
	}

	trail, err := dotenv.ReadAuditTrail(output)
	if err != nil {
		t.Fatalf("ReadAuditTrail: %v", err)
	}

	if len(trail.Events) != 3 {
		t.Fatalf("events len: got %d, want 3", len(trail.Events))
	}

	actions := map[string]string{}
	for _, e := range trail.Events {
		actions[e.Key] = e.Action
	}

	if actions["DB_HOST"] != string(dotenv.ChangeAdded) {
		t.Errorf("DB_HOST action: got %q, want %q", actions["DB_HOST"], dotenv.ChangeAdded)
	}
	if actions["API_KEY"] != string(dotenv.ChangeUpdated) {
		t.Errorf("API_KEY action: got %q", actions["API_KEY"])
	}
	if actions["OLD_VAR"] != string(dotenv.ChangeRemoved) {
		t.Errorf("OLD_VAR action: got %q", actions["OLD_VAR"])
	}
}

func TestBuildAuditTrail_EmptyDiffs(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, ".env")
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	if err := buildAuditTrail(output, "secret/app", "", nil, logger); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	trail, _ := dotenv.ReadAuditTrail(output)
	if len(trail.Events) != 0 {
		t.Errorf("expected 0 events, got %d", len(trail.Events))
	}
}

func TestBuildAuditTrail_ActorPropagated(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, ".env")
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	diffs := []dotenv.DiffEntry{
		{Key: "TOKEN", Change: dotenv.ChangeAdded},
	}

	_ = buildAuditTrail(output, "secret/svc", "deploy-pipeline", diffs, logger)

	trail, _ := dotenv.ReadAuditTrail(output)
	if trail.Events[0].Actor != "deploy-pipeline" {
		t.Errorf("Actor: got %q, want deploy-pipeline", trail.Events[0].Actor)
	}
	if trail.Events[0].Source != "secret/svc" {
		t.Errorf("Source: got %q, want secret/svc", trail.Events[0].Source)
	}
}
