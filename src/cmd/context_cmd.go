package cmd

import "github.com/spf13/cobra"

func newContextCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "context",
		Short: "Scan repo context and seed research",
		RunE: func(c *cobra.Command, args []string) error {
			code := CmdContext(args)
			return codeToErr(code)
		},
	}
	return cmd
}
