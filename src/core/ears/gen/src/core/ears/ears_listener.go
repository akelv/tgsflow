// Code generated from src/core/ears/ears.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // ears

import "github.com/antlr4-go/antlr/v4"

// earsListener is a complete listener for a parse tree produced by earsParser.
type earsListener interface {
	antlr.ParseTreeListener

	// EnterRequirement is called when entering the requirement production.
	EnterRequirement(c *RequirementContext)

	// EnterComplexReq is called when entering the complexReq production.
	EnterComplexReq(c *ComplexReqContext)

	// EnterEventReq is called when entering the eventReq production.
	EnterEventReq(c *EventReqContext)

	// EnterStateReq is called when entering the stateReq production.
	EnterStateReq(c *StateReqContext)

	// EnterUnwantedReq is called when entering the unwantedReq production.
	EnterUnwantedReq(c *UnwantedReqContext)

	// EnterUbiquitousReq is called when entering the ubiquitousReq production.
	EnterUbiquitousReq(c *UbiquitousReqContext)

	// EnterPreconditions is called when entering the preconditions production.
	EnterPreconditions(c *PreconditionsContext)

	// EnterTrigger is called when entering the trigger production.
	EnterTrigger(c *TriggerContext)

	// EnterSystem is called when entering the system production.
	EnterSystem(c *SystemContext)

	// EnterResponse is called when entering the response production.
	EnterResponse(c *ResponseContext)

	// EnterClause is called when entering the clause production.
	EnterClause(c *ClauseContext)

	// EnterToken_word is called when entering the token_word production.
	EnterToken_word(c *Token_wordContext)

	// ExitRequirement is called when exiting the requirement production.
	ExitRequirement(c *RequirementContext)

	// ExitComplexReq is called when exiting the complexReq production.
	ExitComplexReq(c *ComplexReqContext)

	// ExitEventReq is called when exiting the eventReq production.
	ExitEventReq(c *EventReqContext)

	// ExitStateReq is called when exiting the stateReq production.
	ExitStateReq(c *StateReqContext)

	// ExitUnwantedReq is called when exiting the unwantedReq production.
	ExitUnwantedReq(c *UnwantedReqContext)

	// ExitUbiquitousReq is called when exiting the ubiquitousReq production.
	ExitUbiquitousReq(c *UbiquitousReqContext)

	// ExitPreconditions is called when exiting the preconditions production.
	ExitPreconditions(c *PreconditionsContext)

	// ExitTrigger is called when exiting the trigger production.
	ExitTrigger(c *TriggerContext)

	// ExitSystem is called when exiting the system production.
	ExitSystem(c *SystemContext)

	// ExitResponse is called when exiting the response production.
	ExitResponse(c *ResponseContext)

	// ExitClause is called when exiting the clause production.
	ExitClause(c *ClauseContext)

	// ExitToken_word is called when exiting the token_word production.
	ExitToken_word(c *Token_wordContext)
}
