package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/kelvin/tgsflow/src/core/server"
	"github.com/spf13/cobra"
)

const defaultBacklogPath = "tgs/server/backlog.json"

// newServerBacklogCommand creates the `tgs server backlog` command group.
func newServerBacklogCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backlog",
		Short: "Manage thought backlog",
		Long:  `Manage the work queue of approved thoughts ready for implementation.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(
		newBacklogAddCommand(),
		newBacklogListCommand(),
		newBacklogNextCommand(),
		newBacklogClaimCommand(),
		newBacklogCompleteCommand(),
		newBacklogFailCommand(),
		newBacklogRemoveCommand(),
		newBacklogValidateCommand(),
	)

	return cmd
}

// newBacklogAddCommand creates the `tgs server backlog add` command.
func newBacklogAddCommand() *cobra.Command {
	var priority int

	cmd := &cobra.Command{
		Use:   "add <thought-dir>",
		Short: "Add a thought to the backlog",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			thoughtDir := args[0]

			// Get repo root
			repoRoot, err := server.GetRepoRoot()
			if err != nil {
				return fmt.Errorf("get repo root: %w", err)
			}

			backlogPath := filepath.Join(repoRoot, defaultBacklogPath)

			// Load backlog
			backlog, err := server.Load(backlogPath)
			if err != nil {
				return fmt.Errorf("load backlog: %w", err)
			}

			// Extract title from thought
			thoughtPath := filepath.Join(repoRoot, thoughtDir)
			title := extractThoughtTitle(thoughtPath)
			if title == "" {
				title = filepath.Base(thoughtDir)
			}

			// Add thought
			if err := backlog.AddThought(thoughtDir, title, priority); err != nil {
				return fmt.Errorf("add thought: %w", err)
			}

			// Save backlog
			if err := backlog.Save(backlogPath); err != nil {
				return fmt.Errorf("save backlog: %w", err)
			}

			// Commit and push
			gitOps := server.NewGitBacklogOps(repoRoot)
			message := fmt.Sprintf("chore(tgs): add thought %s to backlog", thoughtDir)
			if err := gitOps.Push(message); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to push backlog: %v\n", err)
				fmt.Fprintf(os.Stderr, "Backlog saved locally. Push manually with: git push\n")
			}

			fmt.Printf("Added to backlog: %s\n", thoughtDir)
			return nil
		},
	}

	cmd.Flags().IntVarP(&priority, "priority", "p", 0, "Priority (higher = more urgent)")

	return cmd
}

// newBacklogListCommand creates the `tgs server backlog list` command.
func newBacklogListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all thoughts in backlog",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get repo root
			repoRoot, err := server.GetRepoRoot()
			if err != nil {
				return fmt.Errorf("get repo root: %w", err)
			}

			backlogPath := filepath.Join(repoRoot, defaultBacklogPath)

			// Load backlog
			backlog, err := server.Load(backlogPath)
			if err != nil {
				return fmt.Errorf("load backlog: %w", err)
			}

			thoughts := backlog.ListThoughts()

			if len(thoughts) == 0 {
				fmt.Println("Backlog is empty")
				return nil
			}

			// Print as table
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "DIR\tTITLE\tSTATUS\tPRIORITY\tADDED\tCLAIMED BY")

			for _, t := range thoughts {
				claimedBy := "-"
				if t.ClaimedBy != nil {
					claimedBy = *t.ClaimedBy
				}

				fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\t%s\n",
					t.Dir,
					truncate(t.Title, 40),
					t.Status,
					t.Priority,
					t.AddedAt.Format("2006-01-02"),
					claimedBy,
				)
			}

			w.Flush()

			return nil
		},
	}

	return cmd
}

// newBacklogNextCommand creates the `tgs server backlog next` command.
func newBacklogNextCommand() *cobra.Command {
	var claimedBy string

	cmd := &cobra.Command{
		Use:   "next",
		Short: "Claim next available thought (pull model)",
		Long: `Atomically claims the next highest-priority queued thought.
Returns thought context (paths to research.md, plan.md, etc.).

