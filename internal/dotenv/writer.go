package dotenv

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// Writer writes secrets to a .env file.
type Writer struct {
	OutputPath string
}

// NewWriter creates a new Writer targeting the given file path.
func NewWriter(outputPath string) *Writer {
	return &Writer{OutputPath: outputPath}
}

// Write serialises the provided secrets map into KEY=VALUE lines
// and writes them to the configured output file, replacing any
// existing content.
func (w *Writer) Write(secrets map[string]string) error {
	if len(secrets) == 0 {
		return fmt.Errorf("dotenv: no secrets to write")
	}

	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		v := secrets[k]
		// Quote value if it contains spaces or special characters.
		if strings.ContainsAny(v, " \t\n#") {
			v = fmt.Sprintf("%q", v)
		}
		sb.WriteString(fmt.Sprintf("%s=%s\n", k, v))
	}

	f, err := os.OpenFile(w.OutputPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("dotenv: open %s: %w", w.OutputPath, err)
	}
	defer f.Close()

	if _, err := f.WriteString(sb.String()); err != nil {
		return fmt.Errorf("dotenv: write %s: %w", w.OutputPath, err)
	}
	return nil
}
