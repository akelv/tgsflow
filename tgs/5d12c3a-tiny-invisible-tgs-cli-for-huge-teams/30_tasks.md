# Tasks: Tiny Go-based `tgs` CLI

> Rule of thumb: ≤ 1–2 days per task, ≤ 300 LOC per PR, include tests + docs snippets.
> Labels: `thought:2025-09-11-tgs-go-cli`, `tgs`, `milestone:M#`

---

## M0 — Skeleton & Bootstrap

### T0.1 — Project skeleton & CLI router
- **Scope:** Create `main.go`, basic subcommand router using stdlib `flag`; print version.
- **Estimate:** 0.5d
- **Acceptance:**
  - `tgs --version` prints semver and build info.
  - `tgs help` lists commands (stubs).
- **Deliverables:** `main.go`, `cmd/*` stubs, `internal/version`.

### T0.2 — `init --decorate`
- **Scope:** Create `tgs/` layout, optional CI templates (GH/GitLab), never touch `src/`.
- **Estimate:** 0.5d
- **Acceptance:**
  - Running twice is idempotent.
  - Adds sample `.github/workflows/tgs-approve.yml` when `--ci-template=github`.
- **Deliverables:** `cmd/init.go`, templates under `embed.FS`, tests.

### T0.3 — Config loader (`tgs.yaml`)
- **Scope:** Parse YAML, defaults, env overrides; validate keys.
- **Estimate:** 0.5d
- **Acceptance:** Missing file uses sensible defaults; invalid keys error clearly.
- **Deliverables:** `core/config/loader.go`, unit tests.

### T0.4 — Logging & JSONL option
- **Scope:** Human logs (default) + `--json` JSONL logs; levels.
- **Estimate:** 0.5d
- **Acceptance:** Same messages visible in both modes; stderr vs stdout conventions respected.
- **Deliverables:** `util/log`.

---

## M1 — Core Flow: Context & Thought Files

### T1.1 — `context` (brownfield scan)
- **Scope:** Heuristic scan (langs, modules, schemas), emit `.context.json` + seed `00_research.md`.
- **Estimate:** 1d
- **Acceptance:** On ≤5k files, completes <60s; includes language list & key manifests.
- **Deliverables:** `core/contextscan/*`, `cmd/context.go`, tests with sample repo.

### T1.2 — Thought file helpers
- **Scope:** Read/write helpers for `00/10/20/30/40/50/60`.
- **Estimate:** 0.5d
- **Acceptance:** Append-safe; preserves user edits.
- **Deliverables:** `core/thoughts/*`, unit tests.

### T1.3 — `specify` (Spec Kit proxy + fallback)
- **Scope:** If `specify` on PATH, shell out; else generate minimal sections.
- **Estimate:** 1d
- **Acceptance:** Produces `10_spec.md`; proxy passthrough covered by tests (mock cmd).
- **Deliverables:** `cmd/specify.go`, tests.

### T1.4 — `plan`
- **Scope:** Append technical plan to `20_plan.md`; enforce NFR placeholders.
- **Estimate:** 0.5d
- **Acceptance:** Fails if required NFR sections missing when policy enabled.
- **Deliverables:** `cmd/plan.go`, tests.

### T1.5 — `tasks`
- **Scope:** Generate/validate list; optional `--issues` to open Issues.
- **Estimate:** 0.5d
- **Acceptance:** Creates `30_tasks.md` if missing; validates IDs/format if present.
- **Deliverables:** `cmd/tasks.go`, tests.

### T1.6 — `approve` (governance gate)
- **Scope:** Validate presence of `10/20/30/40` + required roles in `40_approval.md`.
- **Estimate:** 0.5d
- **Acceptance:** `--ci` exits non-zero on failure with actionable message.
- **Deliverables:** `cmd/approve.go`, tests.

---

## M2 — Agents, Transport & Patching

### T2.1 — Protocol v1: HELLO/ACK + JSONL transport
- **Scope:** Define `proto/*` structs, HELLO/ACK, JSONL encoder/decoder.
- **Estimate:** 1d
- **Acceptance:** Golden tests: transcripts round-trip; stdout only protocol, stderr logs.
- **Deliverables:** `agent/proto/*`, tests.

### T2.2 — Adapter discovery & exec bridge
- **Scope:** Discover adapters in `.tgs/agents`, user config dir, PATH; exec with env.
- **Estimate:** 0.5d
- **Acceptance:** First found wins; missing adapter → clear guidance printed.
- **Deliverables:** `agent/execadp/*`, tests.

### T2.3 — Minimal unified diff applier
- **Scope:** Apply text hunks; CRLF/LF normalize; safety checks (path, LOC).
- **Estimate:** 1.5d
- **Acceptance:** Pass patch golden tests; reject forbidden paths & >LOC patches.
- **Deliverables:** `agent/patch/*`, tests.

### T2.4 — `implement` (default patch mode)
- **Scope:** Build CONTEXT/PROMPT; stream PATCH/NOTE; apply; run hooks; open PR.
- **Estimate:** 2d
- **Acceptance:** On dummy adapter, creates branch & PR with labels; hooks run; failures summarized.
- **Deliverables:** `cmd/implement.go`, `core/hooks/*`, `core/pr/*`, tests.

