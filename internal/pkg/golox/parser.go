package golox

// Parser implements a recursive descent parser for the Lox language
type Parser struct {
	tokens     []Token
	statements []Stmt
	current    int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens, statements: []Stmt{}, current: 0}
}

// Parse parses the list of tokens and returns the resulting expression or an error
func (p *Parser) Parse() ([]Stmt, error) {
	for !p.isAtEnd() {
		statement, err := p.declaration()
		if err != nil {
			return nil, err
		}
		p.statements = append(p.statements, statement)
	}
	return p.statements, nil
}

// expression -> assignment ;
func (p *Parser) expression() (Expr, error) {
	return p.assignment()
}

// assignment -> IDENTIFIER "=" assignment | logic_or ;
func (p *Parser) assignment() (Expr, error) {
	expr, err := p.or()
	if err != nil {
		return nil, err
	}

	if p.match(EQUAL) {
		// Every valid assignment target happens to also be valid syntax as a normal expression
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if variable, ok := expr.(*Variable); ok {
			name := variable.Name
			return &Assign{Name: name, Value: value}, nil
		} else if get, ok := expr.(*Get); ok {
			return &Set{Object: get.Object, Name: get.Name, Value: value}, nil
		}

		// TODO: We want to report the error, but continue parsing
		return nil, ParserError{
			Token:   *equals,
			Message: "Invalid assignment target.",
		}
	}

	return expr, nil
}

