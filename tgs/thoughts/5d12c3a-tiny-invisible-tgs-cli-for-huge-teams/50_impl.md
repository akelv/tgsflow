# Implementation Details: Tiny Go-based `tgs` CLI (M0–M1)

- Date: 2025-09-11
- Base Hash: 5d12c3a
- Scope: Minimal single-binary CLI with core commands, governance checks, IDE helpers, and active Thought directory support.

---

## Build & Entry

- go.mod
  - Module: `github.com/kelvin/tgsflow`
  - Go: 1.22
  - Dependency: `gopkg.in/yaml.v3 v3.0.1` (for `tgs.yaml` parsing)

- src/main.go
  - Build-time vars: `version`, `commit`, `date` (defaults: dev/none/unknown). Printed by `--version` or `version` subcommand.
  - Global flags:
    - `--json`: switch logs to JSONL on stderr.
    - `--version`: print version and exit.
  - Subcommand router (stdlib `flag` + switch):
    - Implemented: `help`, `init`, `context`, `specify`, `plan`, `tasks`, `approve`, `verify`, `brief`, `version`.
    - Stubs (announce not implemented): `implement`, `drift-check`, `docs`, `open-pr`, `watch`.

---

## Logging

- src/util/logx/logx.go
  - Levels: DEBUG, INFO, WARN, ERROR. Minimum level defaults to INFO.
  - Output:
    - Human-readable: RFC3339 timestamp, level, message to stderr.
    - JSON mode (`SetJSON(true)`): JSON lines with `ts`, `level`, `msg`.
  - Concurrency: guarded writes via `sync.Mutex`.

---

## Configuration

- src/core/config/loader.go
  - Config:
    - `ApproverRoles []string`
    - `AgentOrder []string`
    - `BranchPrefix string`
    - Policies:
      - `ForbidPaths []string`
      - `MaxPatchLOC int`
      - `EnforceNFR bool`
  - Default(): roles `EM`, `TechLead`; agents `claude`, `gemini`; `branch_prefix: tgs/`; policies: forbid `infra/prod/`, `secrets/`; `max_patch_loc: 300`; `enforce_nfr: false`.
  - Load(repoRoot): reads `tgs.yaml` at repo root; returns defaults if missing.

---

## Thoughts Utilities

- src/core/thoughts/files.go
  - EnsureDir(path): `mkdir -p`.
  - EnsureFile(path, initial): create file if missing.
  - AppendSection(path, title, body): create or append a Markdown section.

- src/core/thoughts/locate.go
  - LocateActiveDir(repoRoot): resolve active Thought directory in priority:
    1) `TGS_THOUGHT_DIR` if points to an existing dir
    2) Most recently modified `tgs/<hash>-*` subdir
    3) Fallback: `tgs/` root
  - SpecFileCandidates(): returns `10_spec.md` and `10_specs.md`.

---

## Commands

- src/cmd/help.go
  - Prints usage and lists available commands.

- src/cmd/init.go
  - Flags: `--decorate` (default true), `--ci-template github|gitlab|none`.
  - Behavior:
    - Ensures `tgs/` exists; seeds common Thought files if missing:
      - `README.md`, `00_research.md`, `10_spec.md`, `20_plan.md`, `30_tasks.md`, `40_approval.md`.
    - CI templates:
      - `github`: writes `.github/workflows/tgs-approve.yml` (Go build + `./tgs approve --ci`).
      - `gitlab`: writes minimal `.gitlab-ci.yml` stub.
      - `none`: no CI files.
    - Idempotent: re-running preserves existing files.

- src/cmd/context.go
  - Heuristics scan: walks repo (skips `.git`, `node_modules`, `vendor`), counts files, detects languages by extension (`.go`, `.ts/.tsx`, `.js`, `.py`).
  - Writes `<active>/.context.json` with `generated_at`, `languages`, `files`.
  - Seeds `<active>/00_research.md` if missing.
  - Logs summary with the active Thought directory.

- src/cmd/specify.go
  - Flags: `--no-spec-kit` to disable proxying.
  - Behavior:
    - If `specify` binary is on PATH and not `--no-spec-kit`: proxies stdin/stdout/stderr.
    - Else: creates minimal `10_spec.md` in active Thought.

- src/cmd/plan.go
  - Appends a `Plan` section to `<active>/20_plan.md` with NFR placeholders (Performance/Security/Reliability).

- src/cmd/tasks.go
  - Flags: `--validate`.
  - Behavior:
    - If `<active>/30_tasks.md` missing: creates starter file with an example task.
    - `--validate`: scans for headings matching `^### T[0-9]+\.[0-9]+ — ` (counts IDs); fails if none found.

- src/cmd/approve.go
  - Flags: `--ci` (return non-zero on failure paths).
  - Required files (in active Thought):
    - Spec: either `10_spec.md` or `10_specs.md`
    - Plus: `20_plan.md`, `30_tasks.md`, `40_approval.md`
  - Role enforcement:
    - If `Config.ApproverRoles` present, each role name must appear in `40_approval.md`.
  - NFR enforcement:
    - If `Policies.EnforceNFR` true, `20_plan.md` must mention "Non-Functional".
  - On success: prints `approve: checks passed`.

- src/cmd/verify.go
  - Flags: `--ci`.
  - Executes repo-local hooks if present: `.tgs/hooks/{fmt,lint,test,perf}` via `os/exec`.
  - On hook failure: prints concise message; returns non-zero when `--ci` is used.

- src/cmd/brief.go
  - Flags: `--task "..."`, `--format md|text`.
  - Assembles a compact brief (≤200 lines):
    - Acceptance criteria (from spec, if present)
    - Non-Functional Requirements (from plan)
    - Tasks (optionally filtered by `--task` string)
    - Constraints (forbidden paths defaults: `infra/prod/`, `secrets/`)
  - Supports spec filename variants (`10_spec.md` or `10_specs.md`).

---

## Active Thought Directory Usage

- Commands resolving and operating within the active Thought directory: `context`, `specify`, `plan`, `tasks`, `approve`, `brief`.
- Resolution: environment override via `TGS_THOUGHT_DIR`, else latest `tgs/<hash>-*`, else `tgs/` root.

---

## Known Gaps / Next Milestones

- Not yet implemented (stubs only in router): `implement`, `drift-check`, `docs`, `open-pr`, `watch`.
- To be added in upcoming milestones (per `30_tasks.md`):
  - Agent protocol (HELLO/ACK, JSONL), adapter exec bridge, unified diff applier and safety gates.
  - Hook scoping (`verify --since <ref>`), drift heuristics, docs generation + attestations, PR creation.
  - CI runner container and workflow templates beyond approve.

---

## Build & Smoke Test

```bash
mkdir -p ./bin
go build -o ./bin/tgs ./src
./bin/tgs --version
./bin/tgs help
```

## Example Workflow

```bash
./bin/tgs init --decorate --ci-template github
./bin/tgs context
./bin/tgs specify
./bin/tgs plan
./bin/tgs tasks --validate
./bin/tgs approve --ci
./bin/tgs verify
./bin/tgs brief --task "T0.1"
```
