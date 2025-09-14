# Implementation Summary: Add agent parent command

## 1. Overview (What & Why)
Implemented a parent `tgs agent` command that groups agent-related functionality and delegates `exec` to the existing adapter execution pipeline. This makes room for future agent subcommands while keeping a thin wrapper over `agent_exec.go`.

## 2. File Changes
- Added `src/cmd/agent.go`: Implements `CmdAgent(args []string) int`, prints usage, and routes `exec` to `NewAgentExecCommand`.
- Added `src/cmd/agent_test.go`: Unit tests for help, unknown subcommand, usage error, and successful delegation using a temp adapter script.

## 3. Commands & Migrations
- Build: `go build ./...`
- Test: `go test ./...`

## 4. How to Test
- Help:
  - `go run ./src agent` → prints agent usage and subcommands.
- Unknown:
  - `go run ./src agent nope` → prints unknown subcommand message, exit code 2.
- Exec (adapter missing):
  - `go run ./src agent exec --task demo --prompt-text "hi" --context README.md` → error about missing adapter unless you provide `--adapter-path`.
- Exec (success path):
  - Create a simple adapter: `echo -e '#!/bin/sh\necho OK' > /tmp/adapter.sh && chmod +x /tmp/adapter.sh`
  - Create a context file: `echo ctx > /tmp/ctx.txt`
  - Run: `go run ./src agent exec --task demo --prompt-text "hi" --context /tmp/ctx.txt --adapter-path /tmp/adapter.sh` → prints `OK`, exit 0.

## 5. Integration Steps
- None beyond build/test; `src/main.go` already dispatches `agent` to `CmdAgent`.

## 6. Rollback
- Revert the commit that added `src/cmd/agent.go` and tests.

## 7. Follow-ups & Next Steps
- Consider adding `tgs agent suggest` or `review` subcommands.
- Document adapters under `adapters/` and provide a default in-repo path.

## 8. Links
- Thought: `tgs/f0d3f9a-add-agent-parent-command/`
