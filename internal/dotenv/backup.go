package dotenv

import (
	"fmt"
	"os"
	"time"
)

// Backup copies src to src.YYYYMMDDHHMMSS.bak.
// If src does not exist, Backup is a no-op and returns ("", nil).
func Backup(src string) (string, error) {
	data, err := os.ReadFile(src)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("backup: read %s: %w", src, err)
	}

	stamp := time.Now().UTC().Format("20060102150405")
	dst := fmt.Sprintf("%s.%s.bak", src, stamp)

	if err := os.WriteFile(dst, data, 0600); err != nil {
		return "", fmt.Errorf("backup: write %s: %w", dst, err)
	}
	return dst, nil
}

// PruneBackups removes all but the newest `keep` backup files for src.
func PruneBackups(src string, keep int) error {
	matches, err := backupFiles(src)
	if err != nil {
		return err
	}
	if len(matches) <= keep {
		return nil
	}
	// matches are lexicographically sorted; oldest first.
	for _, f := range matches[:len(matches)-keep] {
		if err := os.Remove(f); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("prune: remove %s: %w", f, err)
		}
	}
	return nil
}
