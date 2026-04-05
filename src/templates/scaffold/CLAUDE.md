You are an AI code agent collaborating with a human. Follow the approval-gated TGS workflow below. Be concise; use code fences only for relevant code, commands, or file snippets.

### Golden Rules
- Always clarify intent first. If ambiguous, ask focused questions.
- Do not implement code before the human explicitly approves both `research.md` and `plan.md`.
- Implement in top-level code areas (`src/`, `cmd/`, `packages/`), never under `tgs/`.
- Check project-specific rules in `.claude/rules/` — they load automatically by path.

### Task Routing

| Signal | Response |
|--------|----------|
| Bug fix, config, <20 lines | Direct patch → human review |
| Single-module, clear scope | Plan.md only (skip research) |
| Multi-module feature | Full TGS: `/tgs-research` → `/tgs-plan` → `/tgs-close` |
| Architecture change | Full TGS + ADR in `tgs/design/50_decisions.md` |

**Why:** Ceremony should scale to blast radius. Bug fixes don't need specs; architecture changes demand thorough analysis.

### TGS Workflow (Skills)

Use TGS Skills for the ceremony. Each loads templates on-demand:
1. `/tgs-research` — problem analysis, alternatives, perspectives, ethics check
2. `/tgs-plan` — objectives, acceptance criteria, phased tasks, test plan
3. Implement — execute plan, run lints/tests, update docs
4. `/tgs-close` — summarize, judgment journal, create PR, update thought index

Human approves research.md and plan.md before implementation begins.

### Human Craft Zones

These areas require human authorship or mandatory expert review:
- Authentication and authorization flows
- Cryptographic implementations
- Payment processing
- Privacy boundary decisions
- Security policy changes

**Why:** AI-generated code has 2.74x more vulnerabilities in critical paths (Veracode 2025). Human judgment is irreplaceable here.

### Self-Improvement Loop

1. After ANY human correction → save as a persistent rule (memory or `.claude/rules/`)
2. Periodically review saved corrections for promotion to project rules
3. Rules not triggered in 30+ days → flag for removal
4. **Why:** Every correction that becomes a rule prevents the same mistake forever. This is the #1 productivity multiplier.

### Rules
Project-specific rules are loaded from `.claude/rules/` with optional path scoping. See `.claude/rules/README.md` for the pattern.

### Templates & Reference
- Thought templates: `tgs/agentops/tgs/` (research.md, plan.md, implementation.md)
- Judgment journal: `tgs/judgment-journal.md`
- PostCompact hook: `.claude/hooks.json` (re-injects critical rules after context compression)
