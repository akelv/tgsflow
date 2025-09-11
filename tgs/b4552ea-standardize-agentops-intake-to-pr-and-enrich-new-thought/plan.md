# Plan: Standardize AGENTOPS intake-to-PR and enrich new-thought

## 1. Objectives
- Codify intake→research→plan→implement→PR as the default agent procedure.
- Enrich `new-thought` to auto-populate `README.md` with base hash, quick links, and optional idea spec.
- Ensure decoration via `bootstrap.sh` generates identical behavior in `tgs.mk`.
- Document PR creation via `gh` with a clear command in `AGENTOPS.md`.

## 2. Scope / Non-goals
- In-scope: Docs updates, Makefile target changes, bootstrap decorator updates, templates README update.
- Non-goals: Automating `gh auth` or CI workflows; adding app templates beyond existing.

## 3. Acceptance Criteria
- `agentops/AGENTOPS.md` contains the step-by-step procedure including PR creation.
- `make new-thought title="X" spec="Y"` creates a rich `README.md` including base hash and links.
- `bootstrap.sh` decorator writes `tgs.mk` with the same `new-thought` behavior.
- `agentops/tgs/README.md` explains the auto-population behavior.
- All changes committed, branch pushed, PR can be created successfully once `gh` is authenticated.

## 4. Phases & Tasks
- Phase 1: Update documentation
  - [x] Update `agentops/AGENTOPS.md` with the standardized procedure and `gh` command.
  - [x] Update `agentops/tgs/README.md` to reflect enriched README.
- Phase 2: Tooling
  - [x] Enhance root `Makefile` `new-thought` to accept `spec` and write rich README.
  - [x] Sync `bootstrap.sh` to generate equivalent `tgs.mk` logic.
- Phase 3: Thought docs
  - [x] Create new thought via `make new-thought` with the idea spec.
  - [x] Author `research.md`.
  - [x] Author this `plan.md`.
  - [ ] Author `implementation.md` with summary and improvement suggestions.
- Phase 4: Integration
  - [ ] Update `tgs/README.md` index with new entry.
  - [ ] Commit and push thought documents.
  - [ ] Open PR via `gh` and request review.

## 5. File/Module Changes
- Edit: `agentops/AGENTOPS.md` — workflow, PR step, clarifications.
- Edit: `Makefile` — `new-thought` spec parameter and README population.
- Edit: `bootstrap.sh` — `tgs.mk` heredoc to match `new-thought` behavior.
- Edit: `agentops/tgs/README.md` — describe auto-populated README.
- Add: `tgs/b4552ea-standardize-agentops-intake-to-pr-and-enrich-new-thought/*` — research, plan, implementation, README from template.

## 6. Test Plan
- Manual:
  - Run `make new-thought title="Test Thought" spec="Spec"`; verify `README.md` contains base hash, links, and spec.
  - Run `./bootstrap.sh --decorate --dry-run`; verify `tgs.mk` generation step shows updated usage string.
  - After `gh auth login`, run `gh pr create` and verify PR opens with correct title/body.

## 7. Rollout & Rollback
- Rollout: Merge PR; teams use the updated docs and tooling immediately.
- Rollback: Revert the PR commit; no persistent state changes.

## 8. Estimates & Risks
- Estimate: Small change set; <1 day.
- Risks: Divergence between `Makefile` and `tgs.mk` in future; mitigation: keep changes in one place and copy carefully.

---
Approval checkpoint: Please review this plan and reply one of:
- APPROVE plan
- REQUEST CHANGES: <notes>
