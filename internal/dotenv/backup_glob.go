package dotenv

import (
	"fmt"
	"path/filepath"
	"sort"
)

// backupFiles returns sorted backup paths for the given source file.
func backupFiles(src string) ([]string, error) {
	pattern := fmt.Sprintf("%s.*.bak", src)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("backup glob %s: %w", pattern, err)
	}
	sort.Strings(matches)
	return matches, nil
}
