# System Requirements

Use **“The system shall …”** style. One “shall” per requirement.  

## Functional Requirements
- **SR-001**: While executing a thought, the system shall block implementation actions until both `research.md` and `plan.md` are explicitly approved and recorded. (Verification: Demonstration)
- **SR-002**: While coordinating across teams, the system shall record approval metadata (approver and timestamp) within the thought directory for auditability. (Verification: Inspection)
- **SR-003**: While authoring changes, the system shall scaffold a thought directory containing `README.md`, `research.md`, `plan.md`, and `implementation.md` with a base hash and quick links. (Verification: Test)
- **SR-004**: While using AI coding assistants, the system shall provide a canonical system prompt in `agentops/AGENTOPS.md` for assistants to follow the TGS workflow. (Verification: Inspection)
- **SR-005**: When bootstrapping a repository, the system shall apply TGSFlow via a single non-interactive command that initializes required files and structure. (Verification: Demonstration)
- **SR-006**: While installing tooling, the system shall provide installation of the `tgs` CLI via Homebrew and via curl. (Verification: Test)
- **SR-007**: While planning work, the system shall include a system-wide design document set (`tgs/design/00_context.md`, `10_needs.md`, `20_requirements.md`, `30_architecture.md`) to align on context, needs, requirements, and architecture. (Verification: Inspection)
- **SR-008**: While writing needs and requirements, the system shall offer INCOSE/EARS-aligned guidance and optional linting rules. (Verification: Demonstration)
- **SR-009**: When verifying documentation, the system shall provide a `verify` command that evaluates EARS patterns and Markdown bullets and reports pass/fail via exit code. (Verification: Test)
- **SR-010**: While contributing, the system shall scaffold new thoughts via `make new-thought` with provided `title` and optional `spec`. (Verification: Test)
- **SR-011**: While operating in audited contexts, the system shall maintain an audit trail linking approvals and implementation to each thought and reference it in PR content. (Verification: Inspection)
- **SR-012**: While executing tasks, the system shall support non-interactive execution of setup and verification commands. (Verification: Demonstration)
- **SR-013**: While starting new services, the system shall provide ready-to-use templates for React, Python, Go, and CLI under `templates/`. (Verification: Inspection)
- **SR-014**: While maintaining safety for human use, the system shall enforce an approval-gated workflow with phases: Research → Plan → Human Approval → Implement → Document. (Verification: Demonstration)

- **SR-015**: When initializing or decorating a project via the bootstrap script, the system shall scaffold a minimal `tgs/` directory using repository-agnostic templates under `templates/data/tgs/`, excluding project-specific thought history. (Verification: Test)

- **SR-016**: When running `tgs init`, the system shall mirror an organization-approved `tgs/` scaffolding from either embedded templates, a local directory, a remote archive, or a standard git repository (optionally at a specified ref and subdirectory). (Verification: Test)
- **SR-017**: While rendering scaffolding, the system shall apply `.tmpl` templates and copy non-template files as-is, preserving existing files for idempotency. (Verification: Test)
- **SR-018**: While sourcing templates remotely, the system shall clean up temporary files and directories after completion. (Verification: Inspection)
- **SR-019**: While applying guardrails, the system shall include `tgs/agentops/AGENTOPS.md` and design docs under `tgs/design/` to enforce the approval-gated workflow for AI agent collaboration. (Verification: Inspection)

- **SR-020**: While operating in shell AI mode, the system shall provide a Shell Transport that executes the default adapter (`tgs/adapters/claude-code.sh`) with a composed prompt and optional context files, and returns the adapter output as `ChatResp.Text`, honoring timeouts and exit codes. (Verification: Test)

- **SR-027**: While operating in shell AI mode, the system shall support a Gemini adapter script (`tgs/adapters/gemini-code.sh`) with the same interface as the Claude adapter (prompt via `--prompt-text|--prompt-file`, deterministic context file expansion, optional timeout, suggestions routing), returning text to stdout and non-zero on error. (Verification: Test)

- **SR-028**: When running `tgs init`, the system shall also ensure adapter scripts exist under `tgs/adapters/` (at least `claude-code.sh` and `gemini-code.sh`), copying from embedded templates if missing. (Verification: Test)
- **SR-029**: The system shall support `tgs init claude` and `tgs init gemini` subcommands that perform additional decoration such as copying `tgs/agentops/AGENTOPS.md` to the repository root as `CLAUDE.md` or `GEMINI.md` only if absent, otherwise exit with a clear error instructing manual override. (Verification: Test)
- **SR-030**: When initializing a repository, if the root `Makefile` lacks a `new-thought` target, the system shall append or create it with the standard implementation to enable the TGS workflow. (Verification: Test)

- **SR-021**: When running `tgs context pack "<query>"`, the system shall collect relevant sections from `tgs/design/` and the active thought directory into `aibrief.md`. (Verification: Test)
- **SR-022**: The system shall construct the brief using templated prompts to guide a repository-aware search via the brain shell agent. (Verification: Inspection)
- **SR-023**: The system shall include source pointers (path and anchor/line range) for each extracted requirement/context item. (Verification: Inspection)
- **SR-024**: The system shall enforce a configurable token budget for the brief content via `ai.toolpack.budgets.context_pack_tokens`. (Verification: Test)
- **SR-025**: The system shall avoid including secrets or sensitive tokens by applying configured redaction rules. (Verification: Analysis)

- **SR-026**: When verifying design documentation, the system shall lint `tgs/design/10_needs.md` and `tgs/design/20_requirements.md` for EARS compliance (ignoring code fences and allowing bullet response sections), and report findings with `path:line: message`. (Verification: Test)

## Non-Functional Requirements
- **NFR-001**: The system shall ensure traceability such that each implemented change is linked to its originating thought directory. (Verification: Inspection)
- **NFR-002**: The system shall operate on macOS and Linux environments commonly used by developers. (Verification: Test)
- **NFR-003**: The `verify` command shall return exit code 0 on success and non-zero on failure. (Verification: Test)
- **NFR-004**: The system shall require each thought to include `research.md`, `plan.md`, and `implementation.md` before marking it complete. (Verification: Inspection)
- **NFR-005**: The system shall prevent production code from being implemented under `tgs/`, restricting that directory to thought documentation. (Verification: Inspection)
 - **NFR-006**: The `tgs context pack` command shall complete within 30 seconds under default settings on a medium repo. (Verification: Demonstration)

## Interfaces
- **IF-001**: The system shall expose a `make new-thought` target that creates `tgs/<BASE_HASH>-<kebab-title>/` with pre-populated templates. (Verification: Test)
- **IF-002**: The system shall expose a `make ears-gen` target that regenerates the ANTLR parser from `src/core/ears/ears.g4`. (Verification: Demonstration)
- **IF-003**: The system shall provide a `tgs verify --repo <PATH>` command to run documentation linting. (Verification: Test)
- **IF-004**: The system shall provide a repository bootstrap script at `bootstrap.sh` for applying TGSFlow via curl. (Verification: Demonstration)
- **IF-005**: The system shall provide the canonical system prompt at `agentops/AGENTOPS.md` for AI assistant configuration. (Verification: Inspection)

---

### Checklist
- [ ] Each requirement is singular and testable  
- [ ] Uses “shall” (no “should/may”)  
- [ ] Quantified criteria included (units, thresholds)  
- [ ] Verification method assigned (Inspection / Demonstration / Test / Analysis)  
