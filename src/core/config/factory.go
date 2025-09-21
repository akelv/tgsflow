package config

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/kelvin/tgsflow/src/core/thoughts"
	"github.com/kelvin/tgsflow/src/templates"
)

// TemplateData defines fields consumed by templates/data/tgs/tgs.yml.tmpl
type TemplateData struct {
	// Header
	ProjectName string

	// AI
	Provider          string
	Model             string
	APIKeyEnv         string
	TimeoutMS         int
	MaxAttempts       int
	BackoffMS         int
	ContextPackTokens int
	SummaryTokens     int
	ContextPackModel  string
	SummarizeModel    string

	// Triggers & guardrails
	IssueLabels      []string
	PRLabels         []string
	AllowPaths       []string
	DenyPaths        []string
	RequiredChecks   []string
	PRTemplate       string
	CommitConvention string
	EARS             EARSFields

	// Agent
	AgentName         string
	AgentType         string
	AgentEnabled      bool
	AgentCapabilities []string
	AgentLabelsAny    []string
	AgentPathsAny     []string
	AgentBin          string
	AgentArgs         []string
}

type EARSFields struct {
	Enable       bool
	RequireShall bool
}

// DefaultTemplateData returns sane beginner-friendly defaults.
func DefaultTemplateData(projectName string) TemplateData {
	if projectName == "" {
		if wd, err := os.Getwd(); err == nil {
			projectName = filepath.Base(wd)
		} else {
			projectName = "my-project"
		}
	}

	return TemplateData{
		ProjectName:       projectName,
		Provider:          "openai",
		Model:             "gpt-4o-mini",
		APIKeyEnv:         "OPENAI_API_KEY",
		TimeoutMS:         45000,
		MaxAttempts:       2,
		BackoffMS:         800,
		ContextPackTokens: 1200,
		SummaryTokens:     400,
		ContextPackModel:  "gpt-4o-mini",
		SummarizeModel:    "gpt-4o-mini",

		IssueLabels:      []string{"good-first-agent", "agent:fix", "agent:docs"},
		PRLabels:         []string{"agent:auto", "needs:reviewer"},
		AllowPaths:       []string{"cmd/", "internal/", "pkg/", "src/", "docs/", "tgs/"},
		DenyPaths:        []string{"infra/prod/", "deploy/", "secrets/"},
		RequiredChecks:   []string{"lint", "unit", "sast"},
		PRTemplate:       ".github/PULL_REQUEST_TEMPLATE.md",
		CommitConvention: "conventional",
		EARS: EARSFields{
			Enable:       false,
			RequireShall: false,
		},

		AgentName:         "aider-main",
		AgentType:         "cli_editor",
		AgentEnabled:      true,
		AgentCapabilities: []string{"edit_files", "run_tests", "open_pr"},
		AgentLabelsAny:    []string{"agent:aider", "agent:fix"},
		AgentPathsAny:     []string{"cmd/", "internal/"},
		AgentBin:          "aider",
		AgentArgs:         []string{"--no-gitignore", "--read"},
	}
}

// RenderConfigYAML renders the YAML from the embedded template.
func RenderConfigYAML(data TemplateData) (string, error) {
	return templates.Render("tgs/tgs.yml.tmpl", data)
}

// EnsureConfigFile writes tgs/tgs.yml if it does not already exist.
func EnsureConfigFile(repoRoot string, yaml string) (created bool, path string, err error) {
	path = filepath.Join(repoRoot, "tgs", "tgs.yml")
	created, err = thoughts.EnsureFile(path, []byte(yaml))
	return created, path, err
}

// PromptInteractive allows advanced users to customize select fields.
// It never fails hard; invalid input results in keeping defaults.
func PromptInteractive(r io.Reader, w io.Writer, data TemplateData) TemplateData {
	in := bufio.NewScanner(r)

	// Provider
	fmt.Fprintln(w, "Select AI provider [1=openai, 2=anthropic, 3=azure, 4=vertex] (default 1):")
	if in.Scan() {
		switch strings.TrimSpace(in.Text()) {
		case "2":
			data.Provider = "anthropic"
			data.Model = "claude-3-5-sonnet"
			data.APIKeyEnv = "ANTHROPIC_API_KEY"
			data.ContextPackModel = data.Model
			data.SummarizeModel = data.Model
		case "3":
			data.Provider = "azure"
			data.Model = "gpt-4o-mini"
			data.APIKeyEnv = "AZURE_OPENAI_API_KEY"
			data.ContextPackModel = data.Model
			data.SummarizeModel = data.Model
		case "4":
			data.Provider = "vertex"
			data.Model = "gemini-1.5-flash"
			data.APIKeyEnv = "GOOGLE_API_KEY"
			data.ContextPackModel = data.Model
			data.SummarizeModel = data.Model
		default:
			// keep defaults
		}
	}

	// Model (optional)
	fmt.Fprintf(w, "Model [%s]: ", data.Model)
	if in.Scan() {
		txt := strings.TrimSpace(in.Text())
		if txt != "" {
			data.Model = txt
			data.ContextPackModel = txt
			data.SummarizeModel = txt
		}
	}

	// API key env
	fmt.Fprintf(w, "API key env var [%s]: ", data.APIKeyEnv)
	if in.Scan() {
		txt := strings.TrimSpace(in.Text())
		if txt != "" {
			data.APIKeyEnv = txt
		}
	}

	// Enable EARS?
	fmt.Fprintf(w, "Enable EARS linter? (y/N): ")
	if in.Scan() {
		txt := strings.ToLower(strings.TrimSpace(in.Text()))
		if txt == "y" || txt == "yes" {
			data.EARS.Enable = true
		}
	}

	return data
}
