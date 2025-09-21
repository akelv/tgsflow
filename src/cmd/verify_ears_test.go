package cmd

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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

func TestVerify_EARS_DesignDocs_Valid(t *testing.T) {
	dir := t.TempDir()
	// Enable EARS
	writeFile(t, filepath.Join(dir, "tgs", "tgs.yml"), "guardrails:\n  ears:\n    enable: true\n")

	// Create design docs with valid EARS lines and some narrative ignored lines
	needs := "# Stakeholder Needs\n\nWhile planning work, the team shall document needs.\n\nThis is narrative and should be ignored.\n\n```\n- code block should be ignored\n```\n\nWhen a new initiative starts, the team shall capture initial needs.\n"
	reqs := "# System Requirements\n\nThe system shall report verification results.\n\nWhile verifying documentation, the system shall enforce EARS rules.\n"

	writeFile(t, filepath.Join(dir, "tgs", "design", "10_needs.md"), needs)
	writeFile(t, filepath.Join(dir, "tgs", "design", "20_requirements.md"), reqs)

	// Subcommand should also succeed
	code := CmdVerifyEARS([]string{"--repo", dir, "--ci"})
	if code != 0 {
		t.Fatalf("expected code=0, got %d", code)
	}
}

func TestVerify_EARS_DesignDocs_InvalidReportsPathLine(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "tgs", "tgs.yml"), "guardrails:\n  ears:\n    enable: true\n")

	// Valid needs file
	needs := "While planning work, the team shall document needs.\n"
	writeFile(t, filepath.Join(dir, "tgs", "design", "10_needs.md"), needs)

	// Invalid requirements file (missing 'shall')
	reqs := "When a requirement is added, the system must validate it.\n"
	reqPath := filepath.Join(dir, "tgs", "design", "20_requirements.md")
	writeFile(t, reqPath, reqs)

	// Capture stderr to verify message format contains path:line: message
	oldStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stderr = w
	code := CmdVerifyEARS([]string{"--repo", dir, "--ci"})
	w.Close()
	os.Stderr = oldStderr
	out, _ := io.ReadAll(r)
	r.Close()

	if code != 1 {
		t.Fatalf("expected code=1 due to invalid line, got %d", code)
	}
	stderr := string(out)
	// Expect the file path and a colon+line
	if !strings.Contains(stderr, "tgs/design/20_requirements.md:") {
		t.Fatalf("expected stderr to contain path with line prefix, got: %q", stderr)
	}
}
