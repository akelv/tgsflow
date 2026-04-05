# Judgment Journal

Human reasoning captured at TGS checkpoints. This is the highest-value artifact for preserving human craft — it records taste, experience, and domain knowledge that no AI can replicate.

## Why Keep a Judgment Journal?

AI agents handle implementation. But the *decisions* — why option A over B, what tradeoffs were acceptable, which stakeholders mattered most — those are human contributions that compound institutional knowledge over time.

## Template

```markdown
## [Date] — [Thought Hash] — [Title]
**Checkpoint**: Research / Plan / Implementation
**Decision**: APPROVED / APPROVED with modification / REJECTED
**Reasoning**: [Why this decision — the human context AI cannot access]
**Perspectives Weighed**: [Which stakeholders considered, what tradeoffs made]
**Context AI couldn't know**: [Organizational, business, or experience-based factors]
```

## How to Use

1. The `/tgs-close` skill prompts the human for a journal entry at close-out
2. Human provides reasoning (or skips for small changes)
3. Entry is appended to this file
4. Over time, the journal becomes a decision library that informs future TGS research

## Entries

<!-- Add entries below in reverse chronological order -->

## 2026-04-05 — db1ffa7 — Simplify Bootstrap / Upgrade Path

**Checkpoint**: Plan (approved after mid-plan pivot)
**Decision**: APPROVED with major modification — scrap Go `text/template` rendering and `.tmpl` files entirely; ship the scaffold as a tree of plain files copied verbatim.
**Reasoning**: Simplicity is the ultimate goal of this project. Given how elegantly and quickly skills, knowledge, and thinking patterns are evolving in this space (Claude Code slash commands, modular `.claude/rules/`, hooks), maintaining a parallel Go templating layer to stamp out the same files adds complexity for no significant value — it drags. The first plan tried to fix the Go templating gaps; the right move was to delete the templating.
**Perspectives Weighed**:
- **Future-self debt (decisive)**: Every additional abstraction layer is a future maintenance tax. Scrapping `text/template` removes ~500 lines of Go across `render.go`, `factory.go`, and `init.go`, and ~500 lines of bash from `bootstrap.sh` — all while making upgrades actually work.
- **End-user**: Users get a working modern `.claude/` setup in one `curl | bash`; upgrades via `--force` finally function.
- **Operator**: Single source of truth (`src/templates/scaffold/`) instead of a `.tmpl` tree + a rendering engine.
- **Tradeoff accepted**: Per-project variable substitution (project name, AI provider) is no longer automatic — users edit `tgs/tgs.yml` once after install. Worth it for the simplicity gain.
**Context AI couldn't know**: The project's design philosophy — "Simplicity is the ultimate complexity, we should try to solve the complex thoughts to make life simple." The templating was introduced when the scaffold needed variable substitution for `tgs.yml`; but with modern `.claude/` conventions and embedded files, direct copy is strictly better. Knowing when to throw away earlier work is a judgment call AI tends to avoid.
