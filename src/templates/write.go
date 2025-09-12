package templates

import (
	"os"
	"path/filepath"

	"github.com/kelvin/tgsflow/src/core/thoughts"
)

// WriteIfMissing renders tmplName with data and writes it to rootDir/relPath
// if the target file does not already exist. It ensures parent directories exist.
func WriteIfMissing(rootDir string, relPath string, tmplName string, data any) error {
	outPath := filepath.Join(rootDir, relPath)
	if _, err := os.Stat(outPath); err == nil {
		// already exists
		return nil
	}
	content, err := Render(tmplName, data)
	if err != nil {
		return err
	}
	_, err = thoughts.EnsureFile(outPath, []byte(content))
	return err
}
