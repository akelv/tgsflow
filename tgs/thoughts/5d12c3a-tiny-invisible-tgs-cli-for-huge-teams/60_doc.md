# Implementation Summary: Tiny Go-based `tgs` CLI (M0–M1 skeleton)

- Date: 2025-09-11
- Base Hash: 5d12c3a
- Scope: Minimal single-binary CLI with core commands, governance checks, IDE helpers, and active Thought directory support.

---

### What & Why
- Implements a lightweight Go CLI to standardize Spec-driven Development flows for autonomous CI and interactive IDE workflows.
- Emphasizes safety and governance; current milestone delivers structure, gating, and helper utilities.

### Key Capabilities
- **Commands**: `help`, `version`, `init`, `context`, `specify`, `plan`, `tasks`, `approve`, `verify`, `brief`.
- **Active Thought**: Auto-detect most recent `tgs/<hash>-<slug>/` or use `TGS_THOUGHT_DIR`.
- **Governance**: Approval gate validates required files and roles; optional NFR enforcement.
- **IDE helper**: `brief` emits a compact task brief (≤200 lines) with ACs, NFRs, constraints.
- **Hooks**: `verify` runs repo-local `.tgs/hooks/{fmt,lint,test,perf}` if present.
- **Logging**: Human-readable by default; `--json` switches to JSONL on stderr.

### Notable Files
- `src/main.go` (router, global flags, version)
- `src/util/logx/logx.go` (leveled logging, JSONL)
- `src/core/config/loader.go` (`tgs.yaml` loader, defaults, policies)
- `src/core/thoughts/files.go` (ensure dir/file helpers)
- `src/core/thoughts/locate.go` (active Thought resolution; `10_spec(s).md` support)
- `src/cmd/*` command implementations for the features above

### Command Behaviors (current)
- **init**: Idempotently creates `tgs/` skeleton; optional CI workflow (GitHub Actions) when requested.
- **context**: Heuristic scan (langs, file count); writes `<active>/.context.json`; seeds `00_research.md`.
- **specify**: Proxies to Spec Kit if available; fallback creates minimal `10_spec.md` in active Thought.
- **plan**: Appends `Plan` with NFR placeholders to `20_plan.md`.
- **tasks**: Creates `30_tasks.md` or validates task headings (`### T1.2 — Title`).
- **approve**: Requires `10_spec.md` or `10_specs.md`, plus `20_plan.md`, `30_tasks.md`, `40_approval.md`; validates roles and optional NFR presence.
- **verify**: Executes `.tgs/hooks/{fmt,lint,test,perf}` if present; non-zero exit on failures in CI mode recommended.
- **brief**: Emits a concise brief (≤200 lines) from spec/plan/tasks; supports `--task` filter and `--format`.

### Build & Smoke Test
- Build: `go build -o ./bin/tgs ./src`
- Check: `./bin/tgs --version` and `./bin/tgs help`

### Example Workflow
- `tgs init --decorate`
- `tgs context`
- `tgs specify`
- `tgs plan`
- `tgs tasks --validate`
- `tgs approve --ci`
- `tgs verify`
- `tgs brief --task "T0.1"`

### Configuration (`tgs.yaml`)
- **Keys**: `approver_roles`, `agent_order`, `branch_prefix`, `policies{forbid_paths,max_patch_loc,enforce_nfr}`.
- **Defaults**: forbid `infra/prod/`, `secrets/`; `max_patch_loc: 300`; `branch_prefix: tgs/`.

### Limitations / Next Steps
- Not yet implemented: `implement`, `open-pr`, `drift-check`, `docs`, `watch` (currently stubs or TBD).
- Upcoming: diff applier, adapter protocol (HELLO/ACK JSONL), PR creation, CI runner image, scoped verify (`--since`).
