# AI Brief

## Query
"gemini ai brain"

## Context Summary (<= 6 bullets)
- Active thought directory `ff006d7-gemini-ai-brain` exists to create Gemini CLI adapter for TGS AI brain functionality
- Current implementation uses Claude Code adapter (`tgs/adapters/claude-code.sh`) as reference pattern
- Goal is to create equivalent shell adapter that works with Gemini CLI in same manner as Claude option
- Gemini CLI documentation available at https://github.com/google-gemini/gemini-cli/tree/main/docs
- Research, plan, and implementation templates are scaffolded but appear empty (using default templates)
- Shell Transport requirement (SR-020) specifies execution pattern for AI adapters with composed prompts and context files

## Key Needs (EARS-style, with sources)
- [N-004] While using AI coding agents, the Human Engineer needs a standardized system prompt so assistants follow the TGS workflow. (Source: tgs/design/10_needs.md:9)
- [N-012] While executing tasks, the AI Code Agent needs non-interactive commands and explicit, precise instructions. (Source: tgs/design/10_needs.md:17)
- [N-019] While initiating an AI task, the Human Engineer needs a concise, auto-generated brief that packs relevant context from `tgs/design/` and the active thought. (Source: tgs/design/10_needs.md:27)

## Key System Requirements (with sources)
- [SR-020] While operating in shell AI mode, the system shall provide a Shell Transport that executes the `tgs/adapters/claude-code.sh` adapter with a composed prompt and optional context files, and returns the adapter output as `ChatResp.Text`, honoring timeouts and exit codes. (Source: tgs/design/20_requirements.md:28)
- [SR-004] While using AI coding assistants, the system shall provide a canonical system prompt in `agentops/AGENTOPS.md` for assistants to follow the TGS workflow. (Source: tgs/design/20_requirements.md:9)
- [SR-012] While executing tasks, the system shall support non-interactive execution of setup and verification commands. (Source: tgs/design/20_requirements.md:17)

## Links & Pointers
- tgs/thoughts/ff006d7-gemini-ai-brain/README.md – active thought spec and links
- tgs/adapters/claude-code.sh – reference implementation pattern for shell adapters
- tgs/design/20_requirements.md:28 – Shell Transport requirement specifying adapter execution pattern
- https://github.com/google-gemini/gemini-cli/tree/main/docs – Gemini CLI documentation

## Notes
- Token budget: 1200