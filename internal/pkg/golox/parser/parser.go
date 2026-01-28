package parser

import (
	"fmt"

	"github.com/mejroslav/golox/internal/pkg/golox/ast"
	"github.com/mejroslav/golox/internal/pkg/golox/lox_error"
	"github.com/mejroslav/golox/internal/pkg/golox/token"
)

// Parser implements a recursive descent parser for the Lox language
type Parser struct {
	tokens     []token.Token
	statements []ast.Stmt
	current    int
}

func NewParser(tokens []token.Token) *Parser {
	return &Parser{tokens: tokens, statements: []ast.Stmt{}, current: 0}
}

// Parse parses the list of tokens and returns the resulting statements or an error
func (p *Parser) Parse() ([]ast.Stmt, bool) {
	hadError := false
	for !p.isAtEnd() {
		statement, err := p.declaration()
		if err != nil {
			fmt.Println(err.Error())
			p.synchronize()
			hadError = true
			continue
		}
		p.statements = append(p.statements, statement)
	}

	return p.statements, hadError
}

// expression -> assignment ;
func (p *Parser) expression() (ast.Expr, error) {
	return p.assignment()
}

// assignment -> IDENTIFIER "=" assignment | logic_or ;
func (p *Parser) assignment() (ast.Expr, error) {
	expr, err := p.or()
	if err != nil {
		return nil, err
	}

	if p.match(token.EQUAL) {
		// Every valid assignment target happens to also be valid syntax as a normal expression
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if variable, ok := expr.(*ast.Variable); ok {
			name := variable.Name
			return &ast.Assign{Name: name, Value: value}, nil
		} else if get, ok := expr.(*ast.Get); ok {
			return &ast.Set{Object: get.Object, Name: get.Name, Value: value}, nil
		}

		// TODO: We want to report the error, but continue parsing
		return nil, lox_error.ParserError{
			Token:   *equals,
			Message: "Invalid assignment target.",
		}
	}

	return expr, nil
}

