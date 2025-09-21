# Plan: Context Pack command

## 1. Objectives
- Implement `tgs context pack "<query>"` to generate `aibrief.md` under active thought.  
- Use prompt templates for search and brief; respect token budget and include source links.

## 2. Scope / Non-goals
- In-scope: CLI command, prompt templates, brief template, config usage, integration with shell transport.  
- Non-goals: Remote indexing, vector DB, SDK transports.

## 3. Acceptance Criteria
- Command creates/overwrites `aibrief.md` with succinct sections, <= token budget.  
- Each item has a path and anchor/line range.  
- Fails gracefully with exit code and message on adapter error.

## 4. Phases & Tasks
- Phase 1: Templates
  - [ ] Add `templates/data/tgs/context/search_prompt.md.tmpl`
  - [ ] Add `templates/data/tgs/context/brief_template.md.tmpl`
- Phase 2: Command
  - [ ] Implement `src/cmd/context.go` with `pack` subcommand
  - [ ] Wire config, locate thought dir, gather candidate files
  - [ ] Invoke brain transport with budget and prompts
  - [ ] Write `aibrief.md`
- Phase 3: Tests & Docs
  - [ ] Unit tests for file gathering and write path
  - [ ] Update `tgs/design/*` V&V rows

## 5. File/Module Changes
- Add: `templates/data/tgs/context/search_prompt.md.tmpl`  
- Add: `templates/data/tgs/context/brief_template.md.tmpl`  
- Edit: `src/cmd/context.go` implement cobra `context pack`  
- Possibly add helpers under `src/core/thoughts/files.go` if needed.

## 6. Test Plan
- Run `tgs context pack "auth sso"` and inspect `aibrief.md`.  
- Simulate adapter absence to ensure clear error.  
- Validate token budget enforcement by truncation or model instruction.

## 7. Rollout & Rollback
- Feature is additive; can be disabled by not invoking command.  
- Rollback by reverting commit.

## 8. Estimates & Risks
- Estimate: 1 day.  
- Risks: adapter variability; mitigate with clear prompts and budget constraints.

---
Approval checkpoint: Please review this plan and reply one of:
- APPROVE plan
- REQUEST CHANGES: <notes>
