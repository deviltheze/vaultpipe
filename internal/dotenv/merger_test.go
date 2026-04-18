package dotenv

import (
	"sort"
	"testing"
)

func sortedKeys(s []string) []string {
	sort.Strings(s)
	return s
}

func TestMerge_OverwriteAll(t *testing.T) {
	existing := map[string]string{"A": "1", "B": "2"}
	incoming := map[string]string{"B": "99", "C": "3"}
	out, res := Merge(existing, incoming, OverwriteAll)
	if out["B"] != "99" {
		t.Fatalf("expected B=99, got %s", out["B"])
	}
	if out["C"] != "3" {
		t.Fatalf("expected C=3, got %s", out["C"])
	}
	if len(res.Updated) != 1 || res.Updated[0] != "B" {
		t.Fatalf("unexpected Updated: %v", res.Updated)
	}
	if len(res.Added) != 1 || res.Added[0] != "C" {
		t.Fatalf("unexpected Added: %v", res.Added)
	}
}

func TestMerge_KeepExisting(t *testing.T) {
	existing := map[string]string{"A": "1"}
	incoming := map[string]string{"A": "99", "B": "2"}
	out, res := Merge(existing, incoming, KeepExisting)
	if out["A"] != "1" {
		t.Fatalf("expected A=1 (kept), got %s", out["A"])
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "A" {
		t.Fatalf("unexpected Skipped: %v", res.Skipped)
	}
	if len(res.Added) != 1 || res.Added[0] != "B" {
		t.Fatalf("unexpected Added: %v", res.Added)
	}
}

func TestMerge_OverwriteChanged(t *testing.T) {
	existing := map[string]string{"A": "same", "B": "old"}
	incoming := map[string]string{"A": "same", "B": "new", "C": "added"}
	out, res := Merge(existing, incoming, OverwriteChanged)
	if out["A"] != "same" {
		t.Fatalf("expected A unchanged")
	}
	if out["B"] != "new" {
		t.Fatalf("expected B=new")
	}
	if len(sortedKeys(res.Skipped)) != 1 || res.Skipped[0] != "A" {
		t.Fatalf("unexpected Skipped: %v", res.Skipped)
	}
	if len(res.Updated) != 1 || res.Updated[0] != "B" {
		t.Fatalf("unexpected Updated: %v", res.Updated)
	}
	if len(res.Added) != 1 || res.Added[0] != "C" {
		t.Fatalf("unexpected Added: %v", res.Added)
	}
}

func TestMerge_EmptyIncoming(t *testing.T) {
	existing := map[string]string{"A": "1"}
	out, res := Merge(existing, map[string]string{}, OverwriteAll)
	if out["A"] != "1" {
		t.Fatal("existing key should be preserved")
	}
	if len(res.Added)+len(res.Updated)+len(res.Skipped) != 0 {
		t.Fatal("expected empty result")
	}
}
