# Plan: Extend init to include adapters, subcommand, and Makefile target

## 1. Objectives
- Ensure `tgs init` also installs `tgs/adapters/claude-code.sh` and `tgs/adapters/gemini-code.sh` if missing and sets executable bits.
- Add `tgs init claude|gemini` subcommands that safely copy `tgs/agentops/AGENTOPS.md` to root `CLAUDE.md`/`GEMINI.md` only if absent; otherwise, exit with guidance.
- Ensure a standard `new-thought` Makefile target exists in repos after init (append if missing).

## 2. Scope / Non-goals
In-scope:
- Modify `src/cmd/init.go` to implement adapter ensure, subcommand behavior, and Makefile mutation.
- Reuse existing embedded templates; do not change adapter script contents.
Out of scope:
- Changing adapter behavior, CLI interfaces, or adding new vendors beyond `claude|gemini`.
- Interactive prompts or auto-overwrite of existing root docs.

## 3. Acceptance Criteria
- `tgs init` results in `tgs/adapters/claude-code.sh` and `tgs/adapters/gemini-code.sh` present and executable if they were missing.
- `tgs init claude` creates `CLAUDE.md` at repo root using `tgs/agentops/AGENTOPS.md` as source when absent; if present, command exits non-zero with a clear message.
- `tgs init gemini` parallels behavior creating `GEMINI.md`.
- If root `Makefile` lacks a `new-thought:` rule, `tgs init` appends the standard rule; if present, leave unchanged.

## 4. Phases & Tasks
- Phase 1: Implement features in `init`
  - [ ] Add args parsing for `tgs init [adapter]` via Cobra (capture optional arg).
  - [ ] Implement `ensureAdapters()` to copy scripts from repo `tgs/adapters/*.sh` if missing and `chmod +x`.
  - [ ] Implement `ensureMakefileNewThought()` to append standard rule if missing.
  - [ ] Implement `decorateAdapter(adapter)` to copy `tgs/agentops/AGENTOPS.md` to root `CLAUDE.md`/`GEMINI.md` with existence guard.
  - [ ] Wire into `CmdInit` execution flow and logs.
- Phase 2: Tests & docs
  - [ ] Unit tests for Makefile injection and existence checks (temp dirs).
  - [ ] Update docs/messages if needed.

## 5. File/Module Changes
- Edit: `src/cmd/init.go` – new logic for adapter ensure, subcommand, and Makefile rule.
- No new templates required; use existing `tgs/agentops/AGENTOPS.md` and adapter scripts under repo `tgs/adapters/` as sources.

## 6. Test Plan
- Manual:
  - In a temp repo, run `tgs init`; verify `tgs/adapters/*` exist and are executable; check `Makefile` contains `new-thought`.
  - Run `tgs init claude` → creates `CLAUDE.md`; run again → error message, non-zero exit.
  - Run `tgs init gemini` → creates `GEMINI.md` similar behavior.
- Automated (Go tests with temp dirs):
  - Test adapter ensure creates files when absent and preserves when present.
  - Test Makefile append when target missing and preserve when present.
  - Test subcommand error on existing `CLAUDE.md`/`GEMINI.md`.

## 7. Rollout & Rollback
- Rollout: ship as minor feature; no migrations. Idempotent; safe to re-run.
- Rollback: revert file edits; existing files remain.

## 8. Estimates & Risks
- Estimate: 4-6 hours including tests.
- Risks: Makefile detection false negatives; mitigate by simple `grep -q '^new-thought:'` and documenting override if needed.

---
Please review this plan and reply one of:
- APPROVE plan
- REQUEST CHANGES: <notes>
