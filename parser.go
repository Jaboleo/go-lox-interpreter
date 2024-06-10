package main

import (
	"fmt"
	"log"
)

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(tokens []Token) Parser {
	return Parser{tokens, 0}
}

func (p *Parser) Parse() []Stmt {
	var statements []Stmt
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	return statements
}

// HELPER METHODS
func (p *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(type_ TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().ttype == type_
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().ttype == EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) consume(type_ TokenType, message string) (*Token, error) {
	if p.check(type_) {
		token := p.advance()
		return &token, nil
	}
	return nil, p.pError(p.peek(), message)
}

// ACTUAL GRAMMAR
func (p *Parser) expression() Expr {
	return p.assignment()
}

func (p *Parser) declaration() Stmt {
	if p.match(FUN) {
		return p.function("function")
	}
	if p.match(VAR) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = &Expr_Binary{expr, operator, right}
	}
	return expr, err
}

func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = &Expr_Binary{expr, operator, right}
	}
	return expr, err
}

func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = &Expr_Binary{expr, operator, right}
	}
	return expr, err
}

func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = &Expr_Binary{expr, operator, right}
	}
	return expr, err
}

func (p *Parser) unary() (Expr, error) {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right, _ := p.unary()
		return &Expr_Unary{operator, right}, nil
	}
	return p.call(), nil
}

func (p *Parser) call() Expr {
	expr, err := p.primary()
	if err != nil {
		log.Fatal("at call: %w", err)
	}

	for {
		if p.match(LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else {
			break
		}
	}
	return expr
}

func (p *Parser) finishCall(callee Expr) Expr {
	arguments := []Expr{}
	if !p.check(RIGHT_PAREN) {
		if len(arguments) >= 255 {
			log.Fatal(p.peek(), "Can't have more than 255 arguments.")
		}
		arguments = append(arguments, p.expression())
		for p.match(COMMA) {
			if len(arguments) >= 255 {
				log.Fatal(p.peek(), "Can't have more than 255 arguments.")
			}
			arguments = append(arguments, p.expression())
		}
	}
	paren, err := p.consume(RIGHT_PAREN, "Expect ')' after arguments.")
	if err != nil {
		log.Fatal("at finishCall: %w", err)
	}
	return Expr_Call{callee, *paren, arguments}
}

func (p *Parser) primary() (Expr, error) {
	if p.match(FALSE) {
		return &Expr_Literal{false}, nil
	}
	if p.match(TRUE) {
		return &Expr_Literal{true}, nil
	}
	if p.match(NIL) {
		return &Expr_Literal{nil}, nil
	}

	if p.match(NUMBER, STRING) {
		return &Expr_Literal{p.previous().literal}, nil
	}
	if p.match(IDENTIFIER) {
		return &Expr_Variable{p.previous()}, nil
	}
	if p.match(LEFT_PAREN) {
		expr := p.expression()
		_, err := p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}
		return &Expr_Grouping{expr}, nil
	}
	return nil, p.pError(p.peek(), "Expect expression.")
}

type ParseError struct {
	msg string
}

func (p ParseError) Error() string {
	return p.msg
}

func (p *Parser) pError(token Token, message string) error {
	loxError(token, message)
	return ParseError{message}
}

func (p *Parser) statement() Stmt {
	if p.match(FOR) {
		return p.forStatement()
	}
	if p.match(IF) {
		return p.ifStatement()
	}
	if p.match(PRINT) {
		return p.printStatement()
	}
	if p.match(RETURN) {
		return p.returnStatement()
	}
	if p.match(WHILE) {
		return p.whileStatement()
	}
	if p.match(LEFT_BRACE) {
		return Stmt_Block{p.block()}
	}
	return p.expressionStatement()
}

func (p *Parser) ifStatement() Stmt_If {
	p.consume(LEFT_PAREN, "Expect '(' after 'if'.")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' after if condition.")

	thenBranch := p.statement()
	var elseBranch Stmt
	if p.match(ELSE) {
		elseBranch = p.statement()
	}
	return Stmt_If{
		condition:  condition,
		thenBranch: thenBranch,
		elseBranch: elseBranch,
	}
}

func (p *Parser) printStatement() Stmt_Print {
	value := p.expression()
	p.consume(SEMICOLON, "Expect ';' after value.")
	return Stmt_Print{value}
}

