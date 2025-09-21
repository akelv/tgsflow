# Implementation Summary: Verify EARS Subcommand for Design Docs

## 1. Overview (What & Why)
Adds `tgs verify ears` subcommand to lint EARS-style requirements in design docs by default (`tgs/design/10_needs.md`, `tgs/design/20_requirements.md`). Improves signal by:
- Filtering non-EARS lines and code fences
- Special handling: needs require explicit " shall"; requirements are linted even without "shall" to surface missing-shall
- Suppressing noisy ANTLR console errors and emitting only `path:line: message`
- Printing per-file and overall summary counts (captured/valid/invalid)

## 2. File Changes
- `src/cmd/verify.go`
  - Added `verify ears` Cobra subcommand and `CmdVerifyEARS` implementation
  - Implemented targeted linting, ID sanitization for bullets, and summary output (per-file + totals)
- `src/cmd/verify_ears_test.go`
  - Added tests for `CmdVerifyEARS` valid/invalid scenarios and message format
- `src/core/config/loader.go`
  - Extended `Guardrails.EARS` with `paths` and set defaults to design docs
- `src/core/config/factory.go`
  - Added `Paths` to template fields and defaults
- `src/core/ears/lint.go`
  - Suppressed ANTLR default console errors via a silent error listener
- `tgs/design/10_needs.md`
  - Added `N-022` (focused design doc linting)
- `tgs/design/20_requirements.md`
  - Added `SR-026` (verify ears subcommand behavior)
- `tgs/design/40_vnv.md`
  - Added V&V row for `SR-026`

## 3. Commands & Migrations
No migrations.

Build & test:
```bash
make build
go test ./...
```

Run verify ears:
```bash
./bin/tgs verify ears --repo . --ci
./bin/tgs verify ears --repo . --ci --paths tgs/design/10_needs.md,tgs/design/20_requirements.md
```

## 4. How to Test
1) All valid:
```bash
go test ./src/cmd -run TestVerify_EARS_DesignDocs_Valid
```
2) Invalid case shows `path:line: message` and summary counts:
```bash
go test ./src/cmd -run TestVerify_EARS_DesignDocs_InvalidReportsPathLine
./bin/tgs verify ears --repo . --ci
```

## 5. Integration Steps
- Optional config in `tgs/tgs.yml`:
```yaml
guardrails:
  ears:
    enable: true
    paths:
      - tgs/design/10_needs.md
      - tgs/design/20_requirements.md
```
- CI: add a step to run `./bin/tgs verify ears --repo . --ci`

## 6. Rollback
Revert this commit; the subcommand and config field are self-contained. No data migrations.

## 7. Follow-ups & Next Steps
- Consider a config flag to enforce "shall" in needs documents as well
- Add docs to README for `verify ears`

## 8. Links
- Thought: `tgs/thoughts/4b5a2a8-tgs-verify-ears-command/`
