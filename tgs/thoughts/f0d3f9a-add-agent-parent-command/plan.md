# Plan: Add agent parent command

## 1. Objectives
- Provide `tgs agent` parent command with help/usage.
- Implement `CmdAgent(args []string) int` routing subcommands.
- Support `tgs agent exec` by delegating to `NewAgentExecCommand` and mapping output/exit codes.
- Ensure consistent error handling and logging with existing commands.

## 2. Scope / Non-goals
- In-scope: CLI plumbing, help text, delegation, minimal tests, unit tests, docs.
- Out-of-scope: New adapter features, prompt/context semantics, adopting Cobra, additional subcommands beyond `exec`.

## 3. Acceptance Criteria
- Running `tgs agent` prints concise usage with available subcommands.
- Running `tgs agent exec --help` shows flags supported by exec and exits code 2 (flag package behavior) or a custom help path.
- `tgs agent exec` executes adapter per `agent_exec.go` logic and returns appropriate exit codes.
- `src/main.go` successfully builds and dispatches to `CmdAgent`.

## 4. Phases & Tasks
- Phase 1: Implement command
  - [ ] Add `src/cmd/agent.go` with `CmdAgent(args []string) int`.
  - [ ] Implement subcommand parsing: default help; `exec` delegates; unknown => help + exit 2.
  - [ ] Provide `agent` usage/help text and minimal `exec` synopsis.
- Phase 2: Wire and test
  - [ ] Ensure `main.go` compiles and routes to `CmdAgent`.
  - [ ] Manual test: run `tgs agent`, `tgs agent exec --task t --prompt-text p --context README.md` (expect adapter not found unless present; verify error path).
  - [ ] Add lightweight unit test scaffolding if feasible; otherwise manual acceptance.
- Phase 3: Docs
  - [ ] Update `tgs/<dir>/implementation.md` with how to use `tgs agent exec`.

## 5. File/Module Changes
- Add: `src/cmd/agent.go` — new parent command.
- No changes to: `src/cmd/agent_exec.go` (reuse), `src/main.go` (already routes to `CmdAgent`).

## 6. Test Plan
- Build: `go build ./...`.
- Help: `tgs agent` prints usage; `tgs agent foo` prints usage with unknown command message.
- Exec happy path: simulate with a valid adapter path or expect clear "adapter not found" error when default is missing.
- Exit codes: 0 on success; 1 on adapter error; 2 on usage errors.

## 7. Rollout & Rollback
- Rollout: regular merge; no migrations. Backward compatible. New entry point only.
- Rollback: revert the commit adding `src/cmd/agent.go`.

## 8. Estimates & Risks
- Estimate: ~1-2 hours including tests/docs.
- Risks: Help UX divergence; mitigated by mirroring existing command patterns and keeping `CmdAgent` thin.

---
Approval checkpoint: Please review this plan and reply one of:
- APPROVE plan
- REQUEST CHANGES: <notes>
