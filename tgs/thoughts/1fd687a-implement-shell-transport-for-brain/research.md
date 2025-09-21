# Research: Implement Shell Transport for brain

## Problem
The `brain.Transport` interface currently returns a noop stub for shell mode. We need a real Shell Transport that can execute a local adapter (e.g., `tgs/adapters/claude-code.sh`) with prompt and context, collect output, and map it back to `ChatResp` reliably for downstream callers.

## Current State
- `Transport` interface: `Chat(ctx, ChatReq) (ChatResp, error)` exists in `src/core/brain/brain.go`.
- Factory `NewTransport(cfg)` wires `ai.mode` to `NewShellTransport` et al.
- `NewShellTransport` in `transport_stub.go` returns a noop that errors in `Chat`.
- CLI runner `src/cmd/agent_exec.go` already shells out to `tgs/adapters/claude-code.sh` with a robust flag/env contract and timeout handling.
- Adapter `tgs/adapters/claude-code.sh` supports `--prompt-text|--prompt-file`, `--context-list|CONTEXT_FILES`, `--return-mode`, `--out`, `--suggestions-dir`, `--claude-cmd`, and `--timeout`. It returns either stdout text or a suggestions path when content looks like a patch.

## Constraints
- Follow AGENTOPS workflow: implement only after approvals; production code belongs under `src/`, not `tgs/`.
- Prefer absolute paths in commands; use non-interactive flags; propagate deadlines from `context.Context`.
- Do not require network-only SDKs; rely on a local shell adapter and a configurable CLI (`claude`).
- Keep surface minimal: map `ChatReq` → adapter invocation; map adapter output → `ChatResp` with `Text` set and empty `ToolCalls` for now.
- Ensure macOS/Linux compatibility (bash, shasum/sha256sum present). Avoid destructive operations.

## Risks & Security
- Shell execution risks: command injection via prompt or file paths. Mitigation: pass prompt via stdin or environment, not via shell interpolation; avoid `sh -c`. Use `exec.CommandContext` with arg vector.
- Secrets leakage: ensure env redaction if logging; avoid dumping env in errors.
- Timeouts: must respect `ctx` deadline; also allow explicit timeout override.
- Non-deterministic outputs: adapter auto-routes patches to files. Caller expectations should be documented.

## Alternatives
1) Implement an HTTP proxy transport to a server that performs Claude calls. Pros: testable; Cons: additional service.
2) Use an SDK transport (Anthropic/OpenAI). Pros: fewer moving parts; Cons: violates desire to leverage the existing shell adapter and local context handling.
3) MCP transport for tool use. Out of scope for this task.

## Recommendation
Implement `shellTransport` under `src/core/brain/`:
- Config-driven fields: `AdapterPath` (default `tgs/adapters/claude-code.sh`), `ClaudeCmd` (default `claude`), `SuggestionsDir`, `ReturnMode`, optional `OutPath` for temporary capture.
- `Chat(ctx, ChatReq)` will:
  - Synthesize prompt from `System` and `Messages` (role-tagged, newline separated).
  - Create a temporary output file (when needed) to capture adapter output deterministically.
  - Build env: `CLAUDE_CMD`, `RETURN_MODE`, `CONTEXT_FILES` (newline-separated), `PROMPT_TEXT` when inline.
  - Build argv flags mirroring `agent_exec.go` for parity.
  - Respect `ctx` timeout; capture stdout/stderr; on non-zero, return error with stderr.
  - If stdout looks like a suggestions path and exists, read it to `Text`; otherwise use stdout text as `Text`.
- Return `ChatResp{Text: output}` and ignore tool calls for now.

## References
- `src/core/brain/brain.go`, `transport_stub.go`
- `src/cmd/agent_exec.go`
- `tgs/adapters/claude-code.sh`
- `src/core/config/loader.go`
