package dotenv

import (
	"strings"
	"unicode"
)

// SanitizeOptions controls how secret values are sanitized before writing.
type SanitizeOptions struct {
	// StripControlChars removes non-printable control characters from values.
	StripControlChars bool
	// TrimWhitespace trims leading and trailing whitespace from values.
	TrimWhitespace bool
	// MaxValueLength truncates values longer than this. Zero means no limit.
	MaxValueLength int
}

// Sanitize returns a new map with values cleaned according to the given options.
// Keys are never modified. The original map is not mutated.
func Sanitize(secrets map[string]string, opts SanitizeOptions) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if opts.TrimWhitespace {
			v = strings.TrimSpace(v)
		}
		if opts.StripControlChars {
			v = stripControl(v)
		}
		if opts.MaxValueLength > 0 && len(v) > opts.MaxValueLength {
			v = v[:opts.MaxValueLength]
		}
		out[k] = v
	}
	return out
}

// stripControl removes non-printable control characters (except tab and newline
// which are common in multi-line secrets) from s.
func stripControl(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r == '\t' || r == '\n' || r == '\r' {
			b.WriteRune(r)
			continue
		}
		if unicode.IsControl(r) {
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}
