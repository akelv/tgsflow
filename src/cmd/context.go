//CmdContext implements the `tgs context` command.
//It allows the user to pack context files into a single brief for the AI agent.

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/kelvin/tgsflow/src/core/brain"
	"github.com/kelvin/tgsflow/src/core/config"
	"github.com/kelvin/tgsflow/src/core/thoughts"
	"github.com/spf13/cobra"
)

func newContextCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "context",
		Short: "Context-related utilities",
		RunE: func(c *cobra.Command, args []string) error {
			return c.Help()
		},
	}
	cmd.AddCommand(newContextPackCommand())
	return cmd
}

func newContextPackCommand() *cobra.Command {
	var (
		flagOut     string
		flagVerbose bool
	)

	cmd := &cobra.Command{
		Use:   "pack <query>",
		Short: "Pack relevant design and thought context into an AI brief",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			repoRoot, err := os.Getwd()
			if err != nil {
				return err
			}
			cfg, err := config.Load(repoRoot)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			query := strings.TrimSpace(strings.Join(args, " "))
			if query == "" {
				return exitCodeError{code: 2}
			}

			// Locate active thought directory
			activeThought := thoughts.LocateActiveDir(repoRoot)
			if _, err := os.Stat(activeThought); err != nil {
				return fmt.Errorf("active thought dir not found: %s", activeThought)
			}

			// Resolve output path
			outPath := flagOut
			if strings.TrimSpace(outPath) == "" {
				outPath = filepath.Join(activeThought, "aibrief.md")
			}
			if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
				return err
			}

			// Gather candidate context files from design and active thought
			designDir := cfg.Context.PackDir
			if strings.TrimSpace(designDir) == "" {
				designDir = filepath.Join("tgs", "design")
			}
			designDir = filepath.Clean(designDir)
			if !filepath.IsAbs(designDir) {
				designDir = filepath.Join(repoRoot, designDir)
			}

			candidateGlobs := []string{
				filepath.Join(designDir, "*.md"),
			}
			thoughtFiles := []string{
				"README.md", "research.md", "plan.md", "implementation.md",
			}
			for _, spec := range thoughts.SpecFileCandidates() {
				thoughtFiles = append(thoughtFiles, spec)
			}

			var ctxFiles []string
			for _, g := range candidateGlobs {
				matches, _ := filepath.Glob(g)
				ctxFiles = append(ctxFiles, matches...)
			}
			for _, f := range thoughtFiles {
				p := filepath.Join(activeThought, f)
				if st, err := os.Stat(p); err == nil && !st.IsDir() {
					ctxFiles = append(ctxFiles, p)
				}
			}
			// De-duplicate preserving order
			seen := make(map[string]struct{}, len(ctxFiles))
			finalCtx := make([]string, 0, len(ctxFiles))
			for _, p := range ctxFiles {
				if _, ok := seen[p]; ok {
					continue
				}
				seen[p] = struct{}{}
				finalCtx = append(finalCtx, p)
			}
			if len(finalCtx) == 0 {
				return fmt.Errorf("no context files found in %s or %s", designDir, activeThought)
			}

			// Load prompt templates (fallback to embedded defaults)
			searchPrompt := mustLoadPrompt(repoRoot, filepath.Join("tgs", "agentops", "prompts", "context_search.md"), defaultSearchPrompt())
			briefTemplate := mustLoadPrompt(repoRoot, filepath.Join("tgs", "agentops", "prompts", "context_brief.md"), defaultBriefTemplate())

			// Fill variables
			budget := brain.Budget(cfg, "context_pack_tokens", 1200)
			prompt := strings.ReplaceAll(searchPrompt, "{{QUERY}}", query)
			prompt = strings.ReplaceAll(prompt, "{{TOKEN_BUDGET}}", fmt.Sprintf("%d", budget))
			prompt = strings.ReplaceAll(prompt, "{{BRIEF_TEMPLATE}}", briefTemplate)

			// Prepare adapter exec (reuse adapter contract)
			adapterPath := cfg.AI.ShellAdapterPath
			if strings.TrimSpace(adapterPath) == "" {
				adapterPath = filepath.Join("tgs", "adapters", "claude-code.sh")
			}
			if _, err := os.Stat(adapterPath); err != nil {
				return fmt.Errorf("adapter not found: %s", adapterPath)
			}

			execArgs := []string{
				"--return-mode", "text",
				"--claude-cmd", cfg.AI.ShellClaudeCmd,
				"--prompt-text", prompt,
				"--out", outPath,
				"--suggestions-dir", filepath.Join("tgs", "suggestions"),
			}
			// Timeout (seconds) derived from config, if set
			if cfg.AI.TimeoutMS > 0 {
				sec := cfg.AI.TimeoutMS / 1000
				if sec > 0 {
					execArgs = append(execArgs, "--timeout", fmt.Sprintf("%d", sec))
				}
			}

			// Environment: CONTEXT_FILES newline separated
			env := os.Environ()
			env = append(env,
				"CLAUDE_CMD="+cfg.AI.ShellClaudeCmd,
				"RETURN_MODE=text",
				"CONTEXT_FILES="+strings.Join(finalCtx, "\n"),
			)

			execCmd := exec.Command(adapterPath, execArgs...)
			execCmd.Env = env
			execCmd.Stdout = os.Stdout
			execCmd.Stderr = os.Stderr

			if flagVerbose {
				fmt.Fprintf(os.Stderr, "tgs: context pack using %s with %d context files\n", adapterPath, len(finalCtx))
			}

			if err := execCmd.Run(); err != nil {
				return exitCodeError{code: 1}
			}

			if flagVerbose {
				fmt.Fprintf(os.Stderr, "wrote brief: %s\n", outPath)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&flagOut, "out", "", "Output path for brief (default: <active-thought>/aibrief.md)")
	cmd.Flags().BoolVar(&flagVerbose, "verbose", false, "Verbose logs")
	return cmd
}

func mustLoadPrompt(repoRoot, relPath, fallback string) string {
	p := relPath
	if !filepath.IsAbs(p) {
		p = filepath.Join(repoRoot, relPath)
	}
	if b, err := os.ReadFile(p); err == nil {
		if s := strings.TrimSpace(string(b)); s != "" {
			return s
		}
	}
	return fallback
}

func defaultSearchPrompt() string {
	return strings.TrimSpace(`You are an engineering assistant. Given a user query, analyze the provided repository context files to extract only the most relevant information. Focus on:

1) A short problem framing for the query.
2) Key stakeholder needs and system requirements directly related to the query.
3) Pointers to exact sources (file path plus anchor/section or line range) for verification.

Constraints:
- Keep the final brief within {{TOKEN_BUDGET}} tokens.
- Do not include secrets or credentials; if present, redact.
- Prefer EARS-style needs and “The system shall …” requirements.

Output:
Return ONLY the brief using the following structure and style. Do not add any preamble or commentary:

{{BRIEF_TEMPLATE}}

User query: "{{QUERY}}"`)
}

func defaultBriefTemplate() string {
	return strings.TrimSpace(`# AI Brief

## Query
"{{QUERY}}"

## Context Summary (<= 6 bullets)
- ...

## Key Needs (EARS-style, with sources)
- [ID or tag] One-line need statement. (Source: path#anchor)
- ...

## Key System Requirements (with sources)
- [SR-###] The system shall … (Source: path#anchor)
- ...

## Links & Pointers
- path:anchor – why relevant
- ...

## Notes
- Token budget: {{TOKEN_BUDGET}}`)
}
