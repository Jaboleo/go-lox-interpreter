package resolve

import (
	"errors"
	"fmt"
	"log"

	"github.com/kljablon/golox/ast"
	"github.com/kljablon/golox/interpret"
)

type Resolver struct {
	interpreter     interpret.Interpreter
	scopes          []map[string]bool
	currentFunction FunctionType
}

type FunctionType int

const (
	NONE FunctionType = iota
	FUNCTION
)

func NewResover() Resolver {
	interpreter := interpret.NewInterpreter()
	scopes := []map[string]bool{}
	return Resolver{
		interpreter, scopes, NONE,
	}
}

func (r *Resolver) VisitStmt_Block(stmt ast.Stmt_Block) {
	r.beginScope()
	r.ResolveStmts(stmt.Statements)
	r.endScope()
}

func (r *Resolver) VisitStmt_Expression(stmt ast.Stmt_Expression) {
	r.resolveExpr(stmt.Expression)
}

func (r *Resolver) VisitStmt_If(stmt ast.Stmt_If) {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		r.resolveStmt(stmt.ElseBranch)
	}
}

func (r *Resolver) VisitStmt_Print(stmt ast.Stmt_Print) {
	r.resolveExpr(stmt.Expression)
}

func (r *Resolver) VisitStmt_Return(stmt ast.Stmt_Return) {
	if r.currentFunction == NONE {
		log.Fatal("Can't return from top-level code.")
	}
	if stmt.Value != nil {
		r.resolveExpr(stmt.Value)
	}
}

func (r *Resolver) VisitStmt_While(stmt ast.Stmt_While) {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.Body)
}

func (r *Resolver) VisitStmt_Function(stmt ast.Stmt_Function) {
	r.declare(stmt.Name)
	r.define(stmt.Name)
	r.resolveFunction(stmt, FUNCTION)
}

func (r *Resolver) VisitStmt_Var(stmt ast.Stmt_Var) {
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		r.resolveExpr(stmt.Initializer)
	}
	r.define(stmt.Name)
}

func (r *Resolver) VisitExpr_Assign(expr ast.Expr_Assign) any {
	r.resolveExpr(expr.Value)
	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r *Resolver) VisitExpr_Call(expr ast.Expr_Call) any {
	r.resolveExpr(expr.Callee)
	for _, argument := range expr.Arguments {
		r.resolveExpr(argument)
	}
	return nil
}

func (r *Resolver) VisitExpr_Grouping(expr ast.Expr_Grouping) any {
	r.resolveExpr(expr.Expression)
	return nil
}

func (r *Resolver) VisitExpr_Literal(expr ast.Expr_Literal) any {
	return nil
}

func (r *Resolver) VisitExpr_Logical(expr ast.Expr_Logical) any {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitExpr_Binary(expr ast.Expr_Binary) any {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitExpr_Unary(expr ast.Expr_Unary) any {
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitExpr_Variable(expr ast.Expr_Variable) any {
	if len(r.scopes) > 0 {
		scope, err := r.peekScopes()
		if err != nil {
			fmt.Print(err)
		}
		if lexeme, ok := scope[expr.Name.Lexeme]; ok {
			if !lexeme {
				log.Fatal("Can't read local variable in its own initializer.")
			}
		}
	}
	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r *Resolver) ResolveStmts(statements []ast.Stmt) {
	for _, statement := range statements {
		r.resolveStmt(statement)
	}
}

func (r *Resolver) resolveStmt(stmt ast.Stmt) {
	stmt.Accept(r)
}

func (r *Resolver) resolveExpr(expr ast.Expr) {
	expr.Accept(r)
}

func (r *Resolver) resolveFunction(function ast.Stmt_Function, f_type FunctionType) {
	enclosingFunction := r.currentFunction
	r.currentFunction = f_type
	r.beginScope()
	for _, param := range function.Params {
		r.declare(param)
		r.define(param)
	}
	r.ResolveStmts(function.Body)
	r.endScope()
	r.currentFunction = enclosingFunction
}

func (r *Resolver) resolveLocal(expr ast.Expr, name ast.Token) {
	for i := range r.scopes {
		if _, ok := r.scopes[i][name.Lexeme]; ok {
			r.interpreter.Resolve(expr, len(r.scopes)-1-i)
		}
	}
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, (make(map[string]bool)))
}

func (r *Resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) peekScopes() (map[string]bool, error) {
	if len(r.scopes) == 0 {
		return nil, errors.New("index ouf of range")
	}
	return r.scopes[len(r.scopes)-1], nil
}

func (r *Resolver) declare(name ast.Token) any {
	if len(r.scopes) == 0 {
		return nil
	}
	scope, err := r.peekScopes()
	if err != nil {
		fmt.Print(err)
	}
	if _, ok := scope[name.Lexeme]; ok {
		log.Fatal("Already a variable with this name in this scope: ", name)
	}

	scope[name.Lexeme] = false
	return nil
}

func (r *Resolver) define(name ast.Token) any {
	if len(r.scopes) == 0 {
		return nil
	}
	scope, err := r.peekScopes()
	if err != nil {
		fmt.Print(err)
	}
	scope[name.Lexeme] = true
	return nil
}
