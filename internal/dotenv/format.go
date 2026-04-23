package dotenv

import (
	"fmt"
	"sort"
	"strings"
)

// FormatOptions controls how secrets are rendered as .env content.
type FormatOptions struct {
	// SortKeys sorts keys alphabetically when true.
	SortKeys bool
	// GroupByPrefix groups keys under comments by their prefix (e.g. DB_, AWS_).
	GroupByPrefix bool
	// PrefixSeparator is the delimiter used to detect prefixes. Defaults to "_".
	PrefixSeparator string
}

// Format renders a map of secrets into .env file content as a string.
// It respects the provided FormatOptions for sorting and grouping.
func Format(secrets map[string]string, opts FormatOptions) string {
	if len(secrets) == 0 {
		return ""
	}

	sep := opts.PrefixSeparator
	if sep == "" {
		sep = "_"
	}

	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	if opts.SortKeys || opts.GroupByPrefix {
		sort.Strings(keys)
	}

	if !opts.GroupByPrefix {
		var sb strings.Builder
		for _, k := range keys {
			fmt.Fprintf(&sb, "%s=%s\n", k, quoteIfNeeded(secrets[k]))
		}
		return sb.String()
	}

	// Group by prefix
	groups := make(map[string][]string)
	order := []string{}
	seen := map[string]bool{}

	for _, k := range keys {
		prefix := groupPrefix(k, sep)
		if !seen[prefix] {
			seen[prefix] = true
			order = append(order, prefix)
		}
		groups[prefix] = append(groups[prefix], k)
	}

	var sb strings.Builder
	for i, prefix := range order {
		if i > 0 {
			sb.WriteString("\n")
		}
		fmt.Fprintf(&sb, "# %s\n", prefix)
		for _, k := range groups[prefix] {
			fmt.Fprintf(&sb, "%s=%s\n", k, quoteIfNeeded(secrets[k]))
		}
	}
	return sb.String()
}

// groupPrefix returns the prefix portion of a key, e.g. "DB" from "DB_HOST".
// If no separator is found, the full key is used as the group.
func groupPrefix(key, sep string) string {
	if idx := strings.Index(key, sep); idx > 0 {
		return key[:idx]
	}
	return key
}

// quoteIfNeeded wraps values containing spaces or special characters in double quotes.
func quoteIfNeeded(v string) string {
	if strings.ContainsAny(v, " \t\n#$") {
		return `"` + strings.ReplaceAll(v, `"`, `\"`) + `"`
	}
	return v
}
