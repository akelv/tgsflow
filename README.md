# TGSFlow

**Thought-Guided Software development workflow for human-AI collaboration**

TGSFlow enables structured, thoughtful software development through an approval-gated workflow where humans maintain strategic thinking while AI handles implementation. Perfect for use with Claude Code, Cursor, and other AI coding assistants.

## 🚀 Quick Start

Bootstrap a new project with TGSFlow in seconds:

```bash
curl -sSL https://raw.githubusercontent.com/akelv/tgsflow/main/bootstrap.sh | bash
```

Choose from pre-configured templates:
- **React** - Modern React app with TypeScript, Vite, and ESLint
- **Python** - Python project with pyproject.toml and modern tooling
- **Go** - Go application with modules and CLI framework
- **CLI** - Cross-platform CLI tool with build scripts
- **None** - Just the TGS workflow for any project

Each template includes the complete TGS workflow for structured engineering.

## Install the tgs CLI

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
2. Use it as your AI assistant's system prompt
3. The AI will automatically follow the TGS workflow

### Manual TGS Setup

Create a new thought for any feature or change:

```bash
make new-thought title="Add user authentication"
```

This creates a structured directory with templates for research, planning, and implementation documentation.

## Project Templates

Available in `templates/`:
- [React](./templates/react/) - Modern React application with TypeScript
- [Python](./templates/python/) - Python project with modern packaging  
- [Go](./templates/go/) - Go application with standard structure
- [CLI](./templates/cli/) - Cross-platform CLI tool template

## Documentation

- **TGS Workflow Guide**: [agentops/AGENTOPS.md](./agentops/AGENTOPS.md)
- **Thought Organization**: [tgs/README.md](./tgs/README.md)
- **Template Reference**: [templates/README.md](./templates/README.md)

## Why TGSFlow?

- **Reduces AI hallucination** through structured planning
- **Maintains human oversight** on important decisions  
- **Creates audit trail** for all development decisions
- **Scales with team size** - clear handoff points
- **Framework agnostic** - works with any technology stack

## Contributing

TGSFlow follows its own methodology. To contribute:

1. Create a thought: `make new-thought title="Your improvement idea"`
2. Complete research and planning phases
3. Get approval before implementation
4. Submit PR with complete thought documentation

---

**Start building better software with human-AI collaboration** ✨

### EARS Linter (optional)

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