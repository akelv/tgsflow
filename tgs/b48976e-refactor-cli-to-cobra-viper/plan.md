# Plan: Refactor CLI to Cobra/Viper

## 1. Objectives
- Migrate CLI to Cobra for structured commands and consistent UX.
- Integrate Viper for config/env handling without changing current defaults.
- Preserve existing commands, flags, outputs, and exit codes.
- Provide shell completion generation (completion subcommand) and improved help.
- Keep current internal logic and tests working with minimal changes.

## 2. Scope / Non-goals
- In scope: Root and all existing subcommands (init, context, specify, plan, tasks, approve, verify, brief, agent, agent exec).
- In scope: Persistent flags (--json, --version).
- In scope: Viper wiring to read tgs.yaml and env overrides (no behavior change unless env provided).
- Non-goals: Redesign of command semantics; new features beyond completion/help; changing config schema; breaking flag names.

## 3. Acceptance Criteria
- tgs with no args prints help and exits 0.
- tgs --version prints "tgs <version> (commit <sha>, built <date>)" and exits 0.
- tgs --json <cmd> toggles JSON logging as before.
- All existing commands run with same flags and exit codes: 0 success, 1 runtime failure, 2 usage error.
- tgs help <cmd> and <cmd> --help show descriptive usage.
- tgs completion bash|zsh|fish|powershell outputs valid completion scripts.
- tgs agent exec flags/validation unchanged; tests continue to pass.

## 4. Phases & Tasks
- Phase 1: Cobra scaffolding
  - [ ] Add root command with persistent flags and version handling.
  - [ ] Bridge help and default help behavior.
- Phase 2: Subcommand wiring
  - [ ] Add Cobra commands for each existing subcommand that call current Cmd* functions.
  - [ ] Add agent parent and agent exec subcommand wired to NewAgentExecCommand.
- Phase 3: Viper integration
  - [ ] Initialize Viper at root: read tgs.yaml if present; enable env override with prefix TGS_.
  - [ ] Keep config.Load as source of truth; ensure no behavior change if env unset.
- Phase 4: Completions and help polish
  - [ ] Add completion subcommand (bash, zsh, fish, powershell).
  - [ ] Ensure concise, consistent command Short and Long descriptions.
- Phase 5: Tests & validation
  - [ ] Run go test ./... and fix any breakages.
  - [ ] Add/adjust unit tests for all commands per scope below.

## 5. File/Module Changes
- Add dependencies in go.mod:
  - github.com/spf13/cobra
  - github.com/spf13/viper
- New files:
  - src/cmd/root.go: defines RootCmd() returning *cobra.Command; sets persistent --json, --version; Execute(version, commit, date).
  - src/cmd/completion.go: Cobra completion subcommand.
  - src/cmd/*_cmd.go: per-command Cobra wrappers calling existing logic:
    - init_cmd.go → RunE: return codeToErr(CmdInit(args))
    - context_cmd.go, specify_cmd.go, plan_cmd.go, tasks_cmd.go, approve_cmd.go, verify_cmd.go, brief_cmd.go similarly.
    - agent_cmd.go (parent) and agent_exec_cmd.go → call NewAgentExecCommand.
- Modified files:
  - src/main.go: replace manual flag parsing/dispatch with Cobra root Execute(...) call; preserve version vars; maintain logx.SetJSON() via persistent flag handling.
- Compatibility helpers:
  - Add small helper codeToErr(int) error to translate legacy return codes into Cobra errors while ensuring exact process exit codes via SilenceUsage/SilenceErrors and explicit os.Exit after Execute.

## 6. Test Plan
- Run existing tests: go test ./... should pass.
- Manual checks:
  - go build -o bin/tgs ./src then run:
    - ./bin/tgs
    - ./bin/tgs --version
    - ./bin/tgs init --decorate --ci-template none
    - ./bin/tgs plan
    - ./bin/tgs tasks --validate
    - ./bin/tgs approve --ci
    - ./bin/tgs agent exec --prompt-text hi --context README.md --adapter-path /bin/true
    - ./bin/tgs completion zsh | head -20
- Add/adjust tests (if needed): basic root/help/version behavior without over-constraining help text.
- Unit test scope (Cobra wrappers and behavior):
  - Root: no args prints help and exits 0; --version output format; --json toggles logx without error.
  - init: seeds files on empty dir; idempotent; unknown ci-template returns usage code 2.
  - context: creates .context.json and seeds 00_research.md when missing.
  - specify: proxies when specify present; fallback creates 10_spec.md; respects --no-spec-kit.
  - plan: creates new plan or appends section; returns 0; errors surface as code 1.
  - tasks: creates file; --validate counts headings with taskIDRe; invalid/missing yields proper codes/messages.
  - approve: missing required files listed; --ci returns 1 on failures; roles/NFR checks enforced when configured.
  - verify: runs hooks if present; --ci causes non-zero exit when a hook fails.
  - brief: respects --format md|text; caps output lines; task filtering extracts section.
  - agent exec: existing tests maintained; Cobra path delegates to NewAgentExecCommand and preserves exit codes 0/1/2.

## 7. Rollout & Rollback
- Rollout: single PR adding Cobra/Viper and refactor; no migrations.
- CI: ensure build/test pass; optionally add a step to generate and test completion output.
- Rollback: revert PR; code remains functional with previous manual dispatch.

## 8. Risks
- Subtle flag/help differences; mitigated by keeping command logic intact and minimizing UX drift.

---
Approval checkpoint: Please review this plan and reply one of:
- APPROVE plan
- REQUEST CHANGES: <notes>
