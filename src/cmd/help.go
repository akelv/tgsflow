package cmd

import (
	"fmt"
	"os"
)

// CmdHelp prints a concise usage with available commands.
func CmdHelp(_ []string) int {
	fmt.Fprintln(os.Stderr, "Usage: tgs [--json] <command> [args]")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Commands:")
	fmt.Fprintln(os.Stderr, "  help             Show this help")
	fmt.Fprintln(os.Stderr, "  version          Print version")
	fmt.Fprintln(os.Stderr, "  init             Initialize TGS layout (idempotent)")
	fmt.Fprintln(os.Stderr, "  context          Scan repo context and seed research")
	fmt.Fprintln(os.Stderr, "  specify          Generate or proxy to Spec Kit for 10_spec.md")
	fmt.Fprintln(os.Stderr, "  plan             Append or create 20_plan.md with NFR placeholders")
	fmt.Fprintln(os.Stderr, "  tasks            Create or validate 30_tasks.md")
	fmt.Fprintln(os.Stderr, "  approve          Validate approval gate (10/20/30/40)")
	fmt.Fprintln(os.Stderr, "  verify           Run hooks/policy/drift checks")
	fmt.Fprintln(os.Stderr, "  brief            Emit a compact task brief for IDE use")
	return 0
}
