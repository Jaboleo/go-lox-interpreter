package ast

// Expr interface
type Expr interface {
	Accept(Visitor ExprVisitor) any
}

// ExprVisitor interface
type ExprVisitor interface {
	VisitExpr_Assign(e Expr_Assign) any
	VisitExpr_Binary(e Expr_Binary) any
	VisitExpr_Call(e Expr_Call) any
	VisitExpr_Grouping(e Expr_Grouping) any
	VisitExpr_Literal(e Expr_Literal) any
	VisitExpr_Unary(e Expr_Unary) any
	VisitExpr_Variable(e Expr_Variable) any
	VisitExpr_Logical(e Expr_Logical) any
}

// Expr_Binary struct
type Expr_Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (e Expr_Binary) Accept(Visitor ExprVisitor) any {
	return Visitor.VisitExpr_Binary(e)
}

type Expr_Call struct {
	Callee    Expr
	Paren     Token
	Arguments []Expr
}

func (e Expr_Call) Accept(Visitor ExprVisitor) any {
	return Visitor.VisitExpr_Call(e)
}

// Expr_Grouping struct
type Expr_Grouping struct {
	Expression Expr
}

func (e Expr_Grouping) Accept(Visitor ExprVisitor) any {
	return Visitor.VisitExpr_Grouping(e)
}

// Expr_Literal struct
type Expr_Literal struct {
	Value any
}

func (e Expr_Literal) Accept(Visitor ExprVisitor) any {
	return Visitor.VisitExpr_Literal(e)
}

// Expr_Unary struct
type Expr_Unary struct {
	Operator Token
	Right    Expr
}

func (e Expr_Unary) Accept(Visitor ExprVisitor) any {
	return Visitor.VisitExpr_Unary(e)
}

// Expr_Variable struct
type Expr_Variable struct {
	Name Token
}

func (e Expr_Variable) Accept(Visitor ExprVisitor) any {
	return Visitor.VisitExpr_Variable(e)
}

// Expr_Assign struct
type Expr_Assign struct {
	Name  Token
	Value Expr
}

func (e Expr_Assign) Accept(Visitor ExprVisitor) any {
	return Visitor.VisitExpr_Assign(e)
}

type Expr_Logical struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (e Expr_Logical) Accept(Visitor ExprVisitor) any {
	return Visitor.VisitExpr_Logical(e)
}