// or -> and ( "or" and )* ;
func (p *Parser) or() (Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.match(OR) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		expr = &Logical{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

// and -> equality ( "and" equality )* ;
func (p *Parser) and() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(AND) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = &Logical{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

// declaration -> classDecl | varDecl | statement | function ;
func (p *Parser) declaration() (Stmt, error) {
	if p.match(CLASS) {
		return p.classDeclaration()
	}
	if p.match(VAR) {
		return p.varDeclaration()
	}
	if p.match(FUN) {
		return p.function("function")
	}
	return p.statement()
}

// classDecl -> "class" IDENTIFIER ( "<" IDENTIFIER )? "{" function* "}" ;
func (p *Parser) classDeclaration() (Stmt, error) {
	nameToken, err := p.consume(IDENTIFIER, "Expect class name.")
	if err != nil {
		return nil, err
	}

	var superclass *Variable
	if p.match(LESS) {
		superclassToken, err := p.consume(IDENTIFIER, "Expect superclass name.")
		if err != nil {
			return nil, err
		}
		superclass = &Variable{Name: &superclassToken}
	}

	_, err = p.consume(LEFT_BRACE, "Expect '{' before class body.")
	if err != nil {
		return nil, err
	}

	methods := []Function{}
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		functionStmt, err := p.function("method")
		if err != nil {
			return nil, err
		}
		methods = append(methods, *functionStmt.(*Function))
	}

	_, err = p.consume(RIGHT_BRACE, "Expect '}' after class body.")
	if err != nil {
		return nil, err
	}

	return &Class{Name: &nameToken, Superclass: superclass, Methods: methods}, nil
}

// varDecl -> "var" IDENTIFIER ( "=" expression )? ";" ;
func (p *Parser) varDeclaration() (Stmt, error) {
	nameToken, err := p.consume(IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}

	var initializer Expr
	if p.match(EQUAL) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(SEMICOLON, "Expect ';' after variable declaration.")
	if err != nil {
		return nil, err
	}

	return &Var{Name: &nameToken, Initializer: initializer}, nil
}

// statement -> printStmt | forStmt | whileStmt | ifStmt | returnStmt | block | expressionStmt;
func (p *Parser) statement() (Stmt, error) {
	if p.match(IF) {
		return p.ifStatement()
	}
	if p.match(FOR) {
		return p.forStatement()
	}
	if p.match(WHILE) {
		return p.whileStatement()
	}
	if p.match(RETURN) {
		return p.returnStatement()
	}
	if p.match(PRINT) {
		return p.printStatement()
	}
	if p.match(LEFT_BRACE) {
		statements, err := p.block()
		if err != nil {
			return nil, err
		}
		return &Block{Statements: statements}, nil
	}
	return p.expressionStatement()
}

// printStmt -> "print" expression ";" ;
func (p *Parser) printStatement() (Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(SEMICOLON, "Expect ';' after value.")
	if err != nil {
		return nil, err
	}
	return &Print{Expression: value}, nil
}

// whileStmt -> "while" "(" expression ")" statement ;
func (p *Parser) whileStatement() (Stmt, error) {
	_, err := p.consume(LEFT_PAREN, "Expect '(' after 'while'.")
	if err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(RIGHT_PAREN, "Expect ')' after condition.")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return &While{Condition: condition, Body: body}, nil
}

// forStmt -> "for" "(" ( varDecl | expressionStmt | ";" ) expression? ";" expression? ")" statement ;
func (p *Parser) forStatement() (Stmt, error) {
	_, err := p.consume(LEFT_PAREN, "Expect '(' after 'for'.")
	if err != nil {
		return nil, err
	}

	var initializer Stmt
	if p.match(SEMICOLON) {
		initializer = nil
	} else if p.match(VAR) {
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

	var condition Expr
	if !p.check(SEMICOLON) {
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(SEMICOLON, "Expect ';' after loop condition.")
	if err != nil {
		return nil, err
	}

	var increment Expr
	if !p.check(RIGHT_PAREN) {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(RIGHT_PAREN, "Expect ')' after for clauses.")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	// Desugar for loop into while loop
	if increment != nil {
		body = &Block{Statements: []Stmt{
			body,
			&Expression{Expression: increment},
		}}
	}

	if condition == nil {
		condition = &Literal{Value: true}
	}
	body = &While{Condition: condition, Body: body}

	if initializer != nil {
		body = &Block{Statements: []Stmt{
			initializer,
			body,
		}}
	}

	return body, nil
}

// ifStmt -> "if" "(" expression ")" statement ( "else" statement )? ;
func (p *Parser) ifStatement() (Stmt, error) {
	_, err := p.consume(LEFT_PAREN, "Expect '(' after 'if'.")
	if err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(RIGHT_PAREN, "Expect ')' after if condition.")
	if err != nil {
		return nil, err
	}

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}

	var elseBranch Stmt
	if p.match(ELSE) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return &If{Condition: condition, ThenBranch: thenBranch, ElseBranch: elseBranch}, nil
}

// expressionStmt -> expression ";" ;
func (p *Parser) expressionStatement() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(SEMICOLON, "Expect ';' after expression.")
	if err != nil {
		return nil, err
	}
	return &Expression{Expression: expr}, nil
}

// returnStmt -> "return" expression? ";" ;
func (p *Parser) returnStatement() (Stmt, error) {
	keyword := p.previous()

	var value Expr
	var err error
	if !p.check(SEMICOLON) {
		value, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(SEMICOLON, "Expect ';' after return value.")
	if err != nil {
		return nil, err
	}

	return &Return{Keyword: keyword, Value: value}, nil
}

// function -> "fun" IDENTIFIER "(" parameters? ")" block ;
func (p *Parser) function(kind string) (Stmt, error) {
	nameToken, err := p.consume(IDENTIFIER, "Expect "+kind+" name.")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(LEFT_PAREN, "Expect '(' after "+kind+" name.")
	if err != nil {
		return nil, err
	}

	params := []*Token{}
	if !p.check(RIGHT_PAREN) {
		for {
			if len(params) >= 255 {
				return nil, ParserError{
					Token:   *p.peek(),
					Message: "Can't have more than 255 parameters.",
				}
			}
			paramToken, err := p.consume(IDENTIFIER, "Expect parameter name.")
			if err != nil {
				return nil, err
			}
			params = append(params, &paramToken)
			if !p.match(COMMA) {
				break
			}
		}
	}

	_, err = p.consume(RIGHT_PAREN, "Expect ')' after parameters.")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(LEFT_BRACE, "Expect '{' before "+kind+" body.")
	if err != nil {
		return nil, err
	}

	body, err := p.block()
	if err != nil {
		return nil, err
	}

	return &Function{Name: &nameToken, Params: params, Body: body}, nil
}

// block -> "{" declaration* "}" ;
func (p *Parser) block() ([]Stmt, error) {
	statements := []Stmt{}

	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		declaration, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, declaration)
	}

	_, err := p.consume(RIGHT_BRACE, "Expect '}' after block.")
	if err != nil {
		return nil, err
	}

	return statements, nil
}

// equality -> comparison ( ( "!=" | "==" ) comparison )* ;
func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

// comparison -> term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

// term -> factor ( ( "-" | "+" ) factor )* ;
func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

// factor -> unary ( ( "/" | "*" ) unary )* ;
func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

// unary -> ( "!" | "-" ) unary | call ;
func (p *Parser) unary() (Expr, error) {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &Unary{Operator: operator, Right: right}, nil
	}

	return p.call()
}

// call -> primary ( "(" arguments? ")" | "." IDENTIFIER )* ;
func (p *Parser) call() (Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(LEFT_PAREN) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else if p.match(DOT) {
			nameToken, err := p.consume(IDENTIFIER, "Expect property name after '.'.")
			if err != nil {
				return nil, err
			}
			expr = &Get{Object: expr, Name: &nameToken}
		} else {
			break
		}
	}

	return expr, nil
}

