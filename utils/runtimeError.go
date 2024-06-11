package utils

import (
	"fmt"

	"github.com/kljablon/golox/ast"
)

type RuntimeError struct {
	Token   ast.Token
	Message string
}

// Error method to implement the error interface.
func (r RuntimeError) Error() string {
	return fmt.Sprintf("[line %d] Error at '%s': %s", r.Token.Line, r.Token.Lexeme, r.Message)
}

// Constructor-like function to create a new RuntimeError.
func NewRuntimeError(token ast.Token, message string) RuntimeError {
	return RuntimeError{token, message}
}
