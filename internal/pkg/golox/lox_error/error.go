package lox_error

import (
	"fmt"
	"mejroslav/golox/v2/internal/pkg/golox/token"
)

// ScannerError reports an error encountered during scanning
type ScannerError struct {
	File    string
	Line    int
	Column  int
	Context string
	Message string
}

func (s ScannerError) Error() string {
	return fmt.Sprintf("SCANNER ERROR [%s:%d:%d] %s\n%s\n", s.File, s.Line, s.Column, s.Message, s.Context)
}

type ParserError struct {
	Token   token.Token
	Message string
}

// ParserError reports an error encountered during parsing
func (p ParserError) Error() string {
	if p.Token.Type == token.EOF {
		return parserErrorMsg(p.Token.File, p.Token.Line, p.Token.Column, "at end", p.Message)
	} else {
		return parserErrorMsg(p.Token.File, p.Token.Line, p.Token.Column, "at '"+p.Token.Lexeme+"'", p.Message)
	}
}

func parserErrorMsg(file string, line int, column int, where string, message string) string {
	return fmt.Sprintf("PARSER ERROR [%s:%d:%d] %s: %s\n", file, line, column, where, message)
}

type RuntimeError struct {
	Token   token.Token
	Message string
}

func NewRuntimeError(token token.Token, message string) RuntimeError {
	return RuntimeError{
		Token:   token,
		Message: message,
	}
}

func (r RuntimeError) Error() string {
	return fmt.Sprintf("RUNTIME ERROR [%s:%d:%d] %s\n", r.Token.File, r.Token.Line, r.Token.Column, r.Message)
}
