# TGS Close-Out Phase

Implementation must be approved before closing.

## Steps

### 1. Update Thought Index
Add entry to `tgs/README.md` with hash, date, status, description.

### 2. Update Design Docs
Run `/design-sync` if available. Verify ADRs, needs, V&V docs updated if applicable.

### 3. Prompt for Judgment Journal Entry
Ask the human:
> Would you like to add a Judgment Journal entry? What was your reasoning for the key decisions you approved?

If provided, append to `tgs/judgment-journal.md`:
```markdown
## [Date] — [Hash] — [Title]
**Decision**: [What was approved/modified]
**Reasoning**: [Human's reasoning — the context AI cannot access]
**Perspectives**: [Tradeoffs weighed]
```

### 4. Create PR
```bash
gh pr create --fill --title "<type>: <title>" --body-file tgs/<dir>/implementation.md
```

### 5. File Follow-up Backlog Items
Scan for remaining follow-ups. Add to `backlog/pending/`.