This command is intended for pull-model workflows where remote or cloud
Claude Code sessions query the backlog for work.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get repo root
			repoRoot, err := server.GetRepoRoot()
			if err != nil {
				return fmt.Errorf("get repo root: %w", err)
			}

			// Auto-detect claimed_by if not specified
			if claimedBy == "" {
				hostname, _ := os.Hostname()
				claimedBy = fmt.Sprintf("auto@%s", hostname)
			}

			// Atomic claim with git push
			gitOps := server.NewGitBacklogOps(repoRoot)
			claimed, err := gitOps.AtomicClaimAndPush(claimedBy)
			if err != nil {
				if strings.Contains(err.Error(), "no work available") {
					fmt.Println("No work available in backlog")
					return nil
				}
				return fmt.Errorf("claim next: %w", err)
			}

			// Get context files
			context, err := server.GetThoughtContext(repoRoot, claimed.Dir)
			if err != nil {
				return fmt.Errorf("get context: %w", err)
			}

			// Output thought info and context
			fmt.Printf("Claimed: %s\n", claimed.Dir)
			fmt.Printf("Title: %s\n", claimed.Title)
			fmt.Printf("Priority: %d\n", claimed.Priority)
			fmt.Printf("\nContext files:\n")
			for _, file := range context {
				fmt.Printf("  %s\n", file)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&claimedBy, "claimed-by", "", "Identifier for claiming session (default: auto@hostname)")

	return cmd
}

// newBacklogClaimCommand creates the `tgs server backlog claim` command.
func newBacklogClaimCommand() *cobra.Command {
	var claimedBy string

	cmd := &cobra.Command{
		Use:   "claim <thought-dir>",
		Short: "Claim a specific thought",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			thoughtDir := args[0]

			// Get repo root
			repoRoot, err := server.GetRepoRoot()
			if err != nil {
				return fmt.Errorf("get repo root: %w", err)
			}

			if claimedBy == "" {
				hostname, _ := os.Hostname()
				claimedBy = fmt.Sprintf("auto@%s", hostname)
			}

			backlogPath := filepath.Join(repoRoot, defaultBacklogPath)

			// Load backlog
			backlog, err := server.Load(backlogPath)
			if err != nil {
				return fmt.Errorf("load backlog: %w", err)
			}

			// Get thought
			thought, err := backlog.GetThought(thoughtDir)
			if err != nil {
				return fmt.Errorf("get thought: %w", err)
			}

			if thought.Status != server.StatusQueued {
				return fmt.Errorf("thought is not queued (status: %s)", thought.Status)
			}

			// Update to in_progress
			if err := backlog.UpdateStatus(thoughtDir, server.StatusInProgress); err != nil {
				return fmt.Errorf("update status: %w", err)
			}

			// Save and push
			if err := backlog.Save(backlogPath); err != nil {
				return fmt.Errorf("save backlog: %w", err)
			}

			gitOps := server.NewGitBacklogOps(repoRoot)
			message := fmt.Sprintf("chore(tgs): claim thought %s by %s", thoughtDir, claimedBy)
			if err := gitOps.Push(message); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to push: %v\n", err)
			}

			fmt.Printf("Claimed: %s\n", thoughtDir)

			return nil
		},
	}

	cmd.Flags().StringVar(&claimedBy, "claimed-by", "", "Identifier for claiming session")

	return cmd
}

// newBacklogCompleteCommand creates the `tgs server backlog complete` command.
func newBacklogCompleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "complete <thought-dir>",
		Short: "Mark a thought as completed",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			thoughtDir := args[0]

			// Get repo root
			repoRoot, err := server.GetRepoRoot()
			if err != nil {
				return fmt.Errorf("get repo root: %w", err)
			}

			backlogPath := filepath.Join(repoRoot, defaultBacklogPath)

			// Load backlog
			backlog, err := server.Load(backlogPath)
			if err != nil {
				return fmt.Errorf("load backlog: %w", err)
			}

			// Update status
			if err := backlog.UpdateStatus(thoughtDir, server.StatusCompleted); err != nil {
				return fmt.Errorf("update status: %w", err)
			}

			// Save and push
			if err := backlog.Save(backlogPath); err != nil {
				return fmt.Errorf("save backlog: %w", err)
			}

			gitOps := server.NewGitBacklogOps(repoRoot)
			message := fmt.Sprintf("chore(tgs): complete thought %s", thoughtDir)
			if err := gitOps.Push(message); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to push: %v\n", err)
			}

			fmt.Printf("Marked as completed: %s\n", thoughtDir)

			return nil
		},
	}

	return cmd
}

