package main

import "fmt"

type LoxCallable interface {
	arity() int
	call(interpreter Interpreter, arguments []any) any
}

type LoxFunction struct {
	declaration Stmt_Function
	closure     Environment
}

func (l *LoxFunction) call(interpreter Interpreter, arguments []any) (result any) {
	environment := NewEnvironmentWithEnclosing(&l.closure)
	for i, _ := range l.declaration.params {
		if arg, ok := arguments[i].(Expr); ok {
			arguments[i] = interpreter.evaluate(arg)
		}
		environment.define(l.declaration.params[i].lexeme, arguments[i])
	}

	defer func() {
		if r := recover(); r != nil {
			if r, ok := r.(Return); ok {
				result = r.value
			}
		}
	}()

	interpreter.executeBlock(l.declaration.body, &environment)
	return result
}

func (l *LoxFunction) arity() int {
	return len(l.declaration.params)
}

func (l *LoxFunction) toString() string {
	return fmt.Sprintf("<fn %w >", l.declaration.name.lexeme)
}
