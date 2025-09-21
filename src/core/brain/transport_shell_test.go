package brain

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/kelvin/tgsflow/src/core/config"
)

// writeTempAdapter creates a small bash adapter that mimics claude-code.sh behavior
// for the bits we rely on: it reads flags and prints either text or a path.
func writeTempAdapter(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "adapter.sh")
	if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
		t.Fatalf("write temp adapter: %v", err)
	}
	return path
}

func TestShellTransport_HappyPath_Text(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell adapter tests require sh")
	}
	adapter := writeTempAdapter(t, `#!/usr/bin/env bash
set -euo pipefail
echo "Hello from adapter"
`)

	// Inject our adapter path by temporarily moving the real default aside via symlink
	// Instead, rewire NewShellTransport by setting the adapter field post-construction.
	cfg := config.Default()
	tr := NewShellTransport(cfg).(*shellTransport)
	tr.adapterPath = adapter
	tr.claudeCmd = "echo"

	resp, err := tr.Chat(context.Background(), ChatReq{System: "sys", Messages: []Msg{{Role: "user", Content: "hi"}}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got, want := strings.TrimSpace(resp.Text), "Hello from adapter"; got != want {
		t.Fatalf("unexpected text: got %q want %q", got, want)
	}
}

func TestShellTransport_ErrorPath(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell adapter tests require sh")
	}
	adapter := writeTempAdapter(t, `#!/usr/bin/env bash
echo "boom" >&2
exit 7
`)
	cfg := config.Default()
	tr := NewShellTransport(cfg).(*shellTransport)
	tr.adapterPath = adapter
	_, err := tr.Chat(context.Background(), ChatReq{System: "sys"})
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "boom") {
		t.Fatalf("expected stderr in error, got: %v", err)
	}
}

func TestShellTransport_Timeout(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell adapter tests require sh")
	}
	// Adapter sleeps; we set a short context deadline. Note: adapter's internal timeout
	// is best-effort if GNU timeout is present; we rely on CommandContext to kill process.
	adapter := writeTempAdapter(t, `#!/usr/bin/env bash
sleep 5
echo done
`)
	cfg := config.Default()
	tr := NewShellTransport(cfg).(*shellTransport)
	tr.adapterPath = adapter

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	_, err := tr.Chat(ctx, ChatReq{System: "sys"})
	if err == nil {
		t.Fatalf("expected context timeout error")
	}
}
