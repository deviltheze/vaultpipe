package dotenv_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/dotenv"
)

var baseSecrets = map[string]string{
	"APP_HOST":    "localhost",
	"APP_PORT":    "8080",
	"DB_HOST":     "db.local",
	"DB_PASSWORD": "secret",
	"DEBUG":       "true",
}

func TestFilter_NoOptions(t *testing.T) {
	out := dotenv.Filter(baseSecrets, dotenv.FilterOptions{})
	if len(out) != len(baseSecrets) {
		t.Fatalf("expected %d keys, got %d", len(baseSecrets), len(out))
	}
}

func TestFilter_IncludePrefix(t *testing.T) {
	out := dotenv.Filter(baseSecrets, dotenv.FilterOptions{IncludePrefix: []string{"APP_"}})
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
	if _, ok := out["APP_HOST"]; !ok {
		t.Error("expected APP_HOST")
	}
}

func TestFilter_ExcludePrefix(t *testing.T) {
	out := dotenv.Filter(baseSecrets, dotenv.FilterOptions{ExcludePrefix: []string{"DB_"}})
	if _, ok := out["DB_HOST"]; ok {
		t.Error("DB_HOST should be excluded")
	}
	if _, ok := out["APP_HOST"]; !ok {
		t.Error("APP_HOST should remain")
	}
}

func TestFilter_SpecificKeys(t *testing.T) {
	out := dotenv.Filter(baseSecrets, dotenv.FilterOptions{Keys: []string{"DEBUG", "DB_HOST"}})
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
}

func TestFilter_IncludeAndExcludeCombined(t *testing.T) {
	out := dotenv.Filter(baseSecrets, dotenv.FilterOptions{
		IncludePrefix: []string{"DB_"},
		ExcludePrefix: []string{"DB_PASSWORD"},
	})
	if _, ok := out["DB_PASSWORD"]; ok {
		t.Error("DB_PASSWORD should be excluded")
	}
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("DB_HOST should be included")
	}
}

func TestFilter_EmptySecrets(t *testing.T) {
	out := dotenv.Filter(map[string]string{}, dotenv.FilterOptions{IncludePrefix: []string{"APP_"}})
	if len(out) != 0 {
		t.Fatalf("expected empty map")
	}
}
