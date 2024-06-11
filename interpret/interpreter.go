package interpret

import (
	"fmt"
	"log"

	"github.com/kljablon/golox/ast"
	"github.com/kljablon/golox/utils"
)

type Return struct {
	value any
}

type Interpreter struct {
	globals     Environment
	environment Environment
	locals      map[ast.Expr]int
}

func NewInterpreter() Interpreter {
	globals := NewEnvironment()

	globals.define("clock", ClockFunc{})

	locals := make(map[ast.Expr]int)
	return Interpreter{
		globals:     globals,
		environment: globals,
		locals:      locals,
	}
}

func (i *Interpreter) Interpret(statements []ast.Stmt) {
	for _, statement := range statements {
		i.execute(statement)
	}
}

func (i *Interpreter) VisitExpr_Binary(e ast.Expr_Binary) any {
	left := i.evaluate(e.Left)
	right := i.evaluate(e.Right)

	switch e.Operator.TokenType {
	case ast.MINUS:
		i.checkNumberOperand(e.Operator, right)
		return utils.CastToFloat(left) - utils.CastToFloat(right)
	case ast.PLUS:
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
		err := utils.NewRuntimeError(e.Operator, "Operands must be two numbers or two strings.")
		// Print the error message.
		fmt.Println(err.Error())

	case ast.SLASH:
		i.checkNumberOperands(e.Operator, left, right)
		return utils.CastToFloat(left) / utils.CastToFloat(right)
	case ast.STAR:
		i.checkNumberOperands(e.Operator, left, right)
		return utils.CastToFloat(left) * utils.CastToFloat(right)
	case ast.GREATER:
		i.checkNumberOperands(e.Operator, left, right)
		return utils.CastToFloat(left) > utils.CastToFloat(right)
	case ast.GREATER_EQUAL:
		i.checkNumberOperands(e.Operator, left, right)
		return utils.CastToFloat(left) >= utils.CastToFloat(right)
	case ast.LESS:
		i.checkNumberOperands(e.Operator, left, right)
		return utils.CastToFloat(left) < utils.CastToFloat(right)
	case ast.LESS_EQUAL:
		i.checkNumberOperands(e.Operator, left, right)
		return utils.CastToFloat(left) <= utils.CastToFloat(right)
	case ast.BANG_EQUAL:
		i.checkNumberOperands(e.Operator, left, right)
		return !utils.IsEqual(left, right)
	case ast.EQUAL_EQUAL:
		i.checkNumberOperands(e.Operator, left, right)
		return utils.IsEqual(left, right)
	default:
		return nil
	}
	return nil
}

func (i *Interpreter) VisitExpr_Call(e ast.Expr_Call) any {
	callee := i.evaluate(e.Callee)
	arguments := []any{}
	for _, argument := range e.Arguments {
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

func (i *Interpreter) VisitExpr_Grouping(e ast.Expr_Grouping) any {
	return i.evaluate(e.Expression)
}

func (i *Interpreter) VisitExpr_Literal(e ast.Expr_Literal) any {
	return e.Value
}

func (i *Interpreter) VisitExpr_Unary(e ast.Expr_Unary) any {
	right := i.evaluate(e.Right)
	switch e.Operator.TokenType {
	case ast.BANG:
		return !utils.IsTruthy(right)
	case ast.MINUS:
		return -utils.CastToFloat(right)
	default:
		return nil
	}
}

func (i *Interpreter) checkNumberOperand(operator ast.Token, operand any) {
	if _, ok := operand.(float64); ok {
		return
	}
	err := utils.NewRuntimeError(operator, "Operand must be a number.")
	// Print the error message.
	log.Fatal(err.Error())

}

func (i *Interpreter) checkNumberOperands(operator ast.Token, left any, right any) {
	if _, ok := left.(float64); ok {
		if _, ok := right.(float64); ok {
			return
		}
	}
	err := utils.NewRuntimeError(operator, "Operands must be numbers.")
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

func (i *Interpreter) evaluate(expr ast.Expr) any {
	return expr.Accept(i)
}

func (i *Interpreter) execute(stmt ast.Stmt) {
	if stmt == nil {
		panic("nil statement at execute()")
	}
	stmt.Accept(i)
}

func (i *Interpreter) Resolve(expr ast.Expr, depth int) {
	i.locals[expr] = depth
}

func (i *Interpreter) VisitStmt_Block(stmt ast.Stmt_Block) {
	enclosing_env := i.environment
	new_env := NewEnvironmentWithEnclosing(&enclosing_env)
	i.executeBlock(stmt.Statements, &new_env)
}

func (i *Interpreter) executeBlock(statements []ast.Stmt, environment *Environment) {
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

func (i *Interpreter) VisitStmt_Expression(stmt ast.Stmt_Expression) {
	i.evaluate(stmt.Expression)
}

func (i *Interpreter) VisitStmt_Function(stmt ast.Stmt_Function) {
	function := LoxFunction{stmt, i.environment}
	i.environment.define(stmt.Name.Lexeme, function)
}

func (i *Interpreter) VisitStmt_If(stmt ast.Stmt_If) {
	if utils.IsTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		i.execute(stmt.ElseBranch)
	}
}

func (i *Interpreter) VisitStmt_Print(stmt ast.Stmt_Print) {
	value := i.evaluate(stmt.Expression)
	fmt.Println(value)
}

func (i *Interpreter) VisitStmt_Return(stmt ast.Stmt_Return) {
	var value any
	if stmt.Value != nil {
		value = i.evaluate(stmt.Value)
	}
	panic(Return{value})
}

func (i *Interpreter) VisitStmt_While(stmt ast.Stmt_While) {
	for utils.IsTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.Body)
	}
}

func (i *Interpreter) VisitStmt_Var(stmt ast.Stmt_Var) {
	var value any
	if stmt.Initializer != nil {
		value = i.evaluate(stmt.Initializer)
	}
	i.environment.define(stmt.Name.Lexeme, value)
}

func (i *Interpreter) VisitExpr_Assign(expr ast.Expr_Assign) any {
	value := i.evaluate(expr.Value)

	if distance, ok := i.locals[expr]; ok {
		i.environment.assignAt(distance, expr.Name, value)
	} else {
		i.environment.assign(expr.Name, value)
	}

	return value
}

func (i *Interpreter) VisitExpr_Variable(expr ast.Expr_Variable) any {
	return i.lookUpVariable(expr.Name, expr)
}

func (i *Interpreter) lookUpVariable(name ast.Token, expr ast.Expr) any {
	if distance, ok := i.locals[expr]; ok {
		return i.environment.getAt(distance, name.Lexeme)
	} else {
		return i.globals.get(name)
	}
}

func (i *Interpreter) VisitExpr_Logical(expr ast.Expr_Logical) any {
	left := i.evaluate(expr.Left)

	if expr.Operator.TokenType == ast.OR {
		if utils.IsTruthy(left) {
			return left
		} else {
			if !utils.IsTruthy(left) {
				return left
			}
		}
	}
	return i.evaluate(expr.Right)
}
