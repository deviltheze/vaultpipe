package dotenv

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Annotation holds metadata attached to a specific secret key.
type Annotation struct {
	Key       string    `json:"key"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AnnotationFile is the on-disk representation of all annotations.
type AnnotationFile struct {
	Annotations map[string]Annotation `json:"annotations"`
}

// AnnotationPath returns the conventional path for the annotation file
// relative to the given env file path.
func AnnotationPath(envPath string) string {
	dir := filepath.Dir(envPath)
	base := filepath.Base(envPath)
	return filepath.Join(dir, "."+base+".annotations.json")
}

// WriteAnnotations persists annotations for the given env file.
// Existing annotations not present in the provided map are preserved.
func WriteAnnotations(envPath string, incoming map[string]Annotation) error {
	path := AnnotationPath(envPath)

	existing, err := ReadAnnotations(envPath)
	if err != nil {
		return fmt.Errorf("annotate: read existing: %w", err)
	}

	now := time.Now().UTC()
	for k, ann := range incoming {
		if prev, ok := existing[k]; ok {
			ann.CreatedAt = prev.CreatedAt
		} else {
			ann.CreatedAt = now
		}
		ann.UpdatedAt = now
		ann.Key = k
		existing[k] = ann
	}

	af := AnnotationFile{Annotations: existing}
	data, err := json.MarshalIndent(af, "", "  ")
	if err != nil {
		return fmt.Errorf("annotate: marshal: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}

// ReadAnnotations loads annotations for the given env file.
// Returns an empty map if the annotation file does not exist.
func ReadAnnotations(envPath string) (map[string]Annotation, error) {
	path := AnnotationPath(envPath)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return make(map[string]Annotation), nil
	}
	if err != nil {
		return nil, fmt.Errorf("annotate: read file: %w", err)
	}
	var af AnnotationFile
	if err := json.Unmarshal(data, &af); err != nil {
		return nil, fmt.Errorf("annotate: unmarshal: %w", err)
	}
	if af.Annotations == nil {
		af.Annotations = make(map[string]Annotation)
	}
	return af.Annotations, nil
}
