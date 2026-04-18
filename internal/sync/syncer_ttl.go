package sync

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/yourusername/vaultpipe/internal/dotenv"
)

const ttlFileName = ".vaultpipe_ttl.json"

// ttlPath returns the TTL file path alongside the output env file.
func ttlPath(outputFile string) string {
	return filepath.Join(filepath.Dir(outputFile), ttlFileName)
}

// checkTTL returns an error if the last sync is still fresh (not expired).
// Returns nil (proceed) when expired or no TTL file exists.
func checkTTL(outputFile string) error {
	p := ttlPath(outputFile)
	rec, err := dotenv.ReadTTLFile(p)
	if err != nil {
		// No TTL file — first run, proceed.
		return nil
	}
	if !rec.IsExpired() {
		remaining := rec.TTL - time.Since(rec.SyncedAt)
		return fmt.Errorf("ttl: secrets are still fresh (%.0fs remaining), skipping sync", remaining.Seconds())
	}
	return nil
}

// writeTTL persists a new TTL record after a successful sync.
func writeTTL(outputFile, secretPath string, ttl time.Duration) error {
	rec := dotenv.TTLRecord{
		Path:     secretPath,
		SyncedAt: time.Now(),
		TTL:      ttl,
	}
	return dotenv.WriteTTLFile(ttlPath(outputFile), rec)
}
