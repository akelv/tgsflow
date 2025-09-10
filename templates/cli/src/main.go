package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "dev" // Set by build process

var rootCmd = &cobra.Command{
	Use:   "my-cli-tool",
	Short: "A cross-platform CLI tool bootstrapped with TGS workflow",
	Long: `🚀 Welcome to your new CLI tool!

This project was bootstrapped with:
  • TGS workflow for structured engineering
  • Cross-platform build support (Windows, macOS, Linux)
  • Cobra CLI framework for rich command-line interfaces
  • Standard project structure for CLI tools

Next steps:
  1. Run: make install
  2. Build: make build (or make build-all for all platforms)
  3. Start your first thought: make new-thought title="My Feature"`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 Welcome to your new CLI tool!")
		fmt.Println("This project was bootstrapped with:")
		fmt.Println("  • TGS workflow for structured engineering")
		fmt.Println("  • Cross-platform build support")
		fmt.Println("  • Cobra CLI framework")
		fmt.Println("  • Standard CLI project structure")
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Println("  1. Run: make install")
		fmt.Println("  2. Build: make build")
		fmt.Println("  3. Start your first thought: make new-thought title=\"My Feature\"")
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("my-cli-tool version %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}