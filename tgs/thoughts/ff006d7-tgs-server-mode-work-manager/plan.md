# Plan: TGS Server Mode Work Manager

## 1. Objectives

- Provide a `tgs server backlog` command group for managing thought backlog
- Support **push model**: GitHub Actions automatically triggering Claude Code sessions
- Support **pull model**: Cloud-hosted or Claude app sessions querying and claiming work
- Implement atomic claim operations to prevent concurrent processing conflicts
- Track thought lifecycle: `queued` → `in_progress` → `completed`/`failed`
- Enable auto-discovery of approved thoughts (research.md + plan.md committed)
- Provide observability into backlog state and work progress

## 2. Scope / Non-goals

### In-Scope
- Core CLI commands: `add`, `remove`, `list`, `next`, `claim`, `complete`, `fail`, `validate`
- Backlog state file: `tgs/server/backlog.json`
- GitHub Action workflows: `tgs-server-run.yml`, `tgs-server-scan.yml`
- Git-based synchronization with optimistic locking
- Session context building (paths to research.md, plan.md, etc.)
- Basic logging and error handling
- Documentation and examples

### Out-of-Scope (Future Phases)
- Parallel processing (job matrix)
- Advanced observability (dashboards, metrics)
- Web UI for backlog visualization
- External database storage
- Advanced scheduling (priorities, dependencies)
- Automated retry logic
- Integration with issue trackers

## 3. Acceptance Criteria

### CLI Commands
- `tgs server backlog add <thought-dir>` enqueues a thought with validation
- `tgs server backlog list` shows all thoughts with status, timestamps
- `tgs server backlog next` atomically claims next available thought, outputs context
- `tgs server backlog complete <thought-dir>` marks thought completed, commits backlog
- `tgs server backlog fail <thought-dir>` marks thought failed, commits backlog
- `tgs server backlog validate` checks backlog consistency with filesystem

### Push Model (GitHub Actions)
- `tgs-server-run.yml` workflow can be manually triggered
- Workflow claims next thought, invokes `tgs agent exec`, updates status
- Workflow commits backlog.json updates after completion
- Handles failures gracefully (logs error, marks thought as failed)

### Pull Model (Remote Sessions)
- Claude Code session can run `tgs server backlog next` to get work
- Output includes thought dir, title, context file paths
- Session can update status via `complete` or `fail` commands
- Multiple sessions don't claim same thought (atomic operations)

### Data Integrity
- Backlog.json format matches schema (version, thoughts array)
- Git operations handle conflicts with retry logic
- Timestamps are ISO 8601 format
- Thought directories are validated before operations

## 4. Phases & Tasks

### Phase 1: Core Backlog Management (MVP)

#### Task 1.1: Backlog Data Model
- [ ] Define `Backlog` struct in `src/core/server/backlog.go`
  - Fields: Version, Thoughts array
- [ ] Define `ThoughtEntry` struct
  - Fields: Dir, Title, Priority, Status, AddedAt, ClaimedBy, ClaimedAt, StartedAt, CompletedAt
- [ ] Implement JSON marshal/unmarshal
- [ ] Add validation methods (valid paths, status transitions)

#### Task 1.2: Backlog File Operations
- [ ] Implement `Load(path string) (*Backlog, error)` in `src/core/server/backlog.go`
- [ ] Implement `Save(path string) error`
- [ ] Implement atomic read-modify-write with git operations
- [ ] Add file locking for concurrent safety

#### Task 1.3: Backlog Business Logic
- [ ] Implement `AddThought(dir, title string, priority int) error`
- [ ] Implement `RemoveThought(dir string) error`
- [ ] Implement `ClaimNext(claimedBy string) (*ThoughtEntry, error)` with atomic semantics
- [ ] Implement `UpdateStatus(dir string, status string) error`
- [ ] Implement `ListThoughts() []ThoughtEntry`
- [ ] Implement `ValidateBacklog() []ValidationError`

#### Task 1.4: CLI Commands
- [ ] Create `src/cmd/server.go` with parent command
- [ ] Create `src/cmd/server_backlog.go` with subcommands
- [ ] Implement `backlogAddCmd` handler
- [ ] Implement `backlogListCmd` handler (table output)
- [ ] Implement `backlogNextCmd` handler (outputs thought context)
- [ ] Implement `backlogClaimCmd` handler
- [ ] Implement `backlogCompleteCmd` handler
- [ ] Implement `backlogFailCmd` handler
- [ ] Implement `backlogRemoveCmd` handler
- [ ] Implement `backlogValidateCmd` handler
- [ ] Wire commands in `src/cmd/root.go`

