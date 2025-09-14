// cmd/agent_exec.go
package cmd

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// AgentExecOpts holds CLI options.
type AgentExecOpts struct {
	AdapterPath    string
	PromptFile     string
	PromptText     string
	ContextFiles   multiFlag // repeatable --context
	ContextList    string
	ContextGlob    string
	ReturnMode     string // patch_or_text | text
	TimeoutSec     int
	OutPath        string
	SuggestionsDir string
	ClaudeCmd      string
	Verbose        bool
}

// multiFlag is a simple repeatable flag collector.
type multiFlag []string

func (m *multiFlag) String() string { return strings.Join(*m, ",") }
func (m *multiFlag) Set(v string) error {
	*m = append(*m, v)
	return nil
}

// NewAgentExecCommand wires a `tgs agent exec` subcommand-like runner using standard flags.
// If you’re using Cobra, call this from your cobra.Command RunE; otherwise, expose `RunAgentExec()` as `main`.
func NewAgentExecCommand(args []string) (int, error) {
	var opt AgentExecOpts
	fs := flag.NewFlagSet("tgs agent exec", flag.ContinueOnError)
	fs.StringVar(&opt.AdapterPath, "adapter-path", "adapters/claude-code.sh", "Path to claude-code adapter script")
	fs.StringVar(&opt.PromptFile, "prompt-file", "", "Prompt template/text file")
	fs.StringVar(&opt.PromptText, "prompt-text", "", "Prompt text (overrides --prompt-file)")
	fs.Var(&opt.ContextFiles, "context", "Path to a context file (repeatable)")
	fs.StringVar(&opt.ContextList, "context-list", "", "File containing newline-separated context file paths")
	fs.StringVar(&opt.ContextGlob, "context-glob", "", "Shell glob for context files (expanded deterministically)")
	fs.StringVar(&opt.ReturnMode, "return-mode", "patch_or_text", "Return mode: patch_or_text | text")
	fs.IntVar(&opt.TimeoutSec, "timeout", 0, "Timeout in seconds (0 = no timeout)")
	fs.StringVar(&opt.OutPath, "out", "", "Write output to file; default prints to stdout or saves to suggestions")
	fs.StringVar(&opt.SuggestionsDir, "suggestions-dir", "tgs/suggestions", "Directory for suggested patches/text")
	fs.StringVar(&opt.ClaudeCmd, "claude-cmd", "claude", "Claude CLI command (e.g., 'claude')")
	fs.BoolVar(&opt.Verbose, "verbose", false, "Verbose logs")

	// Parse
	if err := fs.Parse(args); err != nil {
		return 2, err
	}

	// Validate minimal input
	if opt.PromptText == "" && opt.PromptFile == "" {
		return 2, errors.New("provide --prompt-text or --prompt-file")
	}
	if _, err := os.Stat(opt.AdapterPath); err != nil {
		return 2, fmt.Errorf("adapter not found: %s (%w)", opt.AdapterPath, err)
	}

	// Build final context file list
	ctxList, err := buildContextList(opt.ContextFiles, opt.ContextList, opt.ContextGlob)
	if err != nil {
		return 2, err
	}
	if len(ctxList) == 0 {
		return 2, errors.New("no context files provided (use --context/--context-list/--context-glob)")
	}
	for _, p := range ctxList {
		if _, err := os.Stat(p); err != nil {
			return 2, fmt.Errorf("context file missing: %s", p)
		}
	}

	// Read prompt if needed
	prompt := opt.PromptText
	if prompt == "" {
		b, err := os.ReadFile(opt.PromptFile)
		if err != nil {
			return 2, fmt.Errorf("read prompt file: %w", err)
		}
		prompt = string(b)
	}
	if strings.TrimSpace(prompt) == "" {
		return 2, errors.New("prompt is empty")
	}

	// Build command & env for adapter
	var cmdArgs []string
	cmdArgs = append(cmdArgs,
		"--return-mode", opt.ReturnMode,
		"--claude-cmd", opt.ClaudeCmd,
	)
	if opt.PromptText == "" && opt.PromptFile != "" {
		cmdArgs = append(cmdArgs, "--prompt-file", opt.PromptFile)
	} else {
		cmdArgs = append(cmdArgs, "--prompt-text", prompt)
	}
	if opt.ContextList != "" {
		cmdArgs = append(cmdArgs, "--context-list", opt.ContextList)
	}
	if opt.ContextGlob != "" {
		cmdArgs = append(cmdArgs, "--context-glob", opt.ContextGlob)
	}
	if opt.OutPath != "" {
		// Ensure parent dir exists
		_ = os.MkdirAll(filepath.Dir(opt.OutPath), 0o755)
		cmdArgs = append(cmdArgs, "--out", opt.OutPath)
	}
	if opt.SuggestionsDir != "" {
		_ = os.MkdirAll(opt.SuggestionsDir, 0o755)
		cmdArgs = append(cmdArgs, "--suggestions-dir", opt.SuggestionsDir)
	}
	if opt.Verbose {
		cmdArgs = append(cmdArgs, "--verbose")
	}
	if opt.TimeoutSec > 0 {
		cmdArgs = append(cmdArgs, "--timeout", fmt.Sprintf("%d", opt.TimeoutSec))
	}

	// Prepare env; pass CONTEXT_FILES as newline-separated list (adapter also supports flags)
	env := os.Environ()
	env = append(env, "CLAUDE_CMD="+opt.ClaudeCmd)
	env = append(env, "RETURN_MODE="+opt.ReturnMode)
	if opt.TimeoutSec > 0 {
		env = append(env, fmt.Sprintf("TIMEOUT_SEC=%d", opt.TimeoutSec))
	}
	if opt.OutPath != "" {
		env = append(env, "OUT_PATH="+opt.OutPath)
	}
	// CONTEXT_FILES
	env = append(env, "CONTEXT_FILES="+strings.Join(ctxList, "\n"))
	// PROMPT_TEXT (only if using inline prompt)
	if opt.PromptText != "" {
		env = append(env, "PROMPT_TEXT="+prompt)
	}

	// Exec
	ctx := context.Background()
	var cancel context.CancelFunc
	if opt.TimeoutSec > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(opt.TimeoutSec)*time.Second)
		defer cancel()
	}
	cmd := exec.CommandContext(ctx, opt.AdapterPath, cmdArgs...)
	cmd.Env = env

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if opt.Verbose {
		fmt.Fprintf(os.Stderr, "tgs: running adapter: %s %s\n", opt.AdapterPath, strings.Join(cmdArgs, " "))
	}
	err = cmd.Run()
	if err != nil {
		// include adapter stderr for diagnostics
		return 1, fmt.Errorf("adapter error: %w\n%s", err, stderr.String())
	}

	// If adapter printed a path (e.g., suggestions/*.patch), show it; otherwise relay stdout.
	out := stdout.String()
	if opt.OutPath == "" && strings.TrimSpace(out) != "" {
		fmt.Print(out)
	}
	return 0, nil
}

func buildContextList(direct multiFlag, listFile, glob string) ([]string, error) {
	ordered := make([]string, 0, len(direct)+8)

	// from repeatable --context
	ordered = append(ordered, direct...)

	// from --context-list file
	if listFile != "" {
		b, err := os.ReadFile(listFile)
		if err != nil {
			return nil, fmt.Errorf("read --context-list: %w", err)
		}
		for _, line := range strings.Split(string(b), "\n") {
			line = strings.TrimSpace(line)
			if line != "" {
				ordered = append(ordered, line)
			}
		}
	}

	// from --context-glob (deterministic: sort via Go’s filepath.Glob order)
	if glob != "" {
		matches, err := filepath.Glob(glob)
		if err != nil {
			return nil, fmt.Errorf("expand --context-glob: %w", err)
		}
		// filepath.Glob returns lexicographic order already
		ordered = append(ordered, matches...)
	}

	// de-dup while preserving order
	seen := make(map[string]struct{}, len(ordered))
	final := make([]string, 0, len(ordered))
	for _, p := range ordered {
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}
		final = append(final, p)
	}
	return final, nil
}
