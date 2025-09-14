package cmd

import "github.com/spf13/cobra"

func newApproveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approve",
		Short: "Validate approval gate (10/20/30/40)",
		RunE: func(c *cobra.Command, args []string) error {
			code := CmdApprove(args)
			return codeToErr(code)
		},
	}
	return cmd
}
