## Problem
The current `src/core/config/loader.go` defines a legacy `Config` schema (`approver_roles`, `policies`, etc.) and reads `tgs.yaml`. The repository now provides a richer canonical config in `tgs/tgs.yml` with sections: `ai`, `triggers`, `guardrails`, `agents`, `steps`, `telemetry`, `context`. The `core/brain/brain.go` expects `cfg.AI` with `Mode`, `Toolpack.Budgets`, etc. We need to implement a new config loader that:

- Maps the new YAML structure into Go structs under `config`.
- Loads from `tgs/tgs.yml` (note: folder and extension differ from legacy `tgs.yaml`).
- Provides sensible defaults when fields are omitted.
- Is stable for future extension and referenced by `brain.NewTransport` and `Budget`.
- Includes table-driven tests to validate decoding, defaults, and error cases.

## Current State
- `src/core/config/loader.go` defines unrelated fields and reads `tgs.yaml` from repo root. Defaults and types do not match new schema.
- `src/core/brain/brain.go` references `cfg.AI.Mode` and `cfg.AI.Toolpack.Budgets` and switches transport by mode.
- A modern YAML example is at `tgs/tgs.yml`.

## Constraints
- Backwards compatibility with nonexistent file: loader should return defaults when file missing.
- Non-destructive behavior; no writes.
- Keep secrets in env; `api_key_env` field remains a string env var name.
- Keep production code outside `tgs/`; only read from it.

## Risks & Security
- Mis-parsing can route AI traffic incorrectly (wrong model/endpoint).
- Leaking secrets: ensure we only read env var names, not resolve or log secret values.
- Unexpected YAML values should not crash; validation should produce clear errors where necessary (e.g., unsupported `ai.mode`). Transport creation already validates mode.

## Alternatives Considered
- Using Viper: flexible but adds heavier dependency and implicit env resolution; plain `yaml.v3` suffices and is already used.
- Keeping legacy fields: increases confusion; instead, migrate `Config` to the new schema and keep loader backward-tolerant by ignoring unknown YAML keys.

## Recommendation
Implement a new `Config` struct modeling `tgs/tgs.yml` with nested types:
- `AI` with `Mode`, `Provider`, `Model`, `Endpoint`, `APIKeyEnv`, `TimeoutMS`, `Retry{MaxAttempts, BackoffMS}`, and `Toolpack{Enabled, AllowFor, Budgets, Routes, Tools[{Name,Desc}], Redaction{RedactEnvKeys, RedactPatterns}}`.
- `Triggers{IssueLabels, PRLabels}`.
- `Guardrails{AllowPaths, DenyPaths, MaxDiffLines, RequiredChecks, PRTemplate, CommitConvention}`.
- `Agents[]` entries with `Name, Type, Enabled, Capabilities, Selector{LabelsAny, PathsAny}, Runtime{Bin, Args, Provider, Events{OnPR, OnCommentCommands}, Endpoint, AuthEnv}` – fields optional.
- `Steps{PlanPrompt, ImplPrompt, ReviewPrompt}`.
- `Telemetry{LogDir, RedactRules}`.
- `Context{PackDir, ThoughtsDir, IncludeGlobs[]}`.

Loader behavior:
- `Load(repoRoot string)` reads `<repoRoot>/tgs/tgs.yml`. If missing, return `Default()`.
- `Default()` sets safe defaults mirroring comments in YAML example.
- YAML unmarshal into the struct; return struct and error.

## References
- Example config: `tgs/tgs.yml`
- Brain usage: `src/core/brain/brain.go`
- YAML lib: `gopkg.in/yaml.v3`

# Research: <Short Title>

- Date: <YYYY-MM-DD>
- Base Hash: <git rev-parse --short HEAD>
- Participants: <Agent/Human>

## 1. Problem Statement
<Clear description of the task and desired outcomes.>

## 2. Current State
<What exists today? Code, tools, versions, constraints.>

## 3. Constraints & Assumptions
<Security, performance, platform, dependencies, compliance, SLAs.>

## 4. Risks & Impact
<Security/privacy, reliability, regressions, scope creep, rollout risk.>

## 5. Alternatives Considered
<Option A, B, C with pros/cons.>

## 6. Recommendation
<Preferred approach and rationale.>

## 7. References & Links
<Docs, tickets, PRs, relevant code paths.>

---
Approval checkpoint: Please review this research and reply one of:
- APPROVE research
- REQUEST CHANGES: <notes>
