// Code generated from src/core/ears/ears.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // ears

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr4-go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type earsParser struct {
	*antlr.BaseParser
}

var EarsParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func earsParserInit() {
	staticData := &EarsParserStaticData
	staticData.LiteralNames = []string{
		"", "", "", "", "", "", "", "", "','",
	}
	staticData.SymbolicNames = []string{
		"", "WHILE", "WHEN", "IF", "THEN", "THE", "SHALL", "CONJ", "COMMA",
		"TEXT_NOCOMMA", "WS", "NEWLINE",
	}
	staticData.RuleNames = []string{
		"requirement", "complexReq", "eventReq", "stateReq", "unwantedReq",
		"ubiquitousReq", "preconditions", "trigger", "system", "response", "clause",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 11, 120, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 3, 0, 38, 8, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2,
		1, 2, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 4, 1, 4, 1, 4,
		1, 4, 3, 4, 71, 8, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1,
		4, 1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 1, 6, 1, 6, 1, 6, 5, 6, 90, 8, 6, 10,
		6, 12, 6, 93, 9, 6, 1, 7, 1, 7, 1, 8, 1, 8, 5, 8, 99, 8, 8, 10, 8, 12,
		8, 102, 9, 8, 1, 9, 1, 9, 5, 9, 106, 8, 9, 10, 9, 12, 9, 109, 9, 9, 3,
		9, 111, 8, 9, 1, 10, 1, 10, 5, 10, 115, 8, 10, 10, 10, 12, 10, 118, 9,
		10, 1, 10, 0, 0, 11, 0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 0, 0, 118,
		0, 37, 1, 0, 0, 0, 2, 39, 1, 0, 0, 0, 4, 50, 1, 0, 0, 0, 6, 58, 1, 0, 0,
		0, 8, 70, 1, 0, 0, 0, 10, 81, 1, 0, 0, 0, 12, 86, 1, 0, 0, 0, 14, 94, 1,
		0, 0, 0, 16, 96, 1, 0, 0, 0, 18, 110, 1, 0, 0, 0, 20, 112, 1, 0, 0, 0,
		22, 23, 3, 2, 1, 0, 23, 24, 5, 0, 0, 1, 24, 38, 1, 0, 0, 0, 25, 26, 3,
		4, 2, 0, 26, 27, 5, 0, 0, 1, 27, 38, 1, 0, 0, 0, 28, 29, 3, 6, 3, 0, 29,
		30, 5, 0, 0, 1, 30, 38, 1, 0, 0, 0, 31, 32, 3, 8, 4, 0, 32, 33, 5, 0, 0,
		1, 33, 38, 1, 0, 0, 0, 34, 35, 3, 10, 5, 0, 35, 36, 5, 0, 0, 1, 36, 38,
		1, 0, 0, 0, 37, 22, 1, 0, 0, 0, 37, 25, 1, 0, 0, 0, 37, 28, 1, 0, 0, 0,
		37, 31, 1, 0, 0, 0, 37, 34, 1, 0, 0, 0, 38, 1, 1, 0, 0, 0, 39, 40, 5, 1,
		0, 0, 40, 41, 3, 12, 6, 0, 41, 42, 5, 8, 0, 0, 42, 43, 5, 2, 0, 0, 43,
		44, 3, 14, 7, 0, 44, 45, 5, 8, 0, 0, 45, 46, 5, 5, 0, 0, 46, 47, 3, 16,
		8, 0, 47, 48, 5, 6, 0, 0, 48, 49, 3, 18, 9, 0, 49, 3, 1, 0, 0, 0, 50, 51,
		5, 2, 0, 0, 51, 52, 3, 14, 7, 0, 52, 53, 5, 8, 0, 0, 53, 54, 5, 5, 0, 0,
		54, 55, 3, 16, 8, 0, 55, 56, 5, 6, 0, 0, 56, 57, 3, 18, 9, 0, 57, 5, 1,
		0, 0, 0, 58, 59, 5, 1, 0, 0, 59, 60, 3, 12, 6, 0, 60, 61, 5, 8, 0, 0, 61,
		62, 5, 5, 0, 0, 62, 63, 3, 16, 8, 0, 63, 64, 5, 6, 0, 0, 64, 65, 3, 18,
		9, 0, 65, 7, 1, 0, 0, 0, 66, 67, 5, 1, 0, 0, 67, 68, 3, 12, 6, 0, 68, 69,
		5, 8, 0, 0, 69, 71, 1, 0, 0, 0, 70, 66, 1, 0, 0, 0, 70, 71, 1, 0, 0, 0,
		71, 72, 1, 0, 0, 0, 72, 73, 5, 3, 0, 0, 73, 74, 3, 14, 7, 0, 74, 75, 5,
		8, 0, 0, 75, 76, 5, 4, 0, 0, 76, 77, 5, 5, 0, 0, 77, 78, 3, 16, 8, 0, 78,
		79, 5, 6, 0, 0, 79, 80, 3, 18, 9, 0, 80, 9, 1, 0, 0, 0, 81, 82, 5, 5, 0,
		0, 82, 83, 3, 16, 8, 0, 83, 84, 5, 6, 0, 0, 84, 85, 3, 18, 9, 0, 85, 11,
		1, 0, 0, 0, 86, 91, 3, 20, 10, 0, 87, 88, 5, 7, 0, 0, 88, 90, 3, 20, 10,
		0, 89, 87, 1, 0, 0, 0, 90, 93, 1, 0, 0, 0, 91, 89, 1, 0, 0, 0, 91, 92,
		1, 0, 0, 0, 92, 13, 1, 0, 0, 0, 93, 91, 1, 0, 0, 0, 94, 95, 3, 20, 10,
		0, 95, 15, 1, 0, 0, 0, 96, 100, 5, 9, 0, 0, 97, 99, 5, 9, 0, 0, 98, 97,
		1, 0, 0, 0, 99, 102, 1, 0, 0, 0, 100, 98, 1, 0, 0, 0, 100, 101, 1, 0, 0,
		0, 101, 17, 1, 0, 0, 0, 102, 100, 1, 0, 0, 0, 103, 107, 5, 9, 0, 0, 104,
		106, 5, 9, 0, 0, 105, 104, 1, 0, 0, 0, 106, 109, 1, 0, 0, 0, 107, 105,
		1, 0, 0, 0, 107, 108, 1, 0, 0, 0, 108, 111, 1, 0, 0, 0, 109, 107, 1, 0,
		0, 0, 110, 103, 1, 0, 0, 0, 110, 111, 1, 0, 0, 0, 111, 19, 1, 0, 0, 0,
		112, 116, 5, 9, 0, 0, 113, 115, 5, 9, 0, 0, 114, 113, 1, 0, 0, 0, 115,
		118, 1, 0, 0, 0, 116, 114, 1, 0, 0, 0, 116, 117, 1, 0, 0, 0, 117, 21, 1,
		0, 0, 0, 118, 116, 1, 0, 0, 0, 7, 37, 70, 91, 100, 107, 110, 116,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// earsParserInit initializes any static state used to implement earsParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewearsParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func EarsParserInit() {
	staticData := &EarsParserStaticData
	staticData.once.Do(earsParserInit)
}

// NewearsParser produces a new parser instance for the optional input antlr.TokenStream.
func NewearsParser(input antlr.TokenStream) *earsParser {
	EarsParserInit()
	this := new(earsParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &EarsParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
	this.GrammarFileName = "ears.g4"

	return this
}

// earsParser tokens.
const (
	earsParserEOF          = antlr.TokenEOF
	earsParserWHILE        = 1
	earsParserWHEN         = 2
	earsParserIF           = 3
	earsParserTHEN         = 4
	earsParserTHE          = 5
	earsParserSHALL        = 6
	earsParserCONJ         = 7
	earsParserCOMMA        = 8
	earsParserTEXT_NOCOMMA = 9
	earsParserWS           = 10
	earsParserNEWLINE      = 11
)

// earsParser rules.
const (
	earsParserRULE_requirement   = 0
	earsParserRULE_complexReq    = 1
	earsParserRULE_eventReq      = 2
	earsParserRULE_stateReq      = 3
	earsParserRULE_unwantedReq   = 4
	earsParserRULE_ubiquitousReq = 5
	earsParserRULE_preconditions = 6
	earsParserRULE_trigger       = 7
	earsParserRULE_system        = 8
	earsParserRULE_response      = 9
	earsParserRULE_clause        = 10
)

// IRequirementContext is an interface to support dynamic dispatch.
type IRequirementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ComplexReq() IComplexReqContext
	EOF() antlr.TerminalNode
	EventReq() IEventReqContext
	StateReq() IStateReqContext
	UnwantedReq() IUnwantedReqContext
	UbiquitousReq() IUbiquitousReqContext

	// IsRequirementContext differentiates from other interfaces.
	IsRequirementContext()
}

type RequirementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRequirementContext() *RequirementContext {
	var p = new(RequirementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = earsParserRULE_requirement
	return p
}

func InitEmptyRequirementContext(p *RequirementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = earsParserRULE_requirement
}

func (*RequirementContext) IsRequirementContext() {}

func NewRequirementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RequirementContext {
	var p = new(RequirementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = earsParserRULE_requirement

	return p
}

func (s *RequirementContext) GetParser() antlr.Parser { return s.parser }

func (s *RequirementContext) ComplexReq() IComplexReqContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IComplexReqContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IComplexReqContext)
}

func (s *RequirementContext) EOF() antlr.TerminalNode {
	return s.GetToken(earsParserEOF, 0)
}

func (s *RequirementContext) EventReq() IEventReqContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventReqContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventReqContext)
}

