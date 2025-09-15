# TGSFlow

**Thought-Guided Software development workflow for human-AI collaboration**

TGSFlow enables structured, thoughtful software development through an approval-gated workflow where humans maintain strategic thinking while AI handles implementation. Perfect for use with Claude Code or any other AI coding agent or interactive vibe code in Cursor or any other IDE. 

TGS aims at helping both small team and big organization with team of teams to apply spec driven developement with heavy use of AI code agent to solve aspect related to guardrails, quality gate approvals, audit trail for enterprise level.

## Introduction

Our fundamental belief is that creating software with AI for human must requires structured, precise language and rigorous verification and validation to ensure that system behaviors are well-tracked, testable, and always traceable back to the original thought's intent. This ensure software continue to be safe for human use. 

To achieve this, our project adopts two proven foundations from the world of systems engineering: **INCOSE** guidelines and the **EARS (Easy Approach to Requirements Syntax)** method.

- INCOSE provides the discipline of writing well-formed needs and requirements that are clear, measurable, and verifiable, ensuring sets of requirements are consistent and complete across the system lifecycle. For more, see the [INCOSE Guide to Writing Requirements](https://www.incose.org/docs/default-source/working-groups/requirements-wg/gtwr/incose_rwg_gtwr_v4_040423_final_drafts.pdf).

- EARS complements this by offering simple yet powerful patterns for expressing requirements in unambiguous, structured natural language. You can review the [EARS resource by Alistair Mavin](https://www.incose.org/docs/default-source/working-groups/requirements-wg/rwg_iw2022/mav_ears_incoserwg_jan22.pdf).

These approaches were shaped in mission-critical domains like aircraft engines, environments where safety at scale, traceability, and team-wide alignment are non-negotiable. 

By bringing them into our AI-driven software development, we treat our work with the same level of care: authoring systems where **precision**, **accountability**, and **human trust** are built in from the very start.


## Quick Start 🚀 

Bootstrap a new project or apply on top of any existing project with TGSFlow in seconds:

```bash
curl -sSL https://raw.githubusercontent.com/akelv/tgsflow/main/bootstrap.sh | bash
```

## Install the tiny invisible tgs cli to improve thought quality 

- Homebrew (macOS/Linux):

```bash
brew tap akelv/tgs
brew install tgs
```

- Curl installer (portable):

```bash
curl -sSL https://raw.githubusercontent.com/akelv/tgsflow/main/scripts/install.sh | bash
```

Once installed, verify:

```bash
tgs --version
```

## The TGS Workflow

**TGS (Thought-Guided Software)** is an approval-gated workflow that ensures thoughtful development:

1. **Research** → Document problem, constraints, alternatives
2. **Plan** → Define implementation strategy and acceptance criteria  
3. **Human Approval** → Review and approve research + plan
4. **Implement** → Execute the approved plan
5. **Document** → Summarize what was built and how to use it

### Key Principles

- **Human thinks, AI implements** - Strategic decisions require human approval
- **Traceable thoughts** - Every change links to its research and planning
- **Approval gates** - No implementation without explicit human approval
- **Documentation-driven** - Clear records of why and how decisions were made

## Using with AI Code Assistants

### Claude Code / Cursor Integration

1. Copy the system prompt from `agentops/AGENTOPS.md`
2. Use it as your AI assistant's system prompt (CLAUDE.md or AGENTS.md) 
3. The AI will automatically follow the TGS workflow

### Manual TGS Setup

Create a new thought for any feature or change:

```bash
make new-thought title="Add user authentication" spec="A requirement specification" 
```

This creates a structured directory with templates for research, planning, and implementation documentation.

## Project Templates

Simple bootstrap project to bootstrap new project available in `templates/`:
- [React](./templates/react/) - Modern React application with TypeScript
- [Python](./templates/python/) - Python project with modern packaging  
- [Go](./templates/go/) - Go application with standard structure
- [CLI](./templates/cli/) - Cross-platform CLI tool template

## Documentation

- **TGS Workflow Guide**: [agentops/AGENTOPS.md](./agentops/AGENTOPS.md)
- **Thought Organization**: [tgs/README.md](./tgs/README.md)
- **Template Reference**: [templates/README.md](./templates/README.md)

## Why TGSFlow?

- **Ensure transparent intention** from every thought to working software 
- **Reduces AI hallucination** through structured planning
- **Maintains human oversight** on important decisions  
- **Creates audit trail** for all development decisions
- **Scales with team size** - clear handoff points
- **Framework agnostic** - works with any technology stack

## Contributing

TGSFlow follows its own methodology. To contribute:

1. Create a thought: `make new-thought`
2. Complete research and planning phases
3. Get approval before implementation
4. Submit PR with complete thoughts documentation in **tgs/**

### EARS Linter grammar update

Generate the ANTLR Go parser for `src/core/ears/ears.g4` (requires Java and ANTLR):

```bash
brew install openjdk antlr
export CLASSPATH="$(brew --prefix)/libexec/antlr-4.13.1-complete.jar:$CLASSPATH"
make ears-gen
```

Enable in `tgs.yaml`:

```yaml
policies:
  ears:
    enable: true
    require_shall: false
```

Run verify (will lint Markdown bullets when enabled):

```bash
./bin/tgs verify --repo .
```
---
**Start engineering serious software for human and AI**

