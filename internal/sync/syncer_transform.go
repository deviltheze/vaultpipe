package sync

import (
	"github.com/your-org/vaultpipe/internal/dotenv"
	"github.com/your-org/vaultpipe/internal/config"
)

// applyTransform applies configured key/value transformations to secrets.
func applyTransform(secrets map[string]string, cfg *config.Config) map[string]string {
	opts := dotenv.TransformOptions{
		UppercaseKeys: cfg.TransformUppercase,
		TrimValues:    cfg.TransformTrimValues,
		Prefix:        cfg.TransformPrefix,
	}
	if !opts.UppercaseKeys && !opts.TrimValues && opts.Prefix == "" {
		return secrets
	}
	return dotenv.Transform(secrets, opts)
}
