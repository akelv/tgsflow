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
		filepath.Join("tgs", "00_research.md"),
		filepath.Join("tgs", "10_spec.md"),
		filepath.Join("tgs", "20_plan.md"),
		filepath.Join("tgs", "30_tasks.md"),
		filepath.Join("tgs", "40_approval.md"),
	}
	for _, p := range checks {
		if _, err := os.Stat(p); err != nil {
			t.Fatalf("expected %s to exist: %v", p, err)
		}
	}
}
