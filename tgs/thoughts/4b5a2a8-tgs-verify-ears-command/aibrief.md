# AI Brief

## Query
"implement tgs verify ears requirement that would review and lint design/10_needs.md and design/20_requirements.md"

## Context Summary (<= 6 bullets)
- TGS verify command already exists with EARS linting capability but needs enhancement for specific design document validation
- EARS parser is fully implemented with ANTLR grammar supporting 5 requirement shapes: ubiquitous, state-driven, event-driven, complex, unwanted
- Current verify function scans all .md files repo-wide but lacks focused design document validation
- Existing thought directory 4b5a2a8-tgs-verify-ears-command contains scaffolded research/plan templates for this exact feature
- Design documents 10_needs.md and 20_requirements.md follow EARS format with N-### and SR-### identifiers
- System already has complete EARS validation infrastructure including ParseRequirement and Lint functions

## Key Needs (EARS-style, with sources)
- [N-008] While writing needs and requirements, the Team needs INCOSE/EARS-aligned guidance and optional linting. (Source: tgs/design/10_needs.md:13)
- [N-009] When verifying docs, the Team needs a `verify` command that checks EARS patterns and Markdown bullets. (Source: tgs/design/10_needs.md:14)

## Key System Requirements (with sources)
- [SR-008] While writing needs and requirements, the system shall offer INCOSE/EARS-aligned guidance and optional linting rules. (Source: tgs/design/20_requirements.md:13)
- [SR-009] When verifying documentation, the system shall provide a `verify` command that evaluates EARS patterns and Markdown bullets and reports pass/fail via exit code. (Source: tgs/design/20_requirements.md:14)
- [IF-003] The system shall provide a `tgs verify --repo <PATH>` command to run documentation linting. (Source: tgs/design/20_requirements.md:47)

## Links & Pointers
- src/cmd/verify.go:86-159 – verifyEARS function implementation with markdown scanning logic
- src/core/ears/lint.go:42-156 – ParseRequirement and EARS validation core functionality  
- tgs/thoughts/4b5a2a8-tgs-verify-ears-command/ – existing thought directory for this exact feature
- tgs/design/40_vnv.md:23-24 – V&V criteria requiring exit code 0/!=0 and EARS issue reporting

## Notes
- Token budget: 1200