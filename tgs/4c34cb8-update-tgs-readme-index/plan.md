# Plan: Update TGS README index

## 1. Objectives
- Add a "Current Thoughts" table to `tgs/README.md` summarizing all thought directories.
- Ensure each row contains: Thought Directory (link), Base Hash, Date, Status, Description.

## 2. Scope / Non-goals
- In scope: Editing `tgs/README.md` only.
- Out of scope: Automating index generation, modifying thought contents.

## 3. Acceptance Criteria
- `tgs/README.md` includes a "## Current Thoughts" section with a markdown table.
- All present thought directories are listed (including `f857d9e-test-tgs-workflow-creation/`).
- Dates are populated from git creation date when available; otherwise today's date.
- Status reflects file presence/content:
  - ✅ Completed if `implementation.md` exists and is non-empty.
  - 🚧 In Progress if `plan.md` exists but implementation missing or empty.
  - 🧭 Research if only `research.md` exists.
- Descriptions are concise, sourced from each thought's `README.md` or `implementation.md` title/overview.

## 4. Phases & Tasks
- Phase 1: Gather metadata
  - [x] List thought directories under `tgs/`.
  - [x] Determine dates via `git log --diff-filter=A`.
  - [x] Infer status by checking files and content length.
  - [x] Extract short descriptions from thought `implementation.md` or `README.md`.
- Phase 2: Update README
  - [ ] Insert/replace a "## Current Thoughts" section with generated table.
  - [ ] Ensure links are relative and correct.

## 5. File/Module Changes
- Edit: `tgs/README.md` — append or update a "## Current Thoughts" section after usage/workflow sections.

## 6. Test Plan
- Render `tgs/README.md` in the IDE preview to confirm table layout.
- Click each link to confirm navigation.

## 7. Rollout & Rollback
- Rollout: Commit the README change.
- Rollback: Revert the README edit.

## 8. Estimates & Risks
- Estimate: ~15 minutes.
- Risks: Description extraction may be ambiguous; if unclear, use directory kebab title as description.

---
Approval checkpoint: Please review this plan and reply one of:
- APPROVE plan
- REQUEST CHANGES: <notes>
