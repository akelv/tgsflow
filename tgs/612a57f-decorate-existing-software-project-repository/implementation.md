# Implementation Summary: Decorate Existing Software Project Repository

## 1. Overview (What & Why)
Add a non-interactive decorate mode to `bootstrap.sh` that injects only the minimal TGS workflow files into an existing repository, avoiding templates and repo-level README by default, and without creating a subdirectory or initializing a new git repository. Includes `--dry-run`, `--force`, and optional `--with-templates` overlay.

## 2. File Changes
- Edited: `bootstrap.sh`
  - Bump `SCRIPT_VERSION` to `1.2.0`
  - Add flags: `--decorate`, `--dry-run`, `--force`, `--with-templates=<type>`
  - Implement decorate flow: copy `agentops/AGENTOPS.md`, `tgs/README.md`, `agentops/tgs/*`; generate `tgs.mk`; ensure `Makefile` includes it
  - Safe copy helpers, idempotent behavior, and logs
  - Detect existing `.git` and, when not explicitly decorating, prompt the user to choose: decorate current repo, create a new project, or quit (respects `--dry-run` messaging)
  - Keep existing bootstrap behavior unchanged by default

## 3. Commands & Migrations
No migrations. New runtime-generated file `tgs.mk` provides `new-thought` target when included.

## 4. How to Test
- Dry run (no changes):
  ```bash
  ./bootstrap.sh --decorate --dry-run
  ```
- Apply minimal workflow into current directory:
  ```bash
  ./bootstrap.sh --decorate
  ```
- Force overwrite existing files if needed:
  ```bash
  ./bootstrap.sh --decorate --force
  ```
- Optional template overlay (e.g., react):
  ```bash
  ./bootstrap.sh --decorate --with-templates=react
  ```

- Run without `--decorate` inside an existing git repository to see the prompt:
  ```bash
  ./bootstrap.sh
  # Expect a prompt: decorate current repo (d), new project (n), or quit (q)
  ```

Expected after decorate (no template):
- `agentops/AGENTOPS.md`, `tgs/README.md`, and `agentops/tgs/*` present
- `tgs.mk` created; `Makefile` includes `tgs.mk` (or minimal `Makefile` created)
- Existing files not overwritten unless `--force`

## 5. Integration Steps
- Commit added files and include `tgs.mk` in your `Makefile` if not already.
- Use `make new-thought title="Your idea"` to create thought directories.

## 6. Rollback
- Remove generated files: `agentops/AGENTOPS.md`, `tgs/README.md`, `agentops/tgs/*` (if newly added), `tgs.mk`, and the `include tgs.mk` line from `Makefile`.

## 7. Follow-ups & Next Steps
- Consider `--only`/`--exclude` patterns for advanced control.
- Add CI smoke test for `--decorate --dry-run` path.

## 8. Links
- Thought: `tgs/612a57f-decorate-existing-software-project-repository/`