#### Task 1.5: Git Integration
- [ ] Implement `GitPullBacklog() error` in `src/core/server/git.go`
- [ ] Implement `GitPushBacklog(message string) error`
- [ ] Implement `GitMergeBacklog() error` with conflict resolution
- [ ] Add retry logic with exponential backoff (max 3 retries)

#### Task 1.6: Thought Discovery
- [ ] Implement `DiscoverApprovedThoughts(repo string) ([]string, error)` in `src/core/server/discovery.go`
- [ ] Logic: Find dirs in `tgs/thoughts/*/` with research.md and plan.md committed, no implementation.md
- [ ] Use `git log` to verify files are committed

### Phase 2: GitHub Actions Workflows

#### Task 2.1: Push Model - Server Run Workflow
- [ ] Create `.github/workflows/tgs-server-run.yml`
- [ ] Trigger: `workflow_dispatch` (manual initially)
- [ ] Inputs: None (picks next from backlog)
- [ ] Steps:
  - Checkout with full history
  - Setup Go
  - Build tgs CLI
  - Run `tgs server backlog next` → capture thought dir
  - Build context file list from thought
  - Run `tgs agent exec` with context
  - On success: `tgs server backlog complete <thought>`
  - On failure: `tgs server backlog fail <thought>`
- [ ] Commit and push backlog updates

#### Task 2.2: Auto-Discovery Workflow
- [ ] Create `.github/workflows/tgs-server-scan.yml`
- [ ] Trigger: `push` to main branch
- [ ] Steps:
  - Checkout
  - Setup Go
  - Build tgs CLI
  - Run `tgs server backlog validate` (sync backlog)
  - Run thought discovery
  - Add new approved thoughts to backlog
  - Commit backlog.json if changes
  - Push or create PR with backlog updates

### Phase 3: Testing & Documentation

#### Task 3.1: Unit Tests
- [ ] Test backlog serialization/deserialization
- [ ] Test AddThought, RemoveThought, ClaimNext logic
- [ ] Test status transitions and validation
- [ ] Test concurrent claim operations (mutex behavior)
- [ ] Test git conflict resolution logic

#### Task 3.2: Integration Tests
- [ ] Test full CLI workflow: add → list → next → complete
- [ ] Test pull model: multiple sessions claiming work
- [ ] Test git push/pull/merge with mock conflicts

#### Task 3.3: Documentation
- [ ] Add `docs/server-mode.md` with architecture overview
- [ ] Document CLI commands with examples
- [ ] Document GitHub Actions setup and configuration
- [ ] Add runbook for common issues (conflicts, failures)
- [ ] Update main README.md with server mode section

## 5. File/Module Changes

### New Files

#### Go Source Files
- `src/core/server/backlog.go`
  - Backlog and ThoughtEntry structs
  - Load, Save, AddThought, RemoveThought, ClaimNext, UpdateStatus, ListThoughts
  - JSON serialization and validation
  - ~300 lines

- `src/core/server/git.go`
  - GitPullBacklog, GitPushBacklog, GitMergeBacklog
  - Retry logic with exponential backoff
  - Conflict detection and resolution
  - ~200 lines

- `src/core/server/discovery.go`
  - DiscoverApprovedThoughts
  - Git log parsing to verify committed files
  - Thought validation (research.md, plan.md present; implementation.md absent)
  - ~150 lines

- `src/cmd/server.go`
  - Parent `tgs server` command using Cobra
  - Help text and subcommand registration
  - ~50 lines

- `src/cmd/server_backlog.go`
  - Subcommands: add, list, next, claim, complete, fail, remove, validate
  - Each command 30-50 lines
  - Total ~400 lines

#### Test Files
- `src/core/server/backlog_test.go` (~300 lines)
- `src/core/server/git_test.go` (~200 lines)
- `src/core/server/discovery_test.go` (~150 lines)
- `src/cmd/server_backlog_test.go` (~200 lines)

#### GitHub Actions Workflows
- `.github/workflows/tgs-server-run.yml`
  - Push model: pick thought, run agent exec, update status
  - ~80 lines

- `.github/workflows/tgs-server-scan.yml`
  - Auto-discovery: scan thoughts, add to backlog
  - ~60 lines

#### Data Files
- `tgs/server/backlog.json`
  - Created on first `tgs server backlog add`
  - Git-tracked, human-readable JSON

#### Documentation
- `docs/server-mode.md` (~500 lines)
- `docs/server-mode-runbook.md` (~300 lines)

### Modified Files

- `src/cmd/root.go`
  - Register `server` parent command
  - +5 lines

- `README.md`
  - Add Server Mode section with quick start
  - Link to docs/server-mode.md
  - +30 lines

- `tgs/README.md`
  - Document tgs/server/ directory purpose
  - +10 lines

