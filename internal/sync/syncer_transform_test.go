package sync

import (
	"testing"

	"github.com/your-org/vaultpipe/internal/config"
)

func baseTransformConfig() *config.Config {
	cfg := newTestConfig()
	return cfg
}

func TestApplyTransform_NoOp(t *testing.T) {
	cfg := baseTransformConfig()
	secrets := map[string]string{"key": " val "}
	out := applyTransform(secrets, cfg)
	if out["key"] != " val " {
		t.Errorf("expected unchanged, got %q", out["key"])
	}
}

func TestApplyTransform_Uppercase(t *testing.T) {
	cfg := baseTransformConfig()
	cfg.TransformUppercase = true
	secrets := map[string]string{"db_host": "localhost"}
	out := applyTransform(secrets, cfg)
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("expected DB_HOST key")
	}
}

func TestApplyTransform_Prefix(t *testing.T) {
	cfg := baseTransformConfig()
	cfg.TransformPrefix = "APP_"
	secrets := map[string]string{"HOST": "localhost"}
	out := applyTransform(secrets, cfg)
	if _, ok := out["APP_HOST"]; !ok {
		t.Error("expected APP_HOST")
	}
}

func TestApplyTransform_TrimValues(t *testing.T) {
	cfg := baseTransformConfig()
	cfg.TransformTrimValues = true
	secrets := map[string]string{"KEY": "  trimmed  "}
	out := applyTransform(secrets, cfg)
	if out["KEY"] != "trimmed" {
		t.Errorf("expected trimmed, got %q", out["KEY"])
	}
}
