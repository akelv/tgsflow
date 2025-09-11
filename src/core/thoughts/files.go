package thoughts

import (
	"fmt"
	"os"
	"path/filepath"
)

// EnsureDir makes sure a directory exists.
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0o755)
}

// EnsureFile ensures a file exists with initial content if it did not exist.
func EnsureFile(path string, initial []byte) (created bool, err error) {
	if _, err = os.Stat(path); err == nil {
		return false, nil
	}
	if err := EnsureDir(filepath.Dir(path)); err != nil {
		return false, err
	}
	if err := os.WriteFile(path, initial, 0o644); err != nil {
		return false, err
	}
	return true, nil
}

// AppendSection appends a Markdown section if the file exists, or creates it with the section.
func AppendSection(path string, title string, body string) error {
	if _, err := os.Stat(path); err != nil {
		// create
		content := fmt.Sprintf("# %s\n\n%s\n", title, body)
		return os.WriteFile(path, []byte(content), 0o644)
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(fmt.Sprintf("\n\n## %s\n\n%s\n", title, body))
	return err
}
