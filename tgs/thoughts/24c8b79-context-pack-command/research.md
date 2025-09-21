# Research: Context Pack command

- Date: 2025-09-21
- Base Hash: 24c8b79
- Participants: Human, Agent

## 1. Problem Statement
Provide a `tgs context pack "<query>"` command that searches repo context (notably `tgs/design/*`) and the active thought directory to generate a concise AI brief (`aibrief.md`) with key requirements and links, honoring a token budget.

## 2. Current State
- CLI uses Cobra; root registers `newContextCommand()` placeholder.  
- Shell brain transport exists (`src/core/brain/transport_shell.go`) to call local adapter `tgs/adapters/claude-code.sh`.
- Config supports `AI.Toolpack` with budgets, tools, and redaction. Default budgets include `context_pack_tokens`.
- Thought discovery exists via `src/core/thoughts/locate.go`.

## 3. Constraints & Assumptions
- Implement production code under `src/`, not `tgs/`.
- Use shell transport; no network SDKs required.
- Limit brief tokens using `ai.toolpack.budgets.context_pack_tokens` with safe default (1200).  
- Redact secrets via configured patterns.
- Do not modify files under `tgs/` except writing the generated brief in the active thought.

## 4. Risks & Impact
- Risk of over-inclusion increasing cost: Mitigate via budget and strict prompts.  
- Risk of stale/irrelevant context: Include explicit source pointers and anchors for validation.  
- Adapter failure: Surface stderr and exit non-zero; do not create partial files.

## 5. Alternatives Considered
- Direct Go summarization without agent: lower cost but weaker quality.  
- SDK transport instead of shell: tighter integration but adds dependencies.  
- Grep-only pack: fast but brittle.

## 6. Recommendation
Implement `tgs context pack` that:  
1) Locates active thought dir; 2) Gathers candidate files from `tgs/design` and active thought; 3) Constructs search and brief prompts from templates; 4) Invokes brain transport with token budget; 5) Writes `aibrief.md` with sources and summary sections.

## 7. References & Links
- `src/core/brain/*`  
- `src/core/config/loader.go`  
- `src/core/thoughts/locate.go`  
- `src/cmd/agent_exec.go`  
- `tgs/agentops/AGENTOPS.md`

---
Approval checkpoint: Please review this research and reply one of:
- APPROVE research
- REQUEST CHANGES: <notes>
