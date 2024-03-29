package main

type AstPrinter struct{}

func (a *AstPrinter) visitExpr_Binary(expr *Expr_Binary) interface{} {
	return a.parenthesize(expr.operator.lexeme, expr.left, expr.right)
}

func (a *AstPrinter) visitExpr_Grouping(expr *Expr_Grouping) interface{} {
	return a.parenthesize("group", expr.expression)
}

func (a *AstPrinter) visitExpr_Literal(expr *Expr_Literal) interface{} {
	if expr.value == nil {
		return "nil"
	}
	return expr.value.(string)
}

func (a *AstPrinter) visitExpr_Unary(expr *Expr_Unary) interface{} {
	return a.parenthesize(expr.operator.lexeme, expr.right)
}

func (a *AstPrinter) print(expr Expr) string {
	return expr.accept(a).(string)
}

func (a *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	builder := "(" + name
	for _, v := range exprs {
		builder += " " + v.accept(a).(string)
	}
	builder += ")"
	return builder
}
