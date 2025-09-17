
# Plan: Tiny Go-based `tgs` CLI (single binary)

> Build a small, dependable CLI that makes SDD “invisible but governed” for both **Autonomous (EM)** and **Interactive (Senior Engineer + IDE)** modes.

---

## 1) Architecture (small, composable, stdlib-first)
  src/
  ├─ cmd/ # subcommand entrypoints (thin)
  │ ├─ init.go
  │ ├─ context.go
  │ ├─ specify.go
  │ ├─ plan.go
  │ ├─ tasks.go
  │ ├─ approve.go
  │ ├─ implement.go
  │ ├─ driftcheck.go
  │ ├─ docs.go
  │ ├─ verify.go
  │ ├─ brief.go
  │ └─ watch.go
  ├─ core/ # domain logic (pure, testable)
  │ ├─ thoughts/ # read/write 00/10/20/30/40/50/60 files
  │ ├─ contextscan/ # brownfield repo scan → .context.json
  │ ├─ hooks/ # fmt/lint/test/perf runners
  │ ├─ policy/ # YAML policies (forbid paths, budgets, coverage)
  │ ├─ drift/ # simple drift heuristics (routes/schemas)
  │ ├─ gitx/ # git apply, branch/PR metadata
  │ ├─ pr/ # PR creation (GH/GitLab)
  │ └─ config/ # tgs.yaml loader/validation
  ├─ agent/ # adapter orchestration (transport, apply)
  │ ├─ proto/ # JSONL v1 structs + HELLO/ACK
  │ ├─ execadp/ # exec adapters (PATH/.tgs/agents discovery)
  │ ├─ transport/ # jsonl (default), lp-cbor (optional)
  │ └─ patch/ # unified-diff apply (safety gate)
  ├─ util/ # logging, fs, timeouts, paths
  └─ main.go # CLI wiring (stdlib flag; no heavy deps)

**Design principles**
- **No framework lock-in**: stdlib `flag` + tiny router for subcommands.
- **Minimize deps**: prefer stdlib; use small, vetted libs only where it saves real time.
- **Pure packages**: `core/*` shouldn’t know about CLI/stdio; easy to test.
- **Replaceable agents**: `agent/execadp` shells to `tgs-agent-*` executables via stdin/out.

---

## 2) Key Data & Files

- **Thought files**
  - `00_research.md`, `10_spec.md`, `20_plan.md`, `30_tasks.md`, `40_approval.md`, `50_impl.md`, `60_docs.md`
- **Repo config**
  - `tgs.yaml` (approvers, agent order, limits, policies, PR prefix)
  - `.tgs/` (hooks, policies, adapters, testdata, optional compose.yml)
- **Context**
  - `tgs/.context.json` (brownfield map: langs, modules, schemas, CI)
- **Attestation**
  - `.tgs/attestations/<slug>-<ts>.json` (who/what/when/hashes)

---

## 3) Command Behaviors (concise contracts)

### 3.1 `tgs init --decorate [--ci-template github|gitlab|none]`
- Create `tgs/` skeleton + optional CI workflows (approval, verify).
- Never touch `src/` or build artifacts.
- Idempotent.

### 3.2 `tgs context`
- Walk repo (≤5k files default), detect languages (via extensions), module hints (package.json, go.mod, pom.xml), DB artifacts (migrations, schema files).
- Emit `00_research.md` (questions + auto-filled facts) and `.context.json`.
- Timeboxed (default 60s), skip large dirs via `.tgsignore`.

### 3.3 `tgs specify|plan|tasks`
- If Spec Kit CLI (`specify`) is present and `--no-spec-kit` not set → proxy.
- Else, generate minimal but structured sections using internal templates.
- Append to existing files; never nuke edits.

### 3.4 `tgs approve [--ci]`
- Validate presence of `10/20/30/40`.
- Validate required roles in `40_approval.md` (configurable).
- Enforce policies present in `tgs.yaml` (e.g., NFR sections in `20_plan.md`).
- Exit non-zero in CI on failure.

### 3.5 `tgs implement [--task|--task-id] [--apply-by-agent] [--open-pr]`
- Select a task from `30_tasks.md` (by ID or text).
- Build **CONTEXT** & **PROMPT**, launch preferred adapter via `execadp`.
- **Default:** receive `PATCH` messages → `patch.Apply()` with preflight guards:
  - max LOC, forbidden paths, no binary writes, path whitelist (optional).
