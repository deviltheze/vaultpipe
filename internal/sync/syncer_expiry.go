package sync

import (
	"fmt"
	"time"

	"github.com/yourusername/vaultpipe/internal/dotenv"
)

// checkExpiry reads the expiry record for the output file and returns an
// error if the file has not yet expired, signalling that a re-sync is
// unnecessary. When no expiry file exists the sync proceeds normally.
func checkExpiry(outputPath string) error {
	rec, err := dotenv.ReadExpiryFile(outputPath)
	if err != nil {
		return fmt.Errorf("expiry: read: %w", err)
	}
	if !rec.ExpiresAt.IsZero() && !rec.IsExpired() {
		remaining := time.Until(rec.ExpiresAt).Round(time.Second)
		return fmt.Errorf("expiry: secrets are still fresh (expires in %s); skipping sync", remaining)
	}
	return nil
}

// writeExpiry persists a new expiry record for outputPath using the
// configured TTL. A zero TTL writes a never-expiring record.
func writeExpiry(outputPath string, ttl time.Duration) error {
	if err := dotenv.WriteExpiryFile(outputPath, ttl); err != nil {
		return fmt.Errorf("expiry: write: %w", err)
	}
	return nil
}
