package main

import (
	"fmt"
	"log"

	"github.com/kljablon/golox/ast"
	"github.com/kljablon/golox/utils"
)

type Parser struct {
	tokens  []ast.Token
	current int
}

func NewParser(tokens []ast.Token) Parser {
	return Parser{tokens, 0}
}

func (p *Parser) Parse() []ast.Stmt {
	var statements []ast.Stmt
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	return statements
}

// HELPER METHODS
func (p *Parser) match(types ...ast.TokenType) bool {
	for _, ttype := range types {
		if p.check(ttype) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(ttype ast.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().TokenType == ttype
}

func (p *Parser) advance() ast.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == ast.EOF
}

func (p *Parser) peek() ast.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() ast.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) consume(ttype ast.TokenType, message string) (*ast.Token, error) {
	if p.check(ttype) {
		token := p.advance()
		return &token, nil
	}
	return nil, p.pError(p.peek(), message)
}

// ACTUAL GRAMMAR
func (p *Parser) expression() ast.Expr {
	return p.assignment()
}

func (p *Parser) declaration() ast.Stmt {
	if p.match(ast.FUN) {
		return p.function("function")
	}
	if p.match(ast.VAR) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) equality() (ast.Expr, error) {
	expr, err := p.comparison()

	for p.match(ast.BANG_EQUAL, ast.EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = &ast.Expr_Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr, err
}

func (p *Parser) comparison() (ast.Expr, error) {
	expr, err := p.term()

	for p.match(ast.GREATER, ast.GREATER_EQUAL, ast.LESS, ast.LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = &ast.Expr_Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr, err
}

func (p *Parser) term() (ast.Expr, error) {
	expr, err := p.factor()

	for p.match(ast.MINUS, ast.PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = &ast.Expr_Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr, err
}

func (p *Parser) factor() (ast.Expr, error) {
	expr, err := p.unary()

	for p.match(ast.SLASH, ast.STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = &ast.Expr_Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr, err
}

func (p *Parser) unary() (ast.Expr, error) {
	if p.match(ast.BANG, ast.MINUS) {
		operator := p.previous()
		right, _ := p.unary()
		return &ast.Expr_Unary{Operator: operator, Right: right}, nil
	}
	return p.call(), nil
}

func (p *Parser) call() ast.Expr {
	expr, err := p.primary()
	if err != nil {
		log.Fatal("at call: %w", err)
	}

	for {
		if p.match(ast.LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else {
			break
		}
	}
	return expr
}

func (p *Parser) finishCall(callee ast.Expr) ast.Expr {
	arguments := []ast.Expr{}
	if !p.check(ast.RIGHT_PAREN) {
		if len(arguments) >= 255 {
			log.Fatal(p.peek(), "Can't have more than 255 arguments.")
		}
		arguments = append(arguments, p.expression())
		for p.match(ast.COMMA) {
			if len(arguments) >= 255 {
				log.Fatal(p.peek(), "Can't have more than 255 arguments.")
			}
			arguments = append(arguments, p.expression())
		}
	}
	paren, err := p.consume(ast.RIGHT_PAREN, "Expect ')' after arguments.")
	if err != nil {
		log.Fatal("at finishCall: %w", err)
	}
	return ast.Expr_Call{Callee: callee, Paren: *paren, Arguments: arguments}
}

func (p *Parser) primary() (ast.Expr, error) {
	if p.match(ast.FALSE) {
		return &ast.Expr_Literal{Value: false}, nil
	}
	if p.match(ast.TRUE) {
		return &ast.Expr_Literal{Value: true}, nil
	}
	if p.match(ast.NIL) {
		return &ast.Expr_Literal{Value: nil}, nil
	}

	if p.match(ast.NUMBER, ast.STRING) {
		return &ast.Expr_Literal{Value: p.previous().Literal}, nil
	}
	if p.match(ast.IDENTIFIER) {
		return &ast.Expr_Variable{Name: p.previous()}, nil
	}
	if p.match(ast.LEFT_PAREN) {
		expr := p.expression()
		_, err := p.consume(ast.RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}
		return &ast.Expr_Grouping{Expression: expr}, nil
	}
	return nil, p.pError(p.peek(), "Expect expression.")
}

type ParseError struct {
	msg string
}

func (p ParseError) Error() string {
	return p.msg
}

func (p *Parser) pError(token ast.Token, message string) error {
	loxError(token, message)
	return ParseError{message}
}

func (p *Parser) statement() ast.Stmt {
	if p.match(ast.FOR) {
		return p.forStatement()
	}
	if p.match(ast.IF) {
		return p.ifStatement()
	}
	if p.match(ast.PRINT) {
		return p.printStatement()
	}
	if p.match(ast.RETURN) {
		return p.returnStatement()
	}
	if p.match(ast.WHILE) {
		return p.whileStatement()
	}
	if p.match(ast.LEFT_BRACE) {
		return ast.Stmt_Block{p.block()}
	}
	return p.expressionStatement()
}

func (p *Parser) ifStatement() ast.Stmt_If {
	p.consume(ast.LEFT_PAREN, "Expect '(' after 'if'.")
	condition := p.expression()
	p.consume(ast.RIGHT_PAREN, "Expect ')' after if condition.")

	thenBranch := p.statement()
	var elseBranch ast.Stmt
	if p.match(ast.ELSE) {
		elseBranch = p.statement()
	}
	return ast.Stmt_If{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}
}

func (p *Parser) printStatement() ast.Stmt_Print {
	value := p.expression()
	p.consume(ast.SEMICOLON, "Expect ';' after value.")
	return ast.Stmt_Print{value}
}

func (p *Parser) returnStatement() ast.Stmt_Return {
	keyword := p.previous()
	var value ast.Expr
	if !p.check(ast.SEMICOLON) {
		value = p.expression()
	}
	p.consume(ast.SEMICOLON, "Expect ';' after return value.")
	return ast.Stmt_Return{keyword, value}
}

func (p *Parser) whileStatement() ast.Stmt_While {
	p.consume(ast.LEFT_PAREN, "Expect '(' after 'while'.")
	condition := p.expression()
	p.consume(ast.RIGHT_PAREN, "Expect ')' after condition.")
	body := p.statement()

	return ast.Stmt_While{condition, body}
}

func (p *Parser) forStatement() ast.Stmt {
	p.consume(ast.LEFT_PAREN, "Expect '(' after 'for'.")
	var initializer ast.Stmt
	if p.match(ast.SEMICOLON) {
		initializer = nil
	} else if p.match(ast.VAR) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}

	var condition ast.Expr
	if !p.check(ast.SEMICOLON) {
		condition = p.expression()
	}
	p.consume(ast.SEMICOLON, "Expect ';' after loop condition.")

	var increment ast.Expr
	if !p.check(ast.RIGHT_PAREN) {
		increment = p.expression()
	}
	p.consume(ast.RIGHT_PAREN, "Expect ')' after for clauses.")

	body := p.statement()
	if increment != nil {
		body = ast.Stmt_Block{
			Statements: []ast.Stmt{body, ast.Stmt_Expression{increment}},
		}
	}

	if condition == nil {
		condition = ast.Expr_Literal{Value: true}
	}
	body = ast.Stmt_While{Condition: condition, Body: body}

	if initializer != nil {
		body = ast.Stmt_Block{
			Statements: []ast.Stmt{initializer, body},
		}
	}

	return body

}

func (p *Parser) expressionStatement() ast.Stmt_Expression {
	expr := p.expression()
	p.consume(ast.SEMICOLON, "Expect ';' after expression.")
	return ast.Stmt_Expression{Expression: expr}
}

func (p *Parser) function(kind string) ast.Stmt_Function {
	name, err := p.consume(ast.IDENTIFIER, fmt.Sprintf("Expect %s name.", kind))
	if err != nil {
		log.Fatalf("%w at varDeclaration()", err)
	}
	p.consume(ast.LEFT_PAREN, fmt.Sprintf("Expect ( after %s name.", kind))
	parameters := []ast.Token{}
	if !p.check(ast.RIGHT_PAREN) {
		for {
			if len(parameters) >= 255 {
				log.Fatal("Can't have more than 255 parameters.")
			}
			param, err := p.consume(ast.IDENTIFIER, "Expect parameter name.")
			if err != nil {
				log.Fatalf("%w at varDeclaration()", err)
			}
			parameters = append(parameters, *param)
			if !p.match(ast.COMMA) {
				break
			}
		}
	}
	p.consume(ast.RIGHT_PAREN, "Expect ')' after parameters.")
	p.consume(ast.LEFT_BRACE, fmt.Sprintf("Expect { before %s body.", kind))
	body := p.block()
	return ast.Stmt_Function{Name: *name, Params: parameters, Body: body}
}

func (p *Parser) block() []ast.Stmt {
	var statements []ast.Stmt
	for !p.check(ast.RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	p.consume(ast.RIGHT_BRACE, "Expect '}' after block.")
	return statements
}

func (p *Parser) varDeclaration() ast.Stmt_Var {
	name, err := p.consume(ast.IDENTIFIER, "Expect variable name.")
	if err != nil {
		fmt.Println("%w at varDeclaration()")
	}
	var initializer ast.Expr
	if p.match(ast.EQUAL) {
		initializer = p.expression()
	}
	p.consume(ast.SEMICOLON, "Expect ';' after variable declaration.")
	return ast.Stmt_Var{Name: *name, Initializer: initializer}
}

func (p *Parser) assignment() ast.Expr {
	expr := p.or()
	if p.match(ast.EQUAL) {
		equals := p.previous()
		value := p.assignment()

		if expr, ok := expr.(*ast.Expr_Variable); ok {
			name := expr.Name
			return &ast.Expr_Assign{Name: name, Value: value}
		}
		err := utils.NewRuntimeError(equals, "Invalid assignment target.")
		log.Fatal(err)
	}
	return expr
}

func (p *Parser) or() ast.Expr {
	expr := p.and()
	for p.match(ast.OR) {
		operator := p.previous()
		right := p.and()
		expr = ast.Expr_Logical{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr
}

func (p *Parser) and() ast.Expr {
	expr, err := p.equality()
	if err != nil {
		log.Fatal("at and(): ", err)
	}

	for p.match(ast.AND) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			log.Fatal("at and(): ", err)
		}
		expr = ast.Expr_Logical{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr
}

// func (p *Parser) synchronize() {
// 	p.advance()

// 	for !p.isAtEnd() {
// 		if p.previous().TokenType == SEMICOLON {
// 			return
// 		}
// 		switch p.peek().TokenType {
// 		case CLASS:
// 		case FUN:
// 		case VAR:
// 		case FOR:
// 		case IF:
// 		case WHILE:
// 		case PRINT:
// 		case RETURN:
// 			return
// 		}
// 		p.advance()
// 	}
// }
