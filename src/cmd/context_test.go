package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFile2(t *testing.T, path, content string, mode os.FileMode) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), mode); err != nil {
		t.Fatal(err)
	}
}

func TestContextPack_MissingAdapter(t *testing.T) {
	cwd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(cwd) })

	repo := t.TempDir()
	if err := os.Chdir(repo); err != nil {
		t.Fatal(err)
	}

	// Minimal repo layout: design files and a thought dir
	design := filepath.Join(repo, "tgs", "design")
	writeFile2(t, filepath.Join(design, "10_needs.md"), "- need", 0o644)

	thought := filepath.Join(repo, "tgs", "thoughts", "abcdef1-context-pack-test")
	writeFile2(t, filepath.Join(thought, "README.md"), "readme", 0o644)
	t.Setenv("TGS_THOUGHT_DIR", thought)

	// Config with non-existent adapter
	tgsYml := `ai:
  mode: shell
  shell_adapter_path: ` + filepath.Join(repo, "no_such_adapter.sh") + `
  shell_claude_cmd: claude
`
	writeFile2(t, filepath.Join(repo, "tgs", "tgs.yml"), tgsYml, 0o644)

	cmd := newContextPackCommand()
	cmd.SetArgs([]string{"auth sso"})
	if err := cmd.Execute(); err == nil {
		t.Fatalf("expected error due to missing adapter")
	}
}

func TestContextPack_HappyPath(t *testing.T) {
	cwd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(cwd) })

	repo := t.TempDir()
	if err := os.Chdir(repo); err != nil {
		t.Fatal(err)
	}

	// Minimal repo layout
	design := filepath.Join(repo, "tgs", "design")
	writeFile2(t, filepath.Join(design, "20_requirements.md"), "- shall", 0o644)

	thought := filepath.Join(repo, "tgs", "thoughts", "abcdef2-context-pack-test")
	writeFile2(t, filepath.Join(thought, "plan.md"), "plan", 0o644)
	t.Setenv("TGS_THOUGHT_DIR", thought)

	// Adapter that parses --out and writes a brief
	adapter := filepath.Join(repo, "fake_adapter.sh")
	script := "#!/bin/sh\n" +
		"out=\"\"\n" +
		"while [ $# -gt 0 ]; do\n" +
		"  if [ \"$1\" = \"--out\" ]; then out=\"$2\"; shift 2; continue; fi\n" +
		"  shift 1\n" +
		"done\n" +
		"[ -n \"$out\" ] || { echo no-out 1>&2; exit 2; }\n" +
		"echo '# AI Brief' > \"$out\"\n" +
		"exit 0\n"
	writeFile2(t, adapter, script, 0o755)

	// Config using our fake adapter
	tgsYml := `ai:
  mode: shell
  shell_adapter_path: ` + adapter + `
  shell_claude_cmd: sh
  timeout_ms: 3000
  toolpack:
    budgets:
      context_pack_tokens: 200
`
	writeFile2(t, filepath.Join(repo, "tgs", "tgs.yml"), tgsYml, 0o644)

	// Run command
	cmd := newContextPackCommand()
	cmd.SetArgs([]string{"authentication and sso"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Expect aibrief.md in thought dir
	brief := filepath.Join(thought, "aibrief.md")
	b, err := os.ReadFile(brief)
	if err != nil {
		t.Fatalf("brief not written: %v", err)
	}
	if len(b) == 0 || string(b[:10]) != "# AI Brief" {
		t.Fatalf("unexpected brief content: %q", string(b))
	}
}
