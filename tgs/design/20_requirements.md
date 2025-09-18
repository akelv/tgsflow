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

## Non-Functional Requirements
- **NFR-001**: The system shall ensure traceability such that each implemented change is linked to its originating thought directory. (Verification: Inspection)
- **NFR-002**: The system shall operate on macOS and Linux environments commonly used by developers. (Verification: Test)
- **NFR-003**: The `verify` command shall return exit code 0 on success and non-zero on failure. (Verification: Test)
- **NFR-004**: The system shall require each thought to include `research.md`, `plan.md`, and `implementation.md` before marking it complete. (Verification: Inspection)
- **NFR-005**: The system shall prevent production code from being implemented under `tgs/`, restricting that directory to thought documentation. (Verification: Inspection)

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
