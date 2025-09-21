package brain

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/kelvin/tgsflow/src/core/config"
)

// shellTransport executes a local shell adapter (e.g., tgs/adapters/claude-code.sh)
// to fulfill Chat requests. It composes a prompt from system + messages and
// returns the adapter's textual output as ChatResp.Text.
type shellTransport struct {
	adapterPath    string
	claudeCmd      string
	returnMode     string
	suggestionsDir string
	timeoutSec     int // optional override; 0 means derive from ctx deadline or no timeout
}

func NewShellTransport(cfg config.Config) Transport {
	// Defaults aligned with agent_exec.go flags, overridable by config
	adapter := cfg.AI.ShellAdapterPath
	if strings.TrimSpace(adapter) == "" {
		adapter = "tgs/adapters/claude-code.sh"
	}
	if st, err := os.Stat(adapter); err != nil || st.Mode()&0111 == 0 {
		// keep as-is; adapter may still be runnable in tests or CI contexts
	}
	claudeCmd := cfg.AI.ShellClaudeCmd
	if strings.TrimSpace(claudeCmd) == "" {
		claudeCmd = "claude"
	}
	return &shellTransport{
		adapterPath:    adapter,
		claudeCmd:      claudeCmd,
		returnMode:     "patch_or_text",
		suggestionsDir: "tgs/suggestions",
		timeoutSec:     0,
	}
}

func (s *shellTransport) Chat(ctx context.Context, req ChatReq) (ChatResp, error) {
	if strings.TrimSpace(req.System) == "" && len(req.Messages) == 0 {
		return ChatResp{}, errors.New("empty prompt")
	}

	// Compose prompt similar to agent_exec: include roles for clarity
	var b strings.Builder
	if strings.TrimSpace(req.System) != "" {
		b.WriteString("[system]\n")
		b.WriteString(req.System)
		b.WriteString("\n\n")
	}
	for _, m := range req.Messages {
		role := strings.TrimSpace(m.Role)
		if role == "" {
			role = "user"
		}
		b.WriteString("[" + role + "]\n")
		b.WriteString(m.Content)
		b.WriteString("\n\n")
	}
	prompt := b.String()

	// Build command arguments mirroring adapter contract
	args := []string{
		"--return-mode", s.returnMode,
		"--claude-cmd", s.claudeCmd,
		"--prompt-text", prompt,
		"--suggestions-dir", s.suggestionsDir,
	}

	// Derive timeout seconds from ctx if set and none explicitly configured
	timeoutSec := s.timeoutSec
	if timeoutSec == 0 {
		if deadline, ok := ctx.Deadline(); ok {
			// Round down to seconds; ensure at least 1 second when positive
			secs := int(time.Until(deadline).Seconds())
			if secs > 0 {
				if secs == 0 {
					secs = 1
				}
				timeoutSec = secs
			}
		}
	}
	if timeoutSec > 0 {
		args = append(args, "--timeout", fmt.Sprintf("%d", timeoutSec))
	}

	// Environment for adapter
	env := os.Environ()
	env = append(env,
		"CLAUDE_CMD="+s.claudeCmd,
		"RETURN_MODE="+s.returnMode,
		// PROMPT_TEXT is set via flag; avoid duplicating in env to keep single source
	)

	// Prepare exec
	cmd := exec.CommandContext(ctx, s.adapterPath, args...)
	cmd.Env = env

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		// Surface adapter stderr for diagnostics
		return ChatResp{}, fmt.Errorf("shell adapter failed: %w\n%s", err, stderr.String())
	}

	out := strings.TrimSpace(stdout.String())
	if out == "" {
		return ChatResp{Text: ""}, nil
	}

	// If adapter returned a path (suggestions destination), read it
	if looksLikePath(out) && fileExists(out) {
		content, err := os.ReadFile(out)
		if err == nil {
			return ChatResp{Text: string(content)}, nil
		}
		// If read fails, fall back to raw output
	}
	return ChatResp{Text: out}, nil
}

func looksLikePath(s string) bool {
	if strings.Contains(s, string(os.PathSeparator)) {
		if strings.HasSuffix(s, ".patch") || strings.HasSuffix(s, ".txt") || strings.HasSuffix(s, ".diff") {
			return true
		}
		// also consider existing file check by caller
	}
	return false
}

func fileExists(p string) bool {
	if p == "" {
		return false
	}
	st, err := os.Stat(p)
	return err == nil && !st.IsDir()
}

// no extra helpers
