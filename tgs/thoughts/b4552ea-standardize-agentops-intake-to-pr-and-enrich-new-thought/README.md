## Thought Directory Template

When you run `make new-thought title="..." [spec="..."]`, a directory `tgs/<BASE_HASH>-<kebab-title>/` is created with:

- `README.md` — auto-populated with Base Hash, quick links, and optional Idea Spec from `spec`.
- `research.md` — problem analysis, options, and recommended approach.
- `plan.md` — objectives, scope, acceptance criteria, phased steps, test/rollout.
- `implementation.md` — summary of changes, how to test, integration, rollback.

Edit `research.md` first, then `plan.md`. Only implement code after both are approved. Production code belongs outside `tgs/` (e.g., `src/`, `cmd/`).

