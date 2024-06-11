package ast

import "fmt"

type Token struct {
	TokenType TokenType
	Lexeme    string
	Literal   any
	Line      int
}

func (t Token) ToString() string {
	return fmt.Sprintf("TYPE: %s, LITERAL: %v", getTokenName(t.TokenType), t.Literal)
}

var tokenNames = []string{
	"LEFT_PAREN", "RIGHT_PAREN", "LEFT_BRACE", "RIGHT_BRACE",
	"COMMA", "DOT", "MINUS", "PLUS", "SEMICOLON", "SLASH", "STAR",
	"BANG", "BANG_EQUAL", "EQUAL", "EQUAL_EQUAL", "GREATER", "GREATER_EQUAL",
	"LESS", "LESS_EQUAL", "IDENTIFIER", "STRING", "NUMBER", "AND", "CLASS",
	"ELSE", "FALSE", "FUN", "FOR", "IF", "NIL", "OR", "PRINT", "RETURN",
	"SUPER", "THIS", "TRUE", "VAR", "WHILE", "EOF",
}

func getTokenName(tokenType TokenType) string {
	index := int(tokenType)
	if index >= 0 && index < len(tokenNames) {
		return tokenNames[index]
	}
	return "UNKNOWN"
}