// finishCall handles parsing the arguments and closing parenthesis of a function call
func (p *Parser) finishCall(callee Expr) (Expr, error) {
	arguments := []Expr{}
	if !p.check(RIGHT_PAREN) {
		for {
			if len(arguments) >= 255 {
				return nil, ParserError{
					Token:   *p.peek(),
					Message: "Can't have more than 255 arguments.",
				}
			}
			arg, err := p.expression()
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, arg)
			if !p.match(COMMA) {
				break
			}
		}
	}

	paren, err := p.consume(RIGHT_PAREN, "Expect ')' after arguments.")
	if err != nil {
		return nil, err
	}

	return &Call{Callee: callee, Paren: &paren, Arguments: arguments}, nil
}

// primary -> "true" | "false" | "nil" | NUMBER | STRING | "(" expression ")" | this | IDENTIFIER | "super" "." IDENTIFIER ;
func (p *Parser) primary() (Expr, error) {
	if p.match(FALSE) {
		return &Literal{Value: false}, nil
	}
	if p.match(TRUE) {
		return &Literal{Value: true}, nil
	}
	if p.match(NIL) {
		return &Literal{Value: nil}, nil
	}
	if p.match(NUMBER, STRING) {
		return &Literal{Value: p.previous().Literal}, nil
	}
	if p.match(LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}
		return &Grouping{Expression: expr}, nil
	}
	if p.match(IDENTIFIER) {
		return &Variable{Name: p.previous()}, nil
	}
	if p.match(THIS) {
		return &This{Keyword: p.previous()}, nil
	}
	if p.match(SUPER) {
		keyword := p.previous()
		_, err := p.consume(DOT, "Expect '.' after 'super'.")
		if err != nil {
			return nil, err
		}
		method, err := p.consume(IDENTIFIER, "Expect superclass method name.")
		if err != nil {
			return nil, err
		}
		return &Super{Keyword: keyword, Method: &method}, nil
	}

	err := ParserError{
		Token:   *p.peek(),
		Message: "Expect expression.",
	}
	return nil, err
}

// Helper methods

// match checks if the current token is of any given types
func (p *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

// consume checks that the current token is of the expected type and advances
func (p *Parser) consume(t TokenType, message string) (Token, error) {
	if p.check(t) {
		return *p.advance(), nil
	}
	return Token{}, ParserError{Token: *p.peek(), Message: message}
}

// check checks if the current token is of the given type
func (p *Parser) check(t TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

// advance moves to the next token and returns the previous one
func (p *Parser) advance() *Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

// isAtEnd checks if we've reached the end of the token list
func (p *Parser) isAtEnd() bool {
	return p.peek().Type == EOF
}

// peek returns the current token
func (p *Parser) peek() *Token {
	return &p.tokens[p.current]
}

// previous returns the most recently consumed token
func (p *Parser) previous() *Token {
	return &p.tokens[p.current-1]
}

// Recovery method to synchronize the parser after an error

// synchronize discards tokens until it thinks it has found a statement boundary
func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == SEMICOLON {
			return
		}

		switch p.peek().Type {
		case CLASS, FUN, VAR, FOR, IF, WHILE, PRINT, RETURN:
			return
		}

		p.advance()
	}
}
