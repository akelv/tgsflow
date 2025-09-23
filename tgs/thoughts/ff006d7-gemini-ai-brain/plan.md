# Plan: Gemini AI Brain Adapter

## 1. Objectives
- Deliver `tgs/adapters/gemini-code.sh` providing a Claude-parity shell adapter for Gemini CLI.
- Maintain deterministic context handling, timeouts, and output routing consistent with SR-020 and SR-027.
- Update tests/docs to validate behavior via existing shell transport tests.

## 2. Scope / Non-goals
- In-scope: New adapter script; minimal docs updates; no core Go refactor.
- Non-goals: Changing `src/core/brain` transport behavior or config schema; vendor SDK integrations.

## 3. Acceptance Criteria
- AC1: `tgs/adapters/gemini-code.sh -h` prints usage including flags: `--prompt-text`, `--prompt-file`, `--context-list`, `--context-glob`, `--timeout`, `--out`, `--suggestions-dir`, `--gemini-cmd`, `--return-mode`.
- AC2: Running with valid prompt and at least one context file returns non-empty output and exit code 0.
- AC3: Missing prompt or no context files causes non-zero exit and clear stderr message.
- AC4: When output is a unified diff, the script writes to `tgs/suggestions/CTX-*.patch` and prints the path; otherwise prints text to stdout.
- AC5: Deterministic snapshot hash is printed in verbose logs and used in suggestions filename prefix.
- AC6: Works on macOS/Linux with bash, `shasum` or `sha256sum`, and optional `timeout`.

## 4. Phases & Tasks
- Phase 1: Adapter Skeleton
  - [ ] Create `tgs/adapters/gemini-code.sh` with CLI parsing and validation.
  - [ ] Implement deterministic context expansion and snapshot hashing.
- Phase 2: Invocation & I/O
  - [ ] Implement Gemini CLI invocation reading prompt from stdin; support `--gemini-cmd` env `GEMINI_CMD`.
  - [ ] Implement timeout wrapper and suggestions routing.
- Phase 3: Docs & Tests
  - [ ] Add brief usage notes in `tgs/README.md` adapters section.
  - [ ] Verify via `go test ./src/core/brain -run TestShellTransport*` overriding adapter path in tests locally.

## 5. File/Module Changes
- Add: `tgs/adapters/gemini-code.sh` (executable bash script).
- Update: `tgs/README.md` (mention Gemini adapter option and parity with Claude).
- No changes to Go code expected.

## 6. Test Plan
- Manual: Run the adapter with `--prompt-text "test" --context-glob "README.md"` and observe output/exit code.
- Automated: Reuse shell transport tests by temporarily pointing `adapterPath` to the new script in a local run; ensure success on happy path, error path, and timeout respects context cancellation.

## 7. Rollout & Rollback
- Rollout: Commit script and docs; no migrations. Default adapter remains Claude; users pass `--adapter-path tgs/adapters/gemini-code.sh` when invoking.
- Rollback: Remove the adapter file; no config/state changes.

## 8. Estimates & Risks
- Estimate: ~2-3 hours including docs and validation.
- Risks: Gemini CLI flag differences; mitigated by designing prompt via stdin and focusing on parity features we control. Fallback to text mode if patch detection is ambiguous.

---
Approval checkpoint: Please review this plan and reply one of:
- APPROVE plan
- REQUEST CHANGES: <notes>
