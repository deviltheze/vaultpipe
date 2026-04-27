package sync

import (
	"fmt"
	"strings"

	"github.com/your-org/vaultpipe/internal/dotenv"
)

// enforcePolicy loads the policy file from the output directory and checks
// secrets against it. If any violations are found the sync is aborted and
// the violations are logged before returning an error.
func enforcePolicy(cfg interface{ GetOutputDir() string }, secrets map[string]string, logger interface {
	LogError(path, msg string)
}) error {
	dir := cfg.GetOutputDir()
	pol, err := dotenv.ReadPolicyFile(dir)
	if err != nil {
		return fmt.Errorf("policy: load: %w", err)
	}
	if len(pol.Rules) == 0 {
		return nil
	}
	violations := dotenv.EnforcePolicy(secrets, pol)
	if len(violations) == 0 {
		return nil
	}
	msgs := make([]string, len(violations))
	for i, v := range violations {
		msgs[i] = v.Error()
		logger.LogError(dir, v.Error())
	}
	return fmt.Errorf("policy enforcement failed:\n  %s", strings.Join(msgs, "\n  "))
}
