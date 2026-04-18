package sync

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/yourusername/vaultpipe/internal/dotenv"
)

// writeTemplate generates a .env.template file alongside the output file
// if the config enables it. The template path is derived from the output path.
func writeTemplate(outputPath string, secrets map[string]string) error {
	tmplPath := templatePath(outputPath)
	if err := dotenv.GenerateTemplate(tmplPath, secrets); err != nil {
		return fmt.Errorf("generate template %s: %w", tmplPath, err)
	}
	return nil
}

// templatePath derives the template file path from the env output path.
// e.g. ".env" -> ".env.template", "config/.env.prod" -> "config/.env.prod.template"
func templatePath(outputPath string) string {
	dir := filepath.Dir(outputPath)
	base := filepath.Base(outputPath)

	// Strip known suffixes before appending .template
	for _, suffix := range []string{".env", ".local"} {
		if strings.HasSuffix(base, suffix) {
			base = strings.TrimSuffix(base, suffix)
			base += ".env.template"
			return filepath.Join(dir, base)
		}
	}
	return filepath.Join(dir, base+".template")
}
