package golox

import "fmt"

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal any
	Line    int
}

func NewToken(tokenType TokenType, lexeme string, literal any, line int) Token {
	return Token{
		Type:    tokenType,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    line,
	}
}

func (t Token) String() string {
	if t.Type == EOF {
		return fmt.Sprintf("%d: %s", t.Line, t.Type)
	}
	return fmt.Sprintf("%d: %s %s %v", t.Line, t.Type, t.Lexeme, t.Literal)
}
