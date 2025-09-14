package cmd

import "github.com/spf13/cobra"

func newVerifyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify",
		Short: "Run hooks/policy/drift checks",
		RunE: func(c *cobra.Command, args []string) error {
			code := CmdVerify(args)
			return codeToErr(code)
		},
	}
	return cmd
}
