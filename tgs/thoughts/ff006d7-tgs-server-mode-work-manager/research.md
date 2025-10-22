# Research: TGS Server Mode Work Manager

- Date: 2025-10-22
- Base Hash: ff006d7
- Participants: Claude Code / Human

## 1. Problem Statement

The current TGSFlow workflow requires manual intervention to:
- Identify approved thoughts (with completed research.md and plan.md committed to git)
- Start Claude Code sessions to implement each thought
- Monitor implementation progress
- Ensure PRs are created and submitted
- Move to the next thought in the backlog

This manual orchestration creates bottlenecks when multiple thoughts are ready for implementation. For cloud-hosted or CI/CD environments, we need an autonomous work manager that can:
- Maintain a backlog of approved, implementation-ready thoughts
- Automatically trigger Claude Code sessions to implement thoughts sequentially
- Track session progress and outcomes (success, failure, partial completion)
- Create PRs for completed implementations
- Dequeue completed thoughts and move to the next item
- Provide visibility into work progress and backlog status

**Desired Outcome:** A "server mode" for TGS that acts as an autonomous work orchestrator, enabling continuous implementation of approved thoughts without human intervention for the execution phase (while maintaining human oversight for research/plan approval and PR review).

## 2. Current State

### Existing Components
- **TGS CLI**: Go-based CLI with Cobra/Viper (`src/cmd/`)
- **Agent Execution**: `tgs agent exec` command bridges to adapter scripts
- **Claude Code Adapter**: `tgs/adapters/claude-code.sh` handles Claude Code invocation
- **Thought Structure**: Standardized directories under `tgs/thoughts/<hash>-<title>/`
- **GitHub Actions**: CI/CD workflows for build/test/release (`.github/workflows/`)

### Current Workflow (Manual)
1. Human creates thought via `make new-thought`
2. AI agent writes `research.md` → human approval
3. AI agent writes `plan.md` → human approval
4. Human commits approved research and plan to git
5. **[Manual]** Human triggers Claude Code session with thought context
6. Claude Code implements according to approved plan
7. **[Manual]** Human monitors progress and reviews PR
8. **[Manual]** Human selects next thought and repeats

### Gaps
- No backlog management system
- No automated detection of implementation-ready thoughts
- No automated Claude Code session orchestration
- No progress tracking or session state management
- No automated queuing/dequeuing of work items

## 3. Constraints & Assumptions

### Technical Constraints
- Must work with existing TGS thought structure
- Must integrate with existing `tgs agent exec` and adapter system
- Must be runnable in GitHub Actions environment (Ubuntu runner)
- Claude Code sessions may be long-running (minutes to hours)
- GitHub Actions has 6-hour job timeout limit

### Assumptions
- Approved thoughts are those with both research.md and plan.md committed and not yet implemented
- A thought is "implemented" when implementation.md exists and PR is created
- Only one thought should be implemented at a time initially (serial processing)
- The server mode will run as a GitHub Actions workflow (not a persistent service)
- Humans will still approve research/plan phases and review PRs

### Architectural Constraints
- Implementation must live outside `tgs/` (per TGSFlow guidelines)
- Must use Go for CLI extensions (consistent with project)
- Backlog state should be git-trackable (file-based, not external DB)
- Must preserve existing `tgs agent exec` interface

## 4. Risks & Impact

### Risks
1. **Concurrency Issues**
   - Risk: Multiple workflows trying to process same thought
   - Mitigation: File-based locking mechanism, workflow concurrency limits

2. **Long-Running Sessions**
   - Risk: GitHub Actions timeout (6 hours max)
   - Mitigation: Break large thoughts into smaller sub-thoughts, add session timeout configuration

3. **Session Failures**
   - Risk: Claude Code crashes, network issues, invalid plans
   - Mitigation: Retry logic with exponential backoff, failure alerting, manual intervention hooks

4. **Backlog Drift**
   - Risk: Backlog file becomes out of sync with actual thought state
   - Mitigation: Backlog validation/sync command, git-based source of truth

5. **Resource Exhaustion**
   - Risk: Too many thoughts queued, runner costs
   - Mitigation: Backlog depth limits, cost monitoring, manual approval gate for server mode runs

### Security Considerations
1. **Credential Handling**
   - Claude Code may need API keys (ANTHROPIC_API_KEY)
   - GitHub token for PR creation (GITHUB_TOKEN)
   - Store in GitHub Secrets, pass via environment variables

2. **Code Injection**
   - Thought titles/specs could contain malicious content
   - Sanitize all inputs, validate thought directory structure

3. **Audit Trail**
   - All server mode actions must be logged with timestamps
   - Preserve session transcripts for debugging and compliance

4. **Access Control**
   - Only authorized workflows can modify backlog
   - PR creation requires proper authentication

## 5. Alternatives Considered

### A) File-Based Backlog with GitHub Actions Dispatcher
**Implementation:**
- Backlog stored as `tgs/server/backlog.json` (or YAML)
- GitHub Action workflow: `tgs-server-dispatch.yml`
- Triggered manually (`workflow_dispatch`) or on schedule (`cron`)
- Workflow reads backlog, picks next thought, invokes `tgs agent exec`
- On completion, updates backlog and commits progress

**Pros:**
- Simple, no external dependencies
- Git-tracked backlog (versioned, reviewable)
- Easy to implement with existing tools
- Runs in standard GitHub Actions runner

**Cons:**
- Manual workflow triggering (unless scheduled)
- Limited to GitHub Actions 6-hour timeout per thought
- Serial processing only (one thought at a time)
- Backlog updates require git commits (may cause conflicts)