### T2.5 — `verify` (one-shot checks)
- **Scope:** Run fmt/lint/test/perf hooks; policy; drift; `--ci` mode.
- **Estimate:** 1d
- **Acceptance:** Clear pass/fail; machine-readable summary in `--json`.
- **Deliverables:** `cmd/verify.go`, tests.

### T2.6 — `open-pr` helper
- **Scope:** Create PR with labels, branch prefix, acceptance criteria link.
- **Estimate:** 0.5d
- **Acceptance:** PR body includes links to `10_spec.md` and selected task.
- **Deliverables:** `cmd/openpr.go`, tests (mock API).

---

## M3 — Governance, Drift & Docs

### T3.1 — `drift-check`
- **Scope:** Heuristics for routes/schemas/public API; suggest spec/plan patch.
- **Estimate:** 1d
- **Acceptance:** Flags added route in sample repo; prints actionable suggestion.
- **Deliverables:** `core/drift/*`, `cmd/driftcheck.go`, tests.

### T3.2 — `docs` + attestation
- **Scope:** Generate `60_docs.md` (user notes, runbook, metrics); attest JSON.
- **Estimate:** 0.5d
- **Acceptance:** Attestation includes approvers, task PRs, hashes, timestamps.
- **Deliverables:** `cmd/docs.go`, `core/thoughts/attest.go`, tests.

### T3.3 — Runner container & GH Actions templates
- **Scope:** Minimal Debian image with `tgs` + reference adapters; publish to GHCR.
- **Estimate:** 1d
- **Acceptance:** Example repo CI passes end-to-end (`approve`, `implement`, `verify`).
- **Deliverables:** `/runner/Dockerfile`, workflows under `/templates/gh/`, docs.

---

## M4 — IDE & Dev UX

### T4.1 — `brief`
- **Scope:** Emit tight task brief (md/text) from spec/plan/tasks + constraints.
- **Estimate:** 0.5d
- **Acceptance:** `tgs brief --task "..."` prints ≤200 lines; includes NFRs + forbidden paths.
- **Deliverables:** `cmd/brief.go`, tests.

### T4.2 — `watch`
- **Scope:** FS watcher; run selected hooks on quiet period; concise failure output.
- **Estimate:** 1d
- **Acceptance:** Detects changes in impacted packages; no runaway processes.
- **Deliverables:** `cmd/watch.go`, tests (fake FS events).

### T4.3 — Policy pack examples
- **Scope:** Ship `.tgs/policies/*.yaml` for NFRs, coverage, perf.
- **Estimate:** 0.5d
- **Acceptance:** `approve` enforces presence; `verify` enforces thresholds (soft → hard).
- **Deliverables:** policy files + docs.

---

## M5 — Optional Enhancements

### T5.1 — Alt transport (length-prefixed CBOR)
- **Scope:** `transport/lpcbor` + HELLO negotiation; behind build tag.
- **Estimate:** 1d
- **Acceptance:** Conformance tests pass with both `jsonl` and `lp-cbor`.
- **Deliverables:** `agent/transport/lpcbor/*`, tests.

### T5.2 — Agent-apply mode (opt-in)
- **Scope:** `--apply-by-agent`: snapshot, allow agent writes, `verify`, auto-revert violators.
- **Estimate:** 1.5d
- **Acceptance:** Violating write to `infra/prod/` is reverted and logged; PR still opens with safe subset.
- **Deliverables:** `cmd/implement.go` extensions, tests.

---

## Cross-cutting Tasks

### TX.1 — Error taxonomy & messages
- **Scope:** Centralize error codes/messages; human + machine forms.
- **Estimate:** 0.5d
- **Acceptance:** All commands use typed errors; CI parsers see stable codes.
- **Deliverables:** `util/errs`.

### TX.2 — Conformance suite for adapters
- **Scope:** Golden JSONL streams + expected outcomes; public docs.
- **Estimate:** 1d
- **Acceptance:** Reference adapters pass; failures produce readable diffs.
- **Deliverables:** `/conformance/*`, docs.

### TX.3 — Docs site (README + HOWTOs)
- **Scope:** Quickstart, autonomous CI, IDE workflow, policy packs, hooks.
- **Estimate:** 1d
- **Acceptance:** New user can bootstrap in <10 minutes.
- **Deliverables:** `docs/*`, updated README.

---

## Nice-to-Haves (later)

- **Monorepo impact analyzer** (git + lang-aware) to further scope tests.
- **Coverage delta gate per package**.
- **Perf smoke integration with k6/autocannon helpers.**
- **GitLab PR & self-hosted GH compatibility matrix.**

---

## Task Execution Notes
- One PR per task.
- Link PRs to this Thought and the task ID in the title (e.g., `M2 T2.4 implement (patch mode)`).
- Include:
  - brief rationale
  - test results snippet
  - acceptance checklist
- Keep commits small and linear; squash on merge.

