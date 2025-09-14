# Research: Add agent parent command

- Date: 2025-09-14
- Base Hash: f0d3f9a
- Participants: Agent (GPT-5) / Human (kelvin)

## 1. Problem Statement
We need a parent `CmdAgent` that provides the `tgs agent` command group and routes to a subcommand that executes model adapters via the existing `agent exec` implementation. Provide help/usage, argument parsing, and dispatch to `NewAgentExecCommand` for `exec`. Keep consistent with the current CLI pattern that uses Go's `flag` package (not Cobra).

Desired outcome:
- `tgs agent` prints help/usage and available subcommands.
- `tgs agent exec [flags...]` delegates to the adapter execution pipeline in `src/cmd/agent_exec.go`.
- Return proper exit codes and error messages via stderr; stdout for command output.

## 2. Current State
- `src/main.go` already includes a switch case for `"agent"` calling `cmdpkg.CmdAgent(subArgs)`, but `CmdAgent` does not exist yet (will not compile until added).
- `src/cmd/agent_exec.go` implements `NewAgentExecCommand(args []string) (int, error)` handling flags, validation, env, and process execution for an adapter script (default `tgs/adapters/claude-code.sh`).
- Other commands (e.g., `CmdBrief`, `CmdPlan`, `CmdApprove`) use the standard library `flag` package and return `int` exit codes.

## 3. Constraints & Assumptions
- Maintain CLI style: stdlib `flag`, simple subcommand parsing, consistent error handling.
- Non-interactive, deterministic behavior; use absolute paths in examples where relevant.
- Do not place production code under `tgs/`.
- Avoid breaking existing commands in `main.go`.
- Keep indentation and formatting consistent with repository style.

## 4. Risks & Impact
- Parsing errors or mismatched flags could cause confusing UX.
- Divergence in help text vs implementation may cause support burden.
- If `agent exec` surface changes in the future, `CmdAgent` must remain a thin router to minimize coupling.

## 5. Alternatives Considered
- A) Implement `CmdAgent` with nested `flag.FlagSet`s and a manual subcommand switch. Pros: minimal deps, consistent with codebase. Cons: manual help formatting.
- B) Adopt Cobra/Viper to model command groups. Pros: richer UX. Cons: new dependency, inconsistent with current codebase, larger refactor.
- C) Keep only `tgs agent-exec` flat command. Pros: simplest. Cons: goes against the desired `agent` parent group UX and future subcommands.

## 6. Recommendation
Implement A) a thin `CmdAgent(args []string) int` using stdlib `flag` with two modes:
- No subcommand or `-h/--help`: print concise `agent` usage and `exec` synopsis with key flags.
- Subcommand `exec`: delegate to `NewAgentExecCommand(subArgs)` and map `(code,error)` to exit code and stderr output. Ensure unknown subcommands print help and return code 2.

## 7. References & Links
- Code paths: `src/main.go`, `src/cmd/agent_exec.go`, other `src/cmd/*.go` for style.
- Adapter script default path: `adapters/claude-code.sh` (ensure existence/runtime checks handled in `agent_exec.go`).

---
Approval checkpoint: Please review this research and reply one of:
- APPROVE research
- REQUEST CHANGES: <notes>
