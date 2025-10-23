# TGS Server Mode

**Autonomous work orchestration for approved thoughts**

## Overview

TGS Server Mode is a work management system that enables autonomous implementation of approved thoughts by Claude Code sessions. It maintains a backlog of implementation-ready thoughts and supports both:

- **Push Model**: GitHub Actions automatically triggering Claude Code sessions
- **Pull Model**: Cloud/remote Claude Code sessions querying the backlog for work

## Architecture

### Components

```
┌─────────────────────────────────────────────────────┐
│                  GitHub Repository                   │
│                                                       │
│  ┌──────────────────────────────────────────────┐  │
│  │         tgs/server/backlog.json              │  │
│  │  (Git-tracked work queue)                     │  │
│  └──────────────────────────────────────────────┘  │
│                         ▲                             │
│                         │                             │
│            ┌────────────┴────────────┐               │
│            │                          │               │
│      Push Model                  Pull Model          │
│            │                          │               │
│  ┌─────────▼────────┐     ┌─────────▼──────────┐   │
│  │ GitHub Actions   │     │ Remote/Cloud        │   │
│  │ Workflow         │     │ Claude Code         │   │
│  │ (tgs-server-run) │     │ Sessions            │   │
│  └──────────────────┘     └─────────────────────┘   │
└─────────────────────────────────────────────────────┘
```

### Backlog State File

Location: `tgs/server/backlog.json`

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
      "claimed_by": null,
      "claimed_at": null,
      "started_at": null,
      "completed_at": null
    }
  ]
}
```

### Status Lifecycle

```
queued → in_progress → completed
                     ↘ failed → (manual retry) → queued
```

## CLI Commands

### Backlog Management

#### Add a thought to backlog

```bash
tgs server backlog add <thought-dir> [--priority N]
```

Example:
```bash
tgs server backlog add tgs/thoughts/abc123-feature-x --priority 10
```

#### List backlog

```bash
tgs server backlog list
```

Output:
```
DIR                               TITLE              STATUS       PRIORITY  ADDED       CLAIMED BY
tgs/thoughts/abc123-feature-x     Implement feature  queued       10        2025-10-22  -
tgs/thoughts/def456-refactor      Refactor code      in_progress  5         2025-10-21  session-123
```

#### Claim next available thought (Pull Model)

```bash
tgs server backlog next [--claimed-by <identifier>]
```

Output:
```
Claimed: tgs/thoughts/abc123-feature-x
Title: Implement feature X
Priority: 10

Context files:
  tgs/thoughts/abc123-feature-x/research.md
  tgs/thoughts/abc123-feature-x/plan.md
  tgs/agentops/AGENTOPS.md
  CLAUDE.md
```

#### Mark thought as completed

```bash
tgs server backlog complete <thought-dir>
```

#### Mark thought as failed

```bash
tgs server backlog fail <thought-dir>
```

#### Remove thought from backlog

```bash
tgs server backlog remove <thought-dir>
```

#### Validate backlog consistency

```bash
tgs server backlog validate
```

## Workflows

### Push Model: GitHub Actions

The `tgs-server-run.yml` workflow automatically processes thoughts from the backlog.

**Trigger**: Manual dispatch via GitHub UI

**Steps**:
1. Claim next highest-priority thought
2. Invoke Claude Code to implement
3. Mark as completed or failed
4. Commit backlog updates

**Usage**:
1. Go to Actions tab in GitHub
2. Select "TGS Server Run" workflow
3. Click "Run workflow"
4. (Optional) Specify claimed_by identifier

### Pull Model: Remote Sessions

Cloud-hosted or Claude app sessions can pull work from the backlog.

**Workflow**:

1. Session queries for next available work:
   ```bash
   tgs server backlog next --claimed-by "cloud-session-123"
   ```

2. Session implements the thought using provided context files

3. Session marks work complete:
   ```bash
   tgs server backlog complete <thought-dir>
   ```

**Example** (Cloud Claude Code session):

```bash
# Clone repo
git clone https://github.com/user/repo.git
cd repo

# Claim next thought
THOUGHT_INFO=$(tgs server backlog next --claimed-by "cloud-$(hostname)")
THOUGHT_DIR=$(echo "$THOUGHT_INFO" | grep "^Claimed:" | cut -d' ' -f2)

# Implement thought (via tgs agent exec or manual implementation)
# ... implementation work ...

