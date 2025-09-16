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
	if parser.HasError() {
		return Result{}, errors.New("syntax error")
	}

	res := Result{}
	req := root.(*earsp.RequirementContext)
	if ctx := req.ComplexReq(); ctx != nil {
		res.Shape = ShapeComplex
		res.Preconditions = extractPreconditions(ctx.Preconditions())
		res.Trigger = extractClauseText(ctx.Trigger())
		res.System = extractSystemText(ctx.System())
		res.Response = extractResponseText(ctx.Response())
		return res, nil
	}
	if ctx := req.EventReq(); ctx != nil {
		res.Shape = ShapeEvent
		res.Trigger = extractClauseText(ctx.Trigger())
		if err := validateTrigger(res.Trigger); err != nil {
			return Result{}, err
		}
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
		res.Trigger = extractClauseText(ctx.Trigger())
		if err := validateTrigger(res.Trigger); err != nil {
			return Result{}, err
		}
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

func textFrom(ctx antlr.ParserRuleContext) string {
	if ctx == nil {
		return ""
	}
	return strings.TrimSpace(ctx.GetText())
}

func extractSystemText(s earsp.ISystemContext) string {
	return textFrom(s)
}
func extractResponseText(r earsp.IResponseContext) string {
	if r == nil {
		return ""
	}
	return textFrom(r)
}
func extractClauseText(c earsp.ITriggerContext) string {
	return textFrom(c)
}
func extractPreconditions(pc earsp.IPreconditionsContext) []string {
	if pc == nil {
		return nil
	}
	// grammar now defines a single clause; return one entry
	return []string{textFrom(pc.Clause())}
}

// validateTrigger adds semantic checks beyond the grammar
func validateTrigger(trigger string) error {
	s := strings.ToLower(" " + trigger + " ")
	if strings.Contains(s, " when ") {
		return errors.New("multiple when clauses in trigger")
	}
	if strings.Contains(s, " while ") {
		return errors.New("mixed 'while' inside trigger")
	}
	return nil
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
