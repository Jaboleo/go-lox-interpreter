package main

import (
	"log"
)

type Environment struct {
	enclosing *Environment
	values    map[string]any
}

func NewEnvironment() Environment {
	new_map := make(map[string]any)
	return Environment{
		enclosing: nil,
		values:    new_map,
	}
}

func NewEnvironmentWithEnclosing(enclosing *Environment) Environment {
	new_map := make(map[string]any)
	return Environment{
		enclosing: enclosing,
		values:    new_map,
	}
}

func (e *Environment) define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) get(name Token) any {
	if v, ok := e.values[name.lexeme]; ok {
		return v
	}
	if e.enclosing != nil && e.enclosing != e {
		return e.enclosing.get(name)
	}

	err := NewRuntimeError(name, "Undefined variable '"+name.lexeme+"'.")
	// Print the error message.
	log.Fatal(err.Error())
	return nil // This return is not necessary because log.Fatal exits the program, but added for clarity.
}

func (e *Environment) assign(name Token, value any) {
	if _, ok := e.values[name.lexeme]; ok {
		e.values[name.lexeme] = value
		return
	}

	if e.enclosing != nil {
		e.enclosing.assign(name, value)
		return
	}

	err := NewRuntimeError(name, "Undefined variable '"+name.lexeme+"'.")
	// Print the error message.
	log.Fatal(err.Error())
}
