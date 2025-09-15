grammar ears;

// --------------------
// Parser rules
// --------------------

// Parse one requirement (one line). You can wrap with (requirement NEWLINE)* in your driver if needed.
requirement
  : complexReq EOF
  | eventReq   EOF
  | stateReq   EOF
  | unwantedReq EOF
  | ubiquitousReq EOF
  ;

// While <preconditions>, when <trigger>, the <system> shall <response>
complexReq
  : WHILE preconditions COMMA WHEN trigger COMMA THE system SHALL response
  ;

// When <trigger>, the <system> shall <response>
eventReq
  : WHEN trigger COMMA THE system SHALL response
  ;

// While <preconditions>, the <system> shall <response>
stateReq
  : WHILE preconditions COMMA THE system SHALL response
  ;

// (While <preconditions>,) if <unwanted trigger>, then the <system> shall <response>
unwantedReq
  : (WHILE preconditions COMMA)? IF trigger COMMA THEN THE system SHALL response
  ;

// The <system> shall <response>
ubiquitousReq
  : THE system SHALL response
  ;

// --------------------
// Components
// --------------------

// One or more precondition clauses joined by "and/or" (comma is the separator between major clauses)
preconditions
  : clause (CONJ clause)*
  ;

// Single trigger clause (no comma inside)
trigger
  : clause
  ;

// System name: one or more TEXT_NOCOMMA tokens (spaces are skipped)
system
  : TEXT_NOCOMMA (TEXT_NOCOMMA)*
  ;

// Response is zero or more TEXT_NOCOMMA tokens to end of line
response
  : (TEXT_NOCOMMA (TEXT_NOCOMMA)*)?
  ;

// A clause is one or more TEXT_NOCOMMA tokens (no comma, spaces skipped)
clause
  : TEXT_NOCOMMA (TEXT_NOCOMMA)*
  ;

// --------------------
// Lexer rules
// --------------------

// Case-insensitive keywords
WHILE : [Ww][Hh][Ii][Ll][Ee] ;
WHEN  : [Ww][Hh][Ee][Nn] ;
IF    : [Ii][Ff] ;
THEN  : [Tt][Hh][Ee][Nn] ;
THE   : [Tt][Hh][Ee] ;
SHALL : [Ss][Hh][Aa][Ll][Ll] ;

// Conjunction between multiple preconditions
CONJ  : (WS? COMMA WS?)? WS+ ([Aa][Nn][Dd] | [Oo][Rr]) WS+ ;

// Punctuation
COMMA : ',' ;

// Free-text tokens
// - TEXT_NOCOMMA: any run of non-space chars that does not include comma or newline
TEXT_NOCOMMA : (~[ \t,\r\n])+ ;

// Whitespace & newlines
WS      : [ \t]+ -> skip ;
NEWLINE : ('\r'? '\n')+ -> skip ;
