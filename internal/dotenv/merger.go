package dotenv

// MergeStrategy controls how existing keys are handled during a merge.
type MergeStrategy int

const (
	// OverwriteAll replaces every existing key with the incoming value.
	OverwriteAll MergeStrategy = iota
	// KeepExisting preserves keys already present in the destination.
	KeepExisting
	// OverwriteChanged updates only keys whose values differ.
	OverwriteChanged
)

// MergeResult summarises what changed after a merge.
type MergeResult struct {
	Added   []string
	Updated []string
	Skipped []string
}

// Merge combines incoming secrets into existing env vars according to the
// chosen strategy and returns a merged map plus a change summary.
func Merge(existing, incoming map[string]string, strategy MergeStrategy) (map[string]string, MergeResult) {
	out := make(map[string]string, len(existing))
	for k, v := range existing {
		out[k] = v
	}

	var result MergeResult

	for k, v := range incoming {
		old, exists := out[k]
		switch {
		case !exists:
			out[k] = v
			result.Added = append(result.Added, k)
		case strategy == KeepExisting:
			result.Skipped = append(result.Skipped, k)
		case strategy == OverwriteChanged && old == v:
			result.Skipped = append(result.Skipped, k)
		default:
			out[k] = v
			result.Updated = append(result.Updated, k)
		}
	}

	return out, result
}
