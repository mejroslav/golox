package token

type TokenType string

const (
	// Single-character tokens.
	LEFT_PAREN  TokenType = "LEFT_PAREN"
	RIGHT_PAREN TokenType = "RIGHT_PAREN"
	LEFT_BRACE  TokenType = "LEFT_BRACE"
	RIGHT_BRACE TokenType = "RIGHT_BRACE"
	COMMA       TokenType = "COMMA"
	DOT         TokenType = "DOT"
	MINUS       TokenType = "MINUS"
	PLUS        TokenType = "PLUS"
	SEMICOLON   TokenType = "SEMICOLON"
	SLASH       TokenType = "SLASH"
	STAR        TokenType = "STAR"

	// One or two character tokens.
	BANG          TokenType = "BANG"
	BANG_EQUAL    TokenType = "BANG_EQUAL"
	EQUAL         TokenType = "EQUAL"
	EQUAL_EQUAL   TokenType = "EQUAL_EQUAL"
	GREATER       TokenType = "GREATER"
	GREATER_EQUAL TokenType = "GREATER_EQUAL"
	LESS          TokenType = "LESS"
	LESS_EQUAL    TokenType = "LESS_EQUAL"

	// Literals.
	IDENTIFIER TokenType = "IDENTIFIER"
	STRING     TokenType = "STRING"
	NUMBER     TokenType = "NUMBER"

	// Keywords.
	NIL   TokenType = "NIL"
	TRUE  TokenType = "TRUE"
	FALSE TokenType = "FALSE"

	// Logical operators.
	AND TokenType = "AND"
	OR  TokenType = "OR"

	// Variable declaration.
	VAR TokenType = "VAR"

	// Control flow.
	IF    TokenType = "IF"
	ELSE  TokenType = "ELSE"
	FOR   TokenType = "FOR"
	WHILE TokenType = "WHILE"
	BREAK TokenType = "BREAK"

	// Functions and methods.
	FUN    TokenType = "FUN"
	RETURN TokenType = "RETURN"

	// Classes and inheritance.
	CLASS TokenType = "CLASS"
	SUPER TokenType = "SUPER"
	THIS  TokenType = "THIS"

	// Output.
	PRINT TokenType = "PRINT"

	// End of file.
	EOF TokenType = "EOF"
)
