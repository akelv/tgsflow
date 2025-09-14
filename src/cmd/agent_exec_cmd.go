package cmd

import (
	"github.com/spf13/cobra"
)

func newAgentExecCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "exec",
		Short:              "Execute an adapter with prompt and context",
		DisableFlagParsing: true, // let NewAgentExecCommand handle all flags
		Args:               cobra.ArbitraryArgs,
		RunE: func(c *cobra.Command, args []string) error {
			code, err := NewAgentExecCommand(args)
			if err != nil {
				// Return code 2 for usage errors detected by flag parsing, else 1
				if code == 2 {
					return exitCodeError{code: 2}
				}
				return exitCodeError{code: 1}
			}
			return codeToErr(code)
		},
	}
	return cmd
}
