package server

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	defaultBacklogPath = "tgs/server/backlog.json"
	maxRetries         = 3
)

// GitBacklogOps provides git operations for backlog synchronization.
type GitBacklogOps struct {
	BacklogPath string
	RepoPath    string
}

// NewGitBacklogOps creates a new GitBacklogOps with default paths.
func NewGitBacklogOps(repoPath string) *GitBacklogOps {
	return &GitBacklogOps{
		BacklogPath: defaultBacklogPath,
		RepoPath:    repoPath,
	}
}

// Pull fetches and merges the latest backlog from remote.
func (g *GitBacklogOps) Pull() error {
	// Get current branch
	branch, err := g.getCurrentBranch()
	if err != nil {
		return fmt.Errorf("get current branch: %w", err)
	}

	// Pull with retry
	return g.retryWithBackoff(func() error {
		cmd := exec.Command("git", "pull", "origin", branch)
		cmd.Dir = g.RepoPath
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("git pull failed: %w\n%s", err, output)
		}
		return nil
	})
}

// Push commits and pushes the backlog to remote.
func (g *GitBacklogOps) Push(message string) error {
	// Add backlog file
	if err := g.gitAdd(g.BacklogPath); err != nil {
		return fmt.Errorf("git add: %w", err)
	}

	// Commit
	if err := g.gitCommit(message); err != nil {
		// Check if there are no changes (nothing to commit)
		if strings.Contains(err.Error(), "nothing to commit") ||
			strings.Contains(err.Error(), "no changes added") {
			return nil // Not an error
		}
		return fmt.Errorf("git commit: %w", err)
	}

	// Push with retry
	branch, err := g.getCurrentBranch()
	if err != nil {
		return fmt.Errorf("get current branch: %w", err)
	}

	return g.retryWithBackoff(func() error {
		cmd := exec.Command("git", "push", "-u", "origin", branch)
		cmd.Dir = g.RepoPath
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("git push failed: %w\n%s", err, output)
		}
		return nil
	})
}

// PullMergeAndPush performs a full sync cycle: pull -> modify -> push.
// The modifyFn is called after pull and before push.
func (g *GitBacklogOps) PullMergeAndPush(message string, modifyFn func() error) error {
	// Pull latest
	if err := g.Pull(); err != nil {
		return fmt.Errorf("pull: %w", err)
	}

	// Apply modification
	if err := modifyFn(); err != nil {
		return fmt.Errorf("modify: %w", err)
	}

	// Push changes
	if err := g.Push(message); err != nil {
		return fmt.Errorf("push: %w", err)
	}

	return nil
}

// AtomicClaimAndPush atomically claims a thought and pushes to remote.
// Uses optimistic locking: pull, claim, push, retry on conflict.
func (g *GitBacklogOps) AtomicClaimAndPush(claimedBy string) (*ThoughtEntry, error) {
	var claimed *ThoughtEntry

	err := g.retryWithBackoff(func() error {
		// Pull latest
		if err := g.Pull(); err != nil {
			return fmt.Errorf("pull: %w", err)
		}

		// Load backlog
		backlog, err := Load(g.BacklogPath)
		if err != nil {
			return fmt.Errorf("load backlog: %w", err)
		}

		// Claim next
		entry, err := backlog.ClaimNext(claimedBy)
		if err != nil {
			return fmt.Errorf("claim next: %w", err)
		}
		if entry == nil {
			return fmt.Errorf("no work available")
		}
		claimed = entry

		// Save
		if err := backlog.Save(g.BacklogPath); err != nil {
			return fmt.Errorf("save backlog: %w", err)
		}

		// Push
		message := fmt.Sprintf("chore(tgs): claim thought %s by %s", entry.Dir, claimedBy)
		if err := g.Push(message); err != nil {
			// If push fails due to conflict, retry the whole operation
			if strings.Contains(err.Error(), "rejected") ||
				strings.Contains(err.Error(), "conflict") {
				return err // Retry
			}
			return fmt.Errorf("push: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return claimed, nil
}

// getCurrentBranch returns the current git branch name.
func (g *GitBacklogOps) getCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = g.RepoPath
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("get branch: %w\n%s", err, out.String())
	}

	branch := strings.TrimSpace(out.String())
	if branch == "" {
		return "", fmt.Errorf("empty branch name")
	}

	return branch, nil
}

// gitAdd runs git add for the specified file.
func (g *GitBacklogOps) gitAdd(path string) error {
	cmd := exec.Command("git", "add", path)
	cmd.Dir = g.RepoPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git add failed: %w\n%s", err, output)
	}
	return nil
}

// gitCommit runs git commit with the given message.
func (g *GitBacklogOps) gitCommit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Dir = g.RepoPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w\n%s", err, output)
	}
	return nil
}

// retryWithBackoff retries a function with exponential backoff.
func (g *GitBacklogOps) retryWithBackoff(fn func() error) error {
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff: 2s, 4s, 8s
			delay := time.Duration(1<<uint(attempt)) * time.Second
			time.Sleep(delay)
		}

		err := fn()
		if err == nil {
			return nil
		}

		lastErr = err

		// Check if error is retryable
		if !isRetryableError(err) {
			return err
		}
	}

	return fmt.Errorf("max retries exceeded: %w", lastErr)
}

// isRetryableError checks if an error should trigger a retry.
func isRetryableError(err error) bool {
	errStr := err.Error()
	retryablePatterns := []string{
		"rejected",
		"conflict",
		"connection",
		"timeout",
		"temporary",
	}

	for _, pattern := range retryablePatterns {
		if strings.Contains(strings.ToLower(errStr), pattern) {
			return true
		}
	}

	return false
}

// IsGitRepo checks if the given path is inside a git repository.
func IsGitRepo(path string) bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = path
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	return err == nil
}

// GetRepoRoot returns the root directory of the git repository.
func GetRepoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("not in a git repo: %w\n%s", err, out.String())
	}

	root := strings.TrimSpace(out.String())
	if root == "" {
		return "", fmt.Errorf("empty git root")
	}

	// Verify the directory exists
	if _, err := os.Stat(root); err != nil {
		return "", fmt.Errorf("git root doesn't exist: %w", err)
	}

	return root, nil
}
