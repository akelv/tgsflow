package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kelvin/tgsflow/src/core/thoughts"
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
	// Seed common files if missing
	seed := map[string]string{
		filepath.Join(tgsDir, "README.md"):      "# TGS Thoughts\n\nThis directory contains Thought artifacts.",
		filepath.Join(tgsDir, "00_research.md"): "# Research\n\n",
		filepath.Join(tgsDir, "10_spec.md"):     "# Specification\n\n",
		filepath.Join(tgsDir, "20_plan.md"):     "# Plan\n\n",
		filepath.Join(tgsDir, "30_tasks.md"):    "# Tasks\n\n",
		filepath.Join(tgsDir, "40_approval.md"): "# Approval\n\n- Approver: \n- Role: \n- Date: \n",
	}
	for p, content := range seed {
		created, err := thoughts.EnsureFile(p, []byte(content))
		if err != nil {
			logx.Errorf("failed to ensure %s: %v", p, err)
			return 1
		}
		if created {
			logx.Infof("created %s", p)
		}
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
		_, err := thoughts.EnsureFile(approve, []byte(githubApproveWorkflow))
		if err != nil {
			logx.Errorf("failed to write workflow: %v", err)
			return 1
		}
		logx.Infof("ensured %s", approve)
	case "gitlab":
		// Minimal stub .gitlab-ci.yml
		_, err := thoughts.EnsureFile(".gitlab-ci.yml", []byte(gitlabCIStub))
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

const githubApproveWorkflow = `name: tgs-approve
on:
  pull_request:
    branches: [ main ]
jobs:
  approve:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - name: Build tgs
        run: |
          go build -o tgs ./...
      - name: Run approve gate
        run: |
          ./tgs approve --ci
`

const gitlabCIStub = `stages: [approve]
approve:
  stage: approve
  image: golang:1.22
  script:
    - go build -o tgs ./...
    - ./tgs approve --ci
`
