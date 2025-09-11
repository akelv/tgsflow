# Research: Standardize AGENTOPS intake-to-PR workflow and enrich new-thought

- Date: 2025-09-11
- Base Hash: b4552ea
- Participants: Agent/Human

## 1. Problem Statement
When given a one-line instruction (e.g., "Implement <idea>"), the agent should consistently:
- Ask clarifying questions if needed
- Create a new thought via `make new-thought` and auto-populate `README.md` with base hash, links, and idea spec
- Produce `research.md` and seek approval
- Produce `plan.md` and seek approval
- Implement changes in top-level code areas (never under `tgs/`)
- Summarize in `implementation.md`, update `tgs/README.md` index, and open a PR with `gh`

## 2. Current State
- `agentops/AGENTOPS.md` already documents a gated workflow but lacked explicit intake clarifications and a PR command.
- `Makefile` had a `new-thought` target that scaffolded the thought and a minimal README line (`# <hash> - <title>`).
- `bootstrap.sh` decorates existing repos and generates a `tgs.mk` that included a similar `new-thought` target.
- The thought `README.md` lacked base hash, quick links, and idea spec; PR step was not standardized.

## 3. Constraints & Assumptions
- Maintain backward compatibility: `title` remains required; any new parameter must be optional.
- Keep implementation code outside `tgs/`.
- Assume `gh` CLI may not be authenticated; docs should mention auth requirements.
- Avoid destructive changes; edits should be small and idempotent.

## 4. Risks & Impact
- Minor behavior change to `new-thought`: creates a richer `README.md`; extremely low risk.
- Divergence risk between root `Makefile` and `bootstrap.sh`-generated `tgs.mk`; mitigated by syncing logic.
- Dependence on `gh` CLI for PR creation; failure mode is documented (authenticate or provide token).

## 5. Alternatives Considered
- Wrapper script to orchestrate thought creation and PR: Higher complexity, less transparent than Makefile.
- GitHub Actions to auto-open PRs: Requires additional secrets and CI setup; out of scope.
- Keep manual README population: Slower and inconsistent; not preferred.

## 6. Recommendation
Update documentation and tooling to codify the intake→research→plan→implement→PR flow, and enhance `new-thought` to accept optional `spec` and auto-write base hash and quick links in `README.md`. Sync the same behavior in `bootstrap.sh` so decorated repos get identical semantics.

## 7. References & Links
- `agentops/AGENTOPS.md`
- `Makefile` (`new-thought`)
- `bootstrap.sh` (decorator-mode `tgs.mk` heredoc)
- `agentops/tgs/README.md`

---
Approval checkpoint: Please review this research and reply one of:
- APPROVE research
- REQUEST CHANGES: <notes>
