package sync

import (
	"os"
	"testing"

	"github.com/your-org/vaultpipe/internal/dotenv"
)

// stubPolicyLogger captures LogError calls.
type stubPolicyLogger struct{ errors []string }

func (s *stubPolicyLogger) LogError(_, msg string) { s.errors = append(s.errors, msg) }

// stubPolicyCfg satisfies the minimal interface required by enforcePolicy.
type stubPolicyCfg struct{ dir string }

func (s *stubPolicyCfg) GetOutputDir() string { return s.dir }

func TestEnforcePolicy_NoPolicyFile_Passes(t *testing.T) {
	dir := t.TempDir()
	log := &stubPolicyLogger{}
	if err := enforcePolicy(&stubPolicyCfg{dir}, map[string]string{"KEY": "val"}, log); err != nil {
		t.Fatalf("expected no error for missing policy file, got %v", err)
	}
	if len(log.errors) != 0 {
		t.Errorf("expected no logged errors")
	}
}

func TestEnforcePolicy_Compliant_Passes(t *testing.T) {
	dir := t.TempDir()
	pol := dotenv.Policy{
		Version: "1",
		Rules:   []dotenv.PolicyRule{{KeyPattern: "SECRET", Required: true, MinLength: 4}},
	}
	if err := dotenv.WritePolicyFile(dir, pol); err != nil {
		t.Fatalf("write policy: %v", err)
	}
	log := &stubPolicyLogger{}
	err := enforcePolicy(&stubPolicyCfg{dir}, map[string]string{"DB_SECRET": "longval"}, log)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestEnforcePolicy_Violation_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	pol := dotenv.Policy{
		Version: "1",
		Rules:   []dotenv.PolicyRule{{KeyPattern: "SECRET", Required: true}},
	}
	if err := dotenv.WritePolicyFile(dir, pol); err != nil {
		t.Fatalf("write policy: %v", err)
	}
	log := &stubPolicyLogger{}
	err := enforcePolicy(&stubPolicyCfg{dir}, map[string]string{"UNRELATED": "value"}, log)
	if err == nil {
		t.Fatal("expected error for policy violation")
	}
	if len(log.errors) == 0 {
		t.Error("expected at least one logged error")
	}
}

func TestEnforcePolicy_BadPolicyFile_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	// Write invalid JSON.
	if err := os.WriteFile(dotenv.PolicyPath(dir), []byte("not-json"), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}
	log := &stubPolicyLogger{}
	if err := enforcePolicy(&stubPolicyCfg{dir}, map[string]string{}, log); err == nil {
		t.Fatal("expected error for bad policy file")
	}
}
