package dotenv

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// RollbackResult describes the outcome of a rollback operation.
type RollbackResult struct {
	RestoredFrom string
	RestoredTo   string
	Timestamp    time.Time
}

// Rollback restores the most recent backup of the given env file.
// It returns an error if no backup is found or the restore fails.
func Rollback(envPath string) (*RollbackResult, error) {
	files, err := backupFiles(envPath)
	if err != nil {
		return nil, fmt.Errorf("rollback: listing backups: %w", err)
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("rollback: no backups found for %s", envPath)
	}

	// Sort descending so the newest backup is first.
	sort.Slice(files, func(i, j int) bool {
		return files[i] > files[j]
	})

	src := files[0]

	data, err := os.ReadFile(src)
	if err != nil {
		return nil, fmt.Errorf("rollback: reading backup %s: %w", src, err)
	}

	if err := os.WriteFile(envPath, data, 0o600); err != nil {
		return nil, fmt.Errorf("rollback: writing %s: %w", envPath, err)
	}

	ts := parseBackupTimestamp(src)

	return &RollbackResult{
		RestoredFrom: src,
		RestoredTo:   envPath,
		Timestamp:    ts,
	}, nil
}

// parseBackupTimestamp extracts the timestamp embedded in a backup filename.
// Backup files follow the pattern: <base>.backup.<timestamp>.
// Returns zero time if parsing fails.
func parseBackupTimestamp(path string) time.Time {
	base := filepath.Base(path)
	const marker = ".backup."
	idx := strings.LastIndex(base, marker)
	if idx == -1 {
		return time.Time{}
	}
	raw := base[idx+len(marker):]
	t, err := time.Parse("20060102150405", raw)
	if err != nil {
		return time.Time{}
	}
	return t
}
