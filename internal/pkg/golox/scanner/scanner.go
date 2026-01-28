package scanner

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/mejroslav/golox/internal/pkg/golox/lox_error"
	"github.com/mejroslav/golox/internal/pkg/golox/token"
)

type CodeScanner struct {
	source  string
	tokens  []token.Token
	start   int
	current int
	line    int
	column  int
	file    string
}

func NewCodeScanner(line int, file string) *CodeScanner {
	return &CodeScanner{line: line, file: file}
}

// Run scans the provided source code and prints the tokens. It returns true if scanning was successful, false otherwise.
// If verbose is true, it prints each token to stdout.
func (s *CodeScanner) Run(source string) ([]token.Token, bool) {
	s.source = source
	s.start = 0
	s.current = 0
	s.tokens = []token.Token{}
	slog.Debug("Starting scan", "file", s.file, "length", len(s.source))
	return s.ScanTokens()
}

func (s *CodeScanner) ScanTokens() ([]token.Token, bool) {
	hadError := false
	for !s.isAtEnd() {
		// We are at the beginning of the next lexeme.
		s.start = s.current
		s.scanToken()
	}

	s.addToken(token.EOF)

	return s.tokens, hadError
}

func (s *CodeScanner) scanToken() {
	var c rune = s.advance()
	switch c {

	// Single-character tokens.
	case '(':
		s.addToken(token.LEFT_PAREN)
	case ')':
		s.addToken(token.RIGHT_PAREN)
	case '{':
		s.addToken(token.LEFT_BRACE)
	case '}':
		s.addToken(token.RIGHT_BRACE)
	case ',':
		s.addToken(token.COMMA)
	case '.':
		s.addToken(token.DOT)
	case '-':
		s.addToken(token.MINUS)
	case '+':
		s.addToken(token.PLUS)
	case ';':
		s.addToken(token.SEMICOLON)
	case '*':
		s.addToken(token.STAR)

	// Operators with one or two characters.
	case '!':
		if s.match('=') {
			s.addToken(token.BANG_EQUAL)
		} else {
			s.addToken(token.BANG)
		}
	case '=':
		if s.match('=') {
			s.addToken(token.EQUAL_EQUAL)
		} else {
			s.addToken(token.EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addToken(token.LESS_EQUAL)
		} else {
			s.addToken(token.LESS)
		}
	case '>':
		if s.match('=') {
			s.addToken(token.GREATER_EQUAL)
		} else {
			s.addToken(token.GREATER)
		}

	case '/':
		if s.match('/') {
			// A comment goes until the end of the line.
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(token.SLASH)
		}

	// String literals.
	case '"':
		s.string()

	// Ignore whitespace.
	case ' ', '\r', '\t':
		// Do nothing.

	// New line.
	case '\n':
		s.newLine()

	default:
		// Here we use the principle of *maximal munch*, trying to consume as many characters as possible.
		if s.isDigit(c) {
			// Number literals.
			s.number()
		} else if s.isAlpha(c) {
			// Identifiers and keywords.
			s.identifier()
		} else {
			// Unexpected character.
			err := lox_error.ScannerError{
				File:   s.file,
				Line:   s.line,
				Column: s.column,

				Context: s.getContextLines(),
				Message: fmt.Sprintf("Unexpected character '%c'.", c),
			}
			fmt.Println(err.Error())
		}
	}
}

func (s *CodeScanner) addToken(tokenType token.TokenType) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, token.NewToken(tokenType, text, nil, s.file, s.line, s.column))
}

func (s *CodeScanner) addTokenWithValue(tokenType token.TokenType, value any) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, token.NewToken(tokenType, text, value, s.file, s.line, s.column))
}

