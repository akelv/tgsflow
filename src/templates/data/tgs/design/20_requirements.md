# System Requirements

Use “The system shall …” style.

## Functional Requirements
- The system shall scaffold a minimal `tgs/` directory from templates.
- The system shall provide `make new-thought` to create thought directories.
- The system shall include AgentOps workflow docs.

## Non-Functional Requirements
- The system shall operate on macOS/Linux.
- The system shall keep production code out of `tgs/`.

## Interfaces
- Make target: `new-thought`
- Script: `scripts/bootstrap.sh` (decorate mode)

---

### Checklist
- [ ] Singular, testable
- [ ] Uses “shall”
- [ ] Verification method assigned
