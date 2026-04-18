package dotenv

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// RotateOptions controls rotation behaviour.
type RotateOptions struct {
	// MaxBackups is the number of rotated files to keep (0 = keep all).
	MaxBackups int
}

// Rotate renames the current env file to a timestamped copy, then writes
// newSecrets as the new env file. The caller is responsible for validation.
func Rotate(envPath string, newSecrets map[string]string, opts RotateOptions) error {
	if _, err := os.Stat(envPath); err == nil {
		timestamp := time.Now().UTC().Format("20060102T150405Z")
		ext := filepath.Ext(envPath)
		base := envPath[:len(envPath)-len(ext)]
		rotated := fmt.Sprintf("%s.%s%s", base, timestamp, ext)
		if err := os.Rename(envPath, rotated); err != nil {
			return fmt.Errorf("rotate: rename current file: %w", err)
		}
	}

	w, err := NewWriter(envPath)
	if err != nil {
		return fmt.Errorf("rotate: create writer: %w", err)
	}
	if err := w.Write(newSecrets); err != nil {
		return fmt.Errorf("rotate: write secrets: %w", err)
	}

	if opts.MaxBackups > 0 {
		if err := PruneBackups(envPath, opts.MaxBackups); err != nil {
			return fmt.Errorf("rotate: prune backups: %w", err)
		}
	}
	return nil
}
