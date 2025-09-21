# Plan: Implement Shell Transport for brain

## Objectives
- Provide a real `shellTransport` implementing `Transport.Chat` to execute the Claude Code adapter.
- Add unit tests covering happy path, error path, and context deadline handling.
- Validate that `tgs/adapters/claude-code.sh` works as the shell-backed Claude executor.

## Scope
- In: `src/core/brain` transport implementation and tests; small additions to `config` if needed (defaults only).
- Out: MCP/Proxy/SDK transports; CLI behavior changes.

## Acceptance Criteria
- `NewShellTransport(cfg)` returns a working transport; `Chat` executes adapter and returns `ChatResp{Text: ...}` on success.
- Respects `context.Context` deadlines; returns on timeout with error.
- Maps `ChatReq{System, Messages, MaxTokens}` to adapter inputs; at minimum prompt composition honored.
- Unit tests pass; coverage added for shell transport logic.

## Tasks
1. Create `src/core/brain/transport_shell.go` with `shellTransport`.
2. Implement `Chat(ctx, ChatReq)`:
   - Compose prompt from system + messages.
   - Build env: `CLAUDE_CMD`, `RETURN_MODE`, `PROMPT_TEXT`, `CONTEXT_FILES` (empty allowed for now).
   - Build args: `--return-mode`, `--claude-cmd`, `--prompt-text`, `--suggestions-dir`.
   - Use `exec.CommandContext` with `cfg.AI.TimeoutMS` fallback to `ctx` deadline.
   - Capture stdout/stderr; if stdout is a path to suggestions, read file content to `Text`.
   - Return `ChatResp` with `Text` and no `ToolCalls`.
3. Wire `NewShellTransport(cfg)` to return the real implementation.
4. Tests in `src/core/brain/transport_shell_test.go`:
   - Fake adapter script written to temp dir, executable, echoing deterministic text.
   - Happy path: returns text.
   - Error path: adapter exits non-zero; ensure error bubbles with stderr.
   - Timeout path: adapter sleeps > timeout; ensure context timeout triggers.
5. Optional: small helper to detect patch path vs text; reuse `agent_exec.go` heuristic if simple.
6. Validate compatibility: basic invocation parity with `agent_exec.go` and `claude-code.sh` flags/env.

## File-by-file changes
- `src/core/brain/transport_shell.go`: new file with implementation.
- `src/core/brain/transport_stub.go`: update `NewShellTransport` to return real transport.
- `src/core/brain/transport_shell_test.go`: new unit tests.

## Test Plan
- `go test ./...` on macOS.
- Inject temp adapter absolute path via cfg or constructor option.

## Rollout/Rollback
- Rollout: add new file; no API changes. Safe.
- Rollback: revert files.

## Estimates
- Implementation: 1-2 hours
- Tests: 1 hour
