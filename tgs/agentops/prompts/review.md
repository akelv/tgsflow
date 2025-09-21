# TGS Review Prompt (phase: REVIEW)

You are reviewing a proposed PR for correctness, safety, and traceability.

## Diff Summary
{{DIFF_SUMMARY}}

## Context
- Relevant EARS: {{RELEVANT_EARS}}
- Requirements: {{RELEVANT_REQUIREMENTS}}

## Checks
- Are edits within allow_paths? If not, block.
- Are tests/docs updated appropriately?
- Do changes satisfy linked needs/requirements?
- Do diffs stay within diff budget?

## Output
- Checklist of blocking vs. non-blocking comments
- Minimal suggested patches if needed

## References
- [Review Checklist](../agentops/review_checklist.md)
- [Testing Recipe](../agentops/testing_recipe.md)
