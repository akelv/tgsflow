package cmd

import (
	"flag"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/kelvin/tgsflow/src/core/thoughts"
	"github.com/kelvin/tgsflow/src/templates"
	"github.com/kelvin/tgsflow/src/util/logx"
	"github.com/spf13/cobra"
)

// CmdSpecify proxies to Spec Kit if available, else creates minimal 10_spec.md.
func CmdSpecify(args []string) int {
	fs := flag.NewFlagSet("tgs specify", flag.ContinueOnError)
	noSpecKit := fs.Bool("no-spec-kit", false, "Do not proxy to Spec Kit even if available")
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		return 2
	}
	_, err := exec.LookPath("specify")
	if err == nil && !*noSpecKit {
		// proxy
		cmd := exec.Command("specify")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			logx.Errorf("specify proxy failed: %v", err)
			return 1
		}
		return 0
	}
	// fallback minimal spec via template
	active := thoughts.LocateActiveDir(".")
	if err := os.MkdirAll(active, 0o755); err != nil {
		logx.Errorf("mkdir thought: %v", err)
		return 1
	}
	path := filepath.Join(active, "10_spec.md")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		content, rerr := templates.Render("thought/10_spec.md.tmpl", map[string]any{"Role": "user", "Outcome": "..."})
		if rerr != nil {
			logx.Errorf("render spec template: %v", rerr)
			return 1
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			logx.Errorf("write 10_spec.md: %v", err)
			return 1
		}
		logx.Infof("created %s", path)
	} else {
		logx.Infof("%s exists; leaving as is", path)
	}
	return 0
}

func newSpecifyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "specify",
		Short: "Generate or proxy to Spec Kit for 10_spec.md",
		RunE: func(c *cobra.Command, args []string) error {
			return codeToErr(CmdSpecify(args))
		},
	}
	return cmd
}
