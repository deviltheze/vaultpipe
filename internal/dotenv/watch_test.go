package dotenv

import (
	"context"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"
)

func TestWatch_DetectsChange(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	if err := os.WriteFile(path, []byte("FOO=bar\n"), 0600); err != nil {
		t.Fatal(err)
	}

	var called atomic.Int32
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cfg := WatchConfig{
		Path:     path,
		Interval: 50 * time.Millisecond,
		OnChange: func(_ string) { called.Add(1); cancel() },
	}

	go func() {
		time.Sleep(100 * time.Millisecond)
		_ = os.WriteFile(path, []byte("FOO=changed\n"), 0600)
	}()

	_ = Watch(ctx, cfg)

	if called.Load() == 0 {
		t.Error("expected OnChange to be called")
	}
}

func TestWatch_NoChangeNoCallback(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	if err := os.WriteFile(path, []byte("FOO=bar\n"), 0600); err != nil {
		t.Fatal(err)
	}

	var called atomic.Int32
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	cfg := WatchConfig{
		Path:     path,
		Interval: 50 * time.Millisecond,
		OnChange: func(_ string) { called.Add(1) },
	}

	_ = Watch(ctx, cfg)

	if called.Load() != 0 {
		t.Errorf("expected no OnChange calls, got %d", called.Load())
	}
}

func TestWatch_MissingFileDoesNotError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "missing.env")

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	cfg := WatchConfig{
		Path:     path,
		Interval: 50 * time.Millisecond,
	}

	if err := Watch(ctx, cfg); err != nil && err != context.DeadlineExceeded {
		t.Errorf("unexpected error: %v", err)
	}
}