### B) GitHub Issues as Backlog with Labels
**Implementation:**
- Each approved thought creates a GitHub Issue with label `tgs:ready`
- GitHub Action watches for issues with this label
- Workflow picks oldest issue, implements thought, closes issue on completion
- Issue comments track progress

**Pros:**
- Native GitHub integration (no custom backlog file)
- Great visibility (issues board shows backlog)
- Built-in commenting/progress tracking
- Can use issue templates for standardization

**Cons:**
- Tight coupling to GitHub (harder to test locally)
- More complex issue/label management
- Potential issue spam if many thoughts queued
- Harder to maintain strict ordering

### C) Separate Persistent Service with API
**Implementation:**
- Standalone service (Go/Python) running on cloud (e.g., GCP Cloud Run)
- Exposes REST API for backlog management
- Persists state in database (PostgreSQL, SQLite)
- GitHub Actions webhook triggers service to enqueue thoughts
- Service spawns Claude Code sessions (docker containers or API calls)

**Pros:**
- No GitHub Actions timeout limits
- True continuous processing
- Scalable to multiple parallel workers
- Rich state management and observability

**Cons:**
- Infrastructure complexity (service hosting, DB, monitoring)
- Operational overhead (deployments, scaling, costs)
- Security complexity (API authentication, secrets management)
- Overengineered for current MVP needs

### D) Hybrid: GitHub Actions Workflow with Job Matrix
**Implementation:**
- GitHub Actions with dynamic matrix generation
- Reads backlog file, generates matrix of thoughts
- Runs multiple thoughts in parallel as separate jobs
- Each job runs `tgs agent exec` for its assigned thought
- Coordinator job merges results and updates backlog

**Pros:**
- Parallel processing of multiple thoughts
- Still uses GitHub Actions (no external service)
- Faster throughput

**Cons:**
- Complex matrix coordination
- Higher chance of conflicts (multiple jobs updating backlog)
- Harder to debug
- Higher GitHub Actions costs

## 6. Recommendation

**Adopt Alternative A: File-Based Backlog with GitHub Actions Dispatcher**

### Rationale
- **Simplicity**: Minimal new infrastructure, leverages existing TGS architecture
- **Git-Native**: Backlog is version-controlled, auditable, and part of the repo
- **Incremental**: Can evolve to Alternative C (persistent service) later if needed
- **Cost-Effective**: Uses included GitHub Actions minutes
- **Testable**: Can run locally using `act` or similar tools

### Phased Approach

#### Phase 1: Core Backlog & Orchestration (MVP)
- Implement `tgs server backlog` command group:
  - `add <thought-dir>`: Enqueue a thought
  - `remove <thought-dir>`: Dequeue a thought
  - `list`: Show backlog with status
  - `validate`: Check backlog consistency with git state
- Backlog file: `tgs/server/backlog.json`
  ```json
  {
    "version": "1",
    "thoughts": [
      {
        "dir": "tgs/thoughts/abc123-feature-x",
        "title": "Implement feature X",
        "priority": 10,
        "status": "queued",
        "added_at": "2025-10-22T10:00:00Z",
        "started_at": null,
        "completed_at": null
      }
    ]
  }
  ```
- GitHub Action: `tgs-server-run.yml` (manual dispatch initially)
  - Reads backlog
  - Picks highest priority thought with status "queued"
  - Invokes `tgs agent exec` with thought context
  - Updates thought status: `queued` → `in_progress` → `completed`/`failed`
  - Commits backlog updates

#### Phase 2: Auto-Discovery & Submission
- GitHub Action: `tgs-server-scan.yml` (runs on push to main)
  - Scans `tgs/thoughts/*/` for approved thoughts
  - Criteria: both research.md and plan.md committed, no implementation.md
  - Auto-adds discovered thoughts to backlog
  - Creates PR or commit to update backlog

#### Phase 3: Observability & Reliability
- Session logging to `tgs/server/logs/<thought>-<timestamp>.log`
- Failure handling: retry logic, alerting via GitHub Issue creation
- Status dashboard: GitHub Pages site showing backlog and progress
- Metrics: CSV file tracking throughput, success rate

#### Phase 4: Parallel Processing (Future)
- Extend to Alternative D (job matrix) if throughput demands
- Resource pooling for API rate limits

### Why Not the Other Alternatives?
- **B (GitHub Issues)**: Good for visibility, but adds complexity for MVP. Could be a future enhancement for user-facing backlog view.
- **C (Persistent Service)**: Overengineered for current scale. Re-evaluate if we need <30min response time or >10 thoughts/day.
- **D (Parallel Matrix)**: Adds complexity without clear immediate benefit. Defer until serial processing is a bottleneck.

## 7. References & Links

### Existing Code
- `src/cmd/agent_exec.go`: Current adapter execution logic
- `tgs/adapters/claude-code.sh`: Claude Code invocation script
- `.github/workflows/ci.yml`: Example GitHub Actions workflow
- `tgs/thoughts/*/`: Existing thought structure and examples

### Documentation
- TGSFlow Workflow: `tgs/agentops/AGENTOPS.md`
- Thought Organization: `tgs/README.md`
- GitHub Actions: https://docs.github.com/en/actions
- Cobra CLI: https://github.com/spf13/cobra

### Related Thoughts
- [f0d3f9a-add-agent-parent-command](../f0d3f9a-add-agent-parent-command/): Created `tgs agent` command group
- [43ec077-automate-releases-with-goreleaser-and-homebrew](../43ec077-automate-releases-with-goreleaser-and-homebrew/): GitHub Actions automation patterns

---

**Approval Checkpoint:** Please review this research and reply one of:
- `APPROVE research`
- `REQUEST CHANGES: <notes>`
