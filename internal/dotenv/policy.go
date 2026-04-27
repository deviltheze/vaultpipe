package dotenv

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// PolicyRule defines a single enforcement rule for a secret key pattern.
type PolicyRule struct {
	KeyPattern  string `json:"key_pattern"`
	Required    bool   `json:"required"`
	MinLength   int    `json:"min_length,omitempty"`
	MaxLength   int    `json:"max_length,omitempty"`
	MustEncrypt bool   `json:"must_encrypt,omitempty"`
}

// Policy holds a set of rules applied during sync.
type Policy struct {
	Version   string       `json:"version"`
	Rules     []PolicyRule `json:"rules"`
	CreatedAt time.Time    `json:"created_at"`
}

// PolicyViolation describes a single rule breach.
type PolicyViolation struct {
	Key     string
	Rule    PolicyRule
	Message string
}

func (v PolicyViolation) Error() string {
	return fmt.Sprintf("policy violation for key %q: %s", v.Key, v.Message)
}

// PolicyPath returns the canonical path for the policy file.
func PolicyPath(dir string) string {
	return filepath.Join(dir, ".vaultpipe-policy.json")
}

// WritePolicyFile serialises p to dir.
func WritePolicyFile(dir string, p Policy) error {
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now().UTC()
	}
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("policy: marshal: %w", err)
	}
	return os.WriteFile(PolicyPath(dir), data, 0o600)
}

// ReadPolicyFile loads a policy from dir. Returns an empty Policy when the
// file does not exist so callers can treat a missing file as "no rules".
func ReadPolicyFile(dir string) (Policy, error) {
	data, err := os.ReadFile(PolicyPath(dir))
	if os.IsNotExist(err) {
		return Policy{}, nil
	}
	if err != nil {
		return Policy{}, fmt.Errorf("policy: read: %w", err)
	}
	var p Policy
	if err := json.Unmarshal(data, &p); err != nil {
		return Policy{}, fmt.Errorf("policy: unmarshal: %w", err)
	}
	return p, nil
}

// EnforcePolicy checks secrets against every rule in p and returns all
// violations found. A nil/empty slice means the secrets are compliant.
func EnforcePolicy(secrets map[string]string, p Policy) []PolicyViolation {
	var violations []PolicyViolation
	for _, rule := range p.Rules {
		matched := keysMatchingPattern(secrets, rule.KeyPattern)
		if rule.Required && len(matched) == 0 {
			violations = append(violations, PolicyViolation{
				Key:     rule.KeyPattern,
				Rule:    rule,
				Message: "required key pattern has no matches",
			})
			continue
		}
		for _, k := range matched {
			v := secrets[k]
			if rule.MinLength > 0 && len(v) < rule.MinLength {
				violations = append(violations, PolicyViolation{Key: k, Rule: rule,
					Message: fmt.Sprintf("value length %d is below minimum %d", len(v), rule.MinLength)})
			}
			if rule.MaxLength > 0 && len(v) > rule.MaxLength {
				violations = append(violations, PolicyViolation{Key: k, Rule: rule,
					Message: fmt.Sprintf("value length %d exceeds maximum %d", len(v), rule.MaxLength)})
			}
		}
	}
	return violations
}

// keysMatchingPattern returns all keys whose name contains pattern as a
// substring (simple, allocation-cheap matching sufficient for env-key rules).
func keysMatchingPattern(secrets map[string]string, pattern string) []string {
	var out []string
	for k := range secrets {
		if containsStr(k, pattern) {
			out = append(out, k)
		}
	}
	return out
}
