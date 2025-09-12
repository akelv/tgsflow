# Research: Tiny Go-based `tgs` CLI — Autonomous & Interactive Modes

- Date: 2025-09-11
- Base Hash: 5d12c3a
- Participants: Agent/Human

## Problem
Teams want to adopt Spec-driven Development (SDD) on brownfield repos with **high leverage and strong governance**:
- **Engineering Manager (EM)**: prefers **full agentic delegation** via CI (GitHub Actions) while retaining **approval gates** and policy enforcement.
- **Senior Engineer**: codes **interactively in Cursor/VS Code** with GPT-like agents; needs **lightweight guardrails** (spec/plan/tasks, verify, drift) without losing IDE speed.

Current shell-based bootstraps add friction (runtime deps, cross-OS quirks) and don’t provide a clean “autonomous pipeline” nor an IDE-friendly “brief/verify” loop out of the box.

**We need a single static Go binary** that is invisible but governs:
- Creates/manages Thoughts (Research → Spec → Plan → Tasks → Approval → Implement → Docs)
- Proxies to code agents (Claude Code, Gemini CLI, local LLM) via a stable adapter protocol
- Enforces approval & policy gates both **locally** and **in CI**
- Works in **two modes**: fully autonomous (EM) and interactive IDE (Senior Engineer)

---

## Goals
1. **Single-binary DX**: No external runtime deps. macOS/Linux/Windows. Simple curl/install.
2. **Brownfield-first**: `tgs init --decorate`, `tgs context` to map legacy repos.
3. **Spec-driven core**: First-class `specify`, `plan`, `tasks`; Spec Kit proxy if present.
4. **Governance**: `approve` gate, policy packs (NFRs, PII, forbidden paths), attestations.
5. **Agent-agnostic**: Exec-adapter protocol over stdin/stdout (JSONL v1), with transport negotiation (lp-cbor ready).
6. **Autonomous EM workflow**: Workflows & runner container to create Thoughts from an Issue, gate, implement tasks into small PRs, and produce docs/attestations.
7. **Interactive IDE workflow**: `brief` → paste into IDE chat; `verify`/`watch` for fast feedback; `open-pr` for neat, labeled PRs.
8. **Safety by default**: default **patch mode** (agents return diffs, `tgs` applies), optional **agent-apply** with strict verification.

---

## Non-Goals
- IDE plugins or UI; we integrate via CLI + prompts.
- Agent vendor lock-in; adapters are replaceable executables.
- Orchestrating deployments; we focus on pre-merge governance and PR hygiene (optional feature-flag rollout docs).

---

## Users & Scenarios

### A) EM — **Fully Agentic Delegation**
- **Intent**: Opens an Issue labeled `tgs:proposal` with outcomes/NFRs.
- **Pipeline**:
  1) `tgs-propose` workflow: runs `context`, `specify`, `plan`, `tasks`; opens Thought PR.
  2) EM approves (`40_approval.md`); `tgs approve --ci` blocks without it.
  3) On merge, `tgs-implement` workflow runs `implement` to open **small PRs per task**, with hooks & policy checks on each PR (`tgs verify`).
  4) Docs & attestation emitted post-merge (`tgs docs`).
- **Why it works**: High leverage with governance and auditability; vendor-agnostic agents.

### B) Senior Engineer — **Interactive IDE (Cursor/GPT-5)**
- **Flow**:
  1) Create Thought locally (`context` → `specify` → `plan` → `tasks` → `approve`).
  2) `tgs brief --task "…"`, paste into IDE chat to get focused agent help.
  3) Code/edit in IDE; run `tgs verify` (fmt/lint/test/perf + drift/policy) before PR.
  4) `tgs open-pr` to create neat PR with labels and acceptance criteria links.
- **Why it works**: Preserves speed while aligning to the same spec/plan/tasks and gates.

---

## Requirements

### Functional
- **F1**: Thought lifecycle commands: `context`, `specify`, `plan`, `tasks`, `approve`, `implement`, `drift-check`, `docs`.
- **F2**: **Autonomous mode** via CI:
  - Create Thought PR from Issue (proposal workflow).
  - Block merges without `40_approval.md` (approval workflow).
  - Implement tasks into small PRs with hooks & policy (implement + verify workflows).
  - Produce docs & attestation.
