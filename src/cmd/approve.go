package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kelvin/tgsflow/src/core/config"
	"github.com/kelvin/tgsflow/src/core/thoughts"
)

// CmdApprove validates presence of required files and optional roles.
func CmdApprove(args []string) int {
	fs := flag.NewFlagSet("tgs approve", flag.ContinueOnError)
	ci := fs.Bool("ci", false, "CI mode (non-zero exit on failure)")
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		return 2
	}
	repoRoot := "."
	active := thoughts.LocateActiveDir(repoRoot)
	required := []string{"20_plan.md", "30_tasks.md", "40_approval.md"}
	missing := []string{}
	for _, f := range required {
		if _, err := os.Stat(filepath.Join(active, f)); os.IsNotExist(err) {
			missing = append(missing, f)
		}
	}
	// support either 10_spec.md or 10_specs.md
	specOK := false
	for _, name := range thoughts.SpecFileCandidates() {
		if _, err := os.Stat(filepath.Join(active, name)); err == nil {
			specOK = true
			break
		}
	}
	if !specOK {
		missing = append(missing, "10_spec.md|10_specs.md")
	}
	if len(missing) > 0 {
		fmt.Fprintf(os.Stderr, "Missing required files: %s\n", strings.Join(missing, ", "))
		if *ci {
			return 1
		}
	}
	cfg, err := config.Load(repoRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "config load failed: %v\n", err)
		if *ci {
			return 1
		}
	}
	// If roles configured, ensure 40_approval.md includes them as words
	if len(cfg.ApproverRoles) > 0 {
		data, err := os.ReadFile(filepath.Join(active, "40_approval.md"))
		if err == nil {
			content := string(data)
			for _, role := range cfg.ApproverRoles {
				if !strings.Contains(content, role) {
					fmt.Fprintf(os.Stderr, "approval missing required role: %s\n", role)
					if *ci {
						return 1
					}
				}
			}
		}
	}
	// If policies enforce NFR, ensure 20_plan.md mentions NFR
	if cfg.Policies.EnforceNFR {
		data, err := os.ReadFile(filepath.Join(active, "20_plan.md"))
		if err == nil {
			if !strings.Contains(strings.ToLower(string(data)), "non-functional") {
				fmt.Fprintln(os.Stderr, "plan missing Non-Functional Requirements section")
				if *ci {
					return 1
				}
			}
		}
	}
	fmt.Fprintln(os.Stderr, "approve: checks passed")
	return 0
}
