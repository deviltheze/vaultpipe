package dotenv

import (
	"testing"
)

func TestMask_SensitivePartiallyMasked(t *testing.T) {
	secrets := map[string]string{
		"DB_PASSWORD": "supersecret",
	}
	out := Mask(secrets, nil)
	got := out["DB_PASSWORD"]
	// default ShowChars=4, value len=11 => 7 stars + last 4 chars
	want := "*******cret"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMask_NonSensitiveUnchanged(t *testing.T) {
	secrets := map[string]string{
		"APP_ENV": "production",
	}
	out := Mask(secrets, nil)
	if out["APP_ENV"] != "production" {
		t.Errorf("expected non-sensitive key to be unchanged")
	}
}

func TestMask_FullyMaskedWhenShowCharsZero(t *testing.T) {
	secrets := map[string]string{
		"API_KEY": "abc123",
	}
	out := Mask(secrets, &MaskOptions{ShowChars: 0, MaskChar: '#'})
	if out["API_KEY"] != "######" {
		t.Errorf("expected fully masked, got %q", out["API_KEY"])
	}
}

func TestMask_EmptyValue(t *testing.T) {
	secrets := map[string]string{
		"DB_PASSWORD": "",
	}
	out := Mask(secrets, nil)
	if out["DB_PASSWORD"] != "" {
		t.Errorf("expected empty string, got %q", out["DB_PASSWORD"])
	}
}

func TestMask_DoesNotMutateOriginal(t *testing.T) {
	secrets := map[string]string{
		"SECRET_TOKEN": "original",
	}
	Mask(secrets, nil)
	if secrets["SECRET_TOKEN"] != "original" {
		t.Error("original map was mutated")
	}
}

func TestMask_CustomMaskChar(t *testing.T) {
	secrets := map[string]string{
		"DB_PASS": "hello",
	}
	out := Mask(secrets, &MaskOptions{ShowChars: 2, MaskChar: 'x'})
	if out["DB_PASS"] != "xxxlo" {
		t.Errorf("got %q, want %q", out["DB_PASS"], "xxxlo")
	}
}
