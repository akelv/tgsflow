package ears

import (
	"errors"
	"regexp"
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

	// Suppress default console error messages from ANTLR; we'll return structured errors instead
	silent := antlr.NewConsoleErrorListener() // placeholder; we'll remove console and add a no-op
	_ = silent
	lexer.RemoveErrorListeners()
	parser.RemoveErrorListeners()
	parser.AddErrorListener(silentNoop{})
	lexer.AddErrorListener(silentNoop{})

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
		if err := validateTriggerCtx(ctx.Trigger()); err != nil {
			return Result{}, err
		}
		// text-level system checks (after second comma, before 'shall')
		if err := validateSystemSegment(line, 2); err != nil {
			return Result{}, err
		}
		if ctx.System() != nil {
			res.System = extractSystemText(ctx.System())
		} else {
			res.System = "it"
		}
		if err := ensureHasShall(line); err != nil {
			return Result{}, err
		}
		res.Response = extractResponseText(ctx.Response())
		return res, nil
	}
	if ctx := req.EventReq(); ctx != nil {
		res.Shape = ShapeEvent
		res.Trigger = extractClauseText(ctx.Trigger())
		if err := validateTriggerCtx(ctx.Trigger()); err != nil {
			return Result{}, err
		}
		if err := validateSystemSegment(line, 1); err != nil {
			return Result{}, err
		}
		if ctx.System() != nil {
			res.System = extractSystemText(ctx.System())
		} else {
			res.System = "it"
		}
		if err := ensureHasShall(line); err != nil {
			return Result{}, err
		}
		res.Response = extractResponseText(ctx.Response())
		return res, nil
	}
	if ctx := req.StateReq(); ctx != nil {
		res.Shape = ShapeState
		res.Preconditions = extractPreconditions(ctx.Preconditions())
		if err := validateSystemSegment(line, 1); err != nil {
			return Result{}, err
		}
		if ctx.System() != nil {
			res.System = extractSystemText(ctx.System())
		} else {
			res.System = "it"
		}
		if err := ensureHasShall(line); err != nil {
			return Result{}, err
		}
		res.Response = extractResponseText(ctx.Response())
		return res, nil
	}
	if ctx := req.UnwantedReq(); ctx != nil {
		res.Shape = ShapeUnwanted
		if pc := ctx.Preconditions(); pc != nil {
			res.Preconditions = extractPreconditions(pc)
		}
		res.Trigger = extractClauseText(ctx.Trigger())
		if err := validateTriggerCtx(ctx.Trigger()); err != nil {
			return Result{}, err
		}
		if err := validateSystemSegment(line, 1); err != nil {
			return Result{}, err
		}
		if ctx.System() != nil {
			res.System = extractSystemText(ctx.System())
		} else {
			res.System = "it"
		}
		if err := ensureHasShall(line); err != nil {
			return Result{}, err
		}
		res.Response = extractResponseText(ctx.Response())
		return res, nil
	}
	if ctx := req.UbiquitousReq(); ctx != nil {
		res.Shape = ShapeUbiquitous
		if err := validateSystemSegment(line, 0); err != nil {
			return Result{}, err
		}
		if ctx.System() != nil {
			res.System = extractSystemText(ctx.System())
		} else {
			res.System = "it"
		}
		if err := ensureHasShall(line); err != nil {
			return Result{}, err
		}
		res.Response = extractResponseText(ctx.Response())
		return res, nil
	}

	return Result{}, errors.New("does not match an allowed EARS form")
}

// silentNoop is an ANTLR error listener that does nothing, preventing noisy console output.
type silentNoop struct{}

func (silentNoop) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
}
func (silentNoop) ReportAmbiguity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, exact bool, ambigAlts *antlr.BitSet, configs *antlr.ATNConfigSet) {
}
func (silentNoop) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, conflictingAlts *antlr.BitSet, configs *antlr.ATNConfigSet) {
}
func (silentNoop) ReportContextSensitivity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex, prediction int, configs *antlr.ATNConfigSet) {
}

