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
	return string(t.Type) + " " + t.Lexeme + " " + fmt.Sprintf("%v", t.Literal)
}
