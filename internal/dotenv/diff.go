package dotenv

// ChangeType describes the kind of change for a secret key.
type ChangeType string

const (
	ChangeAdded   ChangeType = "added"
	ChangeUpdated ChangeType = "updated"
	ChangeRemoved ChangeType = "removed"
	ChangeUnchanged ChangeType = "unchanged"
)

// DiffEntry represents a single key-level change between two env maps.
type DiffEntry struct {
	Key    string
	Change ChangeType
}

// Diff compares an existing env map against an incoming map and returns
// a slice of DiffEntry describing what would change.
func Diff(existing, incoming map[string]string) []DiffEntry {
	seen := make(map[string]bool)
	var entries []DiffEntry

	for k, inVal := range incoming {
		seen[k] = true
		if exVal, ok := existing[k]; !ok {
			entries = append(entries, DiffEntry{Key: k, Change: ChangeAdded})
		} else if exVal != inVal {
			entries = append(entries, DiffEntry{Key: k, Change: ChangeUpdated})
		} else {
			entries = append(entries, DiffEntry{Key: k, Change: ChangeUnchanged})
		}
	}

	for k := range existing {
		if !seen[k] {
			entries = append(entries, DiffEntry{Key: k, Change: ChangeRemoved})
		}
	}

	return entries
}

// Summary returns counts of each change type from a diff result.
func Summary(entries []DiffEntry) map[ChangeType]int {
	m := map[ChangeType]int{
		ChangeAdded:     0,
		ChangeUpdated:   0,
		ChangeRemoved:   0,
		ChangeUnchanged: 0,
	}
	for _, e := range entries {
		m[e.Change]++
	}
	return m
}
