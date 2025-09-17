# Research: Refactor CLI to Cobra/Viper

- Date: 2025-09-14
- Base Hash: b48976e
- Participants: @koderizer

## 1. Problem Statement
Refactor the `tgs` CLI to use Cobra for command structure and Viper for configuration. Goals:
- Preserve existing behavior, flags, exit codes, and outputs for all current commands.
- Improve developer UX: consistent `--help`, auto-generated usage, subcommand hierarchy, completion.
- Enable future extensibility (nested commands like `tgs agent exec`, shared persistent flags, config/env overrides).

## 2. Current State
- Entry point: `src/main.go` uses `flag.FlagSet` for global flags (`--version`, `--json`) and manual subcommand dispatch via `switch`.
- Subcommands implemented as functions in `src/cmd/`: `CmdInit`, `CmdContext`, `CmdSpecify`, `CmdPlan`, `CmdTasks`, `CmdApprove`, `CmdVerify`, `CmdBrief`, plus an `agent` parent with `NewAgentExecCommand` wiring.
- Tests exist for agent exec and some commands; behavior expectations include specific exit codes (2 for usage errors, 1 for runtime failures, 0 for success) and stderr/stdout usage.
- Config loading via YAML at `tgs.yaml` using `core/config/loader.go`; no env var overrides today.

## 3. Constraints & Assumptions
- Must not break existing scripts/CI relying on current commands and flags.
- Keep log behavior: `--json` toggles JSONL logs via `util/logx`.
- Maintain adapter execution semantics for `tgs agent exec` including flag names and validation.
- Add dependencies: `github.com/spf13/cobra` and `github.com/spf13/viper` in `go.mod`.
- Support macOS/Linux shells; no interactive prompts for normal usage.

## 4. Risks & Impact
- Risk of flag incompatibilities or default behavior changes when moving to Cobra.
- Help/usage text changes may affect snapshots if any tests assert exact strings.
- Potential regression in exit codes if error handling is not mapped precisely.
- Build size increase due to Cobra/Viper dependencies.

Mitigations:
- Mirror existing flags and behaviors; keep validation logic in existing functions where feasible.
- Use Cobra’s `SilenceUsage/SilenceErrors` and explicit return codes to match prior behavior.
- Add compatibility shims where necessary; run `go test ./...` to ensure tests still pass.

## 5. Alternatives Considered
- Continue with `flag` and homegrown dispatch: minimal deps but scales poorly; harder to add nested commands and completions.
- Use `urfave/cli`: good UX, but team familiarity and ecosystem around Cobra in Go CLIs is stronger; Viper integration is standard with Cobra.
- Hybrid: Keep `flag` for inner commands, use Cobra only at the root: reduces churn but yields inconsistent UX and partial benefits.

## 6. Recommendation
Adopt Cobra for the root command and all subcommands. Wrap existing command logic within Cobra `RunE` handlers that call current `Cmd*` functions (or extract shared core into reusable functions) to minimize behavioral drift. Integrate Viper at root to load `tgs.yaml` and optionally override with env vars (prefix `TGS_`), but keep current loader as the source of truth for now; expose config to commands via context or package import as today. Provide shell completion generation commands.

## 7. References & Links
- Code paths: `src/main.go`, `src/cmd/*`, `src/core/config/loader.go`, `src/util/logx/logx.go`
- Cobra: https://github.com/spf13/cobra
- Viper: https://github.com/spf13/viper

---
Approval checkpoint: Please review this research and reply one of:
- APPROVE research
- REQUEST CHANGES: <notes>
