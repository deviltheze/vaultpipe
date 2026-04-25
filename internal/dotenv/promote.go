package dotenv

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// PromoteOptions controls how secrets are promoted between environments.
type PromoteOptions struct {
	// SourceEnv is the environment to promote from (e.g. "staging").
	SourceEnv string
	// TargetEnv is the environment to promote to (e.g. "production").
	TargetEnv string
	// Keys restricts promotion to specific keys. Empty means all keys.
	Keys []string
	// DryRun reports what would change without writing.
	DryRun bool
}

// PromoteResult describes the outcome of a promotion.
type PromoteResult struct {
	Promoted []string
	Skipped  []string
	DryRun   bool
}

func (r PromoteResult) String() string {
	var sb strings.Builder
	if r.DryRun {
		sb.WriteString("[dry-run] ")
	}
	fmt.Fprintf(&sb, "promoted=%d skipped=%d", len(r.Promoted), len(r.Skipped))
	return sb.String()
}

// Promote copies selected secrets from a source .env file into a target .env
// file, merging with any existing values in the target.
func Promote(sourceFile, targetFile string, opts PromoteOptions) (PromoteResult, error) {
	src, err := Read(sourceFile)
	if err != nil {
		return PromoteResult{}, fmt.Errorf("promote: read source %q: %w", sourceFile, err)
	}

	var dst map[string]string
	if _, statErr := os.Stat(targetFile); statErr == nil {
		dst, err = Read(targetFile)
		if err != nil {
			return PromoteResult{}, fmt.Errorf("promote: read target %q: %w", targetFile, err)
		}
	} else {
		dst = make(map[string]string)
	}

	keySet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = struct{}{}
	}

	var result PromoteResult
	result.DryRun = opts.DryRun

	for k, v := range src {
		if len(keySet) > 0 {
			if _, ok := keySet[k]; !ok {
				result.Skipped = append(result.Skipped, k)
				continue
			}
		}
		dst[k] = v
		result.Promoted = append(result.Promoted, k)
	}

	sort.Strings(result.Promoted)
	sort.Strings(result.Skipped)

	if opts.DryRun {
		return result, nil
	}

	if err := os.MkdirAll(filepath.Dir(targetFile), 0o755); err != nil {
		return result, fmt.Errorf("promote: mkdir target dir: %w", err)
	}

	w := NewWriter(targetFile)
	if err := w.Write(dst); err != nil {
		return result, fmt.Errorf("promote: write target: %w", err)
	}

	return result, nil
}
