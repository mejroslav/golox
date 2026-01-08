package golox

import "fmt"

// ScannerError reports an error encountered during scanning
type ScannerError struct {
	File    string
	Line    int
	Context string
	Message string
}

func (s ScannerError) Error() string {
	return fmt.Sprintf("SCANNER ERROR [%s:%d] %s\n%s\n", s.File, s.Line, s.Message, s.Context)
}

type ParserError struct {
	Token   Token
	Message string
}

// ParserError reports an error encountered during parsing
func (p ParserError) Error() string {
	if p.Token.Type == EOF {
		return parserErrorMsg(p.Token.File, p.Token.Line, "at end", p.Message)
	} else {
		return parserErrorMsg(p.Token.File, p.Token.Line, "at '"+p.Token.Lexeme+"'", p.Message)
	}
}

func parserErrorMsg(file string, line int, where string, message string) string {
	return fmt.Sprintf("PARSER ERROR [%s:%d] %s: %s\n", file, line, where, message)
}
