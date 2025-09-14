package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newAgentCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Agent-related subcommands",
		RunE: func(c *cobra.Command, args []string) error {
			// Mirror legacy help behavior
			fmt.Fprintln(os.Stderr, "Usage: tgs agent <subcommand> [options]")
			fmt.Fprintln(os.Stderr, "\nSubcommands:")
			fmt.Fprintln(os.Stderr, "  exec   Execute an adapter with prompt and context")
			fmt.Fprintln(os.Stderr, "\nRun 'tgs agent exec -h' to see exec-specific flags.")
			return nil
		},
	}
	cmd.AddCommand(newAgentExecCommand())
	return cmd
}