func (s *RequirementContext) StateReq() IStateReqContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStateReqContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStateReqContext)
}

func (s *RequirementContext) UnwantedReq() IUnwantedReqContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUnwantedReqContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUnwantedReqContext)
}

func (s *RequirementContext) UbiquitousReq() IUbiquitousReqContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUbiquitousReqContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUbiquitousReqContext)
}

func (s *RequirementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RequirementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RequirementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(earsListener); ok {
		listenerT.EnterRequirement(s)
	}
}

func (s *RequirementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(earsListener); ok {
		listenerT.ExitRequirement(s)
	}
}

func (p *earsParser) Requirement() (localctx IRequirementContext) {
	localctx = NewRequirementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, earsParserRULE_requirement)
	p.SetState(37)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 0, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(22)
			p.ComplexReq()
		}
		{
			p.SetState(23)
			p.Match(earsParserEOF)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(25)
			p.EventReq()
		}
		{
			p.SetState(26)
			p.Match(earsParserEOF)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(28)
			p.StateReq()
		}
		{
			p.SetState(29)
			p.Match(earsParserEOF)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(31)
			p.UnwantedReq()
		}
		{
			p.SetState(32)
			p.Match(earsParserEOF)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(34)
			p.UbiquitousReq()
		}
		{
			p.SetState(35)
			p.Match(earsParserEOF)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IComplexReqContext is an interface to support dynamic dispatch.
type IComplexReqContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	WHILE() antlr.TerminalNode
	Preconditions() IPreconditionsContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode
	WHEN() antlr.TerminalNode
	Trigger() ITriggerContext
	THE() antlr.TerminalNode
	System() ISystemContext
	SHALL() antlr.TerminalNode
	Response() IResponseContext

	// IsComplexReqContext differentiates from other interfaces.
	IsComplexReqContext()
}

type ComplexReqContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyComplexReqContext() *ComplexReqContext {
	var p = new(ComplexReqContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = earsParserRULE_complexReq
	return p
}

func InitEmptyComplexReqContext(p *ComplexReqContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = earsParserRULE_complexReq
}

func (*ComplexReqContext) IsComplexReqContext() {}

func NewComplexReqContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ComplexReqContext {
	var p = new(ComplexReqContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = earsParserRULE_complexReq

	return p
}

func (s *ComplexReqContext) GetParser() antlr.Parser { return s.parser }

func (s *ComplexReqContext) WHILE() antlr.TerminalNode {
	return s.GetToken(earsParserWHILE, 0)
}

func (s *ComplexReqContext) Preconditions() IPreconditionsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPreconditionsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPreconditionsContext)
}

func (s *ComplexReqContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(earsParserCOMMA)
}

func (s *ComplexReqContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(earsParserCOMMA, i)
}

func (s *ComplexReqContext) WHEN() antlr.TerminalNode {
	return s.GetToken(earsParserWHEN, 0)
}

func (s *ComplexReqContext) Trigger() ITriggerContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITriggerContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITriggerContext)
}

func (s *ComplexReqContext) THE() antlr.TerminalNode {
	return s.GetToken(earsParserTHE, 0)
}