- Run hooks (`fmt/lint/test/perf`) with timeouts; summarize failures.
- On pass, create **1 PR per task** with labels + acceptance criteria link.

### 3.6 `tgs verify [--since <ref>]`
- One-shot: run hooks (scoped to changed paths), policy checks, drift check.
- Employed locally (pre-PR) and in CI (`--ci`).

### 3.7 `tgs brief [--task "..."] [--format md|text]`
- Emit tight prompt for IDE agents (Cursor): acceptance criteria, NFRs, forbidden paths, deliverables, ≤300 LOC, tests required.

### 3.8 `tgs drift-check [--base <ref>] [--head <ref>]`
- Compare changes vs `10_spec.md`/`20_plan.md`.
- Heuristics: new routes (regex), changed schemas (OpenAPI/GraphQL), added public methods.
- Output actionable suggestions (patch snippet for the spec/plan).

### 3.9 `tgs docs`
- Pull metrics links (if present), create `60_docs.md` (user notes, runbook, dashboards).
- Emit attestation JSON (approvers, task PRs, hashes).

---

## 4) Agent Protocol v1 (JSONL baseline; CBOR ready)

- **HELLO/ACK negotiation**:
  - `HELLO{ schema:"tgs.proto", version:"1.0", encodings:["jsonl","lp-cbor"], compression:["none","gzip"] }`
  - `HELLO_ACK{ accept_encoding:"jsonl", accept_compression:"none" }`
- Message types:
  - `CONTEXT`, `PROMPT` → adapter
  - `PATCH`, `NOTE`, `QUESTION`, `COMPLETE`, `ERROR` ← adapter
  - (optional) `PING/PONG`, `CHUNK_*`
- **Stdout** = protocol; **stderr** = logs
- Size caps: `max_msg_bytes` (default 5MB); chunk if exceeded.
- Safety: reject patches outside repo root; normalize LF; block forbidden paths early.

**Transport package**
- `jsonl` (default)
- `lpcbor` (optional, behind build tag) for future streaming performance

---

## 5) Patching strategy (safety first)

- Prefer **internal unified diff applier**:
  - Minimal implementation to apply hunks (add/modify/delete), CRLF/LF normalization.
  - Verify against base hashes; fail fast on fuzz/offset mismatch.
- **Fallback**: if `git` available → `git apply --3way --whitespace=fix` for complex hunks.
- Preflight:
  - Check LOC cap, path allow/deny rules, binary detection (heuristic by MIME/ext).
- On fail:
  - Write `.rej` files under `.tgs/.rej/`, send summarized failure back to adapter via follow-up `PROMPT` (tight loop).

---

## 6) Hooks & Policies

**Hooks (repo-local, executable)**
- `.tgs/hooks/fmt`, `.tgs/hooks/lint`, `.tgs/hooks/test`, `.tgs/hooks/perf`
- Monorepo aware: pass `--since` or changed file list to scope.
- Compose support: if `.tgs/compose.yml` exists, spin up services for integration tests; teardown after.

**Policies (tgs.yaml / .tgs/policies/*.yaml)**
- `forbid_paths`, `max_patch_loc`, `min_coverage`, `perf_budgets`
- `approve` enforces presence of NFR sections
- `verify` & CI enforce coverage/perf (optional soft-fail thresholds)

---

## 7) CI/CD Runner

**Image**: `ghcr.io/akelv/tgs-runner:stable`
- Includes:
  - `/usr/local/bin/tgs` (static)
  - reference adapters: `tgs-agent-claude`, `tgs-agent-gemini`
  - minimal toolchains commonly used by hooks (curl, git, jq; optionally go/node if needed)
- **Workflows**:
  - `tgs-propose` (from Issue → Thought PR)
  - `tgs-approve` (gate on 10/20/30/40 + policy)
  - `tgs-implement` (per-task PRs)
  - `tgs-verify` (hooks + policy + drift on task PRs)
  - `tgs-docs` (post-merge docs + attestation)

Secrets for agents injected at job runtime.

---

## 8) Error taxonomy (stable, user-friendly)

- Parsing/transport: `BAD_JSON`, `UNSUPPORTED_ENCODING`, `BAD_ORDER`
- Adapter: `AGENT_TIMEOUT`, `AGENT_ERROR`
- Patch: `POLICY_VIOLATION`, `PATCH_APPLY_FAILED`
- Hooks: `HOOKS_FAILED`, `COVERAGE_BELOW_MIN`, `PERF_BUDGET_EXCEEDED`
- Governance: `APPROVAL_MISSING`, `ROLE_MISSING`

All errors return **human summary** + **machine fields** for CI.

---

## 9) Performance & Reliability

- **Timeouts** per step (adapter, hooks). Defaults in `tgs.yaml`.
- **Parallelism**:
  - Local: serial by default (predictability); CLI flag `--parallel` for multiple tasks if safe.
  - CI: let matrix jobs fan out tasks (via workflow strategy).
- **Caching**:
  - Hook-level caching permitted (e.g., `GOMODCACHE`, `~/.cache/pnpm`) on runners.
- **Large diffs**:
  - Gzip compress if `> max_msg_bytes/2`; or chunk.
  - Offer `PATCH_REF` (temp file path) only when both processes share filesystem (CI).

---

## 10) Security

- Never print secrets (agent keys) to stdout/stderr.
- Validate adapter paths (no traversal, no scripts from untrusted dirs).
- Forbid touching `infra/prod/`, `secrets/` by default.
- Provide binary **signing** (cosign) + **SBOM** in releases.

---

## 11) Minimal External Dependencies

- **Must**: `git` (for branch/PR metadata; optional for apply fallback)
- **Optional**: Spec Kit CLI (`specify`) for proxy
- **Libraries** (keep tiny):
  - YAML (`gopkg.in/yaml.v3`)
  - (optional) CBOR (`fxamacker/cbor`) behind build tag
  - (optional) go-diff/patch lib if it saves substantial time; else implement minimal applier

---

## 12) Testing Strategy

- **Unit**: core packages (`thoughts`, `contextscan`, `policy`, `patch`, `transport`), no I/O.
- **Golden tests**: JSONL transcripts for agent protocol; patch apply on edge hunks.
- **Integration**: end-to-end `implement` using the **local dummy adapter**:
  - adapter reads CONTEXT/PROMPT and returns canned PATCH sequences (fixtures).
- **CLI smoke**: `tgs verify` on a sample repo in CI.
- **Conformance**: publish adapter conformance suite (transcripts + expected outcomes).

---

## 13) Rollout Milestones & Tasks

### M0 — Skeleton
- `main.go` + router; `tgs init --decorate`
- Config loader (`tgs.yaml`) with defaults
- Logger (human | `--json`)

### M1 — Core Flow
- `context`, `specify` (proxy + minimal), `plan`, `tasks`
- `approve` (roles + policy presence)
- Thought file read/write helpers

### M2 — Agents & Patching
- Agent protocol (JSONL) + HELLO/ACK
- Exec adapter discovery; STDIO bridger
- `implement` (default patch mode) + hooks runner
- Minimal diff applier + `git apply` fallback

### M3 — Governance & CI
- `drift-check` (v0 heuristics)
- `docs` + attestation
- Runner image + GitHub Actions templates
- `verify` (one-shot) + `open-pr`

### M4 — IDE Helpers
- `brief` (md|text), `watch` (local incremental hooks)
- Polished error messages; concise failure summaries

### M5 — Polish & Optional CBOR
- Transport `lp-cbor` (build tag)
- Policy hardening (coverage/perf gates)
- Adaptation to GitLab PRs

---

## 14) Risks & Mitigations

- **Patch applier edge cases** → Provide `git apply` fallback; log `.rej` with clear fixes.
- **Drift false positives** → Feature flag + per-repo tuning; exclusions.
- **Hook flakiness** → retries & quarantine mode; document best practices.
- **Adapter fragmentation** → ship conformance tests + two reference adapters.

---

## 15) Acceptance Gate for this Plan

- ✅ API shape for each command defined & simple.
- ✅ Agent protocol v1 (JSONL) fixed with HELLO/ACK.
- ✅ Default patch mode + `--apply-by-agent` path verified by `verify`.
- ✅ CI runner & workflows scripted.
- ✅ IDE helpers (`brief`, `verify`, `watch`) specified.

> After M2, we can dogfood on a real brownfield repo (add a Search feature) to validate DX end-to-end before M3.

---
