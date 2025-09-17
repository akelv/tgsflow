# Requirements: Tiny Go-based `tgs` CLI (EARS Notation)

This document restates the acceptance criteria of the spec (`10_spec.md`) using **EARS** syntax:  
- **While** <precondition>, **when** <trigger>, **the <system> shall** <response>.  
- Clauses appear in order; at least a System + Response, with optional preconditions/triggers.

---

## General

- When a developer installs the `tgs` binary, the system shall run on macOS, Linux, and Windows without external runtimes.  
- The system shall enforce that all commands (`init`, `context`, `specify`, `plan`, `tasks`, `approve`, `implement`, `drift-check`, `docs`, `verify`, `brief`, `watch`, `open-pr`) are available from a single binary.  
- While operating in enterprise environments, the system shall allow configuration via `tgs.yaml` to define approver roles, agent order, branch prefix, and policy limits.

---

## Init & Context

- When a user runs `tgs init --decorate`, the system shall create the `tgs/` folder structure and optional CI templates without modifying existing source code.  
- When `tgs context` runs on a repository with ≤5,000 files, the system shall scan languages, modules, and schema files, and the system shall complete within 60 seconds.  
- The system shall output `.context.json` and `00_research.md` as the baseline context.

---

## Specify, Plan, and Tasks

- When `specify` runs and Spec Kit CLI is available, the system shall proxy to Spec Kit; otherwise the system shall generate a minimal spec in `10_spec.md`.  
- When `plan` runs, the system shall create or append a `20_plan.md` with technical design and non-functional requirement sections.  
- When `tasks` runs, the system shall create or validate a `30_tasks.md` file containing implementable units, and the system shall fail if IDs or formatting are invalid.

---

## Approval & Governance

- When `approve --ci` runs, the system shall fail if any of `10_spec.md`, `20_plan.md`, `30_tasks.md`, or `40_approval.md` are missing.  
- While approver roles are configured in `tgs.yaml`, when `approve` runs, the system shall fail if required roles are not signed off in `40_approval.md`.  
- The system shall enforce that `20_plan.md` contains required NFR sections when `policies.enforce_nfr` is true.

---

## Implementation & Agents

- When `implement` runs for a given task, the system shall construct CONTEXT and PROMPT messages, invoke the preferred agent adapter, and receive patches.  
- While apply mode is `patch` (default), the system shall apply diffs itself with safety checks (≤300 LOC, no forbidden paths, no binary writes).  
- While apply mode is `agent`, the system shall snapshot files, allow direct edits, then run verification and revert offending hunks if policies are violated.  
- After patches apply, the system shall run `.tgs/hooks/{fmt,lint,test,perf}` and the system shall fail if any hook exits non-zero.  
- When patches pass hooks, the system shall create a PR on the configured platform with labels (`thought:<slug>`, `agent:<name>`), acceptance criteria references, and branch prefix from config.

---

## Drift Detection

- When `drift-check` runs, the system shall detect new routes, schemas, or public APIs and suggest updates to `10_spec.md` or `20_plan.md`.  
- The system shall fail in CI mode if divergence is detected.

---

## Documentation & Attestation

- When `docs` runs, the system shall produce `60_docs.md` containing runbooks and metrics links.  
- The system shall write an attestation JSON including approvers, task PRs, commit hashes, and timestamps.

---

## Interactive Mode (IDE / Senior Engineer)

- When `brief --task <id>` runs, the system shall output a compact brief (≤200 lines) including acceptance criteria, NFRs, and forbidden paths in Markdown or text.  
- When `verify` runs locally, the system shall run hooks, policy checks, and drift check, and the system shall exit non-zero on failure.  
- When `watch` runs, the system shall detect file changes and run scoped hooks after a quiet period, printing concise failure summaries.  
- When `open-pr` runs, the system shall create a PR labeled `tgs`, linked to the relevant Thought and task.

---

## Non-Functional

- The system shall produce a binary ≤30 MB and include build metadata.  
- The system shall sign releases with cosign and include SBOMs.  
- While operating in CI, the system shall complete `verify` on a medium repo in <4 minutes median.  
- The system shall fail if forbidden paths (e.g., `infra/prod/`, `secrets/`) are modified.  
- The system shall enforce that acceptance criteria use the word “shall” and not “should/may/can” (EARS lint).

