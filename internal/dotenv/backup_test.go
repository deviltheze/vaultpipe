package dotenv

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBackup_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, ".env")
	_ = os.WriteFile(src, []byte("FOO=bar\n"), 0600)

	dst, err := Backup(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dst == "" {
		t.Fatal("expected backup path, got empty string")
	}
	if !strings.HasSuffix(dst, ".bak") {
		t.Errorf("expected .bak suffix, got %s", dst)
	}
	data, _ := os.ReadFile(dst)
	if string(data) != "FOO=bar\n" {
		t.Errorf("backup content mismatch: %q", string(data))
	}
}

func TestBackup_NoOpWhenMissing(t *testing.T) {
	dir := t.TempDir()
	dst, err := Backup(filepath.Join(dir, ".env"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dst != "" {
		t.Errorf("expected empty dst, got %s", dst)
	}
}

func TestPruneBackups_KeepsNewest(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, ".env")
	_ = os.WriteFile(src, []byte("A=1"), 0600)

	// Create 5 backups.
	var created []string
	for i := 0; i < 5; i++ {
		dst, err := Backup(src)
		if err != nil {
			t.Fatalf("backup %d: %v", i, err)
		}
		created = append(created, dst)
		// Ensure unique timestamps by renaming with index suffix.
		_ = os.Rename(dst, dst+string(rune('0'+i)))
	}

	if err := PruneBackups(src, 2); err != nil {
		t.Fatalf("prune error: %v", err)
	}

	matches, _ := backupFiles(src)
	if len(matches) > 2 {
		t.Errorf("expected <=2 backups, got %d", len(matches))
	}
	_ = created
}

func TestPruneBackups_NothingToRemove(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, ".env")
	if err := PruneBackups(src, 3); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
