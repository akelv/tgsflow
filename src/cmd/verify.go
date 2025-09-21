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

	// Optional: EARS linter gate (legacy-compatible via policies.ears.enable)
	if cfg.Policies.EARS.Enable {
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
	return cmd
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
