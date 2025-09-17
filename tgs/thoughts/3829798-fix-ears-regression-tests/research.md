# Research: Fix EARS regression tests

- Date: 2025-09-16
- Base Hash: 3829798
- Participants: Agent

## 1. Problem Statement
The EARS regression test suite under `src/core/ears` fails. Errors show the parser reporting mismatched input at position 1:0 for valid EARS sentences, e.g., “The system shall record events,” indicating tokenization/parsing doesn’t accept initial tokens like THE/WHEN/WHILE.

## 2. Current State
- Grammar: `src/core/ears/ears.g4` defines parser and lexer rules for five EARS shapes and tokens (WHILE, WHEN, IF, THEN, THE, SHALL, PRONOUN, COMMA, TEXT_NOCOMMA, TEXT_EOL, WS, NEWLINE).
- Generated code exists at `src/core/ears/gen/src/core/ears/*`. Parser file shows tokens and rules matching the grammar; ubiquitous rule requires `(THE system | PRONOUN) SHALL response`.
- Linter/parser code: `src/core/ears/lint.go` calls `earsp.NewearsLexer`/`NewearsParser`, then `parser.Requirement()` and maps contexts to results.
- Tests: `lint_test.go` and `regression_test.go` contain positive/negative cases. Many positives fail with “mismatched input ... expecting {WHILE, WHEN, IF, THE, PRONOUN}”.

Observation: Failures happen at 1:0, implying the first token produced by the lexer is not one of the expected keywords, likely `TEXT_NOCOMMA` capturing the first word.

## 3. Constraints & Assumptions
- Keep grammar flexible: allow keywords within clauses and system names.
- Preserve case-insensitivity.
- No trailing newline required in input.
- Do not break existing negative tests.

## 4. Risks & Impact
- Adjusting lexer could make clause parsing too permissive/restrictive.
- Changing tokens requires regeneration and may affect other code.

## 5. Alternatives Considered
- A) Keep `TEXT_NOCOMMA` and use lexer predicates/modes to force keyword tokens at word starts: complex and brittle.
- B) Detect keywords at parser level using string literals: loses clean case-insensitivity and readability.
- C) Replace `TEXT_NOCOMMA` with `WORD` and ensure keywords win ties; replace `TEXT_EOL` with `REST` for response remainder. This keeps simple tokenization and allows keywords to match when equal-length (rule order), avoiding longest-match overshadowing.

## 6. Recommendation
Modify `ears.g4`:
- Replace `TEXT_NOCOMMA` with `WORD : [^,\r\n\s]+ ;` (after keyword tokens).
- Replace `TEXT_EOL` with `REST : ~[\r\n]+ ;` and set `response: REST?`.
- Update `token_word: WORD | THE | WHEN | IF | THEN`.
Then regenerate with `make ears-gen` and rerun the tests. Adjust `lint.go` only if error wording changes are needed.

## 7. References & Links
- Files: `src/core/ears/ears.g4`, `src/core/ears/lint.go`, `src/core/ears/gen/...`
- Tests: `src/core/ears/lint_test.go`, `src/core/ears/regression_test.go`

---
Approval checkpoint: Please review this research and reply one of:
- APPROVE research
- REQUEST CHANGES: <notes>
