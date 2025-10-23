package cmd

import (
	"github.com/spf13/cobra"
)

// newServerCommand creates the `tgs server` parent command.
func newServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Server mode for autonomous work orchestration",
		Long: `Server mode provides work queue management for approved thoughts.

Supports both push (GitHub Actions) and pull (cloud/remote sessions) models
for autonomous implementation of approved thoughts via Claude Code.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Show help if no subcommand
			return cmd.Help()
		},
	}

	// Add subcommands
	cmd.AddCommand(
		newServerBacklogCommand(),
	)

	return cmd
}
