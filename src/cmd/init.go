package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kelvin/tgsflow/src/core/thoughts"
	"github.com/kelvin/tgsflow/src/templates"
	"github.com/kelvin/tgsflow/src/util/logx"
)

// CmdInit implements `tgs init --decorate [--ci-template github|gitlab|none]`.
func CmdInit(args []string) int {
	fs := flag.NewFlagSet("tgs init", flag.ContinueOnError)
	decorate := fs.Bool("decorate", true, "Create tgs/ layout and optional CI templates")
	ciTemplate := fs.String("ci-template", "github", "CI template: github|gitlab|none")
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if !*decorate {
		logx.Infof("Nothing to do (decorate=false)")
		return 0
	}

	// Create tgs/ skeleton
	tgsDir := filepath.Join("tgs")
	if err := thoughts.EnsureDir(tgsDir); err != nil {
		logx.Errorf("failed to create tgs dir: %v", err)
		return 1
	}
	// Seed common files if missing via templates
	renderTo := func(rel string, tmpl string, data any) error {
		outPath := filepath.Join(tgsDir, rel)
		if _, err := os.Stat(outPath); err == nil {
			return nil
		}
		content, err := templates.Render(tmpl, data)
		if err != nil {
			return err
		}
		_, err = thoughts.EnsureFile(outPath, []byte(content))
		if err != nil {
			return err
		}
		logx.Infof("created %s", outPath)
		return nil
	}
	if err := renderTo("00_research.md", "thought/00_research.md.tmpl", nil); err != nil {
		logx.Errorf("seed 00: %v", err)
		return 1
	}
	if err := renderTo("10_spec.md", "thought/10_spec.md.tmpl", map[string]any{"Role": "user", "Outcome": "..."}); err != nil {
		logx.Errorf("seed 10: %v", err)
		return 1
	}
	if err := renderTo("20_plan.md", "thought/20_plan.md.tmpl", nil); err != nil {
		logx.Errorf("seed 20: %v", err)
		return 1
	}
	if err := renderTo("30_tasks.md", "thought/30_tasks.md.tmpl", nil); err != nil {
		logx.Errorf("seed 30: %v", err)
		return 1
	}
	if err := renderTo("40_approval.md", "thought/40_approval.md.tmpl", nil); err != nil {
		logx.Errorf("seed 40: %v", err)
		return 1
	}

	// Optional CI templates
	switch strings.ToLower(*ciTemplate) {
	case "github":
		wfDir := filepath.Join(".github", "workflows")
		if err := thoughts.EnsureDir(wfDir); err != nil {
			logx.Errorf("failed to ensure workflows dir: %v", err)
			return 1
		}
		approve := filepath.Join(wfDir, "tgs-approve.yml")
		content, err := templates.Render("ci/github-approve.yml.tmpl", nil)
		if err != nil {
			logx.Errorf("render workflow: %v", err)
			return 1
		}
		_, err = thoughts.EnsureFile(approve, []byte(content))
		if err != nil {
			logx.Errorf("failed to write workflow: %v", err)
			return 1
		}
		logx.Infof("ensured %s", approve)
	case "gitlab":
		// Minimal stub .gitlab-ci.yml
		content, err := templates.Render("ci/gitlab-ci.yml.tmpl", nil)
		if err != nil {
			logx.Errorf("render gitlab ci: %v", err)
			return 1
		}
		_, err = thoughts.EnsureFile(".gitlab-ci.yml", []byte(content))
		if err != nil && !errors.Is(err, os.ErrExist) {
			logx.Errorf("failed to ensure gitlab ci: %v", err)
			return 1
		}
	case "none":
		// no-op
	default:
		fmt.Fprintln(os.Stderr, "unknown --ci-template value; expected github|gitlab|none")
		return 2
	}
	logx.Infof("tgs init complete (idempotent)")
	return 0
}
