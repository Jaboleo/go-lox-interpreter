package main

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
)

type Interpreter struct {
	environment Environment
}

func newInterpreter() Interpreter {
	return Interpreter{
		environment: NewEnvironment(),
	}
}

func (i *Interpreter) interpret(statements []Stmt) {
	for _, statement := range statements {
		i.execute(statement)
	}
}

func (i *Interpreter) visitExpr_Binary(e Expr_Binary) interface{} {
	left := i.evaluate(e.left)
	right := i.evaluate(e.right)

	switch e.operator.ttype {
	case MINUS:
		i.checkNumberOperand(e.operator, right)
		return i.castToFloat(left) - i.castToFloat(right)
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
		return i.castToFloat(left) / i.castToFloat(right)
	case STAR:
		i.checkNumberOperands(e.operator, left, right)
		return i.castToFloat(left) * i.castToFloat(right)
	case GREATER:
		i.checkNumberOperands(e.operator, left, right)
		return i.castToFloat(left) > i.castToFloat(right)
	case GREATER_EQUAL:
		i.checkNumberOperands(e.operator, left, right)
		return i.castToFloat(left) >= i.castToFloat(right)
	case LESS:
		i.checkNumberOperands(e.operator, left, right)
		return i.castToFloat(left) < i.castToFloat(right)
	case LESS_EQUAL:
		i.checkNumberOperands(e.operator, left, right)
		return i.castToFloat(left) <= i.castToFloat(right)
	case BANG_EQUAL:
		i.checkNumberOperands(e.operator, left, right)
		return !i.isEqual(left, right)
	case EQUAL_EQUAL:
		i.checkNumberOperands(e.operator, left, right)
		return i.isEqual(left, right)
	default:
		return nil
	}
	return nil
}

func (i *Interpreter) visitExpr_Call(e Expr_Call) interface{} {
	callee := i.evaluate(e.callee)
	arguments := []interface{}{}
	for _, argument := range e.arguments {
		arguments = append(arguments, argument)
	}

	if function, ok := callee.(LoxCallable); ok {
		if len(arguments) != function.arity() {
			log.Fatalf("Expected %w arguments but got %w.", function.arity(), len(arguments))
		}
		return function.call(i, arguments)
	}
	return nil
}

func (i *Interpreter) visitExpr_Grouping(e Expr_Grouping) interface{} {
	return i.evaluate(e.expression)
}

func (i *Interpreter) visitExpr_Literal(e Expr_Literal) interface{} {
	return e.value
}

func (i *Interpreter) visitExpr_Unary(e Expr_Unary) interface{} {
	right := i.evaluate(e.right)
	switch e.operator.ttype {
	case BANG:
		return !i.isTruthy(right)
	case MINUS:
		return -i.castToFloat(right)
	default:
		return nil
	}
}

func (i *Interpreter) isTruthy(object interface{}) bool {
	if object == nil {
		return false
	}
	switch v := object.(type) {
	case bool:
		return v
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(object).Int() != 0
	case uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(object).Uint() != 0
	case float32, float64:
		return reflect.ValueOf(object).Float() != 0
	case string:
		return v != ""
		// Default case, return true
	default:
		// For any other type, you can define your own truthy logic
		return true
	}
}

func (i *Interpreter) isEqual(a interface{}, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	return a == b
}

func (i *Interpreter) castToFloat(a interface{}) float64 {
	if val, ok := a.(float64); ok {
		// If `right` is already a float64, simply negate it and return
		return val
	}
	if val, ok := a.(int); ok {
		return float64(val)
	}
	if val, ok := a.(string); ok {
		// Convert the string to float64 using strconv.ParseFloat()
		floatValue, err := strconv.ParseFloat(val, 64)
		if err != nil {
			// Handle error if conversion fails
			// For simplicity, let's return 0 in case of error
			return 0
		}
		return floatValue
		// Return the negation of the float64 value
	}
	return 0
}

func (i *Interpreter) checkNumberOperand(operator Token, operand interface{}) {
	if _, ok := operand.(float64); ok {
		return
	}
	err := NewRuntimeError(operator, "Operand must be a number.")
	// Print the error message.
	log.Fatal(err.Error())

}

func (i *Interpreter) checkNumberOperands(operator Token, left interface{}, right interface{}) {
	if _, ok := left.(float64); ok {
		if _, ok := right.(float64); ok {
			return
		}
	}
	err := NewRuntimeError(operator, "Operands must be numbers.")
	// Print the error message.
	fmt.Println(err.Error())

}

// func (i *Interpreter) castToString(object interface{}) string {
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

func (i *Interpreter) evaluate(expr Expr) interface{} {
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

func (i *Interpreter) visitStmt_If(stmt Stmt_If) {
	if i.isTruthy(i.evaluate(stmt.condition)) {
		i.execute(stmt.thenBranch)
	} else if stmt.elseBranch != nil {
		i.execute(stmt.elseBranch)
	}
}

func (i *Interpreter) visitStmt_Print(stmt Stmt_Print) {
	value := i.evaluate(stmt.expression)
	fmt.Println(value)
}

func (i *Interpreter) visitStmt_While(stmt Stmt_While) {
	for i.isTruthy(i.evaluate(stmt.condition)) {
		i.execute(stmt.body)
	}
}

func (i *Interpreter) visitStmt_Var(stmt Stmt_Var) {
	var value interface{}
	if stmt.initializer != nil {
		value = i.evaluate(stmt.initializer)
	}
	i.environment.define(stmt.name.lexeme, value)
}

func (i *Interpreter) visitExpr_Assign(expr Expr_Assign) interface{} {
	value := i.evaluate(expr.value)
	i.environment.assign(expr.name, value)
	return value
}

func (i *Interpreter) visitExpr_Variable(expr Expr_Variable) interface{} {
	return i.environment.get(expr.name)
}

func (i *Interpreter) visitExpr_Logical(expr Expr_Logical) interface{} {
	left := i.evaluate(expr.left)

	if expr.operator.ttype == OR {
		if i.isTruthy(left) {
			return left
		} else {
			if !i.isTruthy(left) {
				return left
			}
		}
	}
	return i.evaluate(expr.right)
}
