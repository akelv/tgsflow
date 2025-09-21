# Implementation Summary: Context Pack command

## 1. Overview (What & Why)
Adds `tgs context pack "<query>"` to gather relevant context from `tgs/design/*` and the active thought directory, then generate a succinct AI brief (`aibrief.md`) using the shell brain adapter. Honors a configurable token budget and includes source pointers.

## 2. File Changes
- Edited: `src/cmd/context.go`
  - Implemented `context` parent and `context pack` subcommand.
  - Collects design and thought files; loads prompt templates; fills `{{QUERY}}`, `{{TOKEN_BUDGET}}`, and embeds the brief template.
  - Executes `tgs/adapters/claude-code.sh` with `--return-mode text`, passing context via `CONTEXT_FILES` env; writes output to `<thought>/aibrief.md`.
- Added: `tgs/agentops/prompts/context_search.md`, `tgs/agentops/prompts/context_brief.md`
  - Human-editable defaults for search and brief formatting.
- Added: `src/templates/data/tgs/context/search_prompt.md`, `src/templates/data/tgs/context/brief_template.md`
  - Init/decorate templates copied as-is to `tgs/` (no Go templating in these files to avoid parse errors).
- Added: `src/cmd/context_test.go`
  - Tests for missing adapter and happy path using a fake adapter writing `# AI Brief`.
- Updated: `tgs/design/{10_needs.md,20_requirements.md,40_vnv.md}` with needs, requirements, and V&V for context pack.

## 3. Commands & Migrations
- Build and run tests:
  - `go build ./...`
  - `go test ./...`
- Usage:
  - `tgs context pack "authentication and single sign on"`
  - Optional: `--out /path/to/aibrief.md`, `--verbose`

## 4. How to Test
1) Ensure `tgs/adapters/claude-code.sh` exists and `claude` is configured, or point `tgs/tgs.yml` to a test adapter.
2) Create an active thought dir (`TGS_THOUGHT_DIR` or most recent under `tgs/thoughts/`).
3) Run: `tgs context pack "auth sso"`.
4) Verify `<thought>/aibrief.md` exists and contains brief sections.
5) Run unit tests: `go test ./...` (includes adapter-stubbed happy path).

## 5. Integration Steps
- Token budget set via `ai.toolpack.budgets.context_pack_tokens` (default 1200).  
- Model/adapter configured via `ai` fields in `tgs/tgs.yml`.

## 6. Rollback
- Revert changes; command is additive and isolated.

## 7. Follow-ups & Next Steps
- Optional: add anchors/line ranges extraction helpers; richer source mapping.
- Optional: support additional models/routes in config.

## 8. Links
- Thought: `tgs/thoughts/24c8b79-context-pack-command/`
