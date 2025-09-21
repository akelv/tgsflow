package cmd

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func writeFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
}

func TestVerify_EARS_AllValid(t *testing.T) {
	dir := t.TempDir()
	// Enable EARS
	writeFile(t, filepath.Join(dir, "tgs", "tgs.yml"), "guardrails:\n  ears:\n    enable: true\n")
	// Use fixture files from core testdata
	fixtures := []string{
		"src/core/ears/testdata/positive_ubiquitous.md",
		"src/core/ears/testdata/positive_event.md",
		"src/core/ears/testdata/positive_state.md",
		"src/core/ears/testdata/positive_complex.md",
		"src/core/ears/testdata/positive_unwanted.md",
		"src/core/ears/testdata/formatting_bullets_and_skip_blocks.md",
	}
	_, thisFile, _, _ := runtime.Caller(0)
	repoRoot := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", ".."))
	for _, f := range fixtures {
		data, err := os.ReadFile(filepath.Join(repoRoot, f))
		if err != nil {
			t.Fatalf("read fixture: %v", err)
		}
		writeFile(t, filepath.Join(dir, filepath.Base(f)), string(data))
	}

	// CI mode should still succeed (no issues)
	code := CmdVerify([]string{"--repo", dir, "--ci"})
	if code != 0 {
		t.Fatalf("expected code=0, got %d", code)
	}
}

func TestVerify_EARS_WithInvalid(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "tgs", "tgs.yml"), "guardrails:\n  ears:\n    enable: true\n")
	fixtures := []string{
		"src/core/ears/testdata/negative_missing_system.md",
		"src/core/ears/testdata/negative_wrong_order.md",
		"src/core/ears/testdata/negative_multiple_when.md",
	}
	_, thisFile, _, _ := runtime.Caller(0)
	repoRoot := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", ".."))
	for _, f := range fixtures {
		data, err := os.ReadFile(filepath.Join(repoRoot, f))
		if err != nil {
			t.Fatalf("read fixture: %v", err)
		}
		writeFile(t, filepath.Join(dir, filepath.Base(f)), string(data))
	}

	code := CmdVerify([]string{"--repo", dir, "--ci"})
	if code != 1 {
		t.Fatalf("expected code=1 due to invalid line, got %d", code)
	}
}
