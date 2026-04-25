package sync

import (
	"fmt"
	"log/slog"

	"github.com/your-org/vaultpipe/internal/dotenv"
)

// PromoteConfig holds parameters for an environment promotion run.
type PromoteConfig struct {
	SourceFile string
	TargetFile string
	Keys       []string
	DryRun     bool
}

// RunPromotion promotes secrets from a source .env file to a target .env file.
// It logs the result and returns an error if the promotion fails.
//
// When DryRun is true, no changes are written to disk; the result still
// reflects what would have been promoted or skipped.
func RunPromotion(cfg PromoteConfig, logger *slog.Logger) error {
	if cfg.SourceFile == "" {
		return fmt.Errorf("promote: source file must be specified")
	}
	if cfg.TargetFile == "" {
		return fmt.Errorf("promote: target file must be specified")
	}
	if cfg.SourceFile == cfg.TargetFile {
		return fmt.Errorf("promote: source and target files must be different")
	}

	opts := dotenv.PromoteOptions{
		Keys:   cfg.Keys,
		DryRun: cfg.DryRun,
	}

	result, err := dotenv.Promote(cfg.SourceFile, cfg.TargetFile, opts)
	if err != nil {
		logger.Error("promotion failed",
			"source", cfg.SourceFile,
			"target", cfg.TargetFile,
			"error", err,
		)
		return fmt.Errorf("promote: %w", err)
	}

	logger.Info("promotion complete",
		"source", cfg.SourceFile,
		"target", cfg.TargetFile,
		"promoted", len(result.Promoted),
		"skipped", len(result.Skipped),
		"dry_run", result.DryRun,
		"summary", result.String(),
	)

	return nil
}
