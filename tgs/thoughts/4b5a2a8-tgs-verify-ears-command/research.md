# Research: Verify EARS for Core Design Docs

- Date: 2025-09-21
- Base Hash: <auto-populated by thought scaffolding>
- Participants: AI Agent, Human Maintainer

## 1. Problem Statement
We need to implement a targeted enhancement to the `tgs verify` command to lint EARS-style requirement lines specifically within `tgs/design/10_needs.md` and `tgs/design/20_requirements.md`. The linter must ignore code fences, support bullet response sections after EARS headings that end with a colon, and report issues in `path:line: message` format. This fulfills the AI brief request and the new need `N-022` and requirement `SR-026`.

## 2. Current State
- `src/cmd/verify.go` contains `verifyEARS` which currently scans all `*.md` files repo-wide and lints both top-level bullets and lines starting with EARS keywords. It supports skipping code fences and a simple bullet-response mode.
- EARS grammar and parser exist under `src/core/ears/` with `ParseRequirement` and `Lint` functions. Forms supported: ubiquitous, state-driven, event-driven, complex, unwanted.
- Tests exist (`src/cmd/verify_ears_test.go`) validating success/failure on generic fixtures but not focusing on the two design files.
- Config `Guardrails.EARS.Enable` toggles EARS checks; there is no scope option yet.
- V&V (`tgs/design/40_vnv.md`) expects exit code behavior and now includes `SR-026` criteria for the two design files.

## 3. Constraints & Assumptions
- Non-destructive, fast execution; no network needed.
- macOS/Linux supported; Go toolchain only.
- Do not place production code under `tgs/` per policy.
- Keep current generic repo-wide behavior for backward compatibility, but add focused checks for the two design files to satisfy `SR-026`.
- Output format must remain `path:line: message`; CI mode (`--ci`) determines non-zero exit behavior.

## 4. Risks & Impact
- Risk: Over-linting could flag non-requirement narrative lines in design docs. Mitigation: Only lint lines matching EARS starts or bullets, skip headings/blank/fenced blocks, reuse existing logic.
- Risk: Backward-incompatible behavior if we restrict scope. Mitigation: Keep existing repo-wide scan and additionally ensure the two design files are included (i.e., additive focus, not exclusive scope) or gate via a future config scope flag (out of this change).
- Impact: Clearer signal on design documents; improved compliance to EARS.

## 5. Alternatives Considered
- A) Restrict linter scope to design docs only. Pros: Fewer false positives; Cons: Regresses current broad checks. Not chosen.
- B) Add config `guardrails.ears.paths` to define scope. Pros: Flexible; Cons: Extra config work now. Defer.
- C) Keep current behavior and add explicit targeted tests for `10_needs.md` and `20_requirements.md`. Pros: Minimal change; Cons: Still scans all files. Chosen for now.

## 6. Recommendation
Proceed with Alternative C: Leave `verifyEARS` implementation as-is for scanning all `.md` files, but add tests that set up temp repos containing `tgs/design/10_needs.md` and `tgs/design/20_requirements.md` content and assert success/failure per `SR-026`. Optionally, tune `verifyEARS` to ensure bullet-response mode works for colon-terminated EARS lines (already supported) and ensure error messages are formatted `path:line: message` (already implemented).

## 7. References & Links
- `tgs/thoughts/4b5a2a8-tgs-verify-ears-command/aibrief.md`
- `src/cmd/verify.go` (`verifyEARS`)
- `src/core/ears/lint.go`
- `src/cmd/verify_ears_test.go`
- `tgs/design/10_needs.md`, `tgs/design/20_requirements.md`, `tgs/design/40_vnv.md`

---
Approval checkpoint: Please review this research and reply one of:
- APPROVE research
- REQUEST CHANGES: <notes>
