package config

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	AI         AI         `yaml:"ai"`
	Triggers   Triggers   `yaml:"triggers"`
	Guardrails Guardrails `yaml:"guardrails"`
	Agents     []Agent    `yaml:"agents"`
	Steps      Steps      `yaml:"steps"`
	Telemetry  Telemetry  `yaml:"telemetry"`
	Context    Context    `yaml:"context"`
}

func Default() Config {
	return Config{
		AI: AI{
			Mode:      "shell",
			Provider:  "openai",
			Model:     "gpt-4o-mini",
			Endpoint:  "",
			APIKeyEnv: "OPENAI_API_KEY",
			TimeoutMS: 45000,
			// Shell transport specific defaults
			ShellAdapterPath: "tgs/adapters/claude-code.sh",
			ShellClaudeCmd:   "claude",
			Retry: AIRetry{
				MaxAttempts: 2,
				BackoffMS:   800,
			},
			Toolpack: AIToolpack{
				Enabled:  true,
				AllowFor: []string{"context_pack", "need_trace", "plan_summarize", "classify_issue"},
				Budgets: map[string]int{
					"context_pack_tokens": 1200,
					"summary_tokens":      400,
				},
				Routes: map[string]string{
					"context_pack_model": "gpt-4o-mini",
					"summarize_model":    "gpt-4o-mini",
				},
				Tools: []AITool{
					{Name: "fetch_repo_text", Desc: "Read small text slices by path+line range for ranking/quotes."},
					{Name: "list_candidates", Desc: "Return candidate doc sections with path, anchor, and score."},
					{Name: "propose_brief", Desc: "Return ordered brief sections within a token budget."},
				},
				Redaction: AIRedaction{
					RedactEnvKeys:  []string{"API_KEY", "TOKEN", "PASSWORD"},
					RedactPatterns: []string{"(?i)secret\\s*[:=]\\s*['\\\"][^'\\\"]+['\\\"]"},
				},
			},
		},
		Triggers: Triggers{
			IssueLabels: []string{"good-first-agent", "agent:fix", "agent:docs"},
			PRLabels:    []string{"agent:auto", "needs:reviewer"},
		},
		Guardrails: Guardrails{
			AllowPaths:       []string{"cmd/", "internal/", "pkg/", "src/", "docs/", "tgs/"},
			DenyPaths:        []string{"infra/prod/", "deploy/", "secrets/"},
			MaxDiffLines:     800,
			RequiredChecks:   []string{"lint", "unit", "sast"},
			PRTemplate:       ".github/PULL_REQUEST_TEMPLATE.md",
			CommitConvention: "conventional",
			EARS: EARSConfig{
				Enable:       false,
				RequireShall: false,
				Paths:        []string{"tgs/design/10_needs.md", "tgs/design/20_requirements.md"},
			},
		},
		Agents: []Agent{},
		Steps: Steps{
			PlanPrompt:   "tgs/prompts/plan.md",
			ImplPrompt:   "tgs/prompts/impl.md",
			ReviewPrompt: "tgs/prompts/review.md",
		},
		Telemetry: Telemetry{
			LogDir:      "tgs/.tgs/logs",
			RedactRules: []string{"api_key", "token", "password"},
		},
		Context: Context{
			PackDir:      "tgs/design/",
			ThoughtsDir:  "tgs/thoughts/",
			IncludeGlobs: []string{"README.md", "docs/**/*.md", "tgs/**/*.md"},
		},
	}
}

func Load(repoRoot string) (Config, error) {
	cfg := Default()
	// Prefer new location tgs/tgs.yml; fallback to legacy tgs.yaml
	primary := filepath.Join(repoRoot, "tgs", "tgs.yml")
	data, err := os.ReadFile(primary)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			legacy := filepath.Join(repoRoot, "tgs.yaml")
			if data, err = os.ReadFile(legacy); err != nil {
				if errors.Is(err, fs.ErrNotExist) {
					return cfg, nil
				}
				return cfg, err
			}
		} else {
			return cfg, err
		}
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

