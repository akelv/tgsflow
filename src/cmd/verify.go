package cmd

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// CmdVerify runs repo-local hooks if available.
func CmdVerify(args []string) int {
	fs := flag.NewFlagSet("tgs verify", flag.ContinueOnError)
	ci := fs.Bool("ci", false, "CI mode")
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		return 2
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
