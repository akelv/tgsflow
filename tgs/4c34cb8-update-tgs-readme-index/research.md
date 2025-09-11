# Research: Update TGS README index

- Date: 2025-09-11
- Base Hash: 4c34cb8
- Participants: Agent, Human

## 1. Problem Statement
The `tgs/README.md` lacks a "Current Thoughts" index summarizing existing thought directories under `tgs/`. We need to generate and maintain a table listing each thought directory with base hash, date, status, and a concise description, matching the provided example format.

## 2. Current State
- Thought directories detected:
  - `tgs/612a57f-decorate-existing-software-project-repository/`
  - `tgs/b4552ea-standardize-agentops-intake-to-pr-and-enrich-new-thought/`
  - `tgs/f857d9e-test-tgs-workflow-creation/`
- Each has scaffolding; `implementation.md` presence suggests completion, but we should infer status more robustly.
- `tgs/README.md` currently documents workflow and usage but has no index.

## 3. Constraints & Assumptions
- Follow AGENTOPS workflow: author research and plan, get approvals before modifying production docs.
- Use absolute paths in commands and non-interactive flags.
- Do not change content under `tgs/` beyond adding the index section in `tgs/README.md`.
- Dates should be sourced from git (first add date for each dir) when available; if unavailable, default to today's date.
- Status inference:
  - ✅ Completed: `implementation.md` exists and is non-empty.
  - 🚧 In Progress: `plan.md` exists but no `implementation.md` or empty implementation.
  - 🧭 Research: only `research.md` exists.

## 4. Risks & Impact
- Misstating status if file presence is used as proxy. Mitigation: check file content length > 0.
- Git history might not have creation date for new, uncommitted directories. Mitigation: fallback to current date.
- Formatting drift with future example changes. Mitigation: replicate table format exactly and keep logic simple.

## 5. Alternatives Considered
- Manual static list: simple but becomes stale quickly.
- Scripted generator target (e.g., `make update-thought-index`): robust but out of scope for now.
- GitHub Actions to auto-update: overkill for current need.

## 6. Recommendation
Manually compute current entries using repo state and git history, then inject a "## Current Thoughts" table near the bottom of `tgs/README.md`, before the final paragraph, preserving the documented sections.

## 7. References & Links
- `tgs/README.md`
- `tgs/*/README.md`, `implementation.md`
- Provided example table in the user request

---
Approval checkpoint: Please review this research and reply one of:
- APPROVE research
- REQUEST CHANGES: <notes>
