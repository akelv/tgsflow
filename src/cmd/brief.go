package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kelvin/tgsflow/src/core/thoughts"
)

// CmdBrief prints a compact brief for a given task string (best-effort),
// including ACs/NFRs/forbidden paths if available. Max ~200 lines.
func CmdBrief(args []string) int {
	fs := flag.NewFlagSet("tgs brief", flag.ContinueOnError)
	task := fs.String("task", "", "Task ID or search text")
	format := fs.String("format", "md", "Output format: md|text")
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		return 2
	}

	// Collect from spec/plan/tasks
	active := thoughts.LocateActiveDir(".")
	spec := ""
	for _, name := range thoughts.SpecFileCandidates() {
		if s := readIfExists(filepath.Join(active, name)); s != "" {
			spec = s
			break
		}
	}
	plan := readIfExists(filepath.Join(active, "20_plan.md"))
	tasks := readIfExists(filepath.Join(active, "30_tasks.md"))

	var b strings.Builder
	if *format == "md" {
		b.WriteString("# Task Brief\n\n")
	} else {
		b.WriteString("Task Brief\n\n")
	}
	if *task != "" {
		b.WriteString(fmt.Sprintf("Task: %s\n\n", *task))
	}
	if spec != "" {
		b.WriteString("## Acceptance Criteria\n")
		b.WriteString(extractSection(spec, "## Acceptance Criteria"))
		b.WriteString("\n")
	}
	if plan != "" {
		b.WriteString("## Non-Functional Requirements\n")
		b.WriteString(extractSection(plan, "Non-Functional"))
		b.WriteString("\n")
	}
	if tasks != "" {
		b.WriteString("## Tasks\n")
		if *task != "" {
			b.WriteString(extractTask(tasks, *task))
		} else {
			b.WriteString(headLines(tasks, 60))
		}
		b.WriteString("\n")
	}
	// Forbidden paths (defaults)
	b.WriteString("## Constraints\n")
	b.WriteString("- Forbidden paths: infra/prod/, secrets/\n")
	out := headLines(b.String(), 200)
	fmt.Fprint(os.Stdout, out)
	return 0
}

func readIfExists(p string) string {
	data, err := os.ReadFile(p)
	if err != nil {
		return ""
	}
	return string(data)
}

func headLines(s string, n int) string {
	sc := bufio.NewScanner(strings.NewReader(s))
	var lines []string
	for sc.Scan() {
		lines = append(lines, sc.Text())
		if len(lines) >= n {
			break
		}
	}
	return strings.Join(lines, "\n") + "\n"
}

func extractSection(doc string, marker string) string {
	// naive: return lines containing marker or the next ~40 lines
	idx := strings.Index(doc, marker)
	if idx < 0 {
		return headLines(doc, 40)
	}
	tail := doc[idx:]
	return headLines(tail, 80)
}

func extractTask(tasks string, needle string) string {
	// return the section containing the needle heading line
	lines := strings.Split(tasks, "\n")
	var out []string
	capture := false
	for _, ln := range lines {
		if strings.Contains(ln, needle) {
			capture = true
		}
		if capture {
			out = append(out, ln)
			if strings.HasPrefix(ln, "### ") && len(out) > 1 && strings.Contains(ln, "T") {
				// reached next task heading
				break
			}
		}
	}
	if len(out) == 0 {
		return headLines(tasks, 60)
	}
	return strings.Join(out, "\n") + "\n"
}
