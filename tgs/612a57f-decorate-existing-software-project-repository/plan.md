# Plan: Decorate Existing Software Project Repository

## 1. Objectives
- Provide a safe, idempotent "decorate" mode in `bootstrap.sh` to adopt TGSFlow in an existing repository.
- Copy only the essential TGS workflow files; avoid templates and TGSFlow repo `README.md` by default.
- Do not create a subdirectory or initialize a git repo; operate in the current working directory.
- Add `--dry-run` and `--force` flags for visibility and control.

## 2. Scope / Non-goals
In-scope:
- Extend `bootstrap.sh` with `--decorate`, `--dry-run`, `--force`, and optional `--with-templates=<type>`.
- Logic to copy minimal files and wire up a `new-thought` make target via a new `tgs.mk` include file.
- Update thought docs (`implementation.md`) with how to test and integrate.

Out-of-scope:
- Changing existing project templates (react/python/go/cli).
- Overhauling repo docs beyond decorate flow.
- Building a separate distribution artifact; we keep a single script.

## 3. Acceptance Criteria
- `./bootstrap.sh --decorate` runs non-interactively and does NOT prompt for template or project name.
- No new subdirectory is created; no `git init` is executed in decorate mode.
- By default, the following are copied into the current directory (creating parent dirs as needed):
  - `agentops/AGENTOPS.md`
  - `tgs/README.md`
  - `agentops/tgs/*` (file templates for thoughts)
  - `tgs.mk` generated with a `new-thought` target equivalent to the root Makefile's implementation
- Existing files are not overwritten unless `--force` is provided; skipped files are logged.
- If `Makefile` exists and does not already include `tgs.mk` or define `new-thought`, append a single line `include tgs.mk` (idempotent).
- If `Makefile` does not exist, create a minimal one containing `include tgs.mk` and a `help` target that lists `new-thought`.
- The TGSFlow repo `README.md` at root and the `templates/` directory are NOT copied by default.
- `--with-templates=<react|python|go|cli|none>` is supported; default `none`. When set to a template type, files from `templates/<type>/` are copied into the current directory with the same non-destructive semantics.
- `--dry-run` prints all intended actions without writing any files. Works for both decorate and existing bootstrap paths.
- Script version bumped to `1.2.0`; help text updated to document new flags.

## 4. Phases & Tasks
- Phase 1: Design & Flags
  - [x] Define flags and behavior in research/plan
  - [ ] Update help/usage output in script
- Phase 2: Implement Decorate Mode
  - [ ] Add argument parsing and mode dispatch
  - [ ] Implement safe copy functions with `--dry-run` and `--force`
  - [ ] Generate `tgs.mk` and wire `Makefile` include/idempotency
  - [ ] Optional template overlay via `--with-templates`
- Phase 3: Testing & Docs
  - [ ] Manual tests across scenarios (with/without Makefile, with `--force`, idempotency)
  - [ ] Update `tgs/612a57f.../implementation.md` with summary and test steps

## 5. File/Module Changes
- Edit: `bootstrap.sh`
  - Add `--decorate`, `--dry-run`, `--force`, `--with-templates` flags
  - Add decorate execution path; no prompts; no subdir; no git init
  - Implement safe copy helpers and logging; bump `SCRIPT_VERSION` to `1.2.0`
- Add (generated at runtime by script, not committed): `tgs.mk`
  - Contains `new-thought` target copied from current `Makefile` logic
- Edit: `tgs/612a57f.../implementation.md` (docs)

## 6. Test Plan
Manual checks:
1) In a temp dir with no git:
   - Run: `curl -sSL https://raw.githubusercontent.com/akelv/tgsflow/main/bootstrap.sh | bash -s -- --decorate --dry-run`
   - Expect: Logs of intended copies; no files written.
   - Run without `--dry-run`: files `agentops/AGENTOPS.md`, `tgs/README.md`, `agentops/tgs/*`, `tgs.mk`, and a minimal `Makefile` created.
   - Re-run: no changes reported; idempotent.
2) In a repo with existing `Makefile` lacking `new-thought`:
   - Run decorate; expect `include tgs.mk` appended once.
3) With `--with-templates=none` (default): no app template files copied.
4) With `--with-templates=react` in an empty dir: files from `templates/react/` overlayed; existing files preserved unless `--force`.
5) With existing `README.md`: ensure root `README.md` from tgsflow is not copied or modified.

## 7. Rollout & Rollback
- Rollout: merge PR; users can immediately run with new flags.
- Rollback: revert PR; existing installs remain unaffected.

## 8. Estimates & Risks
- Estimate: 2-3 hours including testing and docs.
- Risks: accidental overwrites (mitigated by default skip + `--force`), path edge-cases on different shells (covered by guarded bash, minimal dependencies).

---
Approval checkpoint: Please review this plan and reply one of:
- APPROVE plan
- REQUEST CHANGES: <notes>
