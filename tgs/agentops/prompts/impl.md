# TGS Implementation Prompt (phase: IMPLEMENT)

Use the **approved plan** to propose the smallest possible diff that satisfies the requirements.

## Approved Plan
{{APPROVED_PLAN_JSON}}

## Constraints
- Only modify files under: {{ALLOW_PATHS}}
- Keep total diff <= {{MAX_DIFF}} lines
- Run `make tgs:preflight` locally; ensure it passes

## Instructions
- Implement edits exactly as specified in the plan.
- Propose patches in unified diff format.
- Write **Conventional Commit** messages.

## Style & Patterns
- Follow [Linting Style](../agentops/linting_style.md)
- Add tests using [Testing Recipe](../agentops/testing_recipe.md)
- Use [Commit Message](../agentops/commit_message.md) guidance

## Output
- Patch hunks grouped by file
- Suggested commit message
