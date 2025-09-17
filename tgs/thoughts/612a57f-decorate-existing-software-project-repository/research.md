# Research: Decorate Existing Software Project Repository

- Date: 2025-09-11
- Base Hash: TBD at PR time
- Participants: AI Agent (proposer), Human Maintainer (approver)

## 1. Problem Statement
Engineers with an existing repository want to adopt the TGSFlow workflow without scaffolding a new project or bringing in all templates and repo-level docs from `tgsflow`. The current `bootstrap.sh` always clones the full repo into a new subdirectory, applies a template, copies the root `README.md`, and initializes a new git repository. This behavior is unsuitable for decorating an existing codebase where only the minimal TGS workflow files are desired and no subdirectory should be created.

Desired outcome: a safe "decorate" mode that injects the essential TGS workflow into the current repository without creating a new directory, avoiding template payloads by default, not overwriting project `README.md`, and providing a dry-run for visibility.

## 2. Current State
- `bootstrap.sh` (v1.1.1) flow:
  - Prompts for template and project name.
  - Clones `tgsflow` to a temp dir, then copies everything into `<project_name>/`.
  - Applies selected template from `templates/<type>` into that directory.
  - Updates `README.md` contents, initializes a fresh git repo and commits.
- Minimal TGS workflow components today:
  - `agentops/AGENTOPS.md` (system prompt + workflow guide)
  - `tgs/README.md` (thoughts directory guide)
  - `agentops/tgs/*` (file templates for `make new-thought`)
  - `Makefile` provides `new-thought` target relying on `agentops/tgs/*` existing

Issues for existing repos:
- Creates an extra subdirectory; undesired.
- Pulls in templates and TGSFlow repo `README.md` by default; noisy/irrelevant.
- Initializes git history in a subfolder; conflicts with existing repos.

## 3. Constraints & Assumptions
- Must be safe and non-destructive by default; do not overwrite existing files unless `--force` is specified.
- No additional dependencies beyond `git` and `curl`.
- Should be non-interactive for decorate mode; no prompts for project type/name.
- Idempotent: running decorate multiple times should be safe.
- Cross-platform shell considerations: POSIX-friendly bash; avoid GNU-only utilities where feasible.

## 4. Risks & Impact
- Overwriting `README.md` or `Makefile` in existing repos (mitigated by default skip and `--force`).
- Partial adoption if essential files are omitted (ensure we copy the minimal set reliably).
- Confusion around templates; default to none and require explicit opt-in.

## 5. Alternatives Considered
1) Separate script (e.g., `decorate.sh`): clearer separation but duplicates logic and increases maintenance.
2) Separate branch/tarball containing only minimal TGS files: simple distribution, but diverges from main and risks drift.
3) Enhance `bootstrap.sh` with `--decorate` mode: single entry-point, consistent UX, minimal maintenance.

Chosen: Enhance `bootstrap.sh` with `--decorate`.

## 6. Recommendation
Add `--decorate` mode to `bootstrap.sh` with these properties:
- Operates in the current directory; no project directory creation; no git init.
- Copies only essential files by default:
  - `agentops/AGENTOPS.md`
  - `tgs/README.md`
  - `agentops/tgs/*`
  - A small `tgs.mk` make include containing the `new-thought` target. If a `Makefile` exists, append `include tgs.mk` if absent; otherwise, create a minimal `Makefile` that includes `tgs.mk`.
- Never copy repo-level `README.md` from `tgsflow` and do not copy `templates/` unless explicitly requested.
- Flags:
  - `--decorate`: enable decorate mode
  - `--dry-run`: print planned operations without writing
  - `--force`: allow overwrites when necessary
  - `--with-templates=<react|python|go|cli|none>`: optional; default `none`
- Be idempotent and produce clear logs.

## 7. References & Links
- Current `bootstrap.sh` behavior and structure
- `agentops/AGENTOPS.md`, `tgs/README.md`, `agentops/tgs/*`

---
Approval checkpoint: Please review this research and reply one of:
- APPROVE research
- REQUEST CHANGES: <notes>
