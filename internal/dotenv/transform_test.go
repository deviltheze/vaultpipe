package dotenv

import (
	"testing"
)

func TestTransform_UppercaseKeys(t *testing.T) {
	in := map[string]string{"db_host": "localhost", "api_key": "secret"}
	out := Transform(in, TransformOptions{UppercaseKeys: true})
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("expected DB_HOST")
	}
	if _, ok := out["API_KEY"]; !ok {
		t.Error("expected API_KEY")
	}
}

func TestTransform_TrimValues(t *testing.T) {
	in := map[string]string{"KEY": "  value  "}
	out := Transform(in, TransformOptions{TrimValues: true})
	if out["KEY"] != "value" {
		t.Errorf("expected trimmed value, got %q", out["KEY"])
	}
}

func TestTransform_Prefix(t *testing.T) {
	in := map[string]string{"HOST": "localhost"}
	out := Transform(in, TransformOptions{Prefix: "APP_"})
	if _, ok := out["APP_HOST"]; !ok {
		t.Error("expected APP_HOST")
	}
	if _, ok := out["HOST"]; ok {
		t.Error("original key should not exist")
	}
}

func TestTransform_NoOptions(t *testing.T) {
	in := map[string]string{"key": " val "}
	out := Transform(in, TransformOptions{})
	if out["key"] != " val " {
		t.Errorf("expected unchanged value, got %q", out["key"])
	}
}

func TestTransform_DoesNotMutateOriginal(t *testing.T) {
	in := map[string]string{"key": "value"}
	Transform(in, TransformOptions{UppercaseKeys: true, Prefix: "X_"})
	if _, ok := in["key"]; !ok {
		t.Error("original map should not be mutated")
	}
}

func TestTransform_CombinedOptions(t *testing.T) {
	in := map[string]string{"db_host": "  pg  "}
	out := Transform(in, TransformOptions{UppercaseKeys: true, TrimValues: true, Prefix: "APP_"})
	if out["APP_DB_HOST"] != "pg" {
		t.Errorf("expected APP_DB_HOST=pg, got %v", out)
	}
}
