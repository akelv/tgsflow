package config

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestDefaultTemplateData_WithProjectName(t *testing.T) {
	data := DefaultTemplateData("demo-project")
	if data.ProjectName != "demo-project" {
		t.Fatalf("expected ProjectName 'demo-project', got %q", data.ProjectName)
	}
	if data.Provider != "openai" || data.Model != "gpt-4o-mini" || data.APIKeyEnv != "OPENAI_API_KEY" {
		t.Fatalf("unexpected AI defaults: provider=%q model=%q api=%q", data.Provider, data.Model, data.APIKeyEnv)
	}
	if data.TimeoutMS != 45000 || data.MaxAttempts != 2 || data.BackoffMS != 800 {
		t.Fatalf("unexpected retry/timeouts: %v %v %v", data.TimeoutMS, data.MaxAttempts, data.BackoffMS)
	}
	if len(data.IssueLabels) == 0 || len(data.PRLabels) == 0 {
		t.Fatalf("expected non-empty labels defaults")
	}
	if data.AgentName == "" || data.AgentType == "" || len(data.AgentCapabilities) == 0 {
		t.Fatalf("expected default agent fields populated")
	}
}

func TestRenderConfigYAML_ValidYAMLAndContains(t *testing.T) {
	data := DefaultTemplateData("acme")
	data.Provider = "anthropic"
	data.Model = "claude-3-5-sonnet"
	data.APIKeyEnv = "ANTHROPIC_API_KEY"

	out, err := RenderConfigYAML(data)
	if err != nil {
		t.Fatalf("RenderConfigYAML error: %v", err)
	}
	if !strings.Contains(out, "project: acme") {
		t.Fatalf("expected project name in YAML, got:\n%s", out)
	}
	if !strings.Contains(out, "provider: anthropic") || !strings.Contains(out, "model: claude-3-5-sonnet") {
		t.Fatalf("expected provider/model in YAML, got:\n%s", out)
	}
	// Validate YAML structure (loosely) by unmarshalling to a map
	var m map[string]any
	if err := yaml.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("yaml.Unmarshal failed: %v\nYAML:\n%s", err, out)
	}
	if _, ok := m["ai"]; !ok {
		t.Fatalf("expected 'ai' section in YAML")
	}
}

func TestEnsureConfigFile_CreatesOnce(t *testing.T) {
	root := t.TempDir()
	// render minimal YAML
	yamlStr, err := RenderConfigYAML(DefaultTemplateData("demo"))
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	created, path, err := EnsureConfigFile(root, yamlStr)
	if err != nil {
		t.Fatalf("EnsureConfigFile error: %v", err)
	}
	if !created {
		t.Fatalf("expected created=true on first write")
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file to exist at %s: %v", path, err)
	}
	if want := filepath.Join(root, "tgs", "tgs.yml"); path != want {
		t.Fatalf("unexpected path: got %s want %s", path, want)
	}

	// second call should not create
	created2, path2, err := EnsureConfigFile(root, yamlStr)
	if err != nil {
		t.Fatalf("EnsureConfigFile (second) error: %v", err)
	}
	if created2 {
		t.Fatalf("expected created=false on second write")
	}
	if path2 != path {
		t.Fatalf("unexpected second path: %s vs %s", path2, path)
	}
}

func TestPromptInteractive_ProviderModelAPIAndEARS(t *testing.T) {
	defaults := DefaultTemplateData("demo")
	// User selects provider=anthropic (2), sets model, sets API key env, enables EARS
	in := bytes.NewBufferString("2\nmy-model\nMY_ANTHROPIC_KEY\ny\n")
	var out bytes.Buffer
	got := PromptInteractive(in, &out, defaults)

	if got.Provider != "anthropic" {
		t.Fatalf("provider not applied: %q", got.Provider)
	}
	if got.Model != "my-model" {
		t.Fatalf("model not applied: %q", got.Model)
	}
	if got.APIKeyEnv != "MY_ANTHROPIC_KEY" {
		t.Fatalf("api key env not applied: %q", got.APIKeyEnv)
	}
	if !got.EARS.Enable {
		t.Fatalf("expected EARS.Enable=true")
	}
	// Ensure model routing mirrors model choice
	if got.ContextPackModel != got.Model || got.SummarizeModel != got.Model {
		t.Fatalf("expected routing models to mirror main model, got cp=%q sum=%q main=%q", got.ContextPackModel, got.SummarizeModel, got.Model)
	}
}
