# Thoughts (TGS) Directory

This directory contains organized thought processes, research, planning, and implementation records. Each subdirectory represents a complete thought cycle from research to implementation using the TGSFlow methodology.

## Directory Structure

Each thought is organized in a subdirectory with the naming convention:
```
<BASE_GIT_HASH>-<short-title-description>/
```

Where:
- **BASE_GIT_HASH**: The git commit hash at the moment the thought/research began
- **short-title-description**: A brief description of the thought/improvement

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
2. Create directory: `tgs/<hash>-<short-description>/`
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

## Bootstrapping vs Decorating

- Use bootstrap for greenfield projects: 
  ```bash
  curl -sSL https://raw.githubusercontent.com/akelv/tgsflow/main/bootstrap.sh | bash
  ```
  Follow prompts to select a template and project name.

- Use decorate for existing repositories (adds only the TGS workflow files to the current repo):
  ```bash
  curl -sSL https://raw.githubusercontent.com/akelv/tgsflow/main/bootstrap.sh | bash -s -- --decorate
  ```

Behavior:
- If you run `bootstrap.sh` in a directory that already contains a `.git` folder without `--decorate`, the script will prompt you to choose:
  - Decorate the current repository (recommended to adopt TGS in-place), or
  - Initialize a new project in a subdirectory, or
  - Quit.
- `--dry-run` is supported to preview changes.

## TGSFlow Workflow

This structure supports the TGSFlow methodology:
- **Human oversight**: Research and planning require explicit approval
- **AI implementation**: Detailed execution of approved plans  
- **Documentation**: Complete audit trail for all decisions
- **Traceability**: Every change links back to its thought process

This ensures thoughtful development with clear human-AI collaboration boundaries.

## Current Thoughts

| Thought Directory | Base Hash | Date | Status | Description |
|------------------|-----------|------|--------|-------------|
| [4c34cb8-update-tgs-readme-index](./4c34cb8-update-tgs-readme-index/) | 4c34cb8 | 2025-09-11 | 🚧 In Progress | Update TGS README index with Current Thoughts table |
| [612a57f-decorate-existing-software-project-repository](./612a57f-decorate-existing-software-project-repository/) | 612a57f | 2025-09-11 | ✅ Completed | Add decorate mode to inject TGS workflow into existing repos |
| [b4552ea-standardize-agentops-intake-to-pr-and-enrich-new-thought](./b4552ea-standardize-agentops-intake-to-pr-and-enrich-new-thought/) | b4552ea | 2025-09-11 | ✅ Completed | Standardize AGENTOPS workflow and enrich new-thought scaffolding |
| [f0d3f9a-add-agent-parent-command](./f0d3f9a-add-agent-parent-command/) | f0d3f9a | 2025-09-14 | ✅ Completed | Add `tgs agent` parent command delegating to `agent exec` |
| [5d12c3a-tiny-invisible-tgs-cli-for-huge-teams](./5d12c3a-tiny-invisible-tgs-cli-for-huge-teams/) | 5d12c3a | (prior) | ✅ Completed | Minimal TGS CLI scaffolding for large teams |
| [b48976e-refactor-cli-to-cobra-viper](./b48976e-refactor-cli-to-cobra-viper/) | b48976e | 2025-09-14 | ✅ Completed | Refactor CLI to Cobra/Viper; add completion and preserve behavior |