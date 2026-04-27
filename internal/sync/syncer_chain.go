package sync

import (
	"log/slog"

	"github.com/yourusername/vaultpipe/internal/dotenv"
)

// writeChainRecord persists a chain record describing which keys were synced
// from the Vault path to the output .env file.
func writeChainRecord(cfg chainConfig, keys []string, logger *slog.Logger) {
	if len(keys) == 0 {
		return
	}

	path := dotenv.ChainPath(cfg.outputFile)
	rec := dotenv.ChainRecord{
		Source: cfg.vaultPath,
		Dest:   cfg.outputFile,
		Keys:   keys,
	}

	if err := dotenv.WriteChainRecord(path, rec); err != nil {
		logger.Warn("chain: failed to write record",
			"path", path,
			"error", err,
		)
		return
	}

	logger.Debug("chain: record written",
		"path", path,
		"keys", len(keys),
	)
}

// chainConfig carries the minimal fields needed to write a chain record.
type chainConfig struct {
	vaultPath  string
	outputFile string
}

// sortedKeySlice extracts and sorts keys from a secrets map for chain logging.
func sortedKeySlice(secrets map[string]string) []string {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	// stable sort to keep chain records deterministic
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && keys[j] < keys[j-1]; j-- {
			keys[j], keys[j-1] = keys[j-1], keys[j]
		}
	}
	return keys
}
