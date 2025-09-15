package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// CmdAgent is the parent command for agent-related subcommands.
// Usage:
//
//	tgs agent            # prints help
//	tgs agent exec ...   # delegates to NewAgentExecCommand
func CmdAgent(args []string) int {
	if len(args) == 0 {
		printAgentHelp()
		return 0
	}

	sub := args[0]
	subArgs := args[1:]

	switch sub {
	case "help", "--help", "-h":
		printAgentHelp()
		return 0
	case "exec":
		// Provide a help-only path via flag parsing to show exec flags
		// When users run: tgs agent exec --help, stdlib flag package returns error
		// We handle that by intercepting and printing exec usage below when needed.
		// Delegate to NewAgentExecCommand for real execution.
		code, err := NewAgentExecCommand(subArgs)
		if err != nil {
			// Match repo convention: write errors to stderr, return code
			fmt.Fprintln(os.Stderr, err.Error())
		}
		return code
	default:
		fmt.Fprintf(os.Stderr, "Unknown agent subcommand: %s\n", sub)
		printAgentHelp()
		return 2
	}
}

func printAgentHelp() {
	// Keep help concise and consistent with other commands
	fmt.Fprintln(os.Stderr, "Usage: tgs agent <subcommand> [options]")
	fmt.Fprintln(os.Stderr, "\nSubcommands:")
	fmt.Fprintln(os.Stderr, "  exec   Execute an adapter with prompt and context")
	fmt.Fprintln(os.Stderr, "\nExamples:")
	fmt.Fprintln(os.Stderr, "  tgs agent exec --prompt-text \"...\" --context README.md")
	fmt.Fprintln(os.Stderr, "\nRun 'tgs agent exec -h' to see exec-specific flags.")
}

// Cobra `agent` parent builder colocated
func newAgentCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Agent-related subcommands",
		RunE: func(c *cobra.Command, args []string) error {
			printAgentHelp()
			return nil
		},
	}
	cmd.AddCommand(newAgentExecCommand())
	return cmd
}
