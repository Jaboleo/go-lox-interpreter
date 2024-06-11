package ast

type Stmt interface {
	Accept(Visitor StmtVisitor)
}

type StmtVisitor interface {
	VisitStmt_Block(e Stmt_Block)
	VisitStmt_Expression(e Stmt_Expression)
	VisitStmt_Function(e Stmt_Function)
	VisitStmt_If(e Stmt_If)
	VisitStmt_Print(e Stmt_Print)
	VisitStmt_Return(e Stmt_Return)
	VisitStmt_While(e Stmt_While)
	VisitStmt_Var(e Stmt_Var)
}

type Stmt_Expression struct {
	Expression Expr
}

func (e Stmt_Expression) Accept(Visitor StmtVisitor) {
	Visitor.VisitStmt_Expression(e)
}

type Stmt_Function struct {
	Name   Token
	Params []Token
	Body   []Stmt
}

func (e Stmt_Function) Accept(Visitor StmtVisitor) {
	Visitor.VisitStmt_Function(e)
}

type Stmt_Print struct {
	Expression Expr
}

func (e Stmt_Print) Accept(Visitor StmtVisitor) {
	Visitor.VisitStmt_Print(e)
}

type Stmt_Return struct {
	Keyword Token
	Value   Expr
}

func (e Stmt_Return) Accept(Visitor StmtVisitor) {
	Visitor.VisitStmt_Return(e)
}

type Stmt_Var struct {
	Name        Token
	Initializer Expr
}

func (e Stmt_Var) Accept(Visitor StmtVisitor) {
	Visitor.VisitStmt_Var(e)
}

type Stmt_Block struct {
	Statements []Stmt
}

func (e Stmt_Block) Accept(Visitor StmtVisitor) {
	Visitor.VisitStmt_Block(e)
}

type Stmt_If struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (e Stmt_If) Accept(Visitor StmtVisitor) {
	Visitor.VisitStmt_If(e)
}

type Stmt_While struct {
	Condition Expr
	Body      Stmt
}

func (e Stmt_While) Accept(Visitor StmtVisitor) {
	Visitor.VisitStmt_While(e)
}
