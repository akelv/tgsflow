# Implementation: Extend init to include adapters, subcommand, and Makefile target

## 1. Overview (What & Why)
Enhances `tgs init` to ensure adapters availability, provide vendor decoration (`tgs init claude|gemini`), and guarantee a `new-thought` Makefile target so teams can immediately use the workflow in decorated repos.

## 2. File Changes
- `src/cmd/init.go`
  - Add optional arg passthrough `tgs init [claude|gemini]`.
  - Add `ensureMakefileNewThought()` to append standard `new-thought` target when missing.
  - Add `decorateVendorReadme()` to copy `tgs/agentops/AGENTOPS.md` → root `CLAUDE.md`/`GEMINI.md` if absent.
  - Note: Ensuring adapters is currently disabled by maintainer comment to avoid copying in this repo; logic is present but commented at call-site.
- Docs updated:
  - `tgs/design/10_needs.md` (+N-024..N-026)
  - `tgs/design/20_requirements.md` (+SR-028..SR-030)
  - `tgs/design/40_vnv.md` (V&V for SR-028..SR-030)

## 3. Commands & Migrations
```bash
go build -o ./bin/tgs ./src
./bin/tgs init
./bin/tgs init claude   # creates CLAUDE.md if absent; errors if present
./bin/tgs init gemini   # creates GEMINI.md if absent; errors if present
```

## 4. How to Test
- `tgs init` creates/updates `tgs/` and appends `new-thought` to `Makefile` when missing.
- `tgs init claude` creates root `CLAUDE.md` sourced from `tgs/agentops/AGENTOPS.md` when absent; rerun yields a clear error and non-zero exit.
- `tgs init gemini` mirrors the behavior for `GEMINI.md`.

## 5. Integration Steps
- None required; features are opt-in and idempotent.

## 6. Rollback
- Revert edits to `src/cmd/init.go` and docs.

## 7. Follow-ups & Next Steps
- Consider re-enabling adapter script seeding in `init` (commented call) when policy allows copying from repo seeds.
- Optionally add a guarded `--force` for overwriting existing `CLAUDE.md`/`GEMINI.md`.

## 8. Links
- Thought: `tgs/thoughts/af66921-extend-init-adapters-subcmd-make-target/`