// newBacklogFailCommand creates the `tgs server backlog fail` command.
func newBacklogFailCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fail <thought-dir>",
		Short: "Mark a thought as failed",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			thoughtDir := args[0]

			// Get repo root
			repoRoot, err := server.GetRepoRoot()
			if err != nil {
				return fmt.Errorf("get repo root: %w", err)
			}

			backlogPath := filepath.Join(repoRoot, defaultBacklogPath)

			// Load backlog
			backlog, err := server.Load(backlogPath)
			if err != nil {
				return fmt.Errorf("load backlog: %w", err)
			}

			// Update status
			if err := backlog.UpdateStatus(thoughtDir, server.StatusFailed); err != nil {
				return fmt.Errorf("update status: %w", err)
			}

			// Save and push
			if err := backlog.Save(backlogPath); err != nil {
				return fmt.Errorf("save backlog: %w", err)
			}

			gitOps := server.NewGitBacklogOps(repoRoot)
			message := fmt.Sprintf("chore(tgs): mark thought %s as failed", thoughtDir)
			if err := gitOps.Push(message); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to push: %v\n", err)
			}

			fmt.Printf("Marked as failed: %s\n", thoughtDir)

			return nil
		},
	}

	return cmd
}

// newBacklogRemoveCommand creates the `tgs server backlog remove` command.
func newBacklogRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <thought-dir>",
		Short: "Remove a thought from backlog",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			thoughtDir := args[0]

			// Get repo root
			repoRoot, err := server.GetRepoRoot()
			if err != nil {
				return fmt.Errorf("get repo root: %w", err)
			}

			backlogPath := filepath.Join(repoRoot, defaultBacklogPath)

			// Load backlog
			backlog, err := server.Load(backlogPath)
			if err != nil {
				return fmt.Errorf("load backlog: %w", err)
			}

			// Remove thought
			if err := backlog.RemoveThought(thoughtDir); err != nil {
				return fmt.Errorf("remove thought: %w", err)
			}

			// Save and push
			if err := backlog.Save(backlogPath); err != nil {
				return fmt.Errorf("save backlog: %w", err)
			}

			gitOps := server.NewGitBacklogOps(repoRoot)
			message := fmt.Sprintf("chore(tgs): remove thought %s from backlog", thoughtDir)
			if err := gitOps.Push(message); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to push: %v\n", err)
			}

			fmt.Printf("Removed from backlog: %s\n", thoughtDir)

			return nil
		},
	}

	return cmd
}

// newBacklogValidateCommand creates the `tgs server backlog validate` command.
func newBacklogValidateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate backlog consistency",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get repo root
			repoRoot, err := server.GetRepoRoot()
			if err != nil {
				return fmt.Errorf("get repo root: %w", err)
			}

			backlogPath := filepath.Join(repoRoot, defaultBacklogPath)

			// Load backlog
			backlog, err := server.Load(backlogPath)
			if err != nil {
				return fmt.Errorf("load backlog: %w", err)
			}

			// Validate
			errors := backlog.Validate()

			if len(errors) == 0 {
				fmt.Println("Backlog is valid")
				return nil
			}

			fmt.Println("Validation errors:")
			for _, e := range errors {
				fmt.Printf("  - %s\n", e)
			}

			return fmt.Errorf("backlog validation failed")
		},
	}

	return cmd
}

// extractThoughtTitle extracts title from a thought's README.md.
func extractThoughtTitle(thoughtPath string) string {
	readmePath := filepath.Join(thoughtPath, "README.md")
	data, err := os.ReadFile(readmePath)
	if err != nil {
		return ""
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			title := strings.TrimPrefix(line, "# ")
			// Remove base hash prefix (e.g., "abc123 - Title" -> "Title")
			parts := strings.SplitN(title, " - ", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1])
			}
			return strings.TrimSpace(title)
		}
	}

	return ""
}

// truncate truncates a string to a maximum length.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
