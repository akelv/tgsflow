package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitSeedsFiles(t *testing.T) {
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	code := CmdInit([]string{})
	if code != 0 {
		t.Fatalf("CmdInit returned %d", code)
	}
	checks := []string{
		"CLAUDE.md",
		"Makefile.tgs.mk",
		"Makefile",
		filepath.Join(".claude", "hooks.json"),
		filepath.Join(".claude", "rules", "README.md"),
		filepath.Join(".claude", "commands", "tgs-research.md"),
		filepath.Join(".claude", "commands", "tgs-plan.md"),
		filepath.Join(".claude", "commands", "tgs-close.md"),
		filepath.Join("tgs", "README.md"),
		filepath.Join("tgs", "tgs.yml"),
		filepath.Join("tgs", "judgment-journal.md"),
		filepath.Join("tgs", "design", "00_context.md"),
		filepath.Join("tgs", "design", "10_needs.md"),
		filepath.Join("tgs", "design", "20_requirements.md"),
		filepath.Join("tgs", "design", "30_architecture.md"),
		filepath.Join("tgs", "design", "40_vnv.md"),
		filepath.Join("tgs", "design", "50_decisions.md"),
		filepath.Join("tgs", "agentops", "AGENTOPS.md"),
		filepath.Join("tgs", "agentops", "tgs", "research.md"),
		filepath.Join("tgs", "agentops", "tgs", "plan.md"),
		filepath.Join("tgs", "agentops", "tgs", "implementation.md"),
		filepath.Join("tgs", "adapters", "claude-code.sh"),
		filepath.Join("tgs", "adapters", "gemini-code.sh"),
	}
	for _, p := range checks {
		if _, err := os.Stat(p); err != nil {
			t.Fatalf("expected %s to exist: %v", p, err)
		}
	}

	// Adapter scripts must be executable.
	for _, sh := range []string{
		filepath.Join("tgs", "adapters", "claude-code.sh"),
		filepath.Join("tgs", "adapters", "gemini-code.sh"),
	} {
		info, err := os.Stat(sh)
		if err != nil {
			t.Fatalf("stat %s: %v", sh, err)
		}
		if info.Mode().Perm()&0o111 == 0 {
			t.Fatalf("%s is not executable (mode=%v)", sh, info.Mode().Perm())
		}
	}
}

func TestInit_IdempotentByDefault(t *testing.T) {
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	// Pre-create a scaffold file with sentinel content; a plain init should not overwrite it.
	target := filepath.Join(dir, "CLAUDE.md")
	sentinel := []byte("SENTINEL-USER-EDIT\n")
	if err := os.WriteFile(target, sentinel, 0o644); err != nil {
		t.Fatal(err)
	}
	if code := CmdInit([]string{}); code != 0 {
		t.Fatalf("CmdInit returned %d", code)
	}
	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(sentinel) {
		t.Fatalf("expected sentinel preserved, got: %q", string(got))
	}
}

func TestInit_ForceOverwrites(t *testing.T) {
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(dir, "CLAUDE.md")
	if err := os.WriteFile(target, []byte("STALE\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if code := CmdInit([]string{"--force"}); code != 0 {
		t.Fatalf("CmdInit returned %d", code)
	}
	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) == "STALE\n" {
		t.Fatalf("expected CLAUDE.md to be overwritten with --force")
	}
}

func TestInit_ForcePreservesThoughts(t *testing.T) {
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	// Seed user content under tgs/thoughts/
	thoughtsDir := filepath.Join(dir, "tgs", "thoughts", "abc-my-work")
	if err := os.MkdirAll(thoughtsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	userFile := filepath.Join(thoughtsDir, "research.md")
	userContent := []byte("USER-WORK-IN-PROGRESS\n")
	if err := os.WriteFile(userFile, userContent, 0o644); err != nil {
		t.Fatal(err)
	}
	if code := CmdInit([]string{"--force"}); code != 0 {
		t.Fatalf("CmdInit returned %d", code)
	}
	got, err := os.ReadFile(userFile)
	if err != nil {
		t.Fatalf("user file lost: %v", err)
	}
	if string(got) != string(userContent) {
		t.Fatalf("user file clobbered: got %q", string(got))
	}
}
