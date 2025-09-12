package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTasksCreateAndValidate(t *testing.T) {
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	// Create
	if code := CmdTasks(nil); code != 0 {
		t.Fatalf("CmdTasks returned %d", code)
	}
	tasks := filepath.Join("tgs", "30_tasks.md")
	if _, err := os.Stat(tasks); err != nil {
		t.Fatalf("expected %s: %v", tasks, err)
	}
	// Validate
	if code := CmdTasks([]string{"--validate"}); code != 0 {
		t.Fatalf("CmdTasks --validate returned %d", code)
	}
}
