package main

type Expr interface {
	accept(visitor Visitor) interface{}
}

type Visitor interface {
	visitExpr_Binary(e *Expr_Binary) interface{}
	visitExpr_Grouping(e *Expr_Grouping) interface{}
	visitExpr_Literal(e *Expr_Literal) interface{}
	visitExpr_Unary(e *Expr_Unary) interface{}
}

type Expr_Binary struct {
	left     Expr
	operator Token
	right    Expr
}

func (e *Expr_Binary) accept(visitor Visitor) interface{} {
	return visitor.visitExpr_Binary(e)
}

type Expr_Grouping struct {
	expression Expr
}

func (e *Expr_Grouping) accept(visitor Visitor) interface{} {
	return visitor.visitExpr_Grouping(e)
}

type Expr_Literal struct {
	value interface{}
}

func (e *Expr_Literal) accept(visitor Visitor) interface{} {
	return visitor.visitExpr_Literal(e)
}

type Expr_Unary struct {
	operator Token
	right    Expr
}

func (e *Expr_Unary) accept(visitor Visitor) interface{} {
	return visitor.visitExpr_Unary(e)
}
