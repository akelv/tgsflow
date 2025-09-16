package cmd

import (
	"os"
	"path/filepath"
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
	writeFile(t, filepath.Join(dir, "tgs.yaml"), "policies:\n  ears:\n    enable: true\n")
	// Use fixture files
	fixtures := []string{
		"tgs/suggestions/ears_fixtures/positive_ubiquitous.md",
		"tgs/suggestions/ears_fixtures/positive_event.md",
		"tgs/suggestions/ears_fixtures/positive_state.md",
		"tgs/suggestions/ears_fixtures/positive_complex.md",
		"tgs/suggestions/ears_fixtures/positive_unwanted.md",
		"tgs/suggestions/ears_fixtures/formatting_bullets_and_skip_blocks.md",
	}
	for _, f := range fixtures {
		data, err := os.ReadFile(filepath.Join("/Users/kelvin/github/tgsflow", f))
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
	writeFile(t, filepath.Join(dir, "tgs.yaml"), "policies:\n  ears:\n    enable: true\n")
	fixtures := []string{
		"tgs/suggestions/ears_fixtures/negative_missing_system.md",
		"tgs/suggestions/ears_fixtures/negative_wrong_order.md",
		"tgs/suggestions/ears_fixtures/negative_multiple_when.md",
	}
	for _, f := range fixtures {
		data, err := os.ReadFile(filepath.Join("/Users/kelvin/github/tgsflow", f))
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
