# Plan: EARS linter core and verify integration

## 1. Objectives
- Implement an EARS linter backed by ANTLR4 Go for the grammar in `src/core/ears/ears.g4`.
- Classify and validate requirement lines as one of: ubiquitous, state-driven, event-driven, complex, unwanted; enforce clause order.
- Provide a Go API and CLI integration: `tgs verify` runs the linter when enabled in `tgs.yaml`.
- Emit clear diagnostics with line numbers and shapes; exit non-zero if issues found (when enabled).

## 2. Scope / Non-goals
In-scope:
- ANTLR4 generation to Go and committed parser sources.
- Core linter API and basic semantic rules toggles (e.g., require "shall").
- Config extension to enable/disable linter and rule toggles.
- Verify command integration scanning Markdown bullets as requirement candidates.
- Unit tests for parser wrapper and integration tests for `verify` flow.

Out-of-scope:
- Rich NLP heuristics beyond simple keyword checks.
- Multi-language requirements.
- Deep Markdown AST parsing; we will use simple bullet/line heuristics for now.

## 3. Acceptance Criteria
- With config `policies.ears.enable: true`, running `tgs verify` scans `.md` files and lints candidate lines, printing diagnostics and returning non-zero if any violations.
- Correct classification of EARS shapes per grammar; invalid/mixed forms are reported with helpful messages.
- With linter disabled or no `.md` files, `tgs verify` behavior is unchanged (still runs hooks, prints completion message).
- Unit tests cover: each EARS shape valid example; malformed inputs; unwanted shape; semantic rule toggle.
- CI/lint passes; no new regressions.

## 4. Phases & Tasks
- Phase 1: Parser generation & core API
  - [ ] Add Makefile target `ears-gen` to generate Go lexer/parser into `src/core/ears/gen`.
  - [ ] Commit generated parser; add `go:generate` line in `ears.g4` or a wrapper `gen.go`.
  - [ ] Implement `src/core/ears/lint.go` with `Result`, `Issue`, `ParseRequirement`, `Lint`.
  - [ ] Add custom error listener for parse errors with positions.
- Phase 2: Config & integration
  - [ ] Extend `src/core/config/loader.go` with `Policies.EARS` struct: `{ Enable bool, RequireShall bool }`.
  - [ ] Update `Default()` with sensible defaults and yaml tags.
  - [ ] Update `src/cmd/verify.go` to load config and run linter when enabled.
  - [ ] Implement simple Markdown scan: lines starting with `- ` or `* ` or numbered lists as candidates.
- Phase 3: Testing & docs
  - [ ] Add unit tests for `ears` package covering shapes and errors.
  - [ ] Add integration test for `CmdVerify` with a temp repo and sample markdown.
  - [ ] Update README or new doc section on enabling EARS linter and developer setup for ANTLR.

## 5. File/Module Changes
Add:
- `src/core/ears/gen/` (generated ANTLR Go files).
- `src/core/ears/lint.go` (core linter wrapper and API).
- `src/core/ears/lint_test.go` (unit tests).
- `src/core/ears/gen/doc.go` (package doc and do-not-edit note).
- `src/core/ears/gen/generate.go` or `src/core/ears/generate.go` (optional `go:generate`).

Edit:
- `src/core/config/loader.go` (add `EARS` policy struct, yaml tags, defaults).
- `src/cmd/verify.go` (config load + linter invocation + output formatting).
- `Makefile` (add `ears-gen` target; ensure non-interactive, documented tools).
- `README.md` or a new `docs/ears.md` (usage, config, generation steps).

## 6. Test Plan
- Parser unit tests: valid examples for each shape (ubiquitous/state/event/complex/unwanted);
  malformed ordering; missing keywords; mixed forms; case-insensitivity.
- Semantic rule tests: when `RequireShall` true, ensure sentences lacking "shall" in response are flagged.
- Verify integration tests: temp dir with `tgs.yaml` enabling EARS; markdown file with bullets; ensure exit code non-zero when violations exist and zero when all valid.
- Negative: when disabled, verify runs normally; no parser invoked.

## 7. Rollout & Rollback
- Rollout: gated by `policies.ears.enable` flag (default false). Document setup.
- Regeneration: developers with ANTLR4 can run `make ears-gen`; CI does not require ANTLR since code is committed.
- Rollback: revert config to disable; remove generated code if needed; changes are isolated under `src/core/ears` and `verify` glue.

## 8. Risks
- Risks: ANTLR tool availability and version drift—mitigated by committing generated code and pinning a version in Makefile.

---
Approval checkpoint: Please review this plan and reply one of:
- APPROVE plan
- REQUEST CHANGES: <notes>
