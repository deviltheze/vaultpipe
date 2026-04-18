package dotenv

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"time"
)

// WatchConfig holds configuration for watching a .env file.
type WatchConfig struct {
	Path     string
	Interval time.Duration
	OnChange func(path string)
}

// Watch polls the given .env file for changes and calls OnChange when the
// file's content hash differs from the previously observed hash.
// It blocks until ctx is cancelled.
func Watch(ctx context.Context, cfg WatchConfig) error {
	if cfg.Interval <= 0 {
		cfg.Interval = 5 * time.Second
	}

	last, err := hashFile(cfg.Path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("watch: initial hash: %w", err)
	}

	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			current, err := hashFile(cfg.Path)
			if err != nil {
				if os.IsNotExist(err) {
					continue
				}
				return fmt.Errorf("watch: hash: %w", err)
			}
			if current != last {
				last = current
				if cfg.OnChange != nil {
					cfg.OnChange(cfg.Path)
				}
			}
		}
	}
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
