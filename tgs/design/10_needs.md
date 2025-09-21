# Stakeholder Needs

Format: Use simple, outcome-focused sentences (EARS style if possible).  

## Needs List
- **N-001**: While creating software with AI, the Human Engineer needs a workflow that keeps strategic decisions human-approved and implementation safe.
- **N-002**: While coordinating across teams, the Technical Lead needs explicit approval gates to enforce quality and control change.
- **N-003**: While authoring changes, stakeholders need traceability from thought to implementation and documentation.
- **N-004**: While using AI coding agents, the Human Engineer needs a standardized system prompt so assistants follow the TGS workflow.
- **N-005**: When bootstrapping a project, the Human Engineer needs a one-command setup to apply TGSFlow.
- **N-006**: While installing tooling, the Human Engineer needs an easy `tgs` CLI installation via homebrew or curl.
- **N-007**: While planning work, the Team needs a system wide design document to align on context, needs, requirements, architecture
- **N-008**: While writing needs and requirements, the Team needs INCOSE/EARS-aligned guidance and optional linting.
- **N-009**: When verifying docs, the Team needs a `verify` command that checks EARS patterns and Markdown bullets.
- **N-010**: While contributing, the OSS Contributor needs a clear path to create a thought, get approvals, and submit a PR with linked docs.
- **N-011**: While operating in audited contexts, the Compliance/Auditor needs a reliable audit trail of decisions and approvals.
- **N-012**: While executing tasks, the AI Code Agent needs non-interactive commands and explicit, precise instructions.
- **N-013**: While starting new services, the Team needs ready-to-use templates (React, Python, Go, CLI) to bootrap work.
- **N-014**: While maintaining safety for human use, the Team needs a specification driven development process emphasizing precision, accountability, and rollback.

- **N-015**: When bootstrapping or decorating a repository, the Human Engineer needs the script to initialize a clean, minimal `tgs/` directory from generic templates (not repo-specific documents or history).

- **N-016**: While bootstrapping work with AI agents, the Human Engineer needs a standard, organization-approved scaffolding to start quickly while staying within team guardrails.
- **N-017**: While standardizing across repositories, the Team needs the scaffolding templates to come from a shared source (git repo, local directory, or remote archive) without changing CLI code.
- **N-018**: While re-running initialization, the Human Engineer needs idempotent, non-destructive behavior that preserves any existing files.

- **N-019**: While initiating an AI task, the Human Engineer needs a concise, auto-generated brief that packs relevant context from `tgs/design/` and the active thought.
- **N-020**: While controlling LLM costs, the Team needs a configurable token budget for briefs with safe defaults.
- **N-021**: While reading briefs, the Team needs explicit source links back to original documents for verification and deeper reading.

- **N-022**: When verifying core design docs, the Team needs focused linting of `tgs/design/10_needs.md` and `tgs/design/20_requirements.md` to enforce EARS patterns and ID formatting.

---

### Checklist
- [x] Each need is solution-free (no design baked in)  
- [x] Each need is stated in user/stakeholder voice  
- [x] Needs are validated with actual stakeholders  
- [x] Each need has an ID (N-###)  
