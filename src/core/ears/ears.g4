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
  : WHILE preconditions COMMA WHEN trigger COMMA (THE system | PRONOUN) SHALL response
  ;

// When <trigger>, the <system> shall <response>
eventReq
  : WHEN trigger COMMA (THE system | PRONOUN) SHALL response
  ;

// While <preconditions>, the <system> shall <response>
stateReq
  : WHILE preconditions COMMA (THE system | PRONOUN) SHALL response
  ;

// (While <preconditions>,) if <unwanted trigger>, then the <system> shall <response>
unwantedReq
  : (WHILE preconditions COMMA)? IF trigger COMMA THEN (THE system | PRONOUN) SHALL response
  ;

// The <system> shall <response>
ubiquitousReq
  : (THE system | PRONOUN) SHALL response
  ;

// --------------------
// Components
// --------------------

// One or more precondition clauses joined by "and/or" (comma is the separator between major clauses)
preconditions
  : clause
  ;

// Single trigger clause (no comma inside)
trigger
  : clause
  ;

// System name: one or more token words (allowing keywords inside names)
system
  : token_word+
  ;

// Response is the remainder of the line
response
  : (token_word | COMMA)*
  ;

// A clause is a sequence of token words (allow EARS keywords inside free text)
clause
  : token_word+
  ;

// token_word may be any non-comma word or selected EARS keyword tokens
// (exclude SHALL so it remains the modal before response)
token_word
  : THE
  | WHEN
  | IF
  | THEN
  | WORD
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
PRONOUN : [Ii][Tt] ;

// Punctuation
COMMA : ',' ;

// Free-text tokens
WORD : ~[ \t,\r\n]+ ;

// Whitespace & newlines
WS      : [ \t]+ -> skip ;
NEWLINE : ('\r'? '\n')+ -> skip ;
