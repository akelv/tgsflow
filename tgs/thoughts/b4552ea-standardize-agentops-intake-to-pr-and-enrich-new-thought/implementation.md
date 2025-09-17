# Implementation Summary: Standardize AGENTOPS workflow and enrich new-thought

## 1. Overview (What & Why)
We standardized the end-to-end agent procedure from intake to PR and enriched the `new-thought` scaffolding. This ensures consistent approval gates and removes boilerplate by auto-populating thought README files with base hash, quick links, and optional idea spec.

## 2. File Changes
- Edited `agentops/AGENTOPS.md`
  - Clarified intake and added explicit PR step using `gh`.
  - Emphasized implementing outside `tgs/` and auto-populated README via `new-thought`.
- Edited `Makefile`
  - `new-thought` now supports optional `spec` and generates a richer `README.md` with base hash and quick links.
- Edited `bootstrap.sh`
  - Synced the generated `tgs.mk` heredoc with the updated `new-thought` behavior.
- Edited `agentops/tgs/README.md`
  - Documented auto-populated README and expected workflow.
- Added `tgs/b4552ea-standardize-agentops-intake-to-pr-and-enrich-new-thought/*`
  - `research.md`, `plan.md`, `implementation.md`, `README.md` for this thought.

## 3. Commands & Migrations
- None. Tooling changes only. New behavior is exercised via Make targets and `gh`.

## 4. How to Test
- Thought scaffolding:
  - Run: `make new-thought title="Test Thought" spec="My idea"`
  - Expect: A directory `tgs/<hash>-test-thought/` with `README.md` showing base hash, quick links, and idea spec.
- Decorator flow (dry run):
  - Run: `./bootstrap.sh --decorate --dry-run`
  - Expect: Writes/ensures `agentops/AGENTOPS.md`, `tgs/README.md`, templates, and generates `tgs.mk` with the updated usage string and README logic.
- PR creation:
  - Authenticate: `gh auth login` (or set `GH_TOKEN`).
  - Run: `gh pr create --title "feat(agentops): standardize workflow and enrich new-thought" --body-file tgs/b4552ea-standardize-agentops-intake-to-pr-and-enrich-new-thought/implementation.md`
  - Expect: PR opens with the summary as the body.

## 5. Integration Steps
- Merge this PR. Teams should use the documented procedure in `agentops/AGENTOPS.md`.
- New thoughts should be created with `spec` where appropriate to capture the one-liner idea.

## 6. Rollback
- Revert the PR commit. No data migrations or stateful rollbacks required.

## 7. Follow-ups & Next Steps
- Improvement suggestion: Deduplicate `new-thought` logic by sourcing it from a single script (e.g., `scripts/new_thought.sh`) to avoid drift between `Makefile` and `tgs.mk` in `bootstrap.sh`.
- Improvement suggestion: Add a `make open-pr` convenience target that validates `gh` auth and opens a PR with templated title/body.
- Improvement suggestion: Extend `bootstrap.sh --decorate` to optionally add a minimal `.github/PULL_REQUEST_TEMPLATE.md` referencing `tgs/<dir>/implementation.md`.

## 8. Links
- Branch: `feat/agentops-standard-procedure`
- After auth, PR URL (on creation): `https://github.com/akelv/tgsflow/pull/new/feat/agentops-standard-procedure`
