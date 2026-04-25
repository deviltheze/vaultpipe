package dotenv

import "fmt"

// CompareResult holds the result of comparing two sets of secrets.
type CompareResult struct {
	OnlyInA    map[string]string // keys present in A but not B
	OnlyInB    map[string]string // keys present in B but not A
	Different  map[string][2]string // keys in both but with different values [a, b]
	Identical  map[string]string // keys with identical values in both
}

// Summary returns a human-readable one-line summary of the comparison.
func (r CompareResult) Summary() string {
	return fmt.Sprintf("%d only in source, %d only in target, %d changed, %d identical",
		len(r.OnlyInA), len(r.OnlyInB), len(r.Different), len(r.Identical))
}

// HasDifferences returns true if there is any divergence between the two sets.
func (r CompareResult) HasDifferences() bool {
	return len(r.OnlyInA) > 0 || len(r.OnlyInB) > 0 || len(r.Different) > 0
}

// Compare performs a full key-value comparison between two secret maps.
// a is typically the Vault/source secrets, b is the existing .env contents.
func Compare(a, b map[string]string) CompareResult {
	result := CompareResult{
		OnlyInA:   make(map[string]string),
		OnlyInB:   make(map[string]string),
		Different: make(map[string][2]string),
		Identical: make(map[string]string),
	}

	for k, va := range a {
		if vb, ok := b[k]; ok {
			if va == vb {
				result.Identical[k] = va
			} else {
				result.Different[k] = [2]string{va, vb}
			}
		} else {
			result.OnlyInA[k] = va
		}
	}

	for k, vb := range b {
		if _, ok := a[k]; !ok {
			result.OnlyInB[k] = vb
		}
	}

	return result
}
