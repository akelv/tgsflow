# Context Pack

## Personas
- Human Engineer/Developer — implements with AI under guardrails
- Technical Lead/Reviewer — approves research/plan
- AI Code Agent — executes approved plan precisely
- Compliance/Auditor — needs clear audit trail

## Scope
- In scope: Approval-gated TGS workflow; thought structure; design docs
- Out of scope: Product-specific features; secrets; vendor-specific setups

## Scenarios
- Bootstrap/decorate repo to add minimal TGS workflow
- Create new thought via `make new-thought`
- Use AgentOps prompt to guide AI agents

## Constraints
- macOS/Linux; non-interactive commands; no production code under `tgs/`

---

### Checklist
- [ ] Personas defined
- [ ] Scope explicit
- [ ] Scenarios outcome-focused
- [ ] Constraints documented