## 6. Test Plan

### Unit Tests
- **Backlog Operations**: Test add, remove, claim, update with various inputs
- **Status Transitions**: Test valid/invalid transitions (queued→in_progress, completed→queued should fail)
- **Serialization**: Test JSON round-trip with sample data
- **Validation**: Test thought dir validation, duplicate detection

### Integration Tests
- **CLI Workflow**: End-to-end test: add thought → list → next → complete
- **Git Operations**: Test pull → modify → push cycle with mock conflicts
- **Concurrent Claims**: Simulate multiple sessions claiming (mutex behavior)
- **Discovery**: Test with mock git repo containing approved thoughts

### Manual Testing
- **Push Model**: Trigger `tgs-server-run.yml` manually, verify thought processed
- **Pull Model**:
  - Clone repo in separate directory
  - Run `tgs server backlog next`
  - Verify thought claimed, output shows context
  - Run `tgs server backlog complete <thought>`
  - Verify backlog updated in git
- **Conflict Resolution**:
  - Start two sessions simultaneously
  - First claims thought A
  - Second tries to claim thought A (should get next available or error)

### Edge Cases
- Empty backlog (list, next should handle gracefully)
- Nonexistent thought dir (validation should catch)
- Backlog.json corrupted (should error with clear message)
- Git conflicts during push (should retry with backoff)
- Claim timeout/expiration (future: expired claims released)

## 7. Rollout & Rollback

### Rollout
1. Merge PR with Phase 1 implementation (CLI + core logic)
2. Verify `tgs server backlog` commands work locally
3. Create initial `tgs/server/backlog.json` with test thought
4. Merge PR with Phase 2 (GitHub Actions workflows)
5. Test workflows manually via `workflow_dispatch`
6. Enable auto-discovery workflow for main branch
7. Announce server mode in release notes

### Feature Flags
- None initially; commands are opt-in
- GitHub Actions workflows are manual dispatch initially
- Auto-scan can be disabled by removing workflow file

### Rollback
- Disable workflows: Delete or disable in GitHub UI
- Revert git commits: Standard git revert process
- Remove backlog.json: Not tracked until first use
- No database migrations or external dependencies

### Migration Path
- Existing projects: No changes required
- New adopters: Run `tgs server backlog add <thought>` to start
- Cloud sessions: Install tgs CLI in environment, configure git credentials

## 8. Estimates & Risks

### Time Estimates
- **Phase 1 (Core Backlog)**: 8-12 hours
  - Data model & file ops: 2h
  - Business logic: 3h
  - CLI commands: 3h
  - Git integration: 2h
  - Unit tests: 2h

- **Phase 2 (GitHub Actions)**: 4-6 hours
  - Server run workflow: 2h
  - Auto-scan workflow: 2h
  - Integration testing: 2h

- **Phase 3 (Docs & Polish)**: 4-6 hours
  - Documentation: 3h
  - Manual testing: 2h
  - Runbook: 1h

**Total Estimate**: 16-24 hours

### Risks & Mitigation

1. **Git Merge Conflicts**
   - Risk: High with multiple remote sessions
   - Mitigation: Optimistic locking, retry with backoff, clear error messages
   - Fallback: Claim expiration (future), mutex via GitHub API (future)

2. **Claude Code Session Failures**
   - Risk: Session crashes mid-implementation
   - Mitigation: Mark as failed in backlog, manual retry available
   - Fallback: Logs preserved for debugging

3. **Backlog Corruption**
   - Risk: Invalid JSON from merge conflict or manual edit
   - Mitigation: Validation on load, schema versioning, git history
   - Fallback: Restore from git history

4. **GitHub Actions Costs**
   - Risk: Long-running sessions consume minutes
   - Mitigation: 6-hour timeout, manual dispatch initially
   - Fallback: Disable workflows, use pull model only

5. **Complexity Creep**
   - Risk: Feature requests add complexity (priorities, dependencies)
   - Mitigation: Strict MVP scope, defer enhancements to Phase 4
   - Fallback: Keep simple file-based approach, reject over-engineering

### Dependencies
- **Go 1.23**: Already in use
- **Cobra/Viper**: Already in use for CLI
- **Git**: Required on all systems using server mode
- **GitHub Actions**: Optional (push model only)

### Open Questions
- **Claim Timeout**: Should claims expire? (Future: add `claim_timeout_minutes`)
- **Priority Ordering**: FIFO or priority-based? (Start with FIFO, add priority later)
- **Failure Retry**: Automatic or manual? (Start manual, add auto-retry later)

---

**Approval Checkpoint:** Please review this plan and reply one of:
- `APPROVE plan`
- `REQUEST CHANGES: <notes>`
