package ears

import (
	"errors"
	"strings"

	antlr "github.com/antlr4-go/antlr/v4"
	earsp "github.com/kelvin/tgsflow/src/core/ears/gen/src/core/ears"
)

// Shape represents the recognized EARS requirement shape.
type Shape string

const (
	ShapeUbiquitous Shape = "ubiquitous"
	ShapeState      Shape = "state-driven"
	ShapeEvent      Shape = "event-driven"
	ShapeComplex    Shape = "complex"
	ShapeUnwanted   Shape = "unwanted"
)

// Result is the structured parse result for a requirement line.
type Result struct {
	Shape         Shape
	System        string
	Preconditions []string
	Trigger       string
	Response      string
}

// Issue represents a linting issue.
type Issue struct {
	FilePath string
	Line     int
	Message  string
}

// Available reports whether the generated parser is available.
func Available() bool { return true }

// ParseRequirement parses a single requirement line and returns a structured Result.
func ParseRequirement(line string) (Result, error) {
	input := antlr.NewInputStream(line)
	lexer := earsp.NewearsLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := earsp.NewearsParser(stream)

	root := parser.Requirement()

	res := Result{}
	req := root.(*earsp.RequirementContext)
	if ctx := req.ComplexReq(); ctx != nil {
		res.Shape = ShapeComplex
		res.Preconditions = extractPreconditions(ctx.Preconditions())
		res.Trigger = extractClauseText(ctx.Trigger().Clause())
		res.System = extractSystemText(ctx.System())
		res.Response = extractResponseText(ctx.Response())
		return res, nil
	}
	if ctx := req.EventReq(); ctx != nil {
		res.Shape = ShapeEvent
		res.Trigger = extractClauseText(ctx.Trigger().Clause())
		res.System = extractSystemText(ctx.System())
		res.Response = extractResponseText(ctx.Response())
		return res, nil
	}
	if ctx := req.StateReq(); ctx != nil {
		res.Shape = ShapeState
		res.Preconditions = extractPreconditions(ctx.Preconditions())
		res.System = extractSystemText(ctx.System())
		res.Response = extractResponseText(ctx.Response())
		return res, nil
	}
	if ctx := req.UnwantedReq(); ctx != nil {
		res.Shape = ShapeUnwanted
		if pc := ctx.Preconditions(); pc != nil {
			res.Preconditions = extractPreconditions(pc)
		}
		res.Trigger = extractClauseText(ctx.Trigger().Clause())
		res.System = extractSystemText(ctx.System())
		res.Response = extractResponseText(ctx.Response())
		return res, nil
	}
	if ctx := req.UbiquitousReq(); ctx != nil {
		res.Shape = ShapeUbiquitous
		res.System = extractSystemText(ctx.System())
		res.Response = extractResponseText(ctx.Response())
		return res, nil
	}

	return Result{}, errors.New("does not match an allowed EARS form")
}

func extractSystemText(s earsp.ISystemContext) string {
	if s == nil {
		return ""
	}
	toks := s.(*earsp.SystemContext).AllTEXT_NOCOMMA()
	parts := make([]string, 0, len(toks))
	for _, t := range toks {
		parts = append(parts, t.GetText())
	}
	return strings.TrimSpace(strings.Join(parts, " "))
}
func extractResponseText(r earsp.IResponseContext) string {
	if r == nil {
		return ""
	}
	// Response: zero or more TEXT_NOCOMMA
	// We iterate tokens from the underlying token stream between rule bounds; however,
	// the generated context does not expose them directly, so we can reconstruct via Clause helper
	// by temporarily treating response as a sequence of TEXT_NOCOMMA tokens using the same accessor.
	// Generated context provides no direct AllTEXT_NOCOMMA, so we fallback to rule text minus leading spaces.
	return strings.TrimSpace(r.(*earsp.ResponseContext).GetParser().GetTokenStream().GetTextFromRuleContext(r.(antlr.RuleContext)))
}
func extractClauseText(c earsp.IClauseContext) string {
	if c == nil {
		return ""
	}
	toks := c.(*earsp.ClauseContext).AllTEXT_NOCOMMA()
	parts := make([]string, 0, len(toks))
	for _, t := range toks {
		parts = append(parts, t.GetText())
	}
	return strings.TrimSpace(strings.Join(parts, " "))
}
func extractPreconditions(pc earsp.IPreconditionsContext) []string {
	if pc == nil {
		return nil
	}
	clauses := pc.AllClause()
	out := make([]string, 0, len(clauses))
	for _, c := range clauses {
		out = append(out, extractClauseText(c))
	}
	return out
}

// Lint parses multiple requirement lines and returns issues for those that don't match EARS shapes.
func Lint(lines []string) []Issue {
	var issues []Issue
	for i, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln == "" {
			continue
		}
		if _, err := ParseRequirement(ln); err != nil {
			issues = append(issues, Issue{FilePath: "", Line: i + 1, Message: err.Error()})
		}
	}
	return issues
}
