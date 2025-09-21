# Implementation: Shell Transport for brain

## What & Why
- Implemented a real `shellTransport` that executes `tgs/adapters/claude-code.sh` to fulfill `Transport.Chat` in shell mode.
- Adds tests for happy path, error propagation, and timeout handling.
- Aligns with SR-020 and leverages existing adapter contract used by `tgs agent exec`.
- Expose adapter path/claude cmd via config by consumers.

## Files Changed
- Added `src/core/brain/transport_shell.go` — real implementation.
- Updated `src/core/brain/transport_stub.go` — `NewShellTransport` now returns the real transport.
- Added `src/core/brain/transport_shell_test.go` — unit tests.

## How it works
- Composes prompt from `ChatReq.System` and `ChatReq.Messages` with role headers.
- Invokes adapter with flags: `--return-mode`, `--claude-cmd`, `--prompt-text`, `--suggestions-dir`, and optional `--timeout` derived from `ctx`.
- Reads stdout; if output looks like a suggestions file path and exists, reads file contents into `ChatResp.Text`; otherwise uses text directly.
- How to use in tgs/tgs.yml:
```
    ai:
        mode: shell
        shell_adapter_path: tgs/adapters/claude-code.sh
        shell_claude_cmd: claude
```

## How to test
- Run: `make test`
- Focused tests: `go test ./src/core/brain -run TestShellTransport`

## Adapter compatibility check
- Verified `tgs/adapters/claude-code.sh` is executable and can be driven with a fake Claude command producing `FAKE_OK`.

## Rollback
- Revert `transport_shell.go`, test file, and `NewShellTransport` change in `transport_stub.go`.

## Follow-ups
- Consider mapping `ChatResp.ToolCalls` when tool execution is introduced.
