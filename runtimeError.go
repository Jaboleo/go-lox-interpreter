package main

import "fmt"

type RuntimeError struct {
	token   Token
	message string
}

// Error method to implement the error interface.
func (r RuntimeError) Error() string {
	return fmt.Sprintf("[line %d] Error at '%s': %s", r.token.line, r.token.lexeme, r.message)
}

// Constructor-like function to create a new RuntimeError.
func NewRuntimeError(token Token, message string) RuntimeError {
	return RuntimeError{token, message}
}
