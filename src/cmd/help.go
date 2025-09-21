package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// CmdHelp prints a concise usage with available commands.
func CmdHelp(_ []string) int {
	out := os.Stdout
	fmt.Fprintln(out, "Usage: tgs [--json] <command> [args]")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "Global Flags:")
	fmt.Fprintln(out, "  --json            Emit JSONL logs to stderr")
	fmt.Fprintln(out, "  --version         Print version and exit")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "Commands:")
	fmt.Fprintln(out, "  help              Show this help")
	fmt.Fprintln(out, "  init              Initialize TGS layout (idempotent)")
	fmt.Fprintln(out, "  context           Context tools (e.g., pack)")
	fmt.Fprintln(out, "  verify            Run hooks/policy checks (e.g., ears)")
	fmt.Fprintln(out, "  agent             AI adapter runner (shell adapter)")
	fmt.Fprintln(out, "  version           Print version")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "Settings & Configuration:")
	fmt.Fprintln(out, "  Config file       tgs/tgs.yml (auto-loaded); env prefix TGS_ via Viper")
	fmt.Fprintln(out, "  Key settings      ai.provider, ai.model, ai.api_key_env, ai.shell_adapter_path")
	fmt.Fprintln(out, "                   guardrails.ears.enable, guardrails.ears.paths")
	fmt.Fprintln(out, "  Example (tgs/tgs.yml):")
	fmt.Fprintln(out, "    guardrails:")
	fmt.Fprintln(out, "      ears:")
	fmt.Fprintln(out, "        enable: true")
	fmt.Fprintln(out, "        paths: [\"tgs/design/10_needs.md\", \"tgs/design/20_requirements.md\"]")
	fmt.Fprintln(out, "    ai:")
	fmt.Fprintln(out, "      provider: openai   # or anthropic")
	fmt.Fprintln(out, "      model: gpt-4o-mini # adjust per provider")
	fmt.Fprintln(out, "      api_key_env: OPENAI_API_KEY")
	fmt.Fprintln(out, "      shell_adapter_path: tgs/adapters/claude-code.sh")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "Requirements:")
	fmt.Fprintln(out, "  - Binary installed (brew or scripts/install.sh) or build from source")
	fmt.Fprintln(out, "  - AI API key exported per ai.api_key_env (e.g., OPENAI_API_KEY or ANTHROPIC_API_KEY)")
	fmt.Fprintln(out, "  - For context pack with shell adapter: adapter at ai.shell_adapter_path")
	fmt.Fprintln(out, "  - (Dev) Java + ANTLR only if regenerating EARS grammar (make ears-gen)")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "Quickstart:")
	fmt.Fprintln(out, "  1) Initialize TGS structure (safe, idempotent)")
	fmt.Fprintln(out, "  	tgs init")
	fmt.Fprintln(out, "  2) Author the project context and requirements in tgs/design/ ")
	fmt.Fprintln(out, "  3) Verify design docs with EARS")
	fmt.Fprintln(out, "  	tgs verify ears")
	fmt.Fprintln(out, "  4) To start building a new feature create a thought")
	fmt.Fprintln(out, "  	make new-thought title=\"Add feature\" spec=\"One line spec\"")
	fmt.Fprintln(out, "  5) Pack context into aibrief.md for the active thought")
	fmt.Fprintln(out, "  	tgs context pack \"<your goal>\"")
	fmt.Fprintln(out, "  6) Feed the brief to the AI agent of your choice to research, plan then get your approval before implementation.")
	fmt.Fprintln(out, "  	tgs agent exec --task <taskID> --context aibrief.md")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "Examples:")
	fmt.Fprintln(out, "  tgs verify ears")
	fmt.Fprintln(out, "  tgs context pack \"payment refund flow\" ")
	fmt.Fprintln(out, "")
	return 0
}

func newHelpCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "help",
		Short: "Show help, settings, and usage",
		RunE: func(c *cobra.Command, args []string) error {
			return codeToErr(CmdHelp(args))
		},
	}
	return cmd
}
