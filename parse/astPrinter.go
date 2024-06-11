package parse

import (
	"fmt"

	"github.com/kljablon/golox/ast"
)

type AstPrinter struct{}

func (a *AstPrinter) VisitExpr_Binary(expr ast.Expr_Binary) any {
	return a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (a *AstPrinter) VisitExpr_Grouping(expr ast.Expr_Grouping) any {
	return a.parenthesize("group", expr.Expression)
}

func (a *AstPrinter) VisitExpr_Assign(expr ast.Expr_Assign) any {
	return a.parenthesize("assign", expr.Value)
}

func (a *AstPrinter) VisitExpr_Variable(expr ast.Expr_Variable) any {
	// return a.parenthesize("variable", expr.Name)
	return nil
}

func (a *AstPrinter) VisitExpr_Literal(expr ast.Expr_Literal) any {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", expr.Value)
}

func (a *AstPrinter) VisitExpr_Unary(expr ast.Expr_Unary) any {
	return a.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (a *AstPrinter) VisitExpr_Logical(expr ast.Expr_Logical) any {
	return a.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (a *AstPrinter) VisitExpr_Call(expr ast.Expr_Call) any {
	return a.parenthesize("", expr.Callee)
}

func (a *AstPrinter) Print(expr ast.Expr) string {
	return expr.Accept(a).(string)
}

func (a *AstPrinter) parenthesize(name string, exprs ...ast.Expr) string {
	builder := "(" + name
	for _, v := range exprs {
		builder += " " + v.Accept(a).(string)
	}
	builder += ")"
	return builder
}