// string handles string literals, consuming characters until the closing quote is found.
//
// It also supports multi-line strings.
func (s *CodeScanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			// Strings can span multiple lines, so we need to increment the line counter.
			s.newLine()
		}
		s.advance()
	}

	if s.isAtEnd() {
		err := lox_error.ScannerError{
			File:    s.file,
			Line:    s.line,
			Context: s.getContextLines(),
			Message: "Unterminated string.",
		}
		fmt.Println(err.Error())
		return
	} else {
		// The closing ".
		s.advance()
	}

	// Trim the surrounding quotes.
	value := s.source[s.start+1 : s.current-1]
	s.addTokenWithValue(token.STRING, value)
}

// number handles numeric literals. It supports both integer and floating-point numbers.
func (s *CodeScanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		// Consume the "."
		s.advance()

		for s.isDigit(s.peek()) {
			s.advance()
		}

	}

	value, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
	if err != nil {
		err := lox_error.ScannerError{
			File:    s.file,
			Line:    s.line,
			Context: s.getContextLines(),
			Message: "Invalid number format: " + err.Error(),
		}
		fmt.Println(err.Error())
		return
	}
	s.addTokenWithValue(token.NUMBER, value)
}

// identifier handles identifiers and keywords.
func (s *CodeScanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	tokenType, exists := token.Keywords[text]
	if !exists {
		tokenType = token.IDENTIFIER
	}

	s.addToken(tokenType)
}

func (s *CodeScanner) isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func (s *CodeScanner) isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func (s *CodeScanner) isAlphaNumeric(c rune) bool {
	return s.isAlpha(c) || s.isDigit(c)
}

// advance consumes the current character and returns it.
func (s *CodeScanner) advance() rune {
	c := s.charAt(s.current)
	s.current++
	s.column++
	return c
}

// match checks if the current character matches the expected character.
// If it does, it consumes the character and returns true. Otherwise, it returns false.
func (s *CodeScanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	if s.charAt(s.current) != expected {
		return false
	}
	s.current++
	s.column++
	return true
}

// peek returns the current character without consuming it.
func (s *CodeScanner) peek() rune {
	if s.isAtEnd() {
		return '\000'
	} else {
		return s.charAt(s.current)
	}
}

// peekNext returns the character after the current one without consuming it.
func (s *CodeScanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return '\000'
	} else {
		return s.charAt(s.current + 1)
	}
}

// charAt returns the character at the specified index.
func (s *CodeScanner) charAt(index int) rune {
	return rune(s.source[index])
}

func (s *CodeScanner) newLine() {
	s.line++
	s.column = 0
}

func (s *CodeScanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

// getContextLines returns a string with the lines around the current line for error reporting.
func (s *CodeScanner) getContextLines() string {

	current := s.current
	output := ""
	linesBeforeError := 1
	linesAfterError := 2
	charsBeforeError := 0
	charsAfterError := -1

	// Get lines before the current line.
	lines := 0
	for i := current - 1; i >= 0; i-- {
		c := s.charAt(i)
		if c == '\n' {
			lines++
			if lines > linesBeforeError {
				break
			}
		}
		if lines == 0 {
			charsBeforeError++
		}
		output = string(c) + output
	}

	// Get lines after the current line.
	markerPlaced := false
	lines = 0
	for i := current; i < len(s.source); i++ {
		if s.isAtEnd() {
			break
		}
		c := s.charAt(i)
		if c == '\n' {
			lines++
			if lines == 1 {
				output += marker(charsBeforeError, charsAfterError)
				markerPlaced = true
				charsAfterError = 0
			}
			if lines > linesAfterError {
				break
			}
		} else {
			charsAfterError++
		}
		output += string(c)
	}

	if !markerPlaced {
		output += marker(charsBeforeError, charsAfterError)
	}

	return output
}

func marker(before int, after int) string {
	before = max(before, 1)
	s := make([]rune, before+after+1)
	for i := 0; i < before-1; i++ {
		s[i] = '-'
	}
	s[before-1] = '^'
	for i := before; i < len(s); i++ {
		s[i] = '-'
	}
	return "\n" + string(s)
}
