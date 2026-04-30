package sync

import (
	"fmt"
	"log/slog"
	"os/exec"
	"strings"

	"github.com/your-org/vaultpipe/internal/config"
	"github.com/your-org/vaultpipe/internal/dotenv"
)

// runHook executes a shell hook command if configured and records the invocation.
func runHook(cfg *config.Config, event dotenv.HookEvent, logger *slog.Logger) error {
	var cmd string
	switch event {
	case dotenv.HookPreSync:
		cmd = cfg.PreSyncHook
	case dotenv.HookPostSync:
		cmd = cfg.PostSyncHook
	}

	rec := dotenv.HookRecord{
		Event: event,
		Actor: cfg.Actor,
	}

	var output string
	var hookErr string

	if cmd != "" {
		parts := strings.Fields(cmd)
		out, err := exec.Command(parts[0], parts[1:]...).CombinedOutput() //nolint:gosec
		output = strings.TrimSpace(string(out))
		if err != nil {
			hookErr = err.Error()
			logger.Warn("hook command failed",
				"event", event,
				"cmd", cmd,
				"error", err,
				"output", output,
			)
		} else {
			logger.Info("hook executed",
				"event", event,
				"cmd", cmd,
				"output", output,
			)
		}
	}

	rec.Output = output
	rec.Error = hookErr

	hookPath := dotenv.HookPath(cfg.OutputFile)
	if err := dotenv.WriteHookRecord(hookPath, rec); err != nil {
		return fmt.Errorf("syncer_hook: record %s: %w", event, err)
	}
	return nil
}