// or -> and ( "or" and )* ;
func (p *Parser) or() (ast.Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.match(token.OR) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		expr = &ast.Logical{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

// and -> equality ( "and" equality )* ;
func (p *Parser) and() (ast.Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(token.AND) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = &ast.Logical{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

// declaration -> classDecl | varDecl | statement | function ;
func (p *Parser) declaration() (ast.Stmt, error) {
	if p.match(token.CLASS) {
		return p.classDeclaration()
	}
	if p.match(token.VAR) {
		return p.varDeclaration()
	}
	if p.match(token.FUN) {
		return p.function("function")
	}
	return p.statement()
}

// classDecl -> "class" IDENTIFIER ( "<" IDENTIFIER )? "{" function* "}" ;
func (p *Parser) classDeclaration() (ast.Stmt, error) {
	nameToken, err := p.consume(token.IDENTIFIER, "Expect class name.")
	if err != nil {
		return nil, err
	}

	var superclass *ast.Variable
	if p.match(token.LESS) {
		superclassToken, err := p.consume(token.IDENTIFIER, "Expect superclass name.")
		if err != nil {
			return nil, err
		}
		superclass = &ast.Variable{Name: &superclassToken}
	}

	_, err = p.consume(token.LEFT_BRACE, "Expect '{' before class body.")
	if err != nil {
		return nil, err
	}

	methods := []ast.Function{}
	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		functionStmt, err := p.function("method")
		if err != nil {
			return nil, err
		}
		methods = append(methods, *functionStmt.(*ast.Function))
	}

	_, err = p.consume(token.RIGHT_BRACE, "Expect '}' after class body.")
	if err != nil {
		return nil, err
	}

	return &ast.Class{Name: &nameToken, Superclass: superclass, Methods: methods}, nil
}

// varDecl -> "var" IDENTIFIER ( "=" expression )? ";" ;
func (p *Parser) varDeclaration() (ast.Stmt, error) {
	nameToken, err := p.consume(token.IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}

	var initializer ast.Expr
	if p.match(token.EQUAL) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(token.SEMICOLON, "Expect ';' after variable declaration.")
	if err != nil {
		return nil, err
	}

	return &ast.Var{Name: &nameToken, Initializer: initializer}, nil
}

// statement -> printStmt | forStmt | whileStmt | ifStmt | returnStmt | breakStmt | block | expressionStmt;
func (p *Parser) statement() (ast.Stmt, error) {
	if p.match(token.IF) {
		return p.ifStatement()
	}
	if p.match(token.FOR) {
		return p.forStatement()
	}
	if p.match(token.WHILE) {
		return p.whileStatement()
	}
	if p.match(token.RETURN) {
		return p.returnStatement()
	}
	if p.match(token.PRINT) {
		return p.printStatement()
	}
	if p.match(token.BREAK) {
		return p.breakStatement()
	}
	if p.match(token.LEFT_BRACE) {
		statements, err := p.block()
		if err != nil {
			return nil, err
		}
		return &ast.Block{Statements: statements}, nil
	}
	return p.expressionStatement()
}

// printStmt -> "print" expression ";" ;
func (p *Parser) printStatement() (ast.Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(token.SEMICOLON, "Expect ';' after value.")
	if err != nil {
		return nil, err
	}
	return &ast.Print{Expression: value}, nil
}

// whileStmt -> "while" "(" expression ")" statement ;
func (p *Parser) whileStatement() (ast.Stmt, error) {
	_, err := p.consume(token.LEFT_PAREN, "Expect '(' after 'while'.")
	if err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after condition.")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return &ast.While{Condition: condition, Body: body}, nil
}

// forStmt -> "for" "(" ( varDecl | expressionStmt | ";" ) expression? ";" expression? ")" statement ;
func (p *Parser) forStatement() (ast.Stmt, error) {
	_, err := p.consume(token.LEFT_PAREN, "Expect '(' after 'for'.")
	if err != nil {
		return nil, err
	}

	var initializer ast.Stmt
	if p.match(token.SEMICOLON) {
		initializer = nil
	} else if p.match(token.VAR) {
		initializer, err = p.varDeclaration()
		if err != nil {
			return nil, err
		}
	} else {
		initializer, err = p.expressionStatement()
		if err != nil {
			return nil, err
		}
	}

	var condition ast.Expr
	if !p.check(token.SEMICOLON) {
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(token.SEMICOLON, "Expect ';' after loop condition.")
	if err != nil {
		return nil, err
	}

	var increment ast.Expr
	if !p.check(token.RIGHT_PAREN) {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after for clauses.")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	// Desugar for loop into while loop
	if increment != nil {
		body = &ast.Block{Statements: []ast.Stmt{
			body,
			&ast.Expression{Expression: increment},
		}}
	}

	if condition == nil {
		condition = &ast.Literal{Value: true}
	}
	body = &ast.While{Condition: condition, Body: body}

	if initializer != nil {
		body = &ast.Block{Statements: []ast.Stmt{
			initializer,
			body,
		}}
	}

	return body, nil
}

// breakStmt -> "break" ";" ;
func (p *Parser) breakStatement() (ast.Stmt, error) {
	keyword := p.previous()
	_, err := p.consume(token.SEMICOLON, "Expect ';' after 'break'.")
	if err != nil {
		return nil, err
	}
	return &ast.Break{Keyword: keyword}, nil
}

// ifStmt -> "if" "(" expression ")" statement ( "else" statement )? ;
func (p *Parser) ifStatement() (ast.Stmt, error) {
	_, err := p.consume(token.LEFT_PAREN, "Expect '(' after 'if'.")
	if err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after if condition.")
	if err != nil {
		return nil, err
	}

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}

	var elseBranch ast.Stmt
	if p.match(token.ELSE) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return &ast.If{Condition: condition, ThenBranch: thenBranch, ElseBranch: elseBranch}, nil
}

// expressionStmt -> expression ";" ;
func (p *Parser) expressionStatement() (ast.Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(token.SEMICOLON, "Expect ';' after expression.")
	if err != nil {
		return nil, err
	}
	return &ast.Expression{Expression: expr}, nil
}

// returnStmt -> "return" expression? ";" ;
func (p *Parser) returnStatement() (ast.Stmt, error) {
	keyword := p.previous()

	var value ast.Expr
	var err error
	if !p.check(token.SEMICOLON) {
		value, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(token.SEMICOLON, "Expect ';' after return value.")
	if err != nil {
		return nil, err
	}

	return &ast.Return{Keyword: keyword, Value: value}, nil
}

// function -> "fun" IDENTIFIER "(" parameters? ")" block ;
func (p *Parser) function(kind string) (ast.Stmt, error) {
	nameToken, err := p.consume(token.IDENTIFIER, "Expect "+kind+" name.")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.LEFT_PAREN, "Expect '(' after "+kind+" name.")
	if err != nil {
		return nil, err
	}

	params := []*token.Token{}
	if !p.check(token.RIGHT_PAREN) {
		for {
			if len(params) >= 255 {
				return nil, lox_error.ParserError{
					Token:   *p.peek(),
					Message: "Can't have more than 255 parameters.",
				}
			}
			paramToken, err := p.consume(token.IDENTIFIER, "Expect parameter name.")
			if err != nil {
				return nil, err
			}
			params = append(params, &paramToken)
			if !p.match(token.COMMA) {
				break
			}
		}
	}

	_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after parameters.")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.LEFT_BRACE, "Expect '{' before "+kind+" body.")
	if err != nil {
		return nil, err
	}

	body, err := p.block()
	if err != nil {
		return nil, err
	}

	return &ast.Function{Name: &nameToken, Params: params, Body: body}, nil
}

// block -> "{" declaration* "}" ;
func (p *Parser) block() ([]ast.Stmt, error) {
	statements := []ast.Stmt{}

	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		declaration, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, declaration)
	}

	_, err := p.consume(token.RIGHT_BRACE, "Expect '}' after block.")
	if err != nil {
		return nil, err
	}

	return statements, nil
}

// equality -> comparison ( ( "!=" | "==" ) comparison )* ;
func (p *Parser) equality() (ast.Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

// comparison -> term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
func (p *Parser) comparison() (ast.Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

// term -> factor ( ( "-" | "+" ) factor )* ;
func (p *Parser) term() (ast.Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(token.MINUS, token.PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

// factor -> unary ( ( "/" | "*" ) unary )* ;
func (p *Parser) factor() (ast.Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(token.SLASH, token.STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

// unary -> ( "!" | "-" ) unary | call ;
func (p *Parser) unary() (ast.Expr, error) {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &ast.Unary{Operator: operator, Right: right}, nil
	}

	return p.call()
}

// call -> primary ( "(" arguments? ")" | "." IDENTIFIER )* ;
func (p *Parser) call() (ast.Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(token.LEFT_PAREN) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else if p.match(token.DOT) {
			nameToken, err := p.consume(token.IDENTIFIER, "Expect property name after '.'.")
			if err != nil {
				return nil, err
			}
			expr = &ast.Get{Object: expr, Name: &nameToken}
		} else {
			break
		}
	}

	return expr, nil
}

// finishCall handles parsing the arguments and closing parenthesis of a function call
func (p *Parser) finishCall(callee ast.Expr) (ast.Expr, error) {
	arguments := []ast.Expr{}
	if !p.check(token.RIGHT_PAREN) {
		for {
			if len(arguments) >= 255 {
				return nil, lox_error.ParserError{
					Token:   *p.peek(),
					Message: "Can't have more than 255 arguments.",
				}
			}
			arg, err := p.expression()
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, arg)
			if !p.match(token.COMMA) {
				break
			}
		}
	}

	paren, err := p.consume(token.RIGHT_PAREN, "Expect ')' after arguments.")
	if err != nil {
		return nil, err
	}

	return &ast.Call{Callee: callee, Paren: &paren, Arguments: arguments}, nil
}

// primary -> "true" | "false" | "nil" | NUMBER | STRING | "(" expression ")" | this | IDENTIFIER | "super" "." IDENTIFIER ;
func (p *Parser) primary() (ast.Expr, error) {
	if p.match(token.FALSE) {
		return &ast.Literal{Value: false}, nil
	}
	if p.match(token.TRUE) {
		return &ast.Literal{Value: true}, nil
	}
	if p.match(token.NIL) {
		return &ast.Literal{Value: nil}, nil
	}
	if p.match(token.NUMBER, token.STRING) {
		return &ast.Literal{Value: p.previous().Literal}, nil
	}
	if p.match(token.LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}
		return &ast.Grouping{Expression: expr}, nil
	}
	if p.match(token.IDENTIFIER) {
		return &ast.Variable{Name: p.previous()}, nil
	}
	if p.match(token.THIS) {
		return &ast.This{Keyword: p.previous()}, nil
	}
	if p.match(token.SUPER) {
		keyword := p.previous()
		_, err := p.consume(token.DOT, "Expect '.' after 'super'.")
		if err != nil {
			return nil, err
		}
		method, err := p.consume(token.IDENTIFIER, "Expect superclass method name.")
		if err != nil {
			return nil, err
		}
		return &ast.Super{Keyword: keyword, Method: &method}, nil
	}

	err := lox_error.ParserError{
		Token:   *p.peek(),
		Message: "Expect expression.",
	}
	return nil, err
}

// Helper methods

// match checks if the current token is of any given types
func (p *Parser) match(types ...token.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

// consume checks that the current token is of the expected type and advances
func (p *Parser) consume(t token.TokenType, message string) (token.Token, error) {
	if p.check(t) {
		return *p.advance(), nil
	}
	return token.Token{}, lox_error.ParserError{Token: *p.peek(), Message: message}
}

// check checks if the current token is of the given type
func (p *Parser) check(t token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

// advance moves to the next token and returns the previous one
func (p *Parser) advance() *token.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

// isAtEnd checks if we've reached the end of the token list
func (p *Parser) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

// peek returns the current token
func (p *Parser) peek() *token.Token {
	return &p.tokens[p.current]
}

// previous returns the most recently consumed token
func (p *Parser) previous() *token.Token {
	return &p.tokens[p.current-1]
}

// Recovery method to synchronize the parser after an error

// synchronize discards tokens until it thinks it has found a statement boundary
func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == token.SEMICOLON {
			return
		}

		switch p.peek().Type {
		case token.CLASS, token.FUN, token.VAR, token.FOR, token.IF, token.WHILE, token.PRINT, token.RETURN, token.BREAK:
			return
		}

		p.advance()
	}
}