- **F3**: **Interactive mode**:
  - `brief` emits a compact, copy-paste brief for IDE agents.
  - `verify` runs all local checks in one shot (fmt/lint/test/perf + policy + drift).
  - `watch` supports incremental local feedback.
  - `open-pr` creates labeled PRs.
- **F4**: **Adapter protocol** (v1):
  - JSONL over stdin/stdout with HELLO/ACK (schema/version/encoding).
  - Message types: CONTEXT, PROMPT, PATCH, NOTE, QUESTION, COMPLETE, ERROR (+ optional PING/PONG, CHUNK_*).
  - **Default**: `jsonl`; **Allowed**: `lp-cbor` (configurable).
  - **Default apply**: `patch` (tgs applies diffs); **opt-in**: `--apply-by-agent` (verify after).
- **F5**: **Hooks & Policies**:
  - Hooks: `.tgs/hooks/{fmt,lint,test,perf}`; monorepo-aware.
  - Policies: `min_coverage`, `perf_budgets`, `forbid_paths`, `max_patch_loc`, `max_msg_bytes`.
  - `verify` and CI run these consistently.

### Non-Functional
- **NFR1**: Cross-platform single binary (<20–30MB typical).
- **NFR2**: CI runner image with signed binary; SBOM + cosign attestations for releases.
- **NFR3**: Performance: `verify` on medium repos completes <2–4 min (cache hooks where possible).
- **NFR4**: Security: secrets only via CI secrets; stdout = protocol only for adapters (logs in stderr).
- **NFR5**: Reliability: Timeouts & retries for agent calls; graceful degradation if adapters missing.

---

## Constraints
- Go ≥1.22, Git available in PATH.
- GitHub/GitLab PR APIs supported; GHE compatibility via env overrides.
- Spec Kit CLI proxied when present; otherwise built-in minimal prompts.

---

## Key Design Decisions
- **Default diff-apply in `tgs`** for safety & audit; agent-apply behind flag.
- **Transport negotiation**: JSONL baseline; length-prefixed CBOR ready for large/fast streams.
- **Tight IDE bridge**: `brief`, `verify`, `watch`, `open-pr`—no plugin needed.
- **Small PRs by design**: `max_patch_loc` & per-task PRs reduce review cost and blast radius.

---

## Risks & Mitigations
- **Over-automation in autonomous mode** → Low-quality specs/plans  
  *Mitigation*: Approval gate; policy enforcing NFR sections; EM edits Thought PR.
- **Adapter fragmentation**  
  *Mitigation*: Publish conformance tests & reference adapters (Claude/Gemini).
- **Drift false positives**  
  *Mitigation*: Tunable rules, route/schema awareness, path exclusions.
- **Developer resistance (“too many gates”)**  
  *Mitigation*: Fast `verify`, `watch` for live feedback, policy presets per repo maturity.

---

## Acceptance Criteria
- **AC1**: From an Issue with label `tgs:proposal`, CI opens a Thought PR containing `00/10/20/30` (filled), passing `tgs approve --ci` after reviewers add `40_approval.md`.
- **AC2**: On merging the Thought PR, CI opens **separate task PRs** via `tgs implement`; each PR passes `tgs verify` (hooks + policy + drift).
- **AC3**: Interactive dev can run `tgs brief`, code in Cursor, `tgs verify`, and `tgs open-pr` with no additional tooling.
- **AC4**: Adapters can be swapped (Claude ↔ Gemini) with no changes to `tgs`.
- **AC5**: Default patch mode forbids writes to `infra/prod/` & `secrets/`; violations blocked pre-apply.
- **AC6**: `docs` produces `60_docs.md` with metrics/runbook links; attestation JSON written.

---

## References
- TGSFlow methodology & approval gates
- Spec-driven Development (Spec Kit) concepts
- Internal adapter protocol v1 (HELLO/ACK, JSONL; lp-cbor ready)
