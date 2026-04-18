package dotenv

import (
	"testing"
)

func collectByChange(entries []DiffEntry, ct ChangeType) []string {
	var keys []string
	for _, e := range entries {
		if e.Change == ct {
			keys = append(keys, e.Key)
		}
	}
	return keys
}

func TestDiff_Added(t *testing.T) {
	existing := map[string]string{"A": "1"}
	incoming := map[string]string{"A": "1", "B": "2"}
	entries := Diff(existing, incoming)
	added := collectByChange(entries, ChangeAdded)
	if len(added) != 1 || added[0] != "B" {
		t.Errorf("expected B added, got %v", added)
	}
}

func TestDiff_Updated(t *testing.T) {
	existing := map[string]string{"A": "old"}
	incoming := map[string]string{"A": "new"}
	entries := Diff(existing, incoming)
	updated := collectByChange(entries, ChangeUpdated)
	if len(updated) != 1 || updated[0] != "A" {
		t.Errorf("expected A updated, got %v", updated)
	}
}

func TestDiff_Removed(t *testing.T) {
	existing := map[string]string{"A": "1", "B": "2"}
	incoming := map[string]string{"A": "1"}
	entries := Diff(existing, incoming)
	removed := collectByChange(entries, ChangeRemoved)
	if len(removed) != 1 || removed[0] != "B" {
		t.Errorf("expected B removed, got %v", removed)
	}
}

func TestDiff_Unchanged(t *testing.T) {
	existing := map[string]string{"A": "1"}
	incoming := map[string]string{"A": "1"}
	entries := Diff(existing, incoming)
	unchanged := collectByChange(entries, ChangeUnchanged)
	if len(unchanged) != 1 {
		t.Errorf("expected 1 unchanged, got %v", unchanged)
	}
}

func TestDiff_Summary(t *testing.T) {
	existing := map[string]string{"A": "1", "B": "old", "C": "3"}
	incoming := map[string]string{"A": "1", "B": "new", "D": "4"}
	entries := Diff(existing, incoming)
	summary := Summary(entries)
	if summary[ChangeAdded] != 1 {
		t.Errorf("added: want 1, got %d", summary[ChangeAdded])
	}
	if summary[ChangeUpdated] != 1 {
		t.Errorf("updated: want 1, got %d", summary[ChangeUpdated])
	}
	if summary[ChangeRemoved] != 1 {
		t.Errorf("removed: want 1, got %d", summary[ChangeRemoved])
	}
	if summary[ChangeUnchanged] != 1 {
		t.Errorf("unchanged: want 1, got %d", summary[ChangeUnchanged])
	}
}

func TestDiff_EmptyBoth(t *testing.T) {
	entries := Diff(map[string]string{}, map[string]string{})
	if len(entries) != 0 {
		t.Errorf("expected no entries, got %d", len(entries))
	}
}
