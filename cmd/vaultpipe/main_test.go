package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestRootCmd_MissingConfig(t *testing.T) {
	rootCmd.SetArgs([]string{"--config", "nonexistent.yaml"})
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing config, got nil")
	}
}

func TestRootCmd_ValidConfig(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "vaultpipe.yaml")
	outPath := filepath.Join(dir, ".env")

	cfgContent := "vault_address: http://127.0.0.1:19999\n" +
		"secret_path: secret/data/app\n" +
		"output_file: " + outPath + "\n"

	if err := os.WriteFile(cfgPath, []byte(cfgContent), 0o644); err != nil {
		t.Fatalf("writing config: %v", err)
	}

	rootCmd.SetArgs([]string{"--config", cfgPath})
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)

	err := rootCmd.Execute()
	// Expect an error because vault is not running, but config should load fine.
	if err == nil {
		t.Fatal("expected error connecting to vault")
	}

	if buf.Len() != 0 {
		t.Errorf("unexpected stdout output: %s", buf.String())
	}
}
