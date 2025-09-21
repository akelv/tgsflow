package templates

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/kelvin/tgsflow/src/core/thoughts"
)

//go:embed data/**/*
var tmplFS embed.FS

// Render renders a named template from the embedded FS with the given data.
// name is a path relative to data/ (e.g., "thought/10_spec.md.tmpl").
func Render(name string, data any) (string, error) {
	normalized := filepath.ToSlash(name)
	if !strings.HasPrefix(normalized, "thought/") && !strings.HasPrefix(normalized, "ci/") && !strings.HasPrefix(normalized, "tgs/") {
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

// RenderTGSTree walks the embedded templates under data/tgs and mirrors them into
// destRoot/tgs, rendering any files ending with .tmpl and copying others as-is.
// Existing files are left untouched to keep the operation idempotent.
func RenderTGSTree(destRoot string, data any) error {
	return RenderTGSTreeFromFS(tmplFS, "data/tgs", destRoot, data)
}

// RenderTGSTreeFromDir mirrors templates from a local directory (whose root
// directly contains files like design/, agentops/, README.md.tmpl, tgs.yml.tmpl)
// into destRoot/tgs, applying the same .tmpl rendering rules.
func RenderTGSTreeFromDir(srcDir string, destRoot string, data any) error {
	return RenderTGSTreeFromFS(os.DirFS(srcDir), ".", destRoot, data)
}

// RenderTGSTreeFromFS mirrors templates from any fs.FS starting at srcRoot
// into destRoot/tgs. Files ending with .tmpl are rendered with data and the
// suffix is stripped. Other files are copied as-is. Existing files are left
// untouched to keep the operation idempotent.
func RenderTGSTreeFromFS(source fs.FS, srcRoot string, destRoot string, data any) error {
	return fs.WalkDir(source, srcRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Compute relative path within tgs/
		rel := strings.TrimPrefix(filepath.ToSlash(path), filepath.ToSlash(srcRoot))
		rel = strings.TrimPrefix(rel, "/")
		outPath := filepath.Join(destRoot, "tgs", filepath.FromSlash(rel))
		if d.IsDir() {
			return thoughts.EnsureDir(filepath.Clean(outPath))
		}
		if strings.HasSuffix(outPath, ".tmpl") {
			rendered, rerr := renderTemplateFromFS(source, path, data)
			if rerr != nil {
				return rerr
			}
			outPath = strings.TrimSuffix(outPath, ".tmpl")
			_, werr := thoughts.EnsureFile(outPath, []byte(rendered))
			return werr
		}
		content, rerr := fs.ReadFile(source, path)
		if rerr != nil {
			return rerr
		}
		_, werr := thoughts.EnsureFile(outPath, content)
		return werr
	})
}

// renderEmbeddedTemplate reads and executes a template stored in the embedded FS.
func renderEmbeddedTemplate(path string, data any) (string, error) {
	content, err := tmplFS.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("template read failed for %s: %w", path, err)
	}
	t, err := template.New(filepath.ToSlash(path)).Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("template parse failed for %s: %w", path, err)
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("template execute failed for %s: %w", path, err)
	}
	return buf.String(), nil
}

// renderTemplateFromFS reads a template from an arbitrary fs.FS and executes it.
func renderTemplateFromFS(source fs.FS, path string, data any) (string, error) {
	content, err := fs.ReadFile(source, path)
	if err != nil {
		return "", fmt.Errorf("template read failed for %s: %w", path, err)
	}
	t, err := template.New(filepath.ToSlash(path)).Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("template parse failed for %s: %w", path, err)
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("template execute failed for %s: %w", path, err)
	}
	return buf.String(), nil
}

// end

// import placed at end to avoid circular commentary; actual import above via go fmt
