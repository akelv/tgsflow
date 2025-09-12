package cmd

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/kelvin/tgsflow/src/core/thoughts"
	"github.com/kelvin/tgsflow/src/templates"
	"github.com/kelvin/tgsflow/src/util/logx"
)

// CmdPlan appends minimal plan with NFR placeholders.
func CmdPlan(args []string) int {
	fs := flag.NewFlagSet("tgs plan", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		return 2
	}
	active := thoughts.LocateActiveDir(".")
	if err := os.MkdirAll(active, 0o755); err != nil {
		logx.Errorf("mkdir thought: %v", err)
		return 1
	}
	path := filepath.Join(active, "20_plan.md")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		content, rerr := templates.Render("thought/20_plan.md.tmpl", nil)
		if rerr != nil {
			logx.Errorf("render plan: %v", rerr)
			return 1
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			logx.Errorf("write plan: %v", err)
			return 1
		}
		logx.Infof("created %s", path)
	} else {
		body := "Architecture\n\n- Components...\n\nNon-Functional Requirements\n\n- Performance:\n- Security:\n- Reliability:"
		if err := thoughts.AppendSection(path, "Plan", body); err != nil {
			logx.Errorf("update plan: %v", err)
			return 1
		}
		logx.Infof("updated %s", path)
	}
	return 0
}