func (p *Parser) returnStatement() Stmt_Return {
	keyword := p.previous()
	var value Expr
	if !p.check(SEMICOLON) {
		value = p.expression()
	}
	p.consume(SEMICOLON, "Expect ';' after return value.")
	return Stmt_Return{keyword, value}
}

func (p *Parser) whileStatement() Stmt_While {
	p.consume(LEFT_PAREN, "Expect '(' after 'while'.")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' after condition.")
	body := p.statement()

	return Stmt_While{condition, body}
}

func (p *Parser) forStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'for'.")
	var initializer Stmt
	if p.match(SEMICOLON) {
		initializer = nil
	} else if p.match(VAR) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}

	var condition Expr
	if !p.check(SEMICOLON) {
		condition = p.expression()
	}
	p.consume(SEMICOLON, "Expect ';' after loop condition.")

	var increment Expr
	if !p.check(RIGHT_PAREN) {
		increment = p.expression()
	}
	p.consume(RIGHT_PAREN, "Expect ')' after for clauses.")

	body := p.statement()
	if increment != nil {
		body = Stmt_Block{
			statements: []Stmt{body, Stmt_Expression{increment}},
		}
	}

	if condition == nil {
		condition = Expr_Literal{true}
	}
	body = Stmt_While{condition, body}

	if initializer != nil {
		body = Stmt_Block{
			[]Stmt{initializer, body},
		}
	}

	return body

}

func (p *Parser) expressionStatement() Stmt_Expression {
	expr := p.expression()
	p.consume(SEMICOLON, "Expect ';' after expression.")
	return Stmt_Expression{expr}
}

func (p *Parser) function(kind string) Stmt_Function {
	name, err := p.consume(IDENTIFIER, fmt.Sprintf("Expect %s name.", kind))
	if err != nil {
		log.Fatalf("%w at varDeclaration()", err)
	}
	p.consume(LEFT_PAREN, fmt.Sprintf("Expect ( after %s name.", kind))
	parameters := []Token{}
	if !p.check(RIGHT_PAREN) {
		for {
			if len(parameters) >= 255 {
				log.Fatal("Can't have more than 255 parameters.")
			}
			param, err := p.consume(IDENTIFIER, "Expect parameter name.")
			if err != nil {
				log.Fatalf("%w at varDeclaration()", err)
			}
			parameters = append(parameters, *param)
			if !p.match(COMMA) {
				break
			}
		}
	}
	p.consume(RIGHT_PAREN, "Expect ')' after parameters.")
	p.consume(LEFT_BRACE, fmt.Sprintf("Expect { before %w body.", kind))
	body := p.block()
	return Stmt_Function{*name, parameters, body}
}

func (p *Parser) block() []Stmt {
	var statements []Stmt
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	p.consume(RIGHT_BRACE, "Expect '}' after block.")
	return statements
}

func (p *Parser) varDeclaration() Stmt_Var {
	name, err := p.consume(IDENTIFIER, "Expect variable name.")
	if err != nil {
		fmt.Println("%w at varDeclaration()")
	}
	var initializer Expr
	if p.match(EQUAL) {
		initializer = p.expression()
	}
	p.consume(SEMICOLON, "Expect ';' after variable declaration.")
	return Stmt_Var{*name, initializer}
}

func (p *Parser) assignment() Expr {
	expr := p.or()
	if p.match(EQUAL) {
		equals := p.previous()
		value := p.assignment()

		if expr, ok := expr.(*Expr_Variable); ok {
			name := expr.name
			return &Expr_Assign{name, value}
		}
		err := NewRuntimeError(equals, "Invalid assignment target.")
		log.Fatal(err)
	}
	return expr
}

func (p *Parser) or() Expr {
	expr := p.and()
	for p.match(OR) {
		operator := p.previous()
		right := p.and()
		expr = Expr_Logical{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}
	return expr
}

func (p *Parser) and() Expr {
	expr, err := p.equality()
	if err != nil {
		log.Fatal("at and(): ", err)
	}

	for p.match(AND) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			log.Fatal("at and(): ", err)
		}
		expr = Expr_Logical{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}
	return expr
}

// func (p *Parser) synchronize() {
// 	p.advance()

// 	for !p.isAtEnd() {
// 		if p.previous().ttype == SEMICOLON {
// 			return
// 		}
// 		switch p.peek().ttype {
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

// TODO Czelendż 1 - comma operator
