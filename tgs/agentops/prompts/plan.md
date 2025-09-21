# TGS Plan Prompt (phase: PLAN)

You are preparing a **minimal, verifiable plan** before any code edits.

## Context
{{BRIEF}}

## Task Signal
{{TASK}}

## Guardrails
- Allowed paths: {{ALLOW_PATHS}}
- Max diff lines: {{MAX_DIFF}}
- Required checks: {{REQUIRED_CHECKS}}

## Instructions
1. Identify the exact files to edit (must be within allow_paths).
2. For each file, describe concrete edits (function names, line ranges, new blocks).
3. List tests/docs that should be updated.
4. Stop if the plan would exceed diff budget.

## Output Format
Return JSON:
```json
{
  "files": [{"path":"...", "edits":["..."]}],
  "tests": ["..."],
  "notes": ["risk1", "risk2"]
}
