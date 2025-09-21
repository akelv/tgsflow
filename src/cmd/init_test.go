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
	code := CmdInit([]string{"--decorate", "--ci-template", "none"})
	if code != 0 {
		t.Fatalf("CmdInit returned %d", code)
	}
	checks := []string{
		filepath.Join("tgs", "README.md"),
		filepath.Join("tgs", "tgs.yml"),
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
	}
	for _, p := range checks {
		if _, err := os.Stat(p); err != nil {
			t.Fatalf("expected %s to exist: %v", p, err)
		}
	}
}
