package cmd

import (
	"fmt"

	"github.com/kelvin/tgsflow/src/util/logx"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type exitCodeError struct{ code int }

func (e exitCodeError) Error() string { return fmt.Sprintf("exit code %d", e.code) }

func codeToErr(code int) error {
	if code == 0 {
		return nil
	}
	return exitCodeError{code: code}
}

func exitCodeOf(err error) int {
	if err == nil {
		return 0
	}
	if e, ok := err.(exitCodeError); ok {
		return e.code
	}
	return 1
}

// NewRootCommand builds the Cobra root command tree.
func NewRootCommand(version, commit, date string) *cobra.Command {
	var (
		flagJSON    bool
		flagVersion bool
	)

	root := &cobra.Command{
		Use:           "tgs",
		Short:         "TGSFlow CLI",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// If no subcommand: show help or version
			if flagVersion {
				fmt.Printf("tgs %s (commit %s, built %s)\n", version, commit, date)
				return exitCodeError{code: 0}
			}
			_ = cmd.Help()
			return nil
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if flagJSON {
				logx.SetJSON(true)
			}
		},
	}

	// Persistent flags
	root.PersistentFlags().BoolVar(&flagJSON, "json", false, "Emit JSONL logs to stderr")
	root.PersistentFlags().BoolVar(&flagVersion, "version", false, "Print version and exit")

	// Initialize Viper (non-breaking; config still loaded via config.Load in commands)
	viper.SetConfigName("tgs")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("TGS")
	viper.AutomaticEnv()
	_ = viper.ReadInConfig() // ignore errors; config is optional

	// Subcommands
	root.AddCommand(
		newHelpCommand(),
		newInitCommand(),
		newContextCommand(),
		newVerifyCommand(),
		newAgentCommand(),
	)

	// Use our custom help command
	root.SetHelpCommand(newHelpCommand())

	return root
}

// Execute runs the CLI and maps errors to exit codes.
func Execute(version, commit, date string) int {
	root := NewRootCommand(version, commit, date)
	if err := root.Execute(); err != nil {
		return exitCodeOf(err)
	}
	return 0
}
