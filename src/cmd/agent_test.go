package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCmdAgent_Help(t *testing.T) {
	code := CmdAgent(nil)
	if code != 0 {
		t.Fatalf("expected help exit code 0, got %d", code)
	}
}

func TestCmdAgent_UnknownSubcommand(t *testing.T) {
	code := CmdAgent([]string{"unknown"})
	if code != 2 {
		t.Fatalf("expected exit code 2 for unknown subcommand, got %d", code)
	}
}

func TestCmdAgent_Exec_UsageError(t *testing.T) {
	// Missing required flags should result in usage error (exit code 2)
	code := CmdAgent([]string{"exec"})
	if code != 2 {
		t.Fatalf("expected exit code 2 for usage error, got %d", code)
	}
}

func TestCmdAgent_Exec_DelegatesAndRuns(t *testing.T) {
	// Prepare temp working dir
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	// Create a minimal adapter script that always succeeds
	adapterPath := filepath.Join(dir, "adapter.sh")
	adapter := "#!/bin/sh\necho OK\n"
	if err := os.WriteFile(adapterPath, []byte(adapter), 0o755); err != nil {
		t.Fatalf("write adapter: %v", err)
	}

	// Create a dummy context file to satisfy validation
	ctxPath := filepath.Join(dir, "ctx.txt")
	if err := os.WriteFile(ctxPath, []byte("context"), 0o644); err != nil {
		t.Fatalf("write ctx: %v", err)
	}

	code := CmdAgent([]string{
		"exec",
		"--prompt-text", "hello",
		"--context", ctxPath,
		"--adapter-path", adapterPath,
	})
	if code != 0 {
		t.Fatalf("expected exit code 0 from exec, got %d", code)
	}
}
