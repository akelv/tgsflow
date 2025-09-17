# Implementation Summary: Refactor CLI to Cobra/Viper

## 1. Overview (What & Why)
Migrated the `tgs` CLI to use Cobra for structured commands and Viper for optional config/env handling. This preserves existing behavior and flags while improving UX (consistent help, subcommand hierarchy) and enabling extensibility and shell completions.

## 2. File Changes
- Added:
  - `src/cmd/root.go` — Cobra root with persistent `--json`, `--version`; Viper init; `Execute()`.
  - `src/cmd/completion.go` — `completion` subcommand (bash, zsh, fish, powershell).
  - `src/cmd/init_cmd.go`, `context_cmd.go`, `specify_cmd.go`, `plan_cmd.go`, `tasks_cmd.go`, `approve_cmd.go`, `verify_cmd.go`, `brief_cmd.go` — Cobra wrappers calling existing command logic.
  - `src/cmd/agent_cmd.go`, `agent_exec_cmd.go` — Cobra `agent` parent and `agent exec` delegating to `NewAgentExecCommand`.
- Modified:
  - `src/main.go` — replaced manual flag/dispatch with `cmd.Execute(version, commit, date)`.
- Dependencies:
  - Added `github.com/spf13/cobra`, `github.com/spf13/viper`.

## 3. Commands & Migrations
- No migrations. Optional `tgs.yaml` continues to be read via existing loader; Viper is initialized for env overrides (`TGS_` prefix) without changing defaults.

## 4. How to Test
Build and run checks:
```bash
go build -o bin/tgs ./src
./bin/tgs --version
./bin/tgs
./bin/tgs completion zsh | head -20
```
Smoke tests for existing commands:
```bash
./bin/tgs init --decorate --ci-template none
./bin/tgs plan
./bin/tgs tasks --validate || true
./bin/tgs approve --ci || true
./bin/tgs agent exec --prompt-text hi --context README.md --adapter-path /bin/true || true
go test ./...
```

## 5. Integration Steps
- No config changes required. Existing CI and scripts should continue to work. Optionally source completion output in shell profiles.

## 6. Rollback
- Revert the commits that introduced Cobra/Viper files and restored the previous `src/main.go` dispatch.

## 7. Follow-ups & Next Steps
- Add unit tests for root/help/version and completion command if desired (plan includes test scope).
- Consider exposing more config via Viper with explicit flags mapping.

## 8. Links
- Thought directory: `tgs/b48976e-refactor-cli-to-cobra-viper/`
