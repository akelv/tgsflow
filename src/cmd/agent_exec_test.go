package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAgentExec_MissingPrompt(t *testing.T) {
	code, err := NewAgentExecCommand([]string{
		"--adapter-path", "/bin/true",
		"--context", "/etc/hosts",
	})
	if code != 2 || err == nil {
		t.Fatalf("expected code=2 and error for missing prompt, got code=%d err=%v", code, err)
	}
}

func TestAgentExec_MissingAdapter(t *testing.T) {
	code, err := NewAgentExecCommand([]string{
		"--adapter-path", "/path/does/not/exist.sh",
		"--prompt-text", "hello",
		"--context", "/etc/hosts",
	})
	if code != 2 || err == nil {
		t.Fatalf("expected code=2 and error for missing adapter, got code=%d err=%v", code, err)
	}
}

func TestAgentExec_NoContextProvided(t *testing.T) {
	code, err := NewAgentExecCommand([]string{
		"--adapter-path", "/bin/true",
		"--prompt-text", "hello",
	})
	if code != 2 || err == nil {
		t.Fatalf("expected code=2 and error for missing context, got code=%d err=%v", code, err)
	}
}

func TestAgentExec_ContextFileMissing(t *testing.T) {
	dir := t.TempDir()
	// Adapter that does nothing and succeeds
	adapter := filepath.Join(dir, "adapter.sh")
	if err := os.WriteFile(adapter, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	code, err := NewAgentExecCommand([]string{
		"--adapter-path", adapter,
		"--prompt-text", "hello",
		"--context", filepath.Join(dir, "missing.txt"),
	})
	if code != 2 || err == nil {
		t.Fatalf("expected code=2 and error for missing context file, got code=%d err=%v", code, err)
	}
}

func TestAgentExec_HappyPath(t *testing.T) {
	dir := t.TempDir()
	// Adapter that prints OK and exits 0
	adapter := filepath.Join(dir, "adapter.sh")
	if err := os.WriteFile(adapter, []byte("#!/bin/sh\necho OK\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	ctx := filepath.Join(dir, "ctx.txt")
	if err := os.WriteFile(ctx, []byte("context"), 0o644); err != nil {
		t.Fatal(err)
	}
	code, err := NewAgentExecCommand([]string{
		"--adapter-path", adapter,
		"--prompt-text", "hello",
		"--context", ctx,
	})
	if err != nil || code != 0 {
		t.Fatalf("expected success code=0 err=nil, got code=%d err=%v", code, err)
	}
}

func TestAgentExec_AdapterFailureSurfaceStderr(t *testing.T) {
	dir := t.TempDir()
	// Adapter that writes to stderr and exits non-zero
	adapter := filepath.Join(dir, "adapter.sh")
	script := "#!/bin/sh\necho adapter-problem 1>&2\nexit 3\n"
	if err := os.WriteFile(adapter, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	ctx := filepath.Join(dir, "ctx.txt")
	if err := os.WriteFile(ctx, []byte("context"), 0o644); err != nil {
		t.Fatal(err)
	}
	code, err := NewAgentExecCommand([]string{
		"--adapter-path", adapter,
		"--prompt-text", "hello",
		"--context", ctx,
	})
	if code != 1 || err == nil {
		t.Fatalf("expected adapter failure code=1 with error, got code=%d err=%v", code, err)
	}
}
