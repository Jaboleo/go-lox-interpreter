package interpret

import (
	"log"

	"github.com/kljablon/golox/ast"
	"github.com/kljablon/golox/utils"
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

func (e *Environment) ancestor(distance int) Environment {
	environment := *e
	for i := 0; i < distance; i++ {
		environment = *environment.enclosing
	}
	return environment
}

func (e *Environment) getAt(distance int, name string) any {
	if value, ok := e.ancestor(distance).values[name]; ok {
		return value
	} else {
		return nil
	}
}

func (e *Environment) assignAt(distance int, name ast.Token, value any) {
	e.ancestor(distance).values[name.Lexeme] = value
}

func (e *Environment) get(name ast.Token) any {
	if v, ok := e.values[name.Lexeme]; ok {
		return v
	}
	if e.enclosing != nil && e.enclosing != e {
		return e.enclosing.get(name)
	}

	err := utils.NewRuntimeError(name, "Undefined variable '"+name.Lexeme+"'.")
	// Print the error message.
	log.Fatal(err.Error())
	return nil // This return is not necessary because log.Fatal exits the program, but added for clarity.
}

func (e *Environment) assign(name ast.Token, value any) {
	if _, ok := e.values[name.Lexeme]; ok {
		e.values[name.Lexeme] = value
		return
	}

	if e.enclosing != nil {
		e.enclosing.assign(name, value)
		return
	}

	err := utils.NewRuntimeError(name, "Undefined variable '"+name.Lexeme+"'.")
	// Print the error message.
	log.Fatal(err.Error())
}
