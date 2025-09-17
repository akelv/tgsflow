# Plan: Fix EARS regression tests

## 1. Objectives
- Restore all tests in `src/core/ears` to green.
- Preserve current positive/negative semantics of EARS shapes.
- Keep grammar readable and maintainable; case-insensitive keywords.

## 2. Scope / Non-goals
- In-scope: Grammar and generated parser for `src/core/ears`; minimal `lint.go` adjustments if needed.
- Out-of-scope: Changes to CLI or other packages; expanding EARS shapes beyond current tests.

## 3. Acceptance Criteria
- `go test ./src/core/ears -v` passes: all positives parse; negatives fail as asserted.
- `make test` (repo-wide) has no new failures due to these changes.
- No linter errors introduced.

## 4. Phases & Tasks
- Phase 1: Grammar fix
  - [ ] Replace `TEXT_NOCOMMA` with `WORD : [^,\r\n\s]+ ;` (after keywords)
  - [ ] Replace `TEXT_EOL` with `REST : ~[\r\n]+ ;` and set `response: REST?`
  - [ ] Update `token_word: WORD | THE | WHEN | IF | THEN`
- Phase 2: Regenerate and compile
  - [ ] Run `make ears-gen`
  - [ ] Build to ensure no compile errors
- Phase 3: Verify and iterate
  - [ ] Run `go test ./src/core/ears -v`
  - [ ] If failures remain, adjust grammar minimally (e.g., allow multiple `token_word` in system/clause preserved as is)
- Phase 4: Documentation
  - [ ] Document changes and how to regenerate in `implementation.md`

## 5. File/Module Changes
- Edit: `src/core/ears/ears.g4` (lexer tokens and response rule)
- Regenerated: `src/core/ears/gen/src/core/ears/*`
- Potential small edit: `src/core/ears/lint.go` (no semantic change expected)

## 6. Test Plan
- Run focused tests: `go test ./src/core/ears -run TestParse_ -v`
- Run regression matrix: `go test ./src/core/ears -run Test_EARS_Fixtures_Matrix -v`
- Full package tests: `go test ./src/core/ears -v`
- Full repo tests: `make test`

## 7. Rollout & Rollback
- Rollout: Commit grammar change and generated files in one commit.
- Rollback: Revert commit to restore prior grammar and generated code if issues surface.

## 8. Estimates & Risks
- Estimate: ~30-60 minutes including regeneration and test fixes.
- Risks: Lexer priority subtleties; mitigate by keeping keywords first and narrowing fallback tokens.

---
Approval checkpoint: Please review this plan and reply one of:
- APPROVE plan
- REQUEST CHANGES: <notes>
