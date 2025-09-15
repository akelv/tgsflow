# Research: EARS linter core and verify integration

- Date: 2025-09-15
- Base Hash: f71f872
- Participants: Agent/Human

## 1. Problem Statement
Implement a core EARS linter in Go, generated from the ANTLR4 grammar in `src/core/ears/ears.g4`, and integrate it with the `tgs verify` command behind a configuration flag. The linter must:
- Parse candidate requirements (one line each) and classify them into allowed EARS shapes: ubiquitous, state-driven, event-driven, complex, and unwanted.
- Enforce valid clause order per the grammar; reject invalid or mixed orders.
- Keep System, Preconditions, Trigger, and Response content free-form, with optional semantic checks in Go (e.g., “must contain ‘shall’”, weak verbs, ambiguity heuristics).
- Be invocable by `tgs verify` only when enabled in `tgs.yaml` config.

## 2. Current State
- Grammar `ears.g4` exists under `src/core/ears/ears.g4` with the five EARS forms, tokens, and comments indicating line-oriented parsing.
- No Go parser artifacts exist yet; ANTLR4 generation and Go wrapper are missing.
- `tgs verify` currently runs local hooks from `.tgs/hooks/*` only and has no awareness of repo `tgs.yaml` or EARS linting.
- Config loader `src/core/config/loader.go` supports `tgs.yaml`, but has no fields for ears linting.

## 3. Constraints & Assumptions
- Use ANTLR4 Go target to generate lexer/parser from `ears.g4` during dev and commit generated code for reproducible builds (no Java toolchain requirement at runtime).
- Input unit is one requirement line (driver can pre-split or a simple wrapper can be used; we’ll pre-split lines in Go for simplicity).
- Clause content is free-form; we may run additional semantic checks post-parse. Cardinality (single trigger, ≥1 response) is implied by grammar.
- Config-driven: add `policies.ears.verify` boolean or similar to `tgs.yaml` to gate execution during `tgs verify`.
- Keep public API small and testable: provide `ParseRequirement(line string) (Result, error)` and `Lint(lines []string) []Issue`.

## 4. Risks & Impact
- Tooling friction: Requires ANTLR4 to regenerate parser when grammar changes. Mitigation: provide `make ears-gen` and document setup; commit generated code.
- False negatives/positives from strict grammar on real-world requirements. Mitigation: clear error messages, allow bypass via config, iterative refinement.
- Performance: negligible for line-by-line parsing; batch runs are fast.
- Backward compatibility: gated behind config; no effect unless enabled.

## 5. Alternatives Considered
- Hand-rolled parser: quicker start, but brittle and harder to maintain than formal grammar.
- Simpler regex checks: easier but error-prone; misses structure and ordering accuracy.
- External linter binary: adds dependency and integration complexity; not needed given small scope.

## 6. Recommendation
- Generate Go lexer/parser with ANTLR4 Go target from `ears.g4` into `src/core/ears/gen`.
- Implement `src/core/ears/lint.go` exposing:
  - `Type Shape = string` with constants: `Ubiquitous`, `StateDriven`, `EventDriven`, `Complex`, `Unwanted`.
  - `type Result` capturing shape and text slices for system, preconditions, trigger, response.
  - `ParseRequirement(line string) (Result, error)` using the generated parser and custom error listener.
  - Optional semantic rules (`RequireShall`, weak-verb checks) configurable via options.
- Extend config: add `Policies.EARS` with `{ Enable bool, RequireShall bool }` (extensible).
- Integrate in `cmd/verify.go`: load config, when enabled, locate and read requirement sources (initially scan `*.md` top-level bullets or accept a path flag later), split into lines, run linter, print issues, and exit non-zero if any.
- Provide Makefile targets and docs: `make ears-gen`, `make verify`.

## 7. References & Links
- Grammar: `src/core/ears/ears.g4`
- Verify command: `src/cmd/verify.go`
- Config loader: `src/core/config/loader.go`
- ANTLR4 Go target docs: `https://github.com/antlr/antlr4/blob/master/doc/go-target.md`

---
Approval checkpoint: Please review this research and reply one of:
- APPROVE research
- REQUEST CHANGES: <notes>