// AI defines the intelligence configuration for TGS commands.
type AI struct {
	Mode      string `yaml:"mode"`
	Provider  string `yaml:"provider"`
	Model     string `yaml:"model"`
	Endpoint  string `yaml:"endpoint"`
	APIKeyEnv string `yaml:"api_key_env"`
	TimeoutMS int    `yaml:"timeout_ms"`
	// Shell transport customization
	ShellAdapterPath string     `yaml:"shell_adapter_path"`
	ShellClaudeCmd   string     `yaml:"shell_claude_cmd"`
	Retry            AIRetry    `yaml:"retry"`
	Toolpack         AIToolpack `yaml:"toolpack"`
}

type AIRetry struct {
	MaxAttempts int `yaml:"max_attempts"`
	BackoffMS   int `yaml:"backoff_ms"`
}

type AIToolpack struct {
	Enabled   bool              `yaml:"enabled"`
	AllowFor  []string          `yaml:"allow_for"`
	Budgets   map[string]int    `yaml:"budgets"`
	Routes    map[string]string `yaml:"routes"`
	Tools     []AITool          `yaml:"tools"`
	Redaction AIRedaction       `yaml:"redaction"`
}

type AITool struct {
	Name string `yaml:"name"`
	Desc string `yaml:"desc"`
}

type AIRedaction struct {
	RedactEnvKeys  []string `yaml:"redact_env_keys"`
	RedactPatterns []string `yaml:"redact_patterns"`
}

type Triggers struct {
	IssueLabels []string `yaml:"issue_labels"`
	PRLabels    []string `yaml:"pr_labels"`
}

type Guardrails struct {
	AllowPaths       []string   `yaml:"allow_paths"`
	DenyPaths        []string   `yaml:"deny_paths"`
	MaxDiffLines     int        `yaml:"max_diff_lines"`
	RequiredChecks   []string   `yaml:"required_checks"`
	PRTemplate       string     `yaml:"pr_template"`
	CommitConvention string     `yaml:"commit_convention"`
	EARS             EARSConfig `yaml:"ears"`
}

type Agent struct {
	Name         string        `yaml:"name"`
	Type         string        `yaml:"type"`
	Enabled      bool          `yaml:"enabled"`
	Capabilities []string      `yaml:"capabilities"`
	Selector     AgentSelector `yaml:"selector"`
	Runtime      AgentRuntime  `yaml:"runtime"`
}

type AgentSelector struct {
	LabelsAny []string `yaml:"labels_any"`
	PathsAny  []string `yaml:"paths_any"`
}

type AgentRuntime struct {
	// CLI editor fields
	Bin  string   `yaml:"bin"`
	Args []string `yaml:"args"`
	// Hosted reviewer fields
	Provider string             `yaml:"provider"`
	Events   AgentRuntimeEvents `yaml:"events"`
	// Orchestrated agent fields
	Endpoint string `yaml:"endpoint"`
	AuthEnv  string `yaml:"auth_env"`
}

type AgentRuntimeEvents struct {
	OnPR              bool     `yaml:"on_pr"`
	OnCommentCommands []string `yaml:"on_comment_commands"`
}

type Steps struct {
	PlanPrompt   string `yaml:"plan_prompt"`
	ImplPrompt   string `yaml:"impl_prompt"`
	ReviewPrompt string `yaml:"review_prompt"`
}

type Telemetry struct {
	LogDir      string   `yaml:"log_dir"`
	RedactRules []string `yaml:"redact_rules"`
}

type Context struct {
	PackDir      string   `yaml:"pack_dir"`
	ThoughtsDir  string   `yaml:"thoughts_dir"`
	IncludeGlobs []string `yaml:"include_globs"`
}

type EARSConfig struct {
	Enable       bool     `yaml:"enable"`
	RequireShall bool     `yaml:"require_shall"`
	Paths        []string `yaml:"paths"`
}
