# Plan: Verify EARS for Core Design Docs

## 1. Objectives
- Ensure `tgs verify` lints EARS-style lines in `tgs/design/10_needs.md` and `tgs/design/20_requirements.md` with correct skipping of code fences and bullet response sections.
- Preserve existing repo-wide scanning behavior for backward compatibility.
- Provide deterministic error messages in `path:line: message` format and correct exit codes in CI mode.

## 2. Scope / Non-goals
- In-scope: Tests and docs ensuring focused validation for design docs; minor code adjustments if needed.
- Non-goals: New config options to scope EARS paths; refactors of ANTLR grammar; changing default scan scope.

## 3. Acceptance Criteria
- With valid EARS lines in both `tgs/design/10_needs.md` and `tgs/design/20_requirements.md`, `tgs verify --repo <tempdir> --ci` exits 0.
- Introducing a known invalid line in either file causes `tgs verify --repo <tempdir> --ci` to exit 1 and prints `path:line: message` describing the parse error.
- Code fences and non-EARS narrative lines are ignored; bullet response mode is honored when the EARS line ends with a colon and contains " shall".

## 4. Phases & Tasks
- Phase 1: Tests
  - [ ] Add tests creating temp repos with `tgs/design/10_needs.md` and `tgs/design/20_requirements.md` content (valid and invalid cases).
  - [ ] Assert exit codes and stderr messages per `SR-026`.
- Phase 2: Implementation touch-ups (only if tests reveal issues)
  - [ ] Adjust `verifyEARS` parsing of bullet response mode or EARS line detection as needed.
- Phase 3: Docs
  - [ ] Ensure `tgs/design/10_needs.md`, `20_requirements.md`, and `40_vnv.md` reflect `N-022`, `SR-026`, and V&V.

## 5. File/Module Changes
- `src/cmd/verify_ears_test.go`: Add new tests covering design-doc focus (valid/invalid scenarios using temp repos).
- `src/cmd/verify.go`: No functional changes expected; minor adjustments if tests indicate gaps.
- `tgs/design/*`: Already updated with need/requirement and V&V.

## 6. Test Plan
- Unit tests in `src/cmd/verify_ears_test.go` construct temporary repo roots with design files and run `CmdVerify` with `--repo` and `--ci`.
- Validate both success and failure paths, including error message shape.

## 7. Rollout & Rollback
- Rollout: Merge tests and any minimal code changes; run CI.
- Rollback: Revert commit if regressions are found; feature is additive and test-only unless code changes required.

## 8. Estimates & Risks
- Estimate: 1-2 hours including tests.
- Risks: False positives on narrative lines; mitigated by existing filters and focused tests.

---
Approval checkpoint: Please review this plan and reply one of:
- APPROVE plan
- REQUEST CHANGES: <notes>
