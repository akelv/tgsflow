package cmd

import "github.com/spf13/cobra"

func newPlanCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plan",
		Short: "Append or create 20_plan.md with NFR placeholders",
		RunE: func(c *cobra.Command, args []string) error {
			code := CmdPlan(args)
			return codeToErr(code)
		},
	}
	return cmd
}
