package dotenv

import "strings"

// FilterOptions controls which keys are included or excluded.
type FilterOptions struct {
	IncludePrefix []string
	ExcludePrefix []string
	Keys          []string // if set, only these keys are kept
}

// Filter returns a subset of secrets based on FilterOptions.
func Filter(secrets map[string]string, opts FilterOptions) map[string]string {
	out := make(map[string]string)

	for k, v := range secrets {
		if len(opts.Keys) > 0 && !containsStr(opts.Keys, k) {
			continue
		}
		if len(opts.IncludePrefix) > 0 && !hasAnyPrefix(k, opts.IncludePrefix) {
			continue
		}
		if hasAnyPrefix(k, opts.ExcludePrefix) {
			continue
		}
		out[k] = v
	}
	return out
}

func hasAnyPrefix(s string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}

func containsStr(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
