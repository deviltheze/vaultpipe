package sync

import (
	"github.com/yourusername/vaultpipe/internal/config"
	"github.com/yourusername/vaultpipe/internal/dotenv"
)

// applyFilter applies config-defined filter options to the fetched secrets.
func applyFilter(secrets map[string]string, f config.Filter) map[string]string {
	if len(f.IncludePrefix) == 0 && len(f.ExcludePrefix) == 0 && len(f.Keys) == 0 {
		return secrets
	}
	return dotenv.Filter(secrets, dotenv.FilterOptions{
		IncludePrefix: f.IncludePrefix,
		ExcludePrefix: f.ExcludePrefix,
		Keys:          f.Keys,
	})
}
