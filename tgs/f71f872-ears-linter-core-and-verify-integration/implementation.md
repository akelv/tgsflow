# Implementation Summary: EARS linter core and verify integration

## 1. Overview (What & Why)
- Added an EARS linter in Go powered by ANTLR4 using grammar `src/core/ears/ears.g4`.
- Enforces allowed EARS shapes and clause order; content remains free-form with optional semantic checks later.
- Integrated into `tgs verify` gated by `tgs.yaml` (`policies.ears.enable`).

## 2. File Changes
- Added: `src/core/ears/ears.g4` (updated), generated parser under `src/core/ears/gen/`.
- Added: `src/core/ears/lint.go` (parser wrapper, API) and tests `src/core/ears/lint_test.go`.
- Added: `src/cmd/verify_ears_test.go` (integration tests for verify + EARS).
- Edited: `src/cmd/verify.go` (config load, EARS scan of Markdown bullets, CI exit behavior).
- Edited: `src/core/config/loader.go` (added `Policies.EARS` with `Enable`, `RequireShall`).
- Edited: `Makefile` (added `ears-gen` target).
- Edited: `README.md` (developer notes for ANTLR/Java and usage).

## 3. Commands & Migrations
- Dependencies:
  - Go: `github.com/antlr4-go/antlr/v4` added to `go.mod`.
  - Java + ANTLR (dev only for codegen):
    ```bash
    brew install openjdk antlr
    export CLASSPATH="$(brew --prefix)/libexec/antlr-4.13.2-complete.jar:$CLASSPATH"
    make ears-gen
    ```
- Build & test:
  ```bash
  make build
  make test
  go test -cover ./...
  ```

## 4. How to Test
- Unit tests: `go test ./src/core/ears -v` (coverage >90%).
- Integration tests: `go test ./src/cmd -run Verify_EARS -v`.
- Manual verify run:
  1) Enable in `tgs.yaml`:
  ```yaml
  policies:
    ears:
      enable: true
      require_shall: false
  ```
  2) Ensure Markdown with bullets like:
  ```
  - When user clicks, the service shall log
  - The system shall record events
  ```
  3) Run: `./bin/tgs verify --repo . --ci`
  - Non-zero exit if any invalid lines are found; otherwise 0.

## 5. Integration Steps
- Config gate: `policies.ears.enable` (default false) so existing repos unaffected.
- CI: call `tgs verify --ci` as part of your pipeline to enforce EARS when enabled.

## 6. Rollback
- Disable by setting `policies.ears.enable: false` in `tgs.yaml`.
- Generated code is isolated under `src/core/ears/gen`; remove if decommissioning EARS.

## 7. Follow-ups & Next Steps
- Optional semantic rules (RequireShall, weak verbs, ambiguity) in `ears.Lint`.
- CLI options to target specific files or directories.
- Richer Markdown parsing or file globs.

## 8. Links
- Thought: `tgs/f71f872-ears-linter-core-and-verify-integration/`
- Key sources: `src/core/ears/lint.go`, `src/cmd/verify.go`, `src/core/config/loader.go`
