package cmd

import "github.com/spf13/cobra"

func newInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize TGS layout (idempotent)",
		RunE: func(c *cobra.Command, args []string) error {
			code := CmdInit(args)
			return codeToErr(code)
		},
	}
	return cmd
}