func (s *ComplexReqContext) System() ISystemContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISystemContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISystemContext)
}

func (s *ComplexReqContext) SHALL() antlr.TerminalNode {
	return s.GetToken(earsParserSHALL, 0)
}

func (s *ComplexReqContext) Response() IResponseContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IResponseContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IResponseContext)
}

func (s *ComplexReqContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ComplexReqContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ComplexReqContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(earsListener); ok {
		listenerT.EnterComplexReq(s)
	}
}

func (s *ComplexReqContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(earsListener); ok {
		listenerT.ExitComplexReq(s)
	}
}

func (p *earsParser) ComplexReq() (localctx IComplexReqContext) {
	localctx = NewComplexReqContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, earsParserRULE_complexReq)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(39)
		p.Match(earsParserWHILE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(40)
		p.Preconditions()
	}
	{
		p.SetState(41)
		p.Match(earsParserCOMMA)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(42)
		p.Match(earsParserWHEN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(43)
		p.Trigger()
	}
	{
		p.SetState(44)
		p.Match(earsParserCOMMA)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(45)
		p.Match(earsParserTHE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(46)
		p.System()
	}
	{
		p.SetState(47)
		p.Match(earsParserSHALL)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(48)
		p.Response()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IEventReqContext is an interface to support dynamic dispatch.
type IEventReqContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	WHEN() antlr.TerminalNode
	Trigger() ITriggerContext
	COMMA() antlr.TerminalNode
	THE() antlr.TerminalNode
	System() ISystemContext
	SHALL() antlr.TerminalNode
	Response() IResponseContext

	// IsEventReqContext differentiates from other interfaces.
	IsEventReqContext()
}

type EventReqContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventReqContext() *EventReqContext {
	var p = new(EventReqContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = earsParserRULE_eventReq
	return p
}

func InitEmptyEventReqContext(p *EventReqContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = earsParserRULE_eventReq
}

func (*EventReqContext) IsEventReqContext() {}

func NewEventReqContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventReqContext {
	var p = new(EventReqContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = earsParserRULE_eventReq

	return p
}

func (s *EventReqContext) GetParser() antlr.Parser { return s.parser }

func (s *EventReqContext) WHEN() antlr.TerminalNode {
	return s.GetToken(earsParserWHEN, 0)
}

func (s *EventReqContext) Trigger() ITriggerContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITriggerContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITriggerContext)
}

func (s *EventReqContext) COMMA() antlr.TerminalNode {
	return s.GetToken(earsParserCOMMA, 0)
}

func (s *EventReqContext) THE() antlr.TerminalNode {
	return s.GetToken(earsParserTHE, 0)
}

func (s *EventReqContext) System() ISystemContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISystemContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISystemContext)
}

func (s *EventReqContext) SHALL() antlr.TerminalNode {
	return s.GetToken(earsParserSHALL, 0)
}

func (s *EventReqContext) Response() IResponseContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IResponseContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IResponseContext)
}

func (s *EventReqContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventReqContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventReqContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(earsListener); ok {
		listenerT.EnterEventReq(s)
	}
}

func (s *EventReqContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(earsListener); ok {
		listenerT.ExitEventReq(s)
	}
}

