package cmd

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/kelvin/tgsflow/src/core/config"
	"github.com/kelvin/tgsflow/src/core/ears"
	"github.com/spf13/cobra"
)

// CmdVerify runs repo-local hooks if available.
func CmdVerify(args []string) int {
	fs := flag.NewFlagSet("tgs verify", flag.ContinueOnError)
	ci := fs.Bool("ci", false, "CI mode")
	repoRoot := fs.String("repo", ".", "Repository root path")
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		return 2
	}

	// Load config
	cfg, err := config.Load(*repoRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		if *ci {
			return 1
		}
	}

	// Optional: EARS linter gate (default false)
	if cfg.Guardrails.EARS.Enable {
		issues := verifyEARS(*repoRoot)
		for _, is := range issues {
			fmt.Fprintln(os.Stderr, is)
		}
		if len(issues) > 0 && *ci {
			return 1
		}
	}

	// Minimal placeholder: look for .tgs/hooks/* and execute if present
	hooks := []string{"fmt", "lint", "test", "perf"}
	for _, h := range hooks {
		path := ".tgs/hooks/" + h
		if _, err := os.Stat(path); err == nil {
			// run
			if err := runHook(path); err != nil {
				fmt.Fprintf(os.Stderr, "hook %s failed: %v\n", h, err)
				if *ci {
					return 1
				}
			}
		}
	}
	fmt.Fprintln(os.Stderr, "verify: hooks completed")
	return 0
}

func runHook(path string) error {
	cmd := exec.Command(path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// in future we may add scoped args like --since

func newVerifyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify",
		Short: "Run hooks/policy/drift checks",
		RunE: func(c *cobra.Command, args []string) error {
			return codeToErr(CmdVerify(args))
		},
	}

	// Add subcommand: verify ears
	earsCmd := &cobra.Command{
		Use:   "ears",
		Short: "Lint EARS requirements (design docs by default)",
		RunE: func(c *cobra.Command, args []string) error {
			return codeToErr(CmdVerifyEARS(args))
		},
	}
	cmd.AddCommand(earsCmd)
	return cmd
}

// CmdVerifyEARS lints only configured EARS paths (defaults to design docs).
func CmdVerifyEARS(args []string) int {
	fs := flag.NewFlagSet("tgs verify ears", flag.ContinueOnError)
	repoRoot := fs.String("repo", ".", "Repository root path")
	ci := fs.Bool("ci", false, "CI mode")
	// optional override: --paths comma,separated
	pathsFlag := fs.String("paths", "", "Comma-separated list of paths to lint (defaults from config)")
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		return 2
	}

	cfg, err := config.Load(*repoRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		if *ci {
			return 1
		}
	}
	// Resolve paths
	var paths []string
	if strings.TrimSpace(*pathsFlag) != "" {
		for _, p := range strings.Split(*pathsFlag, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				paths = append(paths, p)
			}
		}
	} else if len(cfg.Guardrails.EARS.Paths) > 0 {
		paths = append(paths, cfg.Guardrails.EARS.Paths...)
	} else {
		paths = []string{"tgs/design/10_needs.md", "tgs/design/20_requirements.md"}
	}

	var (
		issues        []string
		totalCaptured int
		totalValid    int
		totalInvalid  int
	)
	type fileCounts struct{ captured, valid, invalid int }
	perFile := make(map[string]*fileCounts)
	for _, rel := range paths {
		path := filepath.Join(*repoRoot, rel)
		data, err := os.ReadFile(path)
		if err != nil {
			// Missing files should not crash; report and continue
			fmt.Fprintf(os.Stderr, "verify ears: cannot read %s: %v\n", rel, err)
			if *ci {
				return 1
			}
			continue
		}
		lines := strings.Split(string(data), "\n")
		isReqDoc := strings.HasSuffix(rel, "20_requirements.md")
		if _, ok := perFile[rel]; !ok {
			perFile[rel] = &fileCounts{}
		}
		fc := perFile[rel]
		// helper to strip leading bold ID markers like **SR-001**: or **N-001**:
		sanitize := func(s string) string {
			if strings.HasPrefix(s, "**") {
				// find closing '**'
				if idx := strings.Index(s[2:], "**:"); idx >= 0 {
					return strings.TrimSpace(s[2+idx+3:])
				}
				if idx := strings.Index(s[2:], "**:"); idx == -1 {
					if end := strings.Index(s[2:], "**"); end >= 0 {
						rest := s[2+end+2:]
						rest = strings.TrimPrefix(rest, ":")
						return strings.TrimSpace(rest)
					}
				}
			}
			return s
		}
		// Reuse core linter behavior: scan respecting code fences and bullets via minimal adapter
		inFence := false
		bulletResponseMode := false
		for i, ln := range lines {
			raw := ln
			trimmed := strings.TrimSpace(raw)
			if strings.HasPrefix(trimmed, "```") {
				inFence = !inFence
				continue
			}
			if inFence || trimmed == "" || strings.HasPrefix(trimmed, "#") {
				if trimmed == "" {
					bulletResponseMode = false
				}
				continue
			}
			isBullet := strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") || strings.HasPrefix(trimmed, "1. ")
			if isBullet && bulletResponseMode {
				continue
			}
			if !isBullet {
				upper := strings.ToUpper(trimmed)
				// Attempt parse for EARS-shaped lines. For requirements doc, enforce even without 'shall' to surface missing shall.
				if strings.HasPrefix(upper, "WHEN ") || strings.HasPrefix(upper, "WHILE ") || strings.HasPrefix(upper, "IF ") || strings.HasPrefix(upper, "THE ") {
					if !isReqDoc && !strings.Contains(strings.ToLower(trimmed), " shall") {
						// In needs doc, skip lines without 'shall'.
						continue
					}
					totalCaptured++
					fc.captured++
					if _, err := ears.ParseRequirement(trimmed); err != nil {
						issues = append(issues, fmt.Sprintf("%s:%d: %s", rel, i+1, err.Error()))
						totalInvalid++
						fc.invalid++
					} else {
						totalValid++
						fc.valid++
					}
					if strings.HasSuffix(trimmed, ":") && strings.Contains(strings.ToLower(trimmed), " shall") {
						bulletResponseMode = true
					} else {
						bulletResponseMode = false
					}
					continue
				}
			}
			if isBullet {
				candidate := strings.TrimSpace(trimmed[2:])
				candidate = sanitize(candidate)
				upper := strings.ToUpper(candidate)
				hasStarter := strings.HasPrefix(upper, "WHEN ") || strings.HasPrefix(upper, "WHILE ") || strings.HasPrefix(upper, "IF ") || strings.HasPrefix(upper, "THE ")
				if isReqDoc {
					// In requirements doc, lint EARS-shaped bullets regardless of 'shall' to catch missing shall
					if hasStarter {
						totalCaptured++
						fc.captured++
						if _, err := ears.ParseRequirement(candidate); err != nil {
							issues = append(issues, fmt.Sprintf("%s:%d: %s", rel, i+1, err.Error()))
							totalInvalid++
							fc.invalid++
						} else {
							totalValid++
							fc.valid++
						}
						continue
					}
				} else {
					// In needs doc, only lint if explicit 'shall' is present to avoid over-flagging
					if !strings.Contains(strings.ToLower(candidate), " shall") {
						continue
					}
					if hasStarter {
						totalCaptured++
						fc.captured++
						if _, err := ears.ParseRequirement(candidate); err != nil {
							issues = append(issues, fmt.Sprintf("%s:%d: %s", rel, i+1, err.Error()))
							totalInvalid++
							fc.invalid++
						} else {
							totalValid++
							fc.valid++
						}
						continue
					}
				}
				continue
			}
		}
	}

	for _, is := range issues {
		fmt.Fprintln(os.Stderr, is)
	}
	// Per-file summaries in provided order
	for _, rel := range paths {
		if fc := perFile[rel]; fc != nil {
			fmt.Fprintf(os.Stderr, "verify ears: %s captured=%d valid=%d invalid=%d\n", rel, fc.captured, fc.valid, fc.invalid)
		}
	}
	fmt.Fprintf(os.Stderr, "verify ears: captured=%d valid=%d invalid=%d\n", totalCaptured, totalValid, totalInvalid)
	if len(issues) > 0 && *ci {
		return 1
	}
	return 0
}

