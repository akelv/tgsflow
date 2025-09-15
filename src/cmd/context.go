package cmd

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"time"

	"github.com/kelvin/tgsflow/src/core/thoughts"
	"github.com/kelvin/tgsflow/src/templates"
	"github.com/kelvin/tgsflow/src/util/logx"
	"github.com/spf13/cobra"
)

type repoContext struct {
	GeneratedAt string   `json:"generated_at"`
	Languages   []string `json:"languages"`
	Files       int      `json:"files"`
}

// CmdContext performs a lightweight scan and writes tgs/.context.json and seeds 00_research.md.
func CmdContext(args []string) int {
	fs := flag.NewFlagSet("tgs context", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		return 2
	}
	// simple heuristic: count files and detect langs by extension
	langs := map[string]bool{}
	files := 0
	_ = filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if d.Name() == ".git" || d.Name() == "node_modules" || d.Name() == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}
		files++
		ext := filepath.Ext(path)
		switch ext {
		case ".go":
			langs["go"] = true
		case ".ts", ".tsx":
			langs["typescript"] = true
		case ".js":
			langs["javascript"] = true
		case ".py":
			langs["python"] = true
		}
		return nil
	})
	var langsList []string
	for k := range langs {
		langsList = append(langsList, k)
	}
	ctx := repoContext{GeneratedAt: time.Now().Format(time.RFC3339), Languages: langsList, Files: files}
	active := thoughts.LocateActiveDir(".")
	if err := os.MkdirAll(active, 0o755); err != nil {
		logx.Errorf("mkdir thought: %v", err)
		return 1
	}
	f, err := os.Create(filepath.Join(active, ".context.json"))
	if err != nil {
		logx.Errorf("write context: %v", err)
		return 1
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	_ = enc.Encode(ctx)
	_ = f.Close()

	// seed research file if missing
	researchPath := filepath.Join(active, "00_research.md")
	if _, err := os.Stat(researchPath); os.IsNotExist(err) {
		if content, rerr := templates.Render("thought/00_research.md.tmpl", nil); rerr == nil {
			_ = os.WriteFile(researchPath, []byte(content), 0o644)
		}
	}
	logx.Infof("context written to %s/.context.json (%d files)", active, files)
	return 0
}

func newContextCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "context",
		Short: "Scan repo context and seed research",
		RunE: func(c *cobra.Command, args []string) error {
			return codeToErr(CmdContext(args))
		},
	}
	return cmd
}
