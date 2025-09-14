package cmd

import "github.com/spf13/cobra"

func newBriefCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "brief",
		Short: "Emit a compact task brief for IDE use",
		RunE: func(c *cobra.Command, args []string) error {
			code := CmdBrief(args)
			return codeToErr(code)
		},
	}
	return cmd
}
