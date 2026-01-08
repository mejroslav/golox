package golox

import "fmt"

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal any
	File    string
	Line    int
}

func NewToken(tokenType TokenType, lexeme string, literal any, file string, line int) Token {
	return Token{
		Type:    tokenType,
		Lexeme:  lexeme,
		Literal: literal,
		File:    file,
		Line:    line,
	}
}

func (t Token) String() string {
	if t.Type == EOF {
		return fmt.Sprintf("%s:%d: %s", t.File, t.Line, t.Type)
	}
	return fmt.Sprintf("%s:%d: %s %s %v", t.File, t.Line, t.Type, t.Lexeme, t.Literal)
}
