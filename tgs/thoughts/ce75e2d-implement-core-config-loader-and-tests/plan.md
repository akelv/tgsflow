## Objectives
- Implement `config.Config` to mirror `tgs/tgs.yml` schema and defaults.
- Update loader to read `tgs/tgs.yml` and return defaults if missing.
- Add unit tests with table-driven cases for decoding, defaults, and partial configs.
- Keep `brain.NewTransport` working against `cfg.AI` without further change.

## Scope
- In: `src/core/config/loader.go` types and loader implementation; new tests in `src/core/config/loader_test.go`.
- Out: Implementations of AI transports; runtime env resolution; CLI flags.

## Acceptance Criteria
- Tests pass: decoding example YAML matches fields; default values applied when YAML omits fields; missing file returns defaults; invalid YAML yields error.
- `brain.NewTransport` compiles and uses `cfg.AI.Mode` and budgets via `Budget` helper.

## Tasks
1. Define new structs under `config`: `Config`, `AI`, `Retry`, `Toolpack`, `Tool`, `Redaction`, `Triggers`, `Guardrails`, `Agent`, `AgentSelector`, `AgentRuntime`, `AgentRuntimeEvents`, `Steps`, `Telemetry`, `Context`.
2. Implement `Default()` with safe values (mode shell, provider openai, model gpt-4o-mini, timeouts, retries; toolpack disabled budgets set empty, etc.).
3. Implement `Load(repoRoot string)` to read `<repoRoot>/tgs/tgs.yml` and unmarshal.
4. Add tests `loader_test.go` with table cases:
   - missing file → defaults
   - minimal YAML (version+project only) → defaults for others
   - full example YAML → values match
   - partial ai/toolpack budgets and routes → maps populated
   - invalid YAML → error
5. Add small helper to read fixture YAML from temp dirs in tests.

## File-by-file Changes
- `src/core/config/loader.go`: Replace legacy structs with new schema; update loader path; keep package name.
- `src/core/brain/brain.go`: No code change expected; ensure it compiles.
- `src/core/config/loader_test.go`: New test file.
- `tgs/tgs.yml`: Read-only; used as reference for constructing test fixtures.

## Test Plan
- `go test ./src/core/config -run TestLoad` with table cases.
- CI: `make test` runs package tests.

## Rollout / Rollback
- Rollout: Merge alongside brain features; no CLI changes.
- Rollback: Revert changes to `loader.go` and tests.

## Estimates
- Implementation: 1-2 hours
- Tests: 1 hour

# Plan: <Short Title>

## 1. Objectives
<Bulleted outcomes; what success looks like.>

## 2. Scope / Non-goals
<In-scope items>
<Out-of-scope items>

## 3. Acceptance Criteria
<Observable criteria; include tests and docs expectations.>

## 4. Phases & Tasks
- Phase 1: <title>
  - [ ] Task 1
  - [ ] Task 2
- Phase 2: <title>
  - [ ] Task 1

## 5. File/Module Changes
<List files to add/edit/remove and high-level edits expected.>

## 6. Test Plan
<Unit/integration/e2e, manual checks, data migration validation.>

## 7. Rollout & Rollback
<Deployment steps, feature flags, rollback steps.>

## 8. Estimates & Risks
<Time estimates, risk mitigation.>

---
Approval checkpoint: Please review this plan and reply one of:
- APPROVE plan
- REQUEST CHANGES: <notes>
