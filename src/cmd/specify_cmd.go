package cmd

import "github.com/spf13/cobra"

func newSpecifyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "specify",
		Short: "Generate or proxy to Spec Kit for 10_spec.md",
		RunE: func(c *cobra.Command, args []string) error {
			code := CmdSpecify(args)
			return codeToErr(code)
		},
	}
	return cmd
}
