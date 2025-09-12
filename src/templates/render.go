package templates

import (
	"bytes"
	"embed"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed data/**/*
var tmplFS embed.FS

// Render renders a named template from the embedded FS with the given data.
// name is a path relative to data/ (e.g., "thought/10_spec.md.tmpl").
func Render(name string, data any) (string, error) {
	normalized := filepath.ToSlash(name)
	if !strings.HasPrefix(normalized, "thought/") && !strings.HasPrefix(normalized, "ci/") {
		// allow bare names under thought/
		normalized = filepath.ToSlash(filepath.Join("thought", normalized))
	}
	content, err := tmplFS.ReadFile(filepath.Join("data", normalized))
	if err != nil {
		return "", fmt.Errorf("template read failed for %s: %w", normalized, err)
	}
	t, err := template.New(normalized).Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("template parse failed for %s: %w", normalized, err)
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("template execute failed for %s: %w", normalized, err)
	}
	return buf.String(), nil
}
