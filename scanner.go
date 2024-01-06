package main

import (
	"strconv"
	"unicode"
)

type Scanner struct {
	source   string
	tokens   []Token
	start    int
	current  int
	line     int
	keywords map[string]TokenType
}

func NewScanner(source string) Scanner {

	keywords := map[string]TokenType{
		"and":    AND,
		"class":  CLASS,
		"else":   ELSE,
		"false":  FALSE,
		"for":    FOR,
		"fun":    FUN,
		"if":     IF,
		"nil":    NIL,
		"or":     OR,
		"print":  PRINT,
		"return": RETURN,
		"super":  SUPER,
		"this":   THIS,
		"true":   TRUE,
		"var":    VAR,
		"while":  WHILE,
	}

	return Scanner{source, make([]Token, 0), 0, 0, 1, keywords}
}

func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.tokens = append(s.tokens, Token{EOF, "", nil, s.line})
	return s.tokens
}

func (s Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanToken() {
	character := s.advance()
	switch character {
	case '(':
		s.addToken(LEFT_PAREN, nil)
	case ')':
		s.addToken(RIGHT_PAREN, nil)
	case '{':
		s.addToken(LEFT_BRACE, nil)
	case '}':
		s.addToken(RIGHT_BRACE, nil)
	case ',':
		s.addToken(COMMA, nil)
	case '.':
		s.addToken(DOT, nil)
	case '-':
		s.addToken(MINUS, nil)
	case '+':
		s.addToken(PLUS, nil)
	case ';':
		s.addToken(SEMICOLON, nil)
	case '*':
		s.addToken(STAR, nil)
	case '!':
		if s.match('=') {
			s.addToken(BANG_EQUAL, nil)
		} else {
			s.addToken(BANG, nil)
		}
	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL, nil)
		} else {
			s.addToken(EQUAL, nil)
		}
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL, nil)
		} else {
			s.addToken(LESS, nil)
		}
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL, nil)
		} else {
			s.addToken(GREATER, nil)
		}
	case '/':
		if s.match('/') {
			// A comment goes until the end of the line.
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else if s.match('*') {
			s.advance()
			for s.peek() != '*' && s.peekNext() != '/' && !s.isAtEnd() {
				if s.peek() == '\n' {
					s.line++
				}
				s.advance()
			}
			s.advance()
			s.advance()
		} else {
			s.addToken(SLASH, nil)
		}
	case ' ':
	case '\r':
	case '\t':
		// Ignore whitespace.
	case '\n':
		// Increment line number.
		s.line++
	case '"':
		s.stringLit()
	default:
		if unicode.IsDigit(character) {
			s.number()
		} else if unicode.IsLetter(character) || unicode.IsDigit(character) {
			s.identifier()
		} else {
			ReportError(s.line, "Unexpected character.")
		}

	}
}

func (s *Scanner) addToken(ttype TokenType, literal interface{}) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, Token{ttype, text, literal, s.line})
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != byte(expected) {
		return false
	}
	s.current++
	return true
}

// Return current character and proceed
func (s *Scanner) advance() rune {
	char := s.source[s.current]
	s.current++
	return rune(char)
}

// Return current character without proceeding
func (s Scanner) peek() rune {
	if s.isAtEnd() {
		return rune(0)
	}
	return rune(s.source[s.current])
}

// Return next character
func (s Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return '0'
	}
	return rune(s.source[s.current+1])
}

func (s *Scanner) stringLit() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.isAtEnd() {
		ReportError(s.line, "Unterminated string.")
		return
	}
	s.advance()

	// Extract the string value
	value := s.source[s.start+1 : s.current-1]
	s.addToken(STRING, value)
}

func (s *Scanner) number() {
	for unicode.IsDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && unicode.IsDigit(s.peekNext()) {
		s.advance()
		for unicode.IsDigit(s.peek()) {
			s.advance()
		}
	}
	v, err := strconv.ParseFloat(s.source[s.start:s.current], 32)
	if err != nil {
		ReportError(s.line, err.Error())
	}
	s.addToken(NUMBER, v)
}

func (s *Scanner) identifier() {
	for unicode.IsLetter(s.peek()) || unicode.IsDigit(s.peek()) {
		s.advance()
	}
	text := s.source[s.start:s.current]
	type_, ok := s.keywords[text]
	if !ok {
		type_ = IDENTIFIER
	}
	s.addToken(type_, nil)
}
