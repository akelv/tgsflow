package thoughts

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"
)

var thoughtDirRe = regexp.MustCompile(`^[0-9a-f]{7,}-`)

// LocateActiveDir returns the active Thought directory.
// Priority:
// 1) TGS_THOUGHT_DIR env if it exists
// 2) Most recently modified subdir under tgs/ matching <hash>-*
// 3) Fallback: tgs/ root
func LocateActiveDir(repoRoot string) string {
	if p := os.Getenv("TGS_THOUGHT_DIR"); p != "" {
		if st, err := os.Stat(p); err == nil && st.IsDir() {
			return p
		}
	}
	tgsRoot := filepath.Join(repoRoot, "tgs")
	// Prefer new layout tgs/thoughts/*, fallback to legacy tgs/*
	thoughtsRoot := filepath.Join(tgsRoot, "thoughts")
	scanRoots := []string{thoughtsRoot, tgsRoot}
	var entries []os.DirEntry
	var err error
	for _, root := range scanRoots {
		entries, err = os.ReadDir(root)
		if err != nil {
			continue
		}
		type cand struct {
			path string
			mod  time.Time
		}
		var cands []cand
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			name := e.Name()
			if !thoughtDirRe.MatchString(name) {
				continue
			}
			p := filepath.Join(root, name)
			st, err := os.Stat(p)
			if err != nil {
				continue
			}
			cands = append(cands, cand{path: p, mod: st.ModTime()})
		}
		if len(cands) == 1 {
			return cands[0].path
		}
		if len(cands) > 1 {
			sort.Slice(cands, func(i, j int) bool { return cands[i].mod.After(cands[j].mod) })
			return cands[0].path
		}
	}
	return tgsRoot
}

// SpecFileCandidates returns possible spec filenames.
func SpecFileCandidates() []string { return []string{"10_spec.md", "10_specs.md"} }
