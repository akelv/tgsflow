package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPlanCreatesOrAppends(t *testing.T) {
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	// First run should create file from template
	if code := CmdPlan(nil); code != 0 {
		t.Fatalf("CmdPlan returned %d", code)
	}
	plan := filepath.Join("tgs", "20_plan.md")
	if _, err := os.Stat(plan); err != nil {
		t.Fatalf("expected %s: %v", plan, err)
	}
	// Second run should append a section
	if code := CmdPlan(nil); code != 0 {
		t.Fatalf("CmdPlan returned %d on second run", code)
	}
}
