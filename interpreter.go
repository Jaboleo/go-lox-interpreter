package main

import (
	"fmt"
	"log"
)

type Return struct {
	value any
}

type Interpreter struct {
	globals     Environment
	environment Environment
}

func newInterpreter() Interpreter {
	globals := NewEnvironment()

	globals.define("clock", ClockFunc{})

	return Interpreter{
		globals:     globals,
		environment: globals,
	}
}

func (i *Interpreter) interpret(statements []Stmt) {
	for _, statement := range statements {
		i.execute(statement)
	}
}

func (i *Interpreter) visitExpr_Binary(e Expr_Binary) any {
	left := i.evaluate(e.left)
	right := i.evaluate(e.right)

	switch e.operator.ttype {
	case MINUS:
		i.checkNumberOperand(e.operator, right)
		return castToFloat(left) - castToFloat(right)
	case PLUS:
		if valLeft, okLeft := left.(float64); okLeft {
			if valRight, okRight := right.(float64); okRight {
				return valLeft + valRight
			}
		}
		if valLeft, okLeft := left.(string); okLeft {
			if valRight, okRight := right.(string); okRight {
				return valLeft + valRight
			}
		}
		err := NewRuntimeError(e.operator, "Operands must be two numbers or two strings.")
		// Print the error message.
		fmt.Println(err.Error())

	case SLASH:
		i.checkNumberOperands(e.operator, left, right)
		return castToFloat(left) / castToFloat(right)
	case STAR:
		i.checkNumberOperands(e.operator, left, right)
		return castToFloat(left) * castToFloat(right)
	case GREATER:
		i.checkNumberOperands(e.operator, left, right)
		return castToFloat(left) > castToFloat(right)
	case GREATER_EQUAL:
		i.checkNumberOperands(e.operator, left, right)
		return castToFloat(left) >= castToFloat(right)
	case LESS:
		i.checkNumberOperands(e.operator, left, right)
		return castToFloat(left) < castToFloat(right)
	case LESS_EQUAL:
		i.checkNumberOperands(e.operator, left, right)
		return castToFloat(left) <= castToFloat(right)
	case BANG_EQUAL:
		i.checkNumberOperands(e.operator, left, right)
		return !isEqual(left, right)
	case EQUAL_EQUAL:
		i.checkNumberOperands(e.operator, left, right)
		return isEqual(left, right)
	default:
		return nil
	}
	return nil
}

func (i *Interpreter) visitExpr_Call(e Expr_Call) any {
	callee := i.evaluate(e.callee)
	arguments := []any{}
	for _, argument := range e.arguments {
		arguments = append(arguments, argument)
	}

	if function, ok := callee.(LoxFunction); ok {
		if len(arguments) != function.arity() {
			log.Fatalf("Expected %d arguments but got %d.", function.arity(), len(arguments))
		}
		return function.call(*i, arguments)
	}
	log.Fatal("not a LoxCallable")
	return nil
}

func (i *Interpreter) visitExpr_Grouping(e Expr_Grouping) any {
	return i.evaluate(e.expression)
}

func (i *Interpreter) visitExpr_Literal(e Expr_Literal) any {
	return e.value
}

func (i *Interpreter) visitExpr_Unary(e Expr_Unary) any {
	right := i.evaluate(e.right)
	switch e.operator.ttype {
	case BANG:
		return !isTruthy(right)
	case MINUS:
		return -castToFloat(right)
	default:
		return nil
	}
}

func (i *Interpreter) checkNumberOperand(operator Token, operand any) {
	if _, ok := operand.(float64); ok {
		return
	}
	err := NewRuntimeError(operator, "Operand must be a number.")
	// Print the error message.
	log.Fatal(err.Error())

}

func (i *Interpreter) checkNumberOperands(operator Token, left any, right any) {
	if _, ok := left.(float64); ok {
		if _, ok := right.(float64); ok {
			return
		}
	}
	err := NewRuntimeError(operator, "Operands must be numbers.")
	// Print the error message.
	fmt.Println(err.Error())

}

// func (i *Interpreter) castToString(object any) string {
// 	if object == nil {
// 		return "nil"
// 	}
// 	if val, ok := object.(float64); ok {
// 		text := fmt.Sprintf("%f", val)
// 		if strings.HasSuffix(text, ".000000") {
// 			return text[:len(text)-7]
// 		}
// 		return text
// 	}

// 	if val, ok := object.(Token); ok {
// 		return val.ToString()
// 	}

// 	fmt.Println("return nil")
// 	return "nil"
// }

func (i *Interpreter) evaluate(expr Expr) any {
	return expr.accept(i)
}

func (i *Interpreter) execute(stmt Stmt) {
	if stmt == nil {
		panic("nil statement at execute()")
	}
	stmt.accept(i)
}

func (i *Interpreter) visitStmt_Block(stmt Stmt_Block) {
	enclosing_env := i.environment
	new_env := NewEnvironmentWithEnclosing(&enclosing_env)
	i.executeBlock(stmt.statements, &new_env)
}

func (i *Interpreter) executeBlock(statements []Stmt, environment *Environment) {
	previous := i.environment
	i.environment = *environment
	defer func() { i.environment = previous }()
	for _, statement := range statements {
		// if statement != nil {
		// 	i.execute(statement)
		// }
		i.execute(statement)
	}
}

func (i *Interpreter) visitStmt_Expression(stmt Stmt_Expression) {
	i.evaluate(stmt.expression)
}

func (i *Interpreter) visitStmt_Function(stmt Stmt_Function) {
	function := LoxFunction{stmt, i.environment}
	i.environment.define(stmt.name.lexeme, function)
}

func (i *Interpreter) visitStmt_If(stmt Stmt_If) {
	if isTruthy(i.evaluate(stmt.condition)) {
		i.execute(stmt.thenBranch)
	} else if stmt.elseBranch != nil {
		i.execute(stmt.elseBranch)
	}
}

func (i *Interpreter) visitStmt_Print(stmt Stmt_Print) {
	value := i.evaluate(stmt.expression)
	fmt.Println(value)
}

func (i *Interpreter) visitStmt_Return(stmt Stmt_Return) {
	var value any
	if stmt.value != nil {
		value = i.evaluate(stmt.value)
	}
	panic(Return{value})
}

func (i *Interpreter) visitStmt_While(stmt Stmt_While) {
	for isTruthy(i.evaluate(stmt.condition)) {
		i.execute(stmt.body)
	}
}

func (i *Interpreter) visitStmt_Var(stmt Stmt_Var) {
	var value any
	if stmt.initializer != nil {
		value = i.evaluate(stmt.initializer)
	}
	i.environment.define(stmt.name.lexeme, value)
}

func (i *Interpreter) visitExpr_Assign(expr Expr_Assign) any {
	value := i.evaluate(expr.value)
	i.environment.assign(expr.name, value)
	return value
}

func (i *Interpreter) visitExpr_Variable(expr Expr_Variable) any {
	return i.environment.get(expr.name)
}

func (i *Interpreter) visitExpr_Logical(expr Expr_Logical) any {
	left := i.evaluate(expr.left)

	if expr.operator.ttype == OR {
		if isTruthy(left) {
			return left
		} else {
			if !isTruthy(left) {
				return left
			}
		}
	}
	return i.evaluate(expr.right)
}
