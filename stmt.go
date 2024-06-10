package main

type Stmt interface {
	accept(visitor StmtVisitor)
}

type StmtVisitor interface {
	visitStmt_Block(e Stmt_Block)
	visitStmt_Expression(e Stmt_Expression)
	visitStmt_Function(e Stmt_Function)
	visitStmt_If(e Stmt_If)
	visitStmt_Print(e Stmt_Print)
	visitStmt_Return(e Stmt_Return)
	visitStmt_While(e Stmt_While)
	visitStmt_Var(e Stmt_Var)
}

type Stmt_Expression struct {
	expression Expr
}

func (e Stmt_Expression) accept(visitor StmtVisitor) {
	visitor.visitStmt_Expression(e)
}

type Stmt_Function struct {
	name   Token
	params []Token
	body   []Stmt
}

func (e Stmt_Function) accept(visitor StmtVisitor) {
	visitor.visitStmt_Function(e)
}

type Stmt_Print struct {
	expression Expr
}

func (e Stmt_Print) accept(visitor StmtVisitor) {
	visitor.visitStmt_Print(e)
}

type Stmt_Return struct {
	keyword Token
	value   Expr
}

func (e Stmt_Return) accept(visitor StmtVisitor) {
	visitor.visitStmt_Return(e)
}

type Stmt_Var struct {
	name        Token
	initializer Expr
}

func (e Stmt_Var) accept(visitor StmtVisitor) {
	visitor.visitStmt_Var(e)
}

type Stmt_Block struct {
	statements []Stmt
}

func (e Stmt_Block) accept(visitor StmtVisitor) {
	visitor.visitStmt_Block(e)
}

type Stmt_If struct {
	condition  Expr
	thenBranch Stmt
	elseBranch Stmt
}

func (e Stmt_If) accept(visitor StmtVisitor) {
	visitor.visitStmt_If(e)
}

type Stmt_While struct {
	condition Expr
	body      Stmt
}

func (e Stmt_While) accept(visitor StmtVisitor) {
	visitor.visitStmt_While(e)
}
