package dotenv

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// LockFile is the name of the lock file created alongside the .env file.
const lockFileSuffix = ".lock"

// ErrLocked is returned when a lock file already exists and is still valid.
var ErrLocked = errors.New("dotenv: file is locked by another process")

// AcquireLock creates a lock file next to path. Returns ErrLocked if already locked.
func AcquireLock(path string) error {
	lockPath := lockPath(path)

	if data, err := os.ReadFile(lockPath); err == nil {
		parts := strings.SplitN(strings.TrimSpace(string(data)), "\n", 2)
		if len(parts) == 2 {
			ts, err := strconv.ParseInt(parts[1], 10, 64)
			if err == nil && time.Now().Unix()-ts < 30 {
				return fmt.Errorf("%w (pid %s)", ErrLocked, parts[0])
			}
		}
		// stale lock — remove it
		_ = os.Remove(lockPath)
	}

	content := fmt.Sprintf("%d\n%d\n", os.Getpid(), time.Now().Unix())
	if err := os.WriteFile(lockPath, []byte(content), 0600); err != nil {
		return fmt.Errorf("dotenv: could not write lock file: %w", err)
	}
	return nil
}

// ReleaseLock removes the lock file for path.
func ReleaseLock(path string) error {
	lp := lockPath(path)
	if err := os.Remove(lp); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("dotenv: could not release lock: %w", err)
	}
	return nil
}

func lockPath(path string) string {
	return filepath.Join(filepath.Dir(path), filepath.Base(path)+lockFileSuffix)
}
