package dotenv

import (
	"fmt"
	"sort"
	"strings"
)

// DriftEntry describes a single key that has drifted between the live env file
// and the values last synced from Vault.
type DriftEntry struct {
	Key      string
	EnvValue string // value currently in the .env file
	VaultValue string // value from Vault (expected)
	Kind     string // "modified", "missing_in_env", "extra_in_env"
}

// DriftReport holds the complete drift analysis result.
type DriftReport struct {
	Entries []DriftEntry
	Clean   bool
}

// String returns a human-readable summary of the drift report.
func (r DriftReport) String() string {
	if r.Clean {
		return "no drift detected"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "drift detected (%d key(s)):\n", len(r.Entries))
	for _, e := range r.Entries {
		switch e.Kind {
		case "modified":
			fmt.Fprintf(&sb, "  ~ %s (env value differs from vault)\n", e.Key)
		case "missing_in_env":
			fmt.Fprintf(&sb, "  + %s (present in vault, missing from env)\n", e.Key)
		case "extra_in_env":
			fmt.Fprintf(&sb, "  - %s (present in env, not in vault)\n", e.Key)
		}
	}
	return sb.String()
}

// DetectDrift compares the current .env map against the expected Vault secrets
// and returns a DriftReport describing any differences.
func DetectDrift(envSecrets, vaultSecrets map[string]string) DriftReport {
	var entries []DriftEntry

	// Keys in vault — check for missing or modified in env.
	keys := make([]string, 0, len(vaultSecrets))
	for k := range vaultSecrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		envVal, ok := envSecrets[k]
		if !ok {
			entries = append(entries, DriftEntry{Key: k, VaultValue: vaultSecrets[k], Kind: "missing_in_env"})
		} else if envVal != vaultSecrets[k] {
			entries = append(entries, DriftEntry{Key: k, EnvValue: envVal, VaultValue: vaultSecrets[k], Kind: "modified"})
		}
	}

	// Keys only in env — extra keys not tracked by Vault.
	extraKeys := make([]string, 0)
	for k := range envSecrets {
		if _, ok := vaultSecrets[k]; !ok {
			extraKeys = append(extraKeys, k)
		}
	}
	sort.Strings(extraKeys)
	for _, k := range extraKeys {
		entries = append(entries, DriftEntry{Key: k, EnvValue: envSecrets[k], Kind: "extra_in_env"})
	}

	return DriftReport{Entries: entries, Clean: len(entries) == 0}
}
