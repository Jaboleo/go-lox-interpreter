package main

import (
	"fmt"

	"github.com/kljablon/golox/ast"
)

type LoxCallable interface {
	arity() int
	call(interpreter Interpreter, arguments []any) any
}

type LoxFunction struct {
	declaration ast.Stmt_Function
	closure     Environment
}

func (l *LoxFunction) call(interpreter Interpreter, arguments []any) (result any) {
	environment := NewEnvironmentWithEnclosing(&l.closure)
	for i, _ := range l.declaration.Params {
		if arg, ok := arguments[i].(ast.Expr); ok {
			arguments[i] = interpreter.evaluate(arg)
		}
		environment.define(l.declaration.Params[i].Lexeme, arguments[i])
	}

	defer func() {
		if r := recover(); r != nil {
			if r, ok := r.(Return); ok {
				result = r.value
			}
		}
	}()

	interpreter.executeBlock(l.declaration.Body, &environment)
	return result
}

func (l *LoxFunction) arity() int {
	return len(l.declaration.Params)
}

func (l *LoxFunction) toString() string {
	return fmt.Sprintf("<fn %s >", l.declaration.Name.Lexeme)
}
