package cmd

import "github.com/spf13/cobra"

func newTasksCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tasks",
		Short: "Create or validate 30_tasks.md",
		RunE: func(c *cobra.Command, args []string) error {
			code := CmdTasks(args)
			return codeToErr(code)
		},
	}
	return cmd
}
