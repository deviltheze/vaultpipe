package sync

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/config"
)

var filterSecrets = map[string]string{
	"APP_KEY":  "abc",
	"APP_HOST": "localhost",
	"DB_PASS":  "secret",
	"DEBUG":    "1",
}

func TestApplyFilter_NoFilter_ReturnsAll(t *testing.T) {
	out := applyFilter(filterSecrets, config.Filter{})
	if len(out) != len(filterSecrets) {
		t.Fatalf("expected all %d keys, got %d", len(filterSecrets), len(out))
	}
}

func TestApplyFilter_IncludePrefix(t *testing.T) {
	out := applyFilter(filterSecrets, config.Filter{IncludePrefix: []string{"APP_"}})
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["DB_PASS"]; ok {
		t.Error("DB_PASS should not be included")
	}
}

func TestApplyFilter_ExcludePrefix(t *testing.T) {
	out := applyFilter(filterSecrets, config.Filter{ExcludePrefix: []string{"DB_"}})
	if _, ok := out["DB_PASS"]; ok {
		t.Error("DB_PASS should be excluded")
	}
}

func TestApplyFilter_Keys(t *testing.T) {
	out := applyFilter(filterSecrets, config.Filter{Keys: []string{"DEBUG"}})
	if len(out) != 1 {
		t.Fatalf("expected 1 key, got %d", len(out))
	}
	if out["DEBUG"] != "1" {
		t.Error("expected DEBUG=1")
	}
}