// verifyEARS is a temporary placeholder that will be replaced by the real linter integration.
// It scans markdown files for bullet lines and returns issue strings.
func verifyEARS(repoRoot string) []string {
	var issues []string
	filepath.WalkDir(repoRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			// skip vendor/node_modules/.git
			base := filepath.Base(path)
			if base == ".git" || base == "node_modules" || base == "vendor" || strings.HasPrefix(base, ".") {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(strings.ToLower(path), ".md") {
			data, err := os.ReadFile(path)
			if err != nil {
				return nil
			}
			lines := strings.Split(string(data), "\n")
			inFence := false
			bulletResponseMode := false
			for i, ln := range lines {
				raw := ln
				trimmed := strings.TrimSpace(raw)
				// toggle code fence
				if strings.HasPrefix(trimmed, "```") {
					inFence = !inFence
					continue
				}
				if inFence || trimmed == "" || strings.HasPrefix(trimmed, "#") {
					if trimmed == "" {
						bulletResponseMode = false
					}
					continue
				}

				isBullet := strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") || strings.HasPrefix(trimmed, "1. ")
				if isBullet && bulletResponseMode {
					// Treat as continuation of prior requirement response; skip linting this line
					continue
				}

				// Non-bullet candidate lines that start with EARS keywords
				if !isBullet {
					upper := strings.ToUpper(trimmed)
					if strings.HasPrefix(upper, "WHEN ") || strings.HasPrefix(upper, "WHILE ") || strings.HasPrefix(upper, "IF ") || strings.HasPrefix(upper, "THE ") {
						if _, err := ears.ParseRequirement(trimmed); err != nil {
							issues = append(issues, fmt.Sprintf("%s:%d: %s", path, i+1, err.Error()))
						}
						// If this line ends with ":" and contains " shall" before it, enable bullet response mode
						if strings.HasSuffix(trimmed, ":") && strings.Contains(strings.ToLower(trimmed), " shall") {
							bulletResponseMode = true
						} else {
							bulletResponseMode = false
						}
						continue
					}
				}

				// Bullet candidate lines (top-level bullets only)
				if isBullet {
					candidate := strings.TrimSpace(trimmed[2:])
					if _, err := ears.ParseRequirement(candidate); err != nil {
						issues = append(issues, fmt.Sprintf("%s:%d: %s", path, i+1, err.Error()))
					}
					continue
				}
			}
		}
		return nil
	})
	return issues
}
