package main

// Expr interface
type Expr interface {
	accept(visitor ExprVisitor) interface{}
}

// ExprVisitor interface
type ExprVisitor interface {
	visitExpr_Assign(e Expr_Assign) interface{}
	visitExpr_Binary(e Expr_Binary) interface{}
	visitExpr_Call(e Expr_Call) interface{}
	visitExpr_Grouping(e Expr_Grouping) interface{}
	visitExpr_Literal(e Expr_Literal) interface{}
	visitExpr_Unary(e Expr_Unary) interface{}
	visitExpr_Variable(e Expr_Variable) interface{}
	visitExpr_Logical(e Expr_Logical) interface{}
}

// Expr_Binary struct
type Expr_Binary struct {
	left     Expr
	operator Token
	right    Expr
}

func (e Expr_Binary) accept(visitor ExprVisitor) interface{} {
	return visitor.visitExpr_Binary(e)
}

type Expr_Call struct {
	callee    Expr
	paren     Token
	arguments []Expr
}

func (e Expr_Call) accept(visitor ExprVisitor) interface{} {
	return visitor.visitExpr_Call(e)
}

// Expr_Grouping struct
type Expr_Grouping struct {
	expression Expr
}

func (e Expr_Grouping) accept(visitor ExprVisitor) interface{} {
	return visitor.visitExpr_Grouping(e)
}

// Expr_Literal struct
type Expr_Literal struct {
	value interface{}
}

func (e Expr_Literal) accept(visitor ExprVisitor) interface{} {
	return visitor.visitExpr_Literal(e)
}

// Expr_Unary struct
type Expr_Unary struct {
	operator Token
	right    Expr
}

func (e Expr_Unary) accept(visitor ExprVisitor) interface{} {
	return visitor.visitExpr_Unary(e)
}

// Expr_Variable struct
type Expr_Variable struct {
	name Token
}

func (e Expr_Variable) accept(visitor ExprVisitor) interface{} {
	return visitor.visitExpr_Variable(e)
}

// Expr_Assign struct
type Expr_Assign struct {
	name  Token
	value Expr
}

func (e Expr_Assign) accept(visitor ExprVisitor) interface{} {
	return visitor.visitExpr_Assign(e)
}

type Expr_Logical struct {
	left     Expr
	operator Token
	right    Expr
}

func (e Expr_Logical) accept(visitor ExprVisitor) interface{} {
	return visitor.visitExpr_Logical(e)
}
