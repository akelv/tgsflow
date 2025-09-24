# Research: Extend init to include adapters, subcommand, and Makefile target

- Date: 2025-09-24
- Base Hash: af66921

## 1) Problem
Current `tgs init` scaffolds `tgs/` and optional CI but does not:
- Ensure `tgs/adapters/` shell scripts are present.
- Support vendor decoration via `tgs init claude|gemini` that writes a root `CLAUDE.md`/`GEMINI.md` from `tgs/agentops/AGENTOPS.md` safely.
- Ensure repositories have a `make new-thought` target to follow the workflow immediately.

## 2) Current State
- `init` renders from embedded or external templates into `tgs/` and optionally adds CI files.
- Adapter scripts exist in this repo under `tgs/adapters/*.sh`, and requirements (SR-020/027) cover their use but not their initialization into decorated repos.
- Root `Makefile` in this repo already includes a `new-thought` target; decorated repos may not.

## 3) Constraints
- Idempotent and non-destructive: never overwrite existing files unless user opts in; for this change, error out if `CLAUDE.md` (or `GEMINI.md`) exists.
- Keep implementation in `src/`, not under `tgs/`.
- Work offline by default (use embedded templates), but allow external templates as today.

## 4) Risks
- Overwriting user files: mitigated by existence checks and clear errors.
- Makefile mutation conflicts: mitigate by checking for an existing `new-thought:` target before appending; use simple text detection.
- Platform perms: ensure adapter scripts are executable.

## 5) Alternatives
- Separate `tgs adapters install` command: extra step; less streamlined.
- Copy `AGENTOPS.md` always with backup: more invasive; contradicts non-destructive policy.

## 6) Recommendation
Extend `init` to:
1. Always ensure `tgs/adapters/claude-code.sh` and `tgs/adapters/gemini-code.sh` are present and executable if missing.
2. Support `tgs init claude|gemini` to copy `tgs/agentops/AGENTOPS.md` to `CLAUDE.md` or `GEMINI.md` if absent; otherwise, exit non-zero with instructions.
3. Ensure the root `Makefile` contains a `new-thought` target; append the standard implementation if missing.

## 7) References
- `src/cmd/init.go`, `src/templates/render.go`
- Templates: `src/templates/data/tgs/agentops/AGENTOPS.md.tmpl`
- Requirements: SR-028..SR-030; Needs: N-024..N-026; V&V matrix

---
Please review `research.md` in `tgs/thoughts/af66921-extend-init-adapters-subcmd-make-target`. Reply: APPROVE research | REQUEST CHANGES: <notes>.
# Research: <Short Title>

- Date: <YYYY-MM-DD>
- Base Hash: <git rev-parse --short HEAD>
- Participants: <Agent/Human>

## 1. Problem Statement
<Clear description of the task and desired outcomes.>

## 2. Current State
<What exists today? Code, tools, versions, constraints.>

## 3. Constraints & Assumptions
<Security, performance, platform, dependencies, compliance, SLAs.>

## 4. Risks & Impact
<Security/privacy, reliability, regressions, scope creep, rollout risk.>

## 5. Alternatives Considered
<Option A, B, C with pros/cons.>

## 6. Recommendation
<Preferred approach and rationale.>

## 7. References & Links
<Docs, tickets, PRs, relevant code paths.>

---
Approval checkpoint: Please review this research and reply one of:
- APPROVE research
- REQUEST CHANGES: <notes>
