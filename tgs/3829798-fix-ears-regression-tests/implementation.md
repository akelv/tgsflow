# Implementation: Fix EARS regression tests

## What & Why
- Fixed EARS lexer longest-match issue that caused valid sentences to be tokenized as generic text and rejected at position 1:0.
- Simplified tokenization by replacing `TEXT_NOCOMMA/TEXT_EOL` with `WORD`, and broadened `response` to accept tokens to line end. 
- Adjusted semantic validation to allow pronoun `it` as the system and to handle optional `then` before system segment in unwanted form.

## Changes
- Edited `src/core/ears/ears.g4`:
  - Replaced `TEXT_NOCOMMA`/`TEXT_EOL` with `WORD`.
  - `token_word` now: `THE | WHEN | IF | THEN | WORD`.
  - `response` now consumes `(token_word | COMMA)*`.
- Edited `src/core/ears/lint.go`:
  - `validateSystemSegment` now strips optional leading `then` after the comma.
  - Allows `it` (pronoun) in the system segment for event/complex/unwanted/state.
- Regenerated ANTLR artifacts under `src/core/ears/gen/src/core/ears/`.

## Commands
```bash
# Ensure Java and ANTLR are available (example for Homebrew)
export PATH="/usr/local/opt/openjdk/bin:/opt/homebrew/opt/openjdk/bin:$PATH"
export CLASSPATH="$(brew --prefix)/libexec/antlr-4.13.2-complete.jar:$CLASSPATH"

# Regenerate and run tests
make ears-gen
go test ./src/core/ears -v
```

## How to Test
- Unit tests:
  - `go test ./src/core/ears -run TestParse_ -v`
  - `go test ./src/core/ears -run Test_EARS_Fixtures_Matrix -v`
- Full package:
  - `go test ./src/core/ears -v`

## Integration Notes
- No external API changes. Grammar and parser are internal to EARS linter.
- Keep Java/ANTLR toolchain available to regenerate when grammar changes.

## Rollback
- Revert edits to `src/core/ears/ears.g4` and `src/core/ears/lint.go` plus the generated `src/core/ears/gen/...` files.

## Follow-ups
- Consider more nuanced clause splitting for reporting.
- Add CI job to ensure `make ears-gen` output is up-to-date.

## Links
- Thought: `tgs/3829798-fix-ears-regression-tests/`
- Research: `tgs/3829798-fix-ears-regression-tests/research.md`
- Plan: `tgs/3829798-fix-ears-regression-tests/plan.md`
