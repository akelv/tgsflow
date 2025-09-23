# Research: Gemini AI Brain Adapter

- Date: 2025-09-23
- Base Hash: ff006d7
- Participants: AI Agent, Human Reviewer

## 1. Problem Statement
We need a shell adapter for the TGS AI brain that invokes the Gemini CLI in the same manner as the existing Claude adapter (`tgs/adapters/claude-code.sh`). The adapter must accept a composed prompt, include deterministic context files, honor timeouts, and route outputs consistently (text or patch suggestion). This enables teams using Google Gemini to plug into the TGS workflow without changing core code.

## 2. Current State
- Existing Claude adapter: `tgs/adapters/claude-code.sh` implements:
  - `--prompt-text` or `--prompt-file`, deterministic context gathering via `--context-list` / `--context-glob`, snapshot hashing, unique `--add-dir` per context parent, best-effort timeout, stdout or file routing, and exit codes.
  - Used by `shellTransport` to execute with a composed `[system]/[role]` prompt and return `ChatResp.Text`.
- Requirements include SR-020 (shell transport) and newly added SR-027 for a Gemini adapter with the same interface.
- Templates and thought docs exist; Gemini thought scaffold present at `tgs/thoughts/ff006d7-gemini-ai-brain/`.

## 3. Constraints & Assumptions
- macOS/Linux target; bash compatible; prefer absolute paths and non-interactive flags.
- Avoid implementing production code under `tgs/` beyond adapter script.
- Gemini CLI is assumed installed and authenticated via `GOOGLE_API_KEY` or vendor-standard mechanism. Network access may be required for real runs.
- Maintain parity with Claude adapter flags for drop-in replacement where possible.

## 4. Risks & Impact
- CLI flag mismatches could break parity; mitigate by designing adapter to read prompt from stdin and provide flexible flags similar to Claude.
- Token/context size limits differ; mitigate by passing directories with `--add-dir`-like semantics if supported, or fallback to pure prompt text with explicit path listing.
- Secrets exposure in context; mitigate by avoiding raw file content injection unless the model CLI supports safe attachment semantics; keep only file paths in the prompt by default.
- Timeout handling variance; mitigate with external `timeout` when available and rely on process context cancellation from the caller.

## 5. Alternatives Considered
- Direct Go SDK integration for Gemini instead of shell adapter: tighter control but increases vendor lock and complexity; deviates from SR-020 shell transport pattern.
- Generic adapter that conditionally calls different CLIs: reduces duplication but complicates flags and error handling; harder to test in isolation.
- Minimal wrapper that only pipes prompt text: simpler but loses deterministic context handling and suggestions routing.

## 6. Recommendation
Implement `tgs/adapters/gemini-code.sh` mirroring `claude-code.sh` behavior and interface where feasible:
- Flags: `--prompt-text|--prompt-file`, `--context-list`, `--context-glob`, `--timeout`, `--out`, `--suggestions-dir`, `--gemini-cmd`.
- Deterministic context expansion, snapshot hash, unique parent dirs list.
- Compose final prompt: original prompt + "Context files (paths):\n...".
- Invoke Gemini CLI reading prompt from stdin; pass directory hints if supported by the CLI; otherwise rely on plain text mode.
- Exit non-zero on missing prompt/context or CLI failures. Route patch/text similar to Claude.

## 7. References & Links
- Thought README: `tgs/thoughts/ff006d7-gemini-ai-brain/README.md`
- Claude adapter reference: `tgs/adapters/claude-code.sh`
- Shell transport: `src/core/brain/transport_shell.go` and tests `transport_shell_test.go`
- Requirements: `tgs/design/20_requirements.md` (SR-020, SR-027)
- Needs: `tgs/design/10_needs.md` (N-004, N-012, N-019, N-023)
- V&V: `tgs/design/40_vnv.md`
- Gemini CLI docs: github.com/google-gemini/gemini-cli/docs (consult during implementation)

---
Approval checkpoint: Please review this research and reply one of:
- APPROVE research
- REQUEST CHANGES: <notes>
