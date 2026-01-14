package golox

import "fmt"

type AstPrinter struct{}

func NewASTPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (a *AstPrinter) Print(statement Stmt) string {
	result, _ := statement.Accept(a)
	return result.(string)
}

func (a *AstPrinter) VisitBinaryExpr(expr *Binary) (any, error) {
	return a.parenthesizeExprs(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (a *AstPrinter) VisitGroupingExpr(expr *Grouping) (any, error) {
	return a.parenthesizeExprs("group", expr.Expression)
}

func (a *AstPrinter) VisitLiteralExpr(expr *Literal) (any, error) {
	if expr.Value == nil {
		return "nil", nil
	}
	return fmt.Sprintf("%v", expr.Value), nil
}

func (a *AstPrinter) VisitUnaryExpr(expr *Unary) (any, error) {
	return a.parenthesizeExprs(expr.Operator.Lexeme, expr.Right)
}

func (a *AstPrinter) VisitExpressionStmt(expr *Expression) (any, error) {
	return a.parenthesizeExprs("expr", expr.Expression)
}

func (a *AstPrinter) VisitPrintStmt(expr *Print) (any, error) {
	return a.parenthesizeExprs("print", expr.Expression)
}

func (a *AstPrinter) VisitVarStmt(stmt *Var) (any, error) {
	if stmt.Initializer != nil {
		return a.parenthesizeExprs("var "+stmt.Name.Lexeme, stmt.Initializer)
	}
	return "(var " + stmt.Name.Lexeme + ")", nil
}

func (a *AstPrinter) VisitVariableExpr(expr *Variable) (any, error) {
	return expr.Name.Lexeme, nil
}

func (a *AstPrinter) VisitAssignExpr(expr *Assign) (any, error) {
	return a.parenthesizeExprs("assign "+expr.Name.Lexeme, expr.Value)
}

func (a *AstPrinter) VisitBlockStmt(stmt *Block) (any, error) {
	return a.parenthesizeStmts("block", stmt.Statements...)
}

func (a *AstPrinter) VisitIfStmt(stmt *If) (any, error) {
	if stmt.Elsebranch != nil {
		return a.parenthesize("if-else", stmt.Condition, stmt.Thenbranch, stmt.Elsebranch)
	}
	return a.parenthesize("if", stmt.Condition, stmt.Thenbranch)
}

func (a *AstPrinter) VisitLogicalExpr(expr *Logical) (any, error) {
	return a.parenthesizeExprs(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (a *AstPrinter) parenthesizeExprs(name string, exprs ...Expr) (string, error) {
	result := "(" + name
	for _, expr := range exprs {
		subResult, _ := expr.Accept(a)
		result += " " + subResult.(string)
	}
	result += ")"
	return result, nil
}

func (a *AstPrinter) parenthesizeStmts(name string, stmts ...Stmt) (string, error) {
	result := "(" + name
	for _, stmt := range stmts {
		subResult, _ := stmt.Accept(a)
		result += " " + subResult.(string)
	}
	result += ")"
	return result, nil
}

func (a *AstPrinter) parenthesize(name string, parts ...any) (string, error) {
	result := "(" + name
	for _, part := range parts {
		var subResult any
		switch v := part.(type) {
		case Expr:
			subResult, _ = v.Accept(a)
		case Stmt:
			subResult, _ = v.Accept(a)
		}
		result += " " + subResult.(string)
	}
	result += ")"
	return result, nil
}
