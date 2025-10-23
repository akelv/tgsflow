package server

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const thoughtsDir = "tgs/thoughts"

// DiscoveredThought represents a thought discovered as implementation-ready.
type DiscoveredThought struct {
	Dir   string
	Title string
}

// DiscoverApprovedThoughts finds thoughts with research.md and plan.md committed, but no implementation.md.
func DiscoverApprovedThoughts(repoPath string) ([]DiscoveredThought, error) {
	thoughtsPath := filepath.Join(repoPath, thoughtsDir)

	// Check if thoughts directory exists
	if _, err := os.Stat(thoughtsPath); err != nil {
		if os.IsNotExist(err) {
			return []DiscoveredThought{}, nil // No thoughts yet
		}
		return nil, fmt.Errorf("stat thoughts dir: %w", err)
	}

	// List all thought directories
	entries, err := os.ReadDir(thoughtsPath)
	if err != nil {
		return nil, fmt.Errorf("read thoughts dir: %w", err)
	}

	var discovered []DiscoveredThought

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		thoughtDir := filepath.Join(thoughtsDir, entry.Name())
		thoughtPath := filepath.Join(repoPath, thoughtDir)

		// Check if this thought is implementation-ready
		ready, title, err := isImplementationReady(repoPath, thoughtPath)
		if err != nil {
			// Log error but continue with other thoughts
			fmt.Fprintf(os.Stderr, "Warning: failed to check thought %s: %v\n", thoughtDir, err)
			continue
		}

		if ready {
			discovered = append(discovered, DiscoveredThought{
				Dir:   thoughtDir,
				Title: title,
			})
		}
	}

	return discovered, nil
}

// isImplementationReady checks if a thought has research.md and plan.md committed, but no implementation.md.
func isImplementationReady(repoPath, thoughtPath string) (bool, string, error) {
	researchPath := filepath.Join(thoughtPath, "research.md")
	planPath := filepath.Join(thoughtPath, "plan.md")
	implPath := filepath.Join(thoughtPath, "implementation.md")

	// Check if implementation.md exists
	if _, err := os.Stat(implPath); err == nil {
		return false, "", nil // Already implemented
	}

	// Check if research.md exists
	if _, err := os.Stat(researchPath); err != nil {
		return false, "", nil // No research.md
	}

	// Check if plan.md exists
	if _, err := os.Stat(planPath); err != nil {
		return false, "", nil // No plan.md
	}

	// Check if research.md is committed
	researchCommitted, err := isFileCommitted(repoPath, researchPath)
	if err != nil {
		return false, "", fmt.Errorf("check research committed: %w", err)
	}
	if !researchCommitted {
		return false, "", nil
	}

	// Check if plan.md is committed
	planCommitted, err := isFileCommitted(repoPath, planPath)
	if err != nil {
		return false, "", fmt.Errorf("check plan committed: %w", err)
	}
	if !planCommitted {
		return false, "", nil
	}

	// Extract title from README.md if available
	title := extractTitle(thoughtPath)
	if title == "" {
		title = filepath.Base(thoughtPath)
	}

	return true, title, nil
}

// isFileCommitted checks if a file has been committed to git.
func isFileCommitted(repoPath, filePath string) (bool, error) {
	// Make path relative to repo root
	relPath, err := filepath.Rel(repoPath, filePath)
	if err != nil {
		return false, fmt.Errorf("get relative path: %w", err)
	}

	// Check if file is tracked by git (appears in git ls-files)
	cmd := exec.Command("git", "ls-files", "--error-unmatch", relPath)
	cmd.Dir = repoPath
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		// File not tracked or doesn't exist in git
		return false, nil
	}

	// Check if file has uncommitted changes
	cmd = exec.Command("git", "diff", "--quiet", "HEAD", "--", relPath)
	cmd.Dir = repoPath

	err = cmd.Run()
	if err != nil {
		// Exit code 1 means there are differences (uncommitted changes)
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return false, nil
		}
		// Other error
		return false, fmt.Errorf("git diff: %w", err)
	}

	// File is tracked and has no uncommitted changes
	return true, nil
}

// extractTitle attempts to extract the title from a thought's README.md.
func extractTitle(thoughtPath string) string {
	readmePath := filepath.Join(thoughtPath, "README.md")
	data, err := os.ReadFile(readmePath)
	if err != nil {
		return ""
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Look for markdown heading
		if strings.HasPrefix(line, "# ") {
			title := strings.TrimPrefix(line, "# ")
			// Remove base hash suffix if present (e.g., "abc123 - Title" -> "Title")
			parts := strings.SplitN(title, " - ", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1])
			}
			return strings.TrimSpace(title)
		}
	}

	return ""
}

// GetThoughtContext returns paths to key files for a thought.
func GetThoughtContext(repoPath, thoughtDir string) ([]string, error) {
	thoughtPath := filepath.Join(repoPath, thoughtDir)

	// Verify thought directory exists
	if _, err := os.Stat(thoughtPath); err != nil {
		return nil, fmt.Errorf("thought dir not found: %w", err)
	}

	var context []string

	// Add research.md and plan.md
	files := []string{"research.md", "plan.md", "README.md"}
	for _, file := range files {
		path := filepath.Join(thoughtDir, file)
		fullPath := filepath.Join(repoPath, path)
		if _, err := os.Stat(fullPath); err == nil {
			context = append(context, path)
		}
	}

	// Add common design docs
	designFiles := []string{
		"tgs/agentops/AGENTOPS.md",
		"tgs/README.md",
		"CLAUDE.md",
		"README.md",
	}

	for _, file := range designFiles {
		fullPath := filepath.Join(repoPath, file)
		if _, err := os.Stat(fullPath); err == nil {
			context = append(context, file)
		}
	}

	return context, nil
}
