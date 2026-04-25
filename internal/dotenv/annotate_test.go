package dotenv

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAnnotationPath_Convention(t *testing.T) {
	got := AnnotationPath("/tmp/envs/.env")
	want := "/tmp/envs/..env.annotations.json"
	if got != want {
		t.Errorf("AnnotationPath = %q, want %q", got, want)
	}
}

func TestWriteAndReadAnnotations_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	incoming := map[string]Annotation{
		"DB_PASSWORD": {Note: "rotated quarterly"},
		"API_KEY":     {Note: "third-party service"},
	}

	if err := WriteAnnotations(envPath, incoming); err != nil {
		t.Fatalf("WriteAnnotations: %v", err)
	}

	got, err := ReadAnnotations(envPath)
	if err != nil {
		t.Fatalf("ReadAnnotations: %v", err)
	}

	for k, want := range incoming {
		ann, ok := got[k]
		if !ok {
			t.Errorf("missing annotation for key %q", k)
			continue
		}
		if ann.Note != want.Note {
			t.Errorf("key %q note = %q, want %q", k, ann.Note, want.Note)
		}
		if ann.Key != k {
			t.Errorf("key field = %q, want %q", ann.Key, k)
		}
		if ann.CreatedAt.IsZero() {
			t.Errorf("key %q: CreatedAt should not be zero", k)
		}
	}
}

func TestReadAnnotations_MissingFile(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	got, err := ReadAnnotations(envPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestWriteAnnotations_PreservesCreatedAt(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	first := map[string]Annotation{
		"SECRET": {Note: "initial"},
	}
	if err := WriteAnnotations(envPath, first); err != nil {
		t.Fatalf("first write: %v", err)
	}

	time.Sleep(2 * time.Millisecond)

	second := map[string]Annotation{
		"SECRET": {Note: "updated note"},
	}
	if err := WriteAnnotations(envPath, second); err != nil {
		t.Fatalf("second write: %v", err)
	}

	got, _ := ReadAnnotations(envPath)
	ann := got["SECRET"]
	if ann.Note != "updated note" {
		t.Errorf("note = %q, want %q", ann.Note, "updated note")
	}
	if !ann.UpdatedAt.After(ann.CreatedAt) {
		t.Errorf("UpdatedAt (%v) should be after CreatedAt (%v)", ann.UpdatedAt, ann.CreatedAt)
	}
}

func TestWriteAnnotations_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	if err := WriteAnnotations(envPath, map[string]Annotation{
		"K": {Note: "test"},
	}); err != nil {
		t.Fatalf("WriteAnnotations: %v", err)
	}

	info, err := os.Stat(AnnotationPath(envPath))
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("file perm = %o, want 0600", perm)
	}
}