# Mark complete
tgs server backlog complete "$THOUGHT_DIR"
```

## Auto-Discovery

The `tgs-server-scan.yml` workflow automatically discovers approved thoughts and adds them to the backlog.

**Trigger**: Push to main branch (changes to `tgs/thoughts/*/research.md` or `plan.md`)

**Logic**:
- Scans `tgs/thoughts/*/` directories
- Checks if research.md and plan.md are committed
- Checks if implementation.md is absent
- Adds qualifying thoughts to backlog with priority=0

**Manual trigger**:
```bash
# Via GitHub Actions UI, or programmatically:
gh workflow run tgs-server-scan.yml
```

## Concurrency & Conflict Resolution

### Git-Based Locking

Server mode uses **optimistic locking** with git for concurrency control:

1. **Pull** latest backlog.json from remote
2. **Modify** backlog (claim, complete, etc.)
3. **Push** changes to remote
4. **Retry** on conflict (exponential backoff, max 3 retries)

### Atomic Claim Operation

The `tgs server backlog next` command performs atomic claim-and-push:

```go
// Pseudocode
func AtomicClaimAndPush() {
  retry (max 3 times with exponential backoff) {
    git pull origin main
    backlog = load("backlog.json")
    claimed = backlog.ClaimNext(claimedBy)
    save(backlog)
    git push origin main
  }
}
```

If push fails due to conflict (another session claimed simultaneously), the operation retries.

### Best Practices for Remote Sessions

- **Use unique claimed_by identifiers** (e.g., `cloud-session-uuid`)
- **Handle claim failures gracefully** (backlog may be empty)
- **Don't abandon claimed thoughts** (mark as failed if can't complete)
- **Coordinate if running many parallel sessions** (consider rate limiting claims)

## Configuration

### Backlog Path

Default: `tgs/server/backlog.json`

Override via environment variable:
```bash
export TGS_BACKLOG_PATH="custom/path/backlog.json"
```

### Git Branch

Server mode operates on the current git branch. Ensure you're on the correct branch before running commands.

### Priorities

- **Higher number = more urgent**
- Default priority: 0
- Suggested scale: 0 (low) to 10 (critical)

## Security Considerations

### Credentials

- **GitHub Actions**: Use `GITHUB_TOKEN` secret (auto-provided)
- **Remote sessions**: Configure git credentials (SSH keys or tokens)
- **Claude Code**: Set `ANTHROPIC_API_KEY` via environment or secrets

### Access Control

- Only authorized workflows can modify backlog (enforce via branch protection)
- Remote sessions must have git push access
- Consider requiring PR reviews for backlog changes

### Audit Trail

All backlog modifications are git-committed with messages:
```
chore(tgs): claim thought tgs/thoughts/abc123 by session-id
chore(tgs): complete thought tgs/thoughts/abc123
```

Review git history for full audit trail.

## Troubleshooting

### Backlog validation errors

```bash
tgs server backlog validate
```

Common errors:
- Duplicate thought dirs → Remove duplicates manually in backlog.json
- Invalid status values → Fix statuses to: queued, in_progress, completed, failed
- Missing claimedBy for in_progress → Add or reset status to queued

### Git push conflicts

If `tgs server backlog` commands fail with push errors:

1. Pull latest changes:
   ```bash
   git pull origin main
   ```

2. Retry the command (automatic retry should handle this)

3. If persistent conflicts, check for other sessions modifying backlog simultaneously

### Empty backlog

If `tgs server backlog next` returns "No work available":

- Check `tgs server backlog list` to see current state
- Run `tgs-server-scan` workflow to discover approved thoughts
- Manually add thoughts via `tgs server backlog add`

### Session timeout

GitHub Actions has 6-hour timeout. For long-running thoughts:
- Break into smaller sub-thoughts
- Increase timeout in workflow (max 6 hours)
- Use pull model with cloud sessions (no timeout limit)

## Future Enhancements

- Claim timeout/expiration (auto-release abandoned claims)
- Advanced scheduling (dependencies between thoughts)
- Parallel processing (multiple thoughts simultaneously)
- Web UI for backlog visualization
- Metrics and analytics (throughput, success rate)
- Integration with issue trackers

## References

- Research: [tgs/thoughts/ff006d7-tgs-server-mode-work-manager/research.md](../tgs/thoughts/ff006d7-tgs-server-mode-work-manager/research.md)
- Plan: [tgs/thoughts/ff006d7-tgs-server-mode-work-manager/plan.md](../tgs/thoughts/ff006d7-tgs-server-mode-work-manager/plan.md)
- Code: `src/core/server/` and `src/cmd/server*.go`
