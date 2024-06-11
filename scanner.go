package main

import (
	"strconv"
	"unicode"

	"github.com/kljablon/golox/ast"
)

type Scanner struct {
	source   string
	tokens   []ast.Token
	start    int
	current  int
	line     int
	keywords map[string]ast.TokenType
}

func NewScanner(source string) Scanner {

	keywords := map[string]ast.TokenType{
		"and":    ast.AND,
		"class":  ast.CLASS,
		"else":   ast.ELSE,
		"false":  ast.FALSE,
		"for":    ast.FOR,
		"fun":    ast.FUN,
		"if":     ast.IF,
		"nil":    ast.NIL,
		"or":     ast.OR,
		"print":  ast.PRINT,
		"return": ast.RETURN,
		"super":  ast.SUPER,
		"this":   ast.THIS,
		"true":   ast.TRUE,
		"var":    ast.VAR,
		"while":  ast.WHILE,
	}

	return Scanner{source, make([]ast.Token, 0), 0, 0, 1, keywords}
}

func (s *Scanner) ScanTokens() []ast.Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.tokens = append(s.tokens, ast.Token{
		TokenType: ast.EOF,
		Lexeme:    "",
		Literal:   nil,
		Line:      s.line})
	return s.tokens
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanToken() {
	character := s.advance()
	switch character {
	case '(':
		s.addToken(ast.LEFT_PAREN, nil)
	case ')':
		s.addToken(ast.RIGHT_PAREN, nil)
	case '{':
		s.addToken(ast.LEFT_BRACE, nil)
	case '}':
		s.addToken(ast.RIGHT_BRACE, nil)
	case ',':
		s.addToken(ast.COMMA, nil)
	case '.':
		s.addToken(ast.DOT, nil)
	case '-':
		s.addToken(ast.MINUS, nil)
	case '+':
		s.addToken(ast.PLUS, nil)
	case ';':
		s.addToken(ast.SEMICOLON, nil)
	case '*':
		s.addToken(ast.STAR, nil)
	case '!':
		if s.match('=') {
			s.addToken(ast.BANG_EQUAL, nil)
		} else {
			s.addToken(ast.BANG, nil)
		}
	case '=':
		if s.match('=') {
			s.addToken(ast.EQUAL_EQUAL, nil)
		} else {
			s.addToken(ast.EQUAL, nil)
		}
	case '<':
		if s.match('=') {
			s.addToken(ast.LESS_EQUAL, nil)
		} else {
			s.addToken(ast.LESS, nil)
		}
	case '>':
		if s.match('=') {
			s.addToken(ast.GREATER_EQUAL, nil)
		} else {
			s.addToken(ast.GREATER, nil)
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
			s.addToken(ast.SLASH, nil)
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

func (s *Scanner) addToken(ttype ast.TokenType, literal interface{}) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, ast.Token{TokenType: ttype, Lexeme: text, Literal: literal, Line: s.line})
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
func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return rune(0)
	}
	return rune(s.source[s.current])
}

// Return next character
func (s *Scanner) peekNext() rune {
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
	s.addToken(ast.STRING, value)
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
	s.addToken(ast.NUMBER, v)
}

func (s *Scanner) identifier() {
	for unicode.IsLetter(s.peek()) || unicode.IsDigit(s.peek()) {
		s.advance()
	}
	text := s.source[s.start:s.current]
	ttype, ok := s.keywords[text]
	if !ok {
		ttype = ast.IDENTIFIER
	}
	s.addToken(ttype, nil)
}
