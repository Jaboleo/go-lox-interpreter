package main

import (
	"fmt"

	"github.com/kljablon/golox/ast"
)

type RuntimeError struct {
	token   ast.Token
	message string
}

// Error method to implement the error interface.
func (r RuntimeError) Error() string {
	return fmt.Sprintf("[line %d] Error at '%s': %s", r.token.Line, r.token.Lexeme, r.message)
}

// Constructor-like function to create a new RuntimeError.
func NewRuntimeError(token ast.Token, message string) RuntimeError {
	return RuntimeError{token, message}
}
