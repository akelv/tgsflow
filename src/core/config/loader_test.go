package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(dir, rel, content string) (string, error) {
	p := filepath.Join(dir, rel)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		return "", err
	}
	return p, nil
}

func TestLoad_Table(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		yaml    string
		expect  func(t *testing.T, cfg Config)
		wantErr bool
	}{
		{
			name: "missing file returns defaults",
			yaml: "",
			expect: func(t *testing.T, cfg Config) {
				if cfg.AI.Mode != "shell" {
					t.Fatalf("default AI.Mode = %q", cfg.AI.Mode)
				}
				if cfg.AI.Toolpack.Budgets["context_pack_tokens"] == 0 {
					t.Fatalf("default budgets missing")
				}
				if len(cfg.Guardrails.AllowPaths) == 0 {
					t.Fatalf("default guardrails allow_paths missing")
				}
			},
		},
		{
			name: "minimal header only keeps defaults",
			yaml: "version: 0.1\nproject: x\n",
			expect: func(t *testing.T, cfg Config) {
				if cfg.AI.Provider != "openai" {
					t.Fatalf("provider default expected, got %q", cfg.AI.Provider)
				}
			},
		},
		{
			name: "ai fields decode and override defaults",
			yaml: "ai:\n  mode: proxy\n  provider: openai\n  model: gpt-4o-mini\n  endpoint: http://local\n  api_key_env: OPENAI_API_KEY\n  timeout_ms: 1234\n  retry:\n    max_attempts: 9\n    backoff_ms: 77\n  toolpack:\n    enabled: true\n    allow_for: [a,b]\n    budgets: {x: 1}\n    routes: {r1: m1}\n    tools: [{name: n1, desc: d1}]\n    redaction:\n      redact_env_keys: [K1]\n      redact_patterns: [P1]\n",
			expect: func(t *testing.T, cfg Config) {
				if cfg.AI.Mode != "proxy" {
					t.Fatalf("mode decode failed: %q", cfg.AI.Mode)
				}
				if cfg.AI.TimeoutMS != 1234 {
					t.Fatalf("timeout decode failed: %d", cfg.AI.TimeoutMS)
				}
				if cfg.AI.Retry.MaxAttempts != 9 {
					t.Fatalf("retry.max_attempts decode failed: %d", cfg.AI.Retry.MaxAttempts)
				}
				if cfg.AI.Toolpack.Budgets["x"] != 1 {
					t.Fatalf("budgets decode failed: %+v", cfg.AI.Toolpack.Budgets)
				}
				if len(cfg.AI.Toolpack.Tools) != 1 || cfg.AI.Toolpack.Tools[0].Name != "n1" {
					t.Fatalf("tools decode failed: %+v", cfg.AI.Toolpack.Tools)
				}
				if got := cfg.AI.Toolpack.Redaction.RedactEnvKeys; len(got) != 1 || got[0] != "K1" {
					t.Fatalf("redaction keys decode: %+v", got)
				}
			},
		},
		{
			name:    "invalid yaml returns error",
			yaml:    "ai: [oops: 1]",
			wantErr: true,
			expect:  func(t *testing.T, cfg Config) {},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			repo := t.TempDir()
			if tc.yaml != "" {
				if _, err := writeTemp(repo, filepath.Join("tgs", "tgs.yml"), tc.yaml); err != nil {
					t.Fatal(err)
				}
			}
			cfg, err := Load(repo)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			tc.expect(t, cfg)
		})
	}
}
