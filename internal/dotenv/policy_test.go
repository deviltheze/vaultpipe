package dotenv_test

import (
	"os"
	"testing"

	"github.com/your-org/vaultpipe/internal/dotenv"
)

func TestPolicyPath_Convention(t *testing.T) {
	p := dotenv.PolicyPath("/tmp/myenv")
	if p != "/tmp/myenv/.vaultpipe-policy.json" {
		t.Fatalf("unexpected path: %s", p)
	}
}

func TestWriteAndReadPolicyFile_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	pol := dotenv.Policy{
		Version: "1",
		Rules: []dotenv.PolicyRule{
			{KeyPattern: "SECRET", Required: true, MinLength: 8},
		},
	}
	if err := dotenv.WritePolicyFile(dir, pol); err != nil {
		t.Fatalf("write: %v", err)
	}
	got, err := dotenv.ReadPolicyFile(dir)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if got.Version != pol.Version {
		t.Errorf("version: want %s got %s", pol.Version, got.Version)
	}
	if len(got.Rules) != 1 {
		t.Fatalf("rules: want 1 got %d", len(got.Rules))
	}
	if got.Rules[0].KeyPattern != "SECRET" {
		t.Errorf("pattern: want SECRET got %s", got.Rules[0].KeyPattern)
	}
}

func TestReadPolicyFile_MissingFile(t *testing.T) {
	dir := t.TempDir()
	pol, err := dotenv.ReadPolicyFile(dir)
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(pol.Rules) != 0 {
		t.Errorf("expected empty rules")
	}
}

func TestWritePolicyFile_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	if err := dotenv.WritePolicyFile(dir, dotenv.Policy{Version: "1"}); err != nil {
		t.Fatalf("write: %v", err)
	}
	info, err := os.Stat(dotenv.PolicyPath(dir))
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("perm: want 0600 got %o", info.Mode().Perm())
	}
}

func TestEnforcePolicy_NoViolations(t *testing.T) {
	secrets := map[string]string{"DB_SECRET": "supersecret123"}
	pol := dotenv.Policy{
		Rules: []dotenv.PolicyRule{
			{KeyPattern: "SECRET", Required: true, MinLength: 8},
		},
	}
	violations := dotenv.EnforcePolicy(secrets, pol)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %d", len(violations))
	}
}

func TestEnforcePolicy_RequiredMissing(t *testing.T) {
	secrets := map[string]string{"UNRELATED": "value"}
	pol := dotenv.Policy{
		Rules: []dotenv.PolicyRule{
			{KeyPattern: "SECRET", Required: true},
		},
	}
	violations := dotenv.EnforcePolicy(secrets, pol)
	if len(violations) != 1 {
		t.Fatalf("want 1 violation, got %d", len(violations))
	}
}

func TestEnforcePolicy_MinLengthBreached(t *testing.T) {
	secrets := map[string]string{"API_SECRET": "short"}
	pol := dotenv.Policy{
		Rules: []dotenv.PolicyRule{
			{KeyPattern: "SECRET", MinLength: 10},
		},
	}
	violations := dotenv.EnforcePolicy(secrets, pol)
	if len(violations) != 1 {
		t.Fatalf("want 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "API_SECRET" {
		t.Errorf("unexpected key: %s", violations[0].Key)
	}
}

func TestEnforcePolicy_MaxLengthBreached(t *testing.T) {
	secrets := map[string]string{"TOKEN": "this-value-is-way-too-long-for-the-rule"}
	pol := dotenv.Policy{
		Rules: []dotenv.PolicyRule{
			{KeyPattern: "TOKEN", MaxLength: 10},
		},
	}
	violations := dotenv.EnforcePolicy(secrets, pol)
	if len(violations) != 1 {
		t.Fatalf("want 1 violation, got %d", len(violations))
	}
}

func TestEnforcePolicy_EmptyPolicy_NoViolations(t *testing.T) {
	secrets := map[string]string{"ANY_KEY": "any_value"}
	violations := dotenv.EnforcePolicy(secrets, dotenv.Policy{})
	if len(violations) != 0 {
		t.Errorf("empty policy should produce no violations")
	}
}
