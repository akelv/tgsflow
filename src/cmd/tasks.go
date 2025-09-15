package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/kelvin/tgsflow/src/core/thoughts"
	"github.com/kelvin/tgsflow/src/templates"
	"github.com/kelvin/tgsflow/src/util/logx"
	"github.com/spf13/cobra"
)

var taskIDRe = regexp.MustCompile("(?m)^###\\s+T[0-9]+\\.[0-9]+\\s+—\\s+")

// CmdTasks creates or validates 30_tasks.md formatting.
func CmdTasks(args []string) int {
	fs := flag.NewFlagSet("tgs tasks", flag.ContinueOnError)
	validate := fs.Bool("validate", false, "Validate existing 30_tasks.md IDs and formatting")
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		return 2
	}
	active := thoughts.LocateActiveDir(".")
	if err := os.MkdirAll(active, 0o755); err != nil {
		logx.Errorf("mkdir thought: %v", err)
		return 1
	}
	path := filepath.Join(active, "30_tasks.md")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		content, rerr := templates.Render("thought/30_tasks.md.tmpl", nil)
		if rerr != nil {
			logx.Errorf("render tasks: %v", rerr)
			return 1
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			logx.Errorf("write 30_tasks.md: %v", err)
			return 1
		}
		logx.Infof("created %s", path)
		return 0
	}
	if *validate {
		f, err := os.Open(path)
		if err != nil {
			logx.Errorf("open %s: %v", path, err)
			return 1
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		found := 0
		for scanner.Scan() {
			if taskIDRe.MatchString(scanner.Text()) {
				found++
			}
		}
		if err := scanner.Err(); err != nil {
			logx.Errorf("scan %s: %v", path, err)
			return 1
		}
		if found == 0 {
			fmt.Fprintln(os.Stderr, "no task IDs found (expected headings like '### T1.2 — Title')")
			return 1
		}
		logx.Infof("validated %s (%d task IDs)", path, found)
	} else {
		logx.Infof("%s exists; use --validate to check formatting", path)
	}
	return 0
}

func newTasksCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tasks",
		Short: "Create or validate 30_tasks.md",
		RunE: func(c *cobra.Command, args []string) error {
			return codeToErr(CmdTasks(args))
		},
	}
	return cmd
}
