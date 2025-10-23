# Verification & Validation Plan

## Approach
- Verification: prove requirements are built right via inspection, demonstration, test, or analysis.  
- Validation: prove solution meets stakeholder needs in real usage scenarios with the TGS workflow.  

- Methods used:
  - I = Inspection (files, configs, docs exist and are correct)
  - D = Demonstration (observe runtime behavior manually)
  - T = Test (scripted or automated checks with exit codes)
  - A = Analysis (reasoning, static checks, policy reviews)

## V&V Matrix
| Req ID | Method | Acceptance Criteria | Artifact/Test |
|--------|--------|---------------------|---------------|
| SR-001 | I/D    | Implementation actions are blocked until both `research.md` and `plan.md` are approved and recorded (no code merged prior; approvals present). | Thought docs with approvals; process evidence |
| SR-002 | I      | Approval metadata (approver, timestamp or statement) exists within the thought directory. | `tgs/<hash>-*/research.md`, `plan.md` |
| SR-003 | T      | After `make new-thought title="X"`, a directory `tgs/<hash>-x/` exists with `README.md`, `research.md`, `plan.md`, `implementation.md`; README shows base hash and quick links. | `make new-thought` run log, filesystem check |
| SR-004 | I      | `agentops/AGENTOPS.md` exists and contains the canonical system prompt language. | File presence and content |
| SR-005 | D      | Running `bootstrap.sh` initializes TGS structure without interaction. | Run `./bootstrap.sh --dry-run` or local demo |
| SR-006 | T      | `scripts/install.sh` DRY_RUN prints resolved URL and exits 0; Homebrew formula installs `tgs` and `tgs --version` prints version. | `DRY_RUN=1 bash scripts/install.sh`; `brew install tgs` |
| SR-007 | I      | Design docs exist: `tgs/design/00_context.md`, `10_needs.md`, `20_requirements.md`, `30_architecture.md`. | File presence |
| SR-008 | D/T    | With EARS policy enabled, `tgs verify --repo .` analyzes Markdown bullets and reports issues; disabling skips EARS checks. | `tgs.yaml` + `tgs verify` output |
| SR-009 | T      | `tgs verify --repo .` returns exit code 0 when no issues, non-zero when issues are present. | Exit code and stderr messages |
| SR-010 | T      | `make new-thought title=... spec=...` creates a directory with provided spec recorded in `README.md`. | Filesystem check |
| SR-011 | I      | Each PR references a thought directory and includes links to approvals and implementation docs. | PR body template or examples |
| SR-012 | D      | Commands (`tgs init`, `tgs verify`, `tgs agent exec`) run with flags non-interactively (no prompts). | Command runs/logs |
| SR-013 | I      | Templates exist under `templates/{react,python,go,cli}/` and are buildable or runnable per their readmes. | Directory listing |
| SR-014 | I/D    | Workflow phases are followed: Research → Plan → Approval → Implement → Document; evidence in thought directories. | Thought documentation trail |
| SR-015 | T      | Running decorate mode installs only minimal `tgs/` from templates under `templates/data/tgs/`; no repo-specific thought dirs are copied. | `./scripts/bootstrap.sh --decorate --dry-run` output and filesystem check |
| NFR-001 | I     | Every code change is traceable to a thought directory (commit/PR references `tgs/<hash>-*/`). | Repo history & PRs |
| NFR-002 | T     | Build and basic commands succeed on macOS and Linux. | `make build && ./bin/tgs --version` on both OSes |
| NFR-003 | T     | `tgs verify` returns exit code 0/!=0 appropriately. | CI job or local script |
| NFR-004 | I     | Thought directories contain `research.md`, `plan.md`, `implementation.md` before completion. | Directory inspection |
| NFR-005 | A/I   | No production code lives under `tgs/`; `tgs/` contains documentation only, per policy. | Policy review, tree scan |
| IF-001 | T      | `make new-thought` creates `tgs/<hash>-<kebab>/` as specified. | Make output and files |
| IF-002 | D      | `make ears-gen` regenerates parser files under `src/core/ears/gen/` when ANTLR is installed. | Generated files exist |
| IF-003 | T      | `./bin/tgs verify --repo .` runs and prints hook results and EARS issues when enabled. | Command output |
| IF-004 | D      | `bootstrap.sh` can be invoked to apply TGSFlow structure in a clean workspace. | Script run |
| IF-005 | I      | `agentops/AGENTOPS.md` present and referenced by assistants (`AGENTS.md`/`CLAUDE.md`). | File presence and linkage |
| SR-026 | T      | With compliant lines in `tgs/design/10_needs.md` and `tgs/design/20_requirements.md`, `tgs verify --repo . --ci` returns 0; introducing a non-EARS line in either causes a non-zero exit and a `path:line: message` error. | Unit tests in `src/cmd/verify_ears_test.go` using temp repos |
| SR-027 | T      | Executing `tgs/adapters/gemini-code.sh --prompt-text test --context-glob "README.md"` returns text output and exit code 0; passing a missing prompt or context yields non-zero with clear stderr. | `go test ./src/core/brain -run TestShellTransport*` with adapter path overridden |
| SR-020 | T      | `shellTransport.Chat` invokes the adapter and returns output text; errors on non-zero exit; respects context deadline. | `go test ./src/core/brain -run TestShellTransport*` |
| SR-021 | T      | `tgs context pack "auth"` creates `aibrief.md` under active thought with merged context/requirements. | CLI run, file exists and includes sources |
| SR-022 | I      | Prompt templates for search and brief exist under `templates/data/tgs/` and are referenced by the command. | File presence and code reference |
| SR-023 | I      | `aibrief.md` items include path and anchor/line ranges to source material. | Inspect generated brief |
| SR-024 | T      | Token count in generated brief is <= configured `context_pack_tokens`. | Measure tokens or approximate count |
| SR-025 | A/I    | Secrets are redacted per configured patterns; no raw secrets appear in output. | Rule inspection and sample runs |

| SR-028 | T      | After `tgs init`, `tgs/adapters/claude-code.sh` and `tgs/adapters/gemini-code.sh` exist and are executable if previously missing. | Filesystem check and `test -x` |
| SR-029 | T      | Running `tgs init claude` (or `gemini`) creates `CLAUDE.md` (`GEMINI.md`) at repo root if absent; if present, command exits non-zero with an instructive error. | CLI run and file presence/error behavior |
| SR-030 | T      | After `tgs init` on a repo without a `new-thought` target, the root `Makefile` contains the standard `new-thought` rule. | Inspect `Makefile` content |

## Test Environments
- OS: macOS (darwin), Linux (ubuntu-latest).  
- Tooling: Go toolchain (per `go.mod`), Java + ANTLR for `make ears-gen` (optional), Homebrew for CLI distribution, curl for installer.  
- Network access: Required for installer GitHub API requests and release downloads (DRY_RUN avoids download).  
- CI: GitHub Actions or GitLab CI (optional) to run `make build`, `go test ./...`, and `./bin/tgs verify --repo .`.  

---

### Checklist
- [x] Each requirement has a V&V method  
- [x] Acceptance criteria are objective/measurable  
- [x] Traceability from requirement → test case is clear  
- [x] Validation plan covers real stakeholder scenarios  
