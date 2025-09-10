package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "my-go-project",
	Short: "A Go project bootstrapped with TGS workflow",
	Long: `🚀 Welcome to your new Go project!

This project was bootstrapped with:
  • TGS workflow for structured engineering
  • Go modules for dependency management
  • Cobra CLI framework for command-line interfaces
  • Standard Go project structure

Next steps:
  1. Run: go mod tidy
  2. Build: go build -o bin/my-go-project cmd/main/main.go
  3. Start your first thought: make new-thought title="My Feature"`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 Welcome to your new Go project!")
		fmt.Println("This project was bootstrapped with:")
		fmt.Println("  • TGS workflow for structured engineering")
		fmt.Println("  • Go modules for dependency management")
		fmt.Println("  • Cobra CLI framework")
		fmt.Println("  • Standard Go project structure")
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Println("  1. Run: go mod tidy")
		fmt.Println("  2. Build: go build -o bin/my-go-project cmd/main/main.go")
		fmt.Println("  3. Start your first thought: make new-thought title=\"My Feature\"")
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}