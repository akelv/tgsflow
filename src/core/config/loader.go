package config

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Policies struct {
	ForbidPaths []string   `yaml:"forbid_paths"`
	MaxPatchLOC int        `yaml:"max_patch_loc"`
	EnforceNFR  bool       `yaml:"enforce_nfr"`
	EARS        EARSConfig `yaml:"ears"`
}

// EARSConfig controls EARS linter behavior.
type EARSConfig struct {
	Enable       bool `yaml:"enable"`
	RequireShall bool `yaml:"require_shall"`
}

type Config struct {
	ApproverRoles []string `yaml:"approver_roles"`
	AgentOrder    []string `yaml:"agent_order"`
	BranchPrefix  string   `yaml:"branch_prefix"`
	Policies      Policies `yaml:"policies"`
}

func Default() Config {
	return Config{
		ApproverRoles: []string{"EM", "TechLead"},
		AgentOrder:    []string{"claude", "gemini"},
		BranchPrefix:  "tgs/",
		Policies: Policies{
			ForbidPaths: []string{"infra/prod/", "secrets/"},
			MaxPatchLOC: 300,
			EnforceNFR:  false,
			EARS: EARSConfig{
				Enable:       false,
				RequireShall: false,
			},
		},
	}
}

func Load(repoRoot string) (Config, error) {
	cfg := Default()
	path := filepath.Join(repoRoot, "tgs.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return cfg, nil
		}
		return cfg, err
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
