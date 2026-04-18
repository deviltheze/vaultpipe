package dotenv

import "strings"

// MaskOptions controls how values are masked in output.
type MaskOptions struct {
	// ShowChars is the number of trailing characters to reveal. 0 means fully masked.
	ShowChars int
	// MaskChar is the character used for masking. Defaults to '*'.
	MaskChar rune
}

var defaultMaskOptions = MaskOptions{
	ShowChars: 4,
	MaskChar:  '*',
}

// Mask returns a copy of secrets where sensitive values are partially masked.
// Non-sensitive values are returned as-is.
func Mask(secrets map[string]string, opts *MaskOptions) map[string]string {
	if opts == nil {
		opts = &defaultMaskOptions
	}
	mc := opts.MaskChar
	if mc == 0 {
		mc = '*'
	}

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if isSensitive(k) {
			out[k] = maskValue(v, opts.ShowChars, mc)
		} else {
			out[k] = v
		}
	}
	return out
}

func maskValue(v string, showChars int, mc rune) string {
	if len(v) == 0 {
		return ""
	}
	if showChars <= 0 || showChars >= len(v) {
		return strings.Repeat(string(mc), len(v))
	}
	masked := strings.Repeat(string(mc), len(v)-showChars)
	return masked + v[len(v)-showChars:]
}