func textFrom(ctx antlr.ParserRuleContext) string {
	if ctx == nil {
		return ""
	}
	return strings.TrimSpace(ctx.GetText())
}

func extractSystemText(s earsp.ISystemContext) string { return textFrom(s) }
func extractResponseText(r earsp.IResponseContext) string {
	if r == nil {
		return ""
	}
	return textFrom(r)
}
func extractClauseText(c earsp.ITriggerContext) string { return textFrom(c) }

func extractPreconditions(pc earsp.IPreconditionsContext) []string {
	if pc == nil {
		return nil
	}
	clause := pc.Clause()
	words := wordsFromRule(clause.(antlr.RuleContext))
	var segments [][]string
	var cur []string
	for _, w := range words {
		lw := strings.ToLower(w)
		if lw == "and" || lw == "or" {
			if len(cur) > 0 {
				segments = append(segments, cur)
				cur = nil
			}
			continue
		}
		cur = append(cur, w)
	}
	if len(cur) > 0 {
		segments = append(segments, cur)
	}
	out := make([]string, 0, len(segments))
	for _, seg := range segments {
		p := strings.TrimSpace(strings.Join(seg, " "))
		if p != "" {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		if s := textFrom(clause); s != "" {
			return []string{s}
		}
	}
	return out
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

func validateTriggerCtx(t earsp.ITriggerContext) error {
	words := wordsFromRule(t.(antlr.RuleContext))
	countWhen := 0
	for _, w := range words {
		lw := strings.ToLower(w)
		if lw == "when" {
			countWhen++
		}
		if lw == "while" {
			return errors.New("mixed 'while' inside trigger")
		}
	}
	if countWhen > 0 {
		return errors.New("multiple when clauses in trigger")
	}
	return nil
}

// validate the segment between the Nth comma and the word 'shall'
// nComma: 0 = start of line, 1 = after first comma, 2 = after second comma
func validateSystemSegment(line string, nComma int) error {
	low := strings.ToLower(line)
	idx := 0
	for i := 0; i < nComma; i++ {
		p := strings.Index(low[idx:], ",")
		if p < 0 {
			break
		}
		idx += p + 1
	}
	rest := strings.TrimSpace(low[idx:])
	// Allow optional leading "then" (for unwanted form: ", then the <system> shall ...")
	if strings.HasPrefix(rest, "then ") {
		rest = strings.TrimSpace(strings.TrimPrefix(rest, "then "))
	}
	shallPos := strings.Index(rest, " shall")
	seg := rest
	if shallPos >= 0 {
		seg = strings.TrimSpace(rest[:shallPos])
	}
	if nComma > 0 {
		// must start with 'the <name>' or pronoun 'it'
		if strings.HasPrefix(seg, "the ") {
			// ok, 'the' followed by a name
		} else if seg == "it" || strings.HasPrefix(seg, "it ") {
			// ok, pronoun only or pronoun followed by words (rare)
		} else {
			return errors.New("missing system name")
		}
	}
	// disallow keywords in system segment
	if strings.Contains(seg, " when ") || strings.Contains(seg, " while ") {
		return errors.New("invalid keyword in system")
	}
	return nil
}

func ensureHasShall(line string) error {
	if !strings.Contains(strings.ToLower(line), " shall") {
		return errors.New("missing shall")
	}
	return nil
}

var conjSplit = regexp.MustCompile(`(?i)\s+(and|or)\s+`)

func splitPreconditions(text string) []string {
	parts := conjSplit.Split(text, -1)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		if s := strings.TrimSpace(text); s != "" {
			return []string{s}
		}
	}
	return out
}

func wordsFromRule(rc antlr.RuleContext) []string {
	var words []string
	var walk func(n antlr.Tree)
	walk = func(n antlr.Tree) {
		switch t := n.(type) {
		case antlr.TerminalNode:
			words = append(words, t.GetText())
		case antlr.RuleContext:
			for _, ch := range t.GetChildren() {
				walk(ch)
			}
		}
	}
	walk(rc)
	return words
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