func (p *earsParser) EventReq() (localctx IEventReqContext) {
	localctx = NewEventReqContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, earsParserRULE_eventReq)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(50)
		p.Match(earsParserWHEN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(51)
		p.Trigger()
	}
	{
		p.SetState(52)
		p.Match(earsParserCOMMA)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(53)
		p.Match(earsParserTHE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(54)
		p.System()
	}
	{
		p.SetState(55)
		p.Match(earsParserSHALL)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(56)
		p.Response()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IStateReqContext is an interface to support dynamic dispatch.
type IStateReqContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	WHILE() antlr.TerminalNode
	Preconditions() IPreconditionsContext
	COMMA() antlr.TerminalNode
	THE() antlr.TerminalNode
	System() ISystemContext
	SHALL() antlr.TerminalNode
	Response() IResponseContext

	// IsStateReqContext differentiates from other interfaces.
	IsStateReqContext()
}

type StateReqContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStateReqContext() *StateReqContext {
	var p = new(StateReqContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = earsParserRULE_stateReq
	return p
}

func InitEmptyStateReqContext(p *StateReqContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = earsParserRULE_stateReq
}

func (*StateReqContext) IsStateReqContext() {}

func NewStateReqContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StateReqContext {
	var p = new(StateReqContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = earsParserRULE_stateReq

	return p
}

func (s *StateReqContext) GetParser() antlr.Parser { return s.parser }

func (s *StateReqContext) WHILE() antlr.TerminalNode {
	return s.GetToken(earsParserWHILE, 0)
}

func (s *StateReqContext) Preconditions() IPreconditionsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPreconditionsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPreconditionsContext)
}

func (s *StateReqContext) COMMA() antlr.TerminalNode {
	return s.GetToken(earsParserCOMMA, 0)
}

func (s *StateReqContext) THE() antlr.TerminalNode {
	return s.GetToken(earsParserTHE, 0)
}

func (s *StateReqContext) System() ISystemContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISystemContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISystemContext)
}

func (s *StateReqContext) SHALL() antlr.TerminalNode {
	return s.GetToken(earsParserSHALL, 0)
}

func (s *StateReqContext) Response() IResponseContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IResponseContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IResponseContext)
}

func (s *StateReqContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StateReqContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StateReqContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(earsListener); ok {
		listenerT.EnterStateReq(s)
	}
}

func (s *StateReqContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(earsListener); ok {
		listenerT.ExitStateReq(s)
	}
}

func (p *earsParser) StateReq() (localctx IStateReqContext) {
	localctx = NewStateReqContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, earsParserRULE_stateReq)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(58)
		p.Match(earsParserWHILE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(59)
		p.Preconditions()
	}
	{
		p.SetState(60)
		p.Match(earsParserCOMMA)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(61)
		p.Match(earsParserTHE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(62)
		p.System()
	}
	{
		p.SetState(63)
		p.Match(earsParserSHALL)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(64)
		p.Response()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IUnwantedReqContext is an interface to support dynamic dispatch.
type IUnwantedReqContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IF() antlr.TerminalNode
	Trigger() ITriggerContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode
	THEN() antlr.TerminalNode
	THE() antlr.TerminalNode
	System() ISystemContext
	SHALL() antlr.TerminalNode
	Response() IResponseContext
	WHILE() antlr.TerminalNode
	Preconditions() IPreconditionsContext

	// IsUnwantedReqContext differentiates from other interfaces.
	IsUnwantedReqContext()
}

type UnwantedReqContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUnwantedReqContext() *UnwantedReqContext {
	var p = new(UnwantedReqContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = earsParserRULE_unwantedReq
	return p
}

func InitEmptyUnwantedReqContext(p *UnwantedReqContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = earsParserRULE_unwantedReq
}

func (*UnwantedReqContext) IsUnwantedReqContext() {}

func NewUnwantedReqContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UnwantedReqContext {
	var p = new(UnwantedReqContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = earsParserRULE_unwantedReq

	return p
}

func (s *UnwantedReqContext) GetParser() antlr.Parser { return s.parser }

func (s *UnwantedReqContext) IF() antlr.TerminalNode {
	return s.GetToken(earsParserIF, 0)
}

func (s *UnwantedReqContext) Trigger() ITriggerContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITriggerContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITriggerContext)
}

func (s *UnwantedReqContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(earsParserCOMMA)
}

func (s *UnwantedReqContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(earsParserCOMMA, i)
}

func (s *UnwantedReqContext) THEN() antlr.TerminalNode {
	return s.GetToken(earsParserTHEN, 0)
}

func (s *UnwantedReqContext) THE() antlr.TerminalNode {
	return s.GetToken(earsParserTHE, 0)
}

func (s *UnwantedReqContext) System() ISystemContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISystemContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISystemContext)
}

func (s *UnwantedReqContext) SHALL() antlr.TerminalNode {
	return s.GetToken(earsParserSHALL, 0)
}

func (s *UnwantedReqContext) Response() IResponseContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IResponseContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IResponseContext)
}

func (s *UnwantedReqContext) WHILE() antlr.TerminalNode {
	return s.GetToken(earsParserWHILE, 0)
}

func (s *UnwantedReqContext) Preconditions() IPreconditionsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPreconditionsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPreconditionsContext)
}

func (s *UnwantedReqContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UnwantedReqContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *UnwantedReqContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(earsListener); ok {
		listenerT.EnterUnwantedReq(s)
	}
}

func (s *UnwantedReqContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(earsListener); ok {
		listenerT.ExitUnwantedReq(s)
	}
}

func (p *earsParser) UnwantedReq() (localctx IUnwantedReqContext) {
	localctx = NewUnwantedReqContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, earsParserRULE_unwantedReq)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(70)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == earsParserWHILE {
		{
			p.SetState(66)
			p.Match(earsParserWHILE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(67)
			p.Preconditions()
		}
		{
			p.SetState(68)
			p.Match(earsParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(72)
		p.Match(earsParserIF)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(73)
		p.Trigger()
	}
	{
		p.SetState(74)
		p.Match(earsParserCOMMA)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(75)
		p.Match(earsParserTHEN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(76)
		p.Match(earsParserTHE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(77)
		p.System()
	}
	{
		p.SetState(78)
		p.Match(earsParserSHALL)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(79)
		p.Response()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IUbiquitousReqContext is an interface to support dynamic dispatch.
type IUbiquitousReqContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	THE() antlr.TerminalNode
	System() ISystemContext
	SHALL() antlr.TerminalNode
	Response() IResponseContext

	// IsUbiquitousReqContext differentiates from other interfaces.
	IsUbiquitousReqContext()
}

type UbiquitousReqContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUbiquitousReqContext() *UbiquitousReqContext {
	var p = new(UbiquitousReqContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = earsParserRULE_ubiquitousReq
	return p
}

func InitEmptyUbiquitousReqContext(p *UbiquitousReqContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = earsParserRULE_ubiquitousReq
}

func (*UbiquitousReqContext) IsUbiquitousReqContext() {}

func NewUbiquitousReqContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UbiquitousReqContext {
	var p = new(UbiquitousReqContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = earsParserRULE_ubiquitousReq

	return p
}

func (s *UbiquitousReqContext) GetParser() antlr.Parser { return s.parser }

func (s *UbiquitousReqContext) THE() antlr.TerminalNode {
	return s.GetToken(earsParserTHE, 0)
}

func (s *UbiquitousReqContext) System() ISystemContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISystemContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISystemContext)
}

func (s *UbiquitousReqContext) SHALL() antlr.TerminalNode {
	return s.GetToken(earsParserSHALL, 0)
}

func (s *UbiquitousReqContext) Response() IResponseContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IResponseContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IResponseContext)
}

func (s *UbiquitousReqContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UbiquitousReqContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *UbiquitousReqContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(earsListener); ok {
		listenerT.EnterUbiquitousReq(s)
	}
}

func (s *UbiquitousReqContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(earsListener); ok {
		listenerT.ExitUbiquitousReq(s)
	}
}

func (p *earsParser) UbiquitousReq() (localctx IUbiquitousReqContext) {
	localctx = NewUbiquitousReqContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, earsParserRULE_ubiquitousReq)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(81)
		p.Match(earsParserTHE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(82)
		p.System()
	}
	{
		p.SetState(83)
		p.Match(earsParserSHALL)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(84)
		p.Response()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IPreconditionsContext is an interface to support dynamic dispatch.
type IPreconditionsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllClause() []IClauseContext
	Clause(i int) IClauseContext
	AllCONJ() []antlr.TerminalNode
	CONJ(i int) antlr.TerminalNode

	// IsPreconditionsContext differentiates from other interfaces.
	IsPreconditionsContext()
}

type PreconditionsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPreconditionsContext() *PreconditionsContext {
	var p = new(PreconditionsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = earsParserRULE_preconditions
	return p
}

func InitEmptyPreconditionsContext(p *PreconditionsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = earsParserRULE_preconditions
}

func (*PreconditionsContext) IsPreconditionsContext() {}

func NewPreconditionsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PreconditionsContext {
	var p = new(PreconditionsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = earsParserRULE_preconditions

	return p
}

func (s *PreconditionsContext) GetParser() antlr.Parser { return s.parser }

func (s *PreconditionsContext) AllClause() []IClauseContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IClauseContext); ok {
			len++
		}
	}

	tst := make([]IClauseContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IClauseContext); ok {
			tst[i] = t.(IClauseContext)
			i++
		}
	}

	return tst
}

func (s *PreconditionsContext) Clause(i int) IClauseContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IClauseContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IClauseContext)
}

