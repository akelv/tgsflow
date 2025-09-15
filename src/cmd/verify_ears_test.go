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
	// Markdown with valid EARS bullets
	md := "# Requirements\n\n- The payment service shall record transactions\n- When button is pressed, the controller shall start\n- While in armed mode, the system shall alert\n- While battery is low, when charger is connected, the device shall charge\n- If overheating, then the system shall shutdown\n"
	writeFile(t, filepath.Join(dir, "REQS.md"), md)

	// CI mode should still succeed (no issues)
	code := CmdVerify([]string{"--repo", dir, "--ci"})
	if code != 0 {
		t.Fatalf("expected code=0, got %d", code)
	}
}

func TestVerify_EARS_WithInvalid(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "tgs.yaml"), "policies:\n  ears:\n    enable: true\n")
	md := "# Requirements\n\n- The system shall record events\n- Because of X the system might respond\n"
	writeFile(t, filepath.Join(dir, "REQS.md"), md)

	code := CmdVerify([]string{"--repo", dir, "--ci"})
	if code != 1 {
		t.Fatalf("expected code=1 due to invalid line, got %d", code)
	}
}
