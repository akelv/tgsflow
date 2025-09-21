## What & Why
Implemented new configuration loader to mirror `tgs/tgs.yml` schema, enabling `core/brain` to access AI settings and adding backward compatibility for legacy `tgs.yaml` EARS gating used in verify tests.

## File Changes
- `src/core/config/loader.go`: Replaced legacy schema with new `Config` struct containing `AI`, `Triggers`, `Guardrails`, `Agents`, `Steps`, `Telemetry`, `Context`; added legacy `Policies` for compatibility; default values; loader reads `tgs/tgs.yml` with fallback to `tgs.yaml`.
- `src/core/config/loader_test.go`: Added table-driven tests for defaults, minimal YAML, overrides, invalid YAML.
- `src/core/brain/brain.go`: Fixed import path to module `github.com/kelvin/tgsflow/src/core/config`.
- `src/core/brain/transport_stub.go`: Added stub transports to satisfy interface during tests.
- `src/cmd/verify.go`: Restored EARS gating via `cfg.Policies.EARS.Enable`; compile fix for unused variable.

## Commands
- `go test ./...`

## How to Test
1. Ensure `tgs/tgs.yml` present; run `go test ./...` — all tests pass.
2. Remove config file or point `--repo` to empty dir; loader returns defaults.
3. Modify AI fields in a temp `tgs/tgs.yml`; confirm overrides via unit tests.

## Rollback
Revert changes to the files listed above.

## Follow-ups
- Implement real transports for `shell`, `proxy`, `sdk`, `mcp`.
- Consider a new lint section replacing legacy `Policies`.