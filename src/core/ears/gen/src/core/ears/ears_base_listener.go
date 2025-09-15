// Code generated from src/core/ears/ears.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // ears

import "github.com/antlr4-go/antlr/v4"

// BaseearsListener is a complete listener for a parse tree produced by earsParser.
type BaseearsListener struct{}

var _ earsListener = &BaseearsListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseearsListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseearsListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseearsListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseearsListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterRequirement is called when production requirement is entered.
func (s *BaseearsListener) EnterRequirement(ctx *RequirementContext) {}

// ExitRequirement is called when production requirement is exited.
func (s *BaseearsListener) ExitRequirement(ctx *RequirementContext) {}

// EnterComplexReq is called when production complexReq is entered.
func (s *BaseearsListener) EnterComplexReq(ctx *ComplexReqContext) {}

// ExitComplexReq is called when production complexReq is exited.
func (s *BaseearsListener) ExitComplexReq(ctx *ComplexReqContext) {}

// EnterEventReq is called when production eventReq is entered.
func (s *BaseearsListener) EnterEventReq(ctx *EventReqContext) {}

// ExitEventReq is called when production eventReq is exited.
func (s *BaseearsListener) ExitEventReq(ctx *EventReqContext) {}

// EnterStateReq is called when production stateReq is entered.
func (s *BaseearsListener) EnterStateReq(ctx *StateReqContext) {}

// ExitStateReq is called when production stateReq is exited.
func (s *BaseearsListener) ExitStateReq(ctx *StateReqContext) {}

// EnterUnwantedReq is called when production unwantedReq is entered.
func (s *BaseearsListener) EnterUnwantedReq(ctx *UnwantedReqContext) {}

// ExitUnwantedReq is called when production unwantedReq is exited.
func (s *BaseearsListener) ExitUnwantedReq(ctx *UnwantedReqContext) {}

// EnterUbiquitousReq is called when production ubiquitousReq is entered.
func (s *BaseearsListener) EnterUbiquitousReq(ctx *UbiquitousReqContext) {}

// ExitUbiquitousReq is called when production ubiquitousReq is exited.
func (s *BaseearsListener) ExitUbiquitousReq(ctx *UbiquitousReqContext) {}

// EnterPreconditions is called when production preconditions is entered.
func (s *BaseearsListener) EnterPreconditions(ctx *PreconditionsContext) {}

// ExitPreconditions is called when production preconditions is exited.
func (s *BaseearsListener) ExitPreconditions(ctx *PreconditionsContext) {}

// EnterTrigger is called when production trigger is entered.
func (s *BaseearsListener) EnterTrigger(ctx *TriggerContext) {}

// ExitTrigger is called when production trigger is exited.
func (s *BaseearsListener) ExitTrigger(ctx *TriggerContext) {}

// EnterSystem is called when production system is entered.
func (s *BaseearsListener) EnterSystem(ctx *SystemContext) {}

// ExitSystem is called when production system is exited.
func (s *BaseearsListener) ExitSystem(ctx *SystemContext) {}

// EnterResponse is called when production response is entered.
func (s *BaseearsListener) EnterResponse(ctx *ResponseContext) {}

// ExitResponse is called when production response is exited.
func (s *BaseearsListener) ExitResponse(ctx *ResponseContext) {}

// EnterClause is called when production clause is entered.
func (s *BaseearsListener) EnterClause(ctx *ClauseContext) {}

// ExitClause is called when production clause is exited.
func (s *BaseearsListener) ExitClause(ctx *ClauseContext) {}
