# Thoughts (TGS) Directory

This directory contains organized thought processes, research, planning, and implementation records. Each subdirectory represents a complete thought cycle from research to implementation using the TGSFlow methodology.

## Directory Structure

This `tgs/` area is organized into two primary subdirectories:

- `thoughts/` — Per-thought working directories created by `make new-thought`.
  - Naming: `<BASE_GIT_HASH>-<short-title-description>/`
  - Where:
    - **BASE_GIT_HASH**: The git commit hash at the moment the thought/research began
    - **short-title-description**: A brief description of the thought/improvement

- `design/` — Long-lived system-level design documentation.
  - Files: `00_context.md`, `10_needs.md`, `20_requirements.md`, `30_architecture.md`, `40_vnv.md`, `50_decisions.md`
  - Purpose: Architecture, requirements, verification/validation, and decision history maintained outside individual thoughts.

- `agentops/` — AgentOps workflow guide and thought templates used by `make new-thought`.
  - Files: `AGENTOPS.md` (system prompt and workflow), `tgs/*` (scaffold templates)
  - Purpose: Canonical workflow and the files used to scaffold new thought directories.

- `adapters/` — External model/tool adapters invoked by `tgs agent exec`.
  - Files: `claude-code.sh` (default Claude Code adapter), `gemini-code.sh` (Gemini CLI adapter)
  - Purpose: Bridge between the TGS CLI and external AI/automation tools. Default path used by the CLI: `tgs/adapters/claude-code.sh` (override with `--adapter-path`, e.g. `tgs/adapters/gemini-code.sh`).

## Thought Structure

Each thought directory contains:

- **`research.md`** - Problem analysis, constraint identification, and solution exploration
- **`plan.md`** - Detailed implementation plan with phases and technical specifications  
- **`implementation.md`** - Complete implementation summary and integration guide
- **`README.md`** - Navigation index and quick links to related files
## Purpose

This organizational structure provides:

1. **Traceability**: Each thought is linked to its originating git state
2. **Completeness**: Full research → plan → implementation → summary cycle
3. **Organization**: Related documentation grouped together
4. **History**: Clear evolution of ideas and implementations
5. **Context**: Preserved decision-making context for future reference

## Usage

When starting a new thought/improvement in a decorated or bootstrapped repo:

1. Get the current git HEAD hash: `git rev-parse --short HEAD`
2. Create directory: `tgs/thoughts/<hash>-<short-description>/`
3. Conduct research and create `research.md`
4. Develop plan and create `plan.md`
5. **Get human approval** for both research and plan
6. Implement changes according to approved plan
7. Document implementation in `implementation.md`
8. Update this index with the new thought entry

Or use the helper:
```bash
make new-thought title="Your idea here" spec="One-line idea spec (optional)"
```
## TGSFlow Workflow

This structure supports the TGSFlow methodology:
- **Human oversight**: Research and planning require explicit approval
- **AI implementation**: Detailed execution of approved plans  
- **Documentation**: Complete audit trail for all decisions
- **Traceability**: Every change links back to its thought process

This ensures thoughtful development with clear human-AI collaboration boundaries.

## Current Thoughts
- Base `be75a85` — Fix EARS regression tests — 2025-09-16 — Status: Implemented
  - Dir: `tgs/thoughts/3829798-fix-ears-regression-tests/`
  - Docs: `research.md`, `plan.md`, `implementation.md`

| Thought Directory | Base Hash | Date | Status | Description |
|------------------|-----------|------|--------|-------------|
| [612a57f-decorate-existing-software-project-repository](./thoughts/612a57f-decorate-existing-software-project-repository/) | 612a57f | 2025-09-11 | ✅ Completed | Add decorate mode to inject TGS workflow into existing repos |
| [b4552ea-standardize-agentops-intake-to-pr-and-enrich-new-thought](./thoughts/b4552ea-standardize-agentops-intake-to-pr-and-enrich-new-thought/) | b4552ea | 2025-09-11 | ✅ Completed | Standardize AGENTOPS workflow and enrich new-thought scaffolding |
| [f0d3f9a-add-agent-parent-command](./thoughts/f0d3f9a-add-agent-parent-command/) | f0d3f9a | 2025-09-14 | ✅ Completed | Add `tgs agent` parent command delegating to `agent exec` |
| [5d12c3a-tiny-invisible-tgs-cli-for-huge-teams](./thoughts/5d12c3a-tiny-invisible-tgs-cli-for-huge-teams/) | 5d12c3a | (prior) | ✅ Completed | Minimal TGS CLI scaffolding for large teams |
| [b48976e-refactor-cli-to-cobra-viper](./thoughts/b48976e-refactor-cli-to-cobra-viper/) | b48976e | 2025-09-14 | ✅ Completed | Refactor CLI to Cobra/Viper; add completion and preserve behavior |
| [43ec077-automate-releases-with-goreleaser-and-homebrew](./thoughts/43ec077-automate-releases-with-goreleaser-and-homebrew/) | 43ec077 | 2025-09-14 | ✅ Completed | Automate releases with GoReleaser, GitHub Actions and Homebrew |
| [f71f872-ears-linter-core-and-verify-integration](./thoughts/f71f872-ears-linter-core-and-verify-integration/) | f71f872 | 2025-09-15 | ✅ Completed | Implement EARS linter and integrate with verify |
| [1fd687a-implement-shell-transport-for-brain](./thoughts/1fd687a-implement-shell-transport-for-brain/) | 1fd687a | 2025-09-21 | ✅ Completed | Implement Shell Transport (Claude shell adapter) with tests |