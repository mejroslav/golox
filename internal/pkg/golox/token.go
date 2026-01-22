package golox

import "fmt"

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal any
	File    string
	Line    int
	Column  int
}

func NewToken(tokenType TokenType, lexeme string, literal any, file string, line int, column int) Token {
	return Token{
		Type:    tokenType,
		Lexeme:  lexeme,
		Literal: literal,
		File:    file,
		Line:    line,
		Column:  column,
	}
}

func (t Token) String() string {
	if t.Type == EOF {
		return fmt.Sprintf("%s:%d:%d: %s", t.File, t.Line, t.Column, t.Type)
	}
	return fmt.Sprintf("%s:%d:%d: %s %s %v", t.File, t.Line, t.Column, t.Type, t.Lexeme, t.Literal)
}
