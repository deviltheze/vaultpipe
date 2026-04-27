package sync

import (
	"log/slog"
	"time"

	"github.com/yourusername/vaultpipe/internal/dotenv"
)

// buildAuditTrail converts a diff result into an AuditTrail and writes it.
func buildAuditTrail(outputFile, source, actor string, diffs []dotenv.DiffEntry, logger *slog.Logger) error {
	events := make([]dotenv.AuditEvent, 0, len(diffs))
	now := time.Now().UTC()

	for _, d := range diffs {
		events = append(events, dotenv.AuditEvent{
			Timestamp: now,
			Key:       d.Key,
			Action:    string(d.Change),
			Source:    source,
			Actor:     actor,
		})
	}

	trail := dotenv.AuditTrail{
		SyncedAt: now,
		Events:   events,
	}

	if err := dotenv.WriteAuditTrail(outputFile, trail); err != nil {
		logger.Warn("audit_trail: failed to write", "error", err)
		return err
	}

	logger.Info("audit_trail: written",
		"output", outputFile,
		"events", len(events),
	)
	return nil
}
