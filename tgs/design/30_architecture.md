# Architecture Overview

## System Context (C4 Level 1)
- TGSFlow is a CLI-first workflow tool used by humans and AI code agents to enforce an approval-gated software process.
- External actors/systems:
  - Human Engineer/Developer (runs the CLI, authors thoughts)
  - Technical Lead/Approver (approves `research.md` and `plan.md`)
  - AI Code Agent (invoked via `tgs agent exec` through adapter)
  - Compliance/Auditor (reads thought docs and audit trail)
  - GitHub/GitLab CI (optional workflows and verification)
  - Homebrew/GitHub Releases (CLI distribution)
  - Java/ANTLR toolchain (only when regenerating EARS parser)

## Containers (C4 Level 2)
- CLI binary `tgs` (Go, Cobra/Viper) — entry point for commands (`init`, `context`, `verify`, `agent exec`).
- Repository workspace — contains `tgs/` thought directories, design docs, templates, and source.
- Templates (embedded via `embed.FS`) — Markdown and CI stubs used by `init`/writers.
- EARS Linter and Parser — generated ANTLR artifacts in `src/core/ears/gen/...` consumed at runtime.
- Adapter process (`adapters/claude-code.sh`) — spawned by `tgs agent exec` to call external AI.
- CI pipelines — optional GitHub/GitLab workflows seeded by `init`.

## Components (C4 Level 3)
- CLI Commands (`src/cmd/`):
  - `root.go` — builds root command, flags, and wires subcommands.
  - `init.go` — `tgs init` scaffolds TGS layout and optional CI workflows.
  - `context.go` — `tgs context` writes `.context.json` and seeds research.
  - `verify.go` — `tgs verify` loads `tgs.yaml`, runs EARS lint and repo hooks.
  - `agent.go` — parent for agent commands; prints help and routes.
  - `agent_exec.go` — `tgs agent exec` runs adapter with prompt/context and returns output.
- Core (`src/core/`):
  - `thoughts/*` — ensure/locate thought dirs, append sections, file ops.
  - `ears/*` — EARS grammar integration, requirement parsing and linting helpers.
  - `config/loader.go` — loads `tgs.yaml` with defaults and policies.
- Templates (`src/templates/*`):
  - `render.go` — renders embedded templates.
  - `write.go` — writes rendered templates if missing.
- Utils (`src/util/logx/logx.go`) — simple logging with optional JSONL.

## Interfaces
- CLI entrypoints:
  - `tgs init [--decorate] [--ci-template github|gitlab|none]`
  - `tgs context`
  - `tgs verify [--ci] [--repo PATH]`
  - `tgs agent exec --prompt-text|--prompt-file --context ... [--return-mode patch_or_text|text] [--timeout N]`
- Make targets:
  - `make new-thought title="..." [spec="..."]` → `tgs/<BASE_HASH>-<kebab-title>/`
  - `make ears-gen` → regenerates ANTLR parser from `src/core/ears/ears.g4`
  - `make build|test|tidy` — developer convenience.
- Scripts:
  - `scripts/install.sh` — install `tgs` from GitHub releases (Homebrew alternative).
  - `bootstrap.sh` — project bootstrap via curl (optional).
- Files/artifacts:
  - `tgs.yaml` — policies and defaults (EARS, NFR enforcement, forbid paths).
  - Thought docs (`tgs/<hash>-<slug>/...`) — `research.md`, `plan.md`, `implementation.md`.
  - CI templates under `.github/workflows/` or `.gitlab-ci.yml` (optional).

## Data & Models (if AI/ML is used)
- No first-class ML models in-repo. The adapter may call external LLMs.
- EARS grammar: `src/core/ears/ears.g4` with generated parser under `src/core/ears/gen/...`.
- Configuration: `tgs.yaml` for policies, approver roles, agent order, branch prefix.

---

### Checklist
- [x] Context diagram shows external actors/systems  
- [x] Containers cover all runtime elements  
- [x] Major components identified with responsibilities  
- [x] Interfaces documented with data flow/protocol  