func (s *PreconditionsContext) AllCONJ() []antlr.TerminalNode {
	return s.GetTokens(earsParserCONJ)
}

func (s *PreconditionsContext) CONJ(i int) antlr.TerminalNode {
	return s.GetToken(earsParserCONJ, i)
}

func (s *PreconditionsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PreconditionsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PreconditionsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(earsListener); ok {
		listenerT.EnterPreconditions(s)
	}
}

func (s *PreconditionsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(earsListener); ok {
		listenerT.ExitPreconditions(s)
	}
}

func (p *earsParser) Preconditions() (localctx IPreconditionsContext) {
	localctx = NewPreconditionsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, earsParserRULE_preconditions)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(86)
		p.Clause()
	}
	p.SetState(91)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == earsParserCONJ {
		{
			p.SetState(87)
			p.Match(earsParserCONJ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(88)
			p.Clause()
		}

		p.SetState(93)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITriggerContext is an interface to support dynamic dispatch.
type ITriggerContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Clause() IClauseContext

	// IsTriggerContext differentiates from other interfaces.
	IsTriggerContext()
}

type TriggerContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTriggerContext() *TriggerContext {
	var p = new(TriggerContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = earsParserRULE_trigger
	return p
}

func InitEmptyTriggerContext(p *TriggerContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = earsParserRULE_trigger
}

func (*TriggerContext) IsTriggerContext() {}

func NewTriggerContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TriggerContext {
	var p = new(TriggerContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = earsParserRULE_trigger

	return p
}

func (s *TriggerContext) GetParser() antlr.Parser { return s.parser }

func (s *TriggerContext) Clause() IClauseContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IClauseContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IClauseContext)
}

func (s *TriggerContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TriggerContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TriggerContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(earsListener); ok {
		listenerT.EnterTrigger(s)
	}
}

func (s *TriggerContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(earsListener); ok {
		listenerT.ExitTrigger(s)
	}
}

func (p *earsParser) Trigger() (localctx ITriggerContext) {
	localctx = NewTriggerContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, earsParserRULE_trigger)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(94)
		p.Clause()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ISystemContext is an interface to support dynamic dispatch.
type ISystemContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllTEXT_NOCOMMA() []antlr.TerminalNode
	TEXT_NOCOMMA(i int) antlr.TerminalNode

	// IsSystemContext differentiates from other interfaces.
	IsSystemContext()
}

type SystemContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySystemContext() *SystemContext {
	var p = new(SystemContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = earsParserRULE_system
	return p
}

func InitEmptySystemContext(p *SystemContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = earsParserRULE_system
}

func (*SystemContext) IsSystemContext() {}

func NewSystemContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SystemContext {
	var p = new(SystemContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = earsParserRULE_system

	return p
}

func (s *SystemContext) GetParser() antlr.Parser { return s.parser }

func (s *SystemContext) AllTEXT_NOCOMMA() []antlr.TerminalNode {
	return s.GetTokens(earsParserTEXT_NOCOMMA)
}

func (s *SystemContext) TEXT_NOCOMMA(i int) antlr.TerminalNode {
	return s.GetToken(earsParserTEXT_NOCOMMA, i)
}

func (s *SystemContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SystemContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SystemContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(earsListener); ok {
		listenerT.EnterSystem(s)
	}
}

func (s *SystemContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(earsListener); ok {
		listenerT.ExitSystem(s)
	}
}

func (p *earsParser) System() (localctx ISystemContext) {
	localctx = NewSystemContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, earsParserRULE_system)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(96)
		p.Match(earsParserTEXT_NOCOMMA)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(100)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == earsParserTEXT_NOCOMMA {
		{
			p.SetState(97)
			p.Match(earsParserTEXT_NOCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(102)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IResponseContext is an interface to support dynamic dispatch.
type IResponseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllTEXT_NOCOMMA() []antlr.TerminalNode
	TEXT_NOCOMMA(i int) antlr.TerminalNode

	// IsResponseContext differentiates from other interfaces.
	IsResponseContext()
}

type ResponseContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyResponseContext() *ResponseContext {
	var p = new(ResponseContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = earsParserRULE_response
	return p
}

func InitEmptyResponseContext(p *ResponseContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = earsParserRULE_response
}

func (*ResponseContext) IsResponseContext() {}

func NewResponseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ResponseContext {
	var p = new(ResponseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = earsParserRULE_response

	return p
}

func (s *ResponseContext) GetParser() antlr.Parser { return s.parser }

func (s *ResponseContext) AllTEXT_NOCOMMA() []antlr.TerminalNode {
	return s.GetTokens(earsParserTEXT_NOCOMMA)
}

func (s *ResponseContext) TEXT_NOCOMMA(i int) antlr.TerminalNode {
	return s.GetToken(earsParserTEXT_NOCOMMA, i)
}

func (s *ResponseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ResponseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ResponseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(earsListener); ok {
		listenerT.EnterResponse(s)
	}
}

func (s *ResponseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(earsListener); ok {
		listenerT.ExitResponse(s)
	}
}

func (p *earsParser) Response() (localctx IResponseContext) {
	localctx = NewResponseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, earsParserRULE_response)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(110)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == earsParserTEXT_NOCOMMA {
		{
			p.SetState(103)
			p.Match(earsParserTEXT_NOCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(107)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == earsParserTEXT_NOCOMMA {
			{
				p.SetState(104)
				p.Match(earsParserTEXT_NOCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

			p.SetState(109)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IClauseContext is an interface to support dynamic dispatch.
type IClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllTEXT_NOCOMMA() []antlr.TerminalNode
	TEXT_NOCOMMA(i int) antlr.TerminalNode

	// IsClauseContext differentiates from other interfaces.
	IsClauseContext()
}

type ClauseContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyClauseContext() *ClauseContext {
	var p = new(ClauseContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = earsParserRULE_clause
	return p
}

func InitEmptyClauseContext(p *ClauseContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = earsParserRULE_clause
}

func (*ClauseContext) IsClauseContext() {}

func NewClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ClauseContext {
	var p = new(ClauseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = earsParserRULE_clause

	return p
}

func (s *ClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *ClauseContext) AllTEXT_NOCOMMA() []antlr.TerminalNode {
	return s.GetTokens(earsParserTEXT_NOCOMMA)
}

func (s *ClauseContext) TEXT_NOCOMMA(i int) antlr.TerminalNode {
	return s.GetToken(earsParserTEXT_NOCOMMA, i)
}

func (s *ClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(earsListener); ok {
		listenerT.EnterClause(s)
	}
}

func (s *ClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(earsListener); ok {
		listenerT.ExitClause(s)
	}
}

func (p *earsParser) Clause() (localctx IClauseContext) {
	localctx = NewClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, earsParserRULE_clause)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(112)
		p.Match(earsParserTEXT_NOCOMMA)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(116)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == earsParserTEXT_NOCOMMA {
		{
			p.SetState(113)
			p.Match(earsParserTEXT_NOCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(118)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}
