# Context Pack

## Personas
- **Human Engineer/Developer** — uses AI code assistants with TGS to implement safely; goals: faster delivery with guardrails, clear approvals and audit trail; pain: unstructured AI output, unclear ownership.
- **Technical Lead/Reviewer** — approves research/plan, enforces quality gates; goals: traceability and alignment; pain: inconsistent processes across teams.
- **AI Code Agent** — executes implementation under human-approved plan; needs precise instructions, non-interactive commands, and clear constraints.
- **Compliance/Auditor** — needs transparent intention and traceability from thought to change; goals: reliable audit trail.
- **OSS Contributor/New Team Member** — needs quick onboarding via templates and thought docs; goals: clear contributing steps.

## Scope
- **In scope:** Approval-gated TGS workflow (Research → Plan → Human Approval → Implement → Document); thought directories and templates; `tgs` CLI usage; AI assistant integration via `agentops/AGENTOPS.md`; EARS linter enablement and verification; project templates and docs.
- **Out of scope:** Product-specific feature implementations; ML model training/selection; secrets management and vendor-specific setups; bespoke CI beyond provided templates.

## Scenarios
- When bootstrapping a project, the user runs the bootstrap script to apply TGSFlow, resulting in a repo with workflow docs and structure.
- When installing the CLI via Homebrew or curl, the user gets `tgs` available and verifies with `tgs --version`.
- When creating a new thought, the user runs `make new-thought title="..." spec="..."`, resulting in a scaffolded directory with research/plan/implementation docs.
- When integrating with an AI code assistant, the user sets the system prompt from `agentops/AGENTOPS.md`, and the assistant follows the approval-gated workflow.
- When enabling EARS linting, the user updates `tgs.yaml` (e.g., `policies.ears.enable: true`) and runs `./bin/tgs verify --repo .` to lint Markdown bullets.
- When contributing, the user follows TGS methodology, ensures approvals, and submits a PR with complete thought documentation in `tgs/`.

## Constraints
- Technical: Primary platforms macOS/Linux; `tgs` CLI via Homebrew/curl; Java + ANTLR required only to regenerate EARS parser (`make ears-gen`); prefer absolute paths and non-interactive flags; do not implement production code under `tgs/`.
- Business: Human approval gates are mandatory; maintain audit trail and documentation-driven process.
- Regulatory/Safety mindset: Emphasize traceability, precision, and INCOSE/EARS-informed documentation to keep software safe for human use.

## Contracts / Interfaces
- CLI `tgs`: `make new-thought ...`; verification via `./bin/tgs verify --repo .`.
- Make targets: `make new-thought`, `make ears-gen`.
- EARS grammar: `src/core/ears/ears.g4` and generated parser; optional toolchain (Java/ANTLR) for regeneration.
- AI assistant integration: system prompt at `agentops/AGENTOPS.md` (also referenced by `CLAUDE.md`/`AGENTS.md`).
- Project templates: `templates/` for React, Python, Go, CLI.

---

### Checklist
- [x] Personas cover all primary user types  
- [x] Scope is explicit (what’s in / what’s out)  
- [x] Key scenarios written in outcome-focused language  
- [x] Constraints documented and understood  
