package golox

import "strconv"

type Scanner struct {
	source  string
	tokens  []Token
	start   int
	current int
	line    int
}

func NewCodeScanner() *Scanner {
	return &Scanner{}
}

func (s *Scanner) Run() bool {
	var codeScanner Scanner = *NewCodeScanner()
	codeScanner.ScanTokens("dummy source")
	return false
}

func (s *Scanner) ScanTokens(source string) []Token {
	s.source = source
	eof := Token{Type: EOF, Lexeme: "", Literal: nil, Line: 0}
	s.tokens = append(s.tokens, eof)
	return s.tokens
}

func (s *Scanner) isAtAnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanToken() {
	var c rune = s.advance()
	switch c {

	// Single-character tokens.
	case '(':
		s.addToken(LEFT_PAREN)
	case ')':
		s.addToken(RIGHT_PAREN)
	case '{':
		s.addToken(LEFT_BRACE)
	case '}':
		s.addToken(RIGHT_BRACE)
	case ',':
		s.addToken(COMMA)
	case '.':
		s.addToken(DOT)
	case '-':
		s.addToken(MINUS)
	case '+':
		s.addToken(PLUS)
	case ';':
		s.addToken(SEMICOLON)
	case '*':
		s.addToken(STAR)

	// Operators with one or two characters.
	case '!':
		if s.match('=') {
			s.addToken(BANG_EQUAL)
		} else {
			s.addToken(BANG)
		}
	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL)
		} else {
			s.addToken(EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL)
		} else {
			s.addToken(LESS)
		}
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL)
		} else {
			s.addToken(GREATER)
		}

	case '/':
		if s.match('/') {
			// A comment goes until the end of the line.
			for s.peek() != '\n' && !s.isAtAnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH)
		}

	// String literals.
	case '"':
		s.string()

	// Ignore whitespace.
	case ' ', '\r', '\t':
		// Do nothing.

	// New line.
	case '\n':
		s.line++

	// Unexpected character.
	default:
		if s.isDigit(c) {
			s.number()
		} else {
			Error(s.line, "Unexpected character.")
		}

	}
}

func (s *Scanner) addToken(tokenType TokenType) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, NewToken(tokenType, text, nil, s.line))
}

func (s *Scanner) addTokenWithValue(tokenType TokenType, value any) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, NewToken(tokenType, text, value, s.line))
}

// string handles string literals, consuming characters until the closing quote is found.
//
// It also supports multi-line strings.
func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtAnd() {
		if s.peek() == '\n' {
			// Strings can span multiple lines, so we need to increment the line counter.
			s.line++
		}
		s.advance()
	}

	if s.isAtAnd() {
		Error(s.line, "Unterminated string.")
		return
	} else {
		// The closing ".
		s.advance()
	}

	// Trim the surrounding quotes.
	value := s.source[s.start+1 : s.current-1]
	s.addTokenWithValue(STRING, value)
}

// number handles numeric literals. It supports both integer and floating-point numbers.
func (s *Scanner) number() {
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
		Error(s.line, "Invalid number format.")
		return
	}
	s.addTokenWithValue(NUMBER, value)
}

func (s *Scanner) isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) advance() rune {
	c := s.charAt(s.current)
	s.current++
	return c
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtAnd() {
		return false
	}
	if s.charAt(s.current) != expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) peek() rune {
	if s.isAtAnd() {
		return '\000'
	} else {
		return s.charAt(s.current)
	}
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return '\000'
	} else {
		return s.charAt(s.current + 1)
	}
}

func (s *Scanner) charAt(index int) rune {
	return rune(s.source[index])
}
