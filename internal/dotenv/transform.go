package dotenv

import (
	"strings"
)

// TransformOptions controls how secret values are transformed before writing.
type TransformOptions struct {
	UppercaseKeys bool
	TrimValues    bool
	Prefix        string
}

// Transform applies transformations to a secrets map and returns a new map.
func Transform(secrets map[string]string, opts TransformOptions) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if opts.TrimValues {
			v = strings.TrimSpace(v)
		}
		if opts.UppercaseKeys {
			k = strings.ToUpper(k)
		}
		if opts.Prefix != "" {
			k = opts.Prefix + k
		}
		out[k] = v
	}
	return out
}
