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
	return a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (a *AstPrinter) VisitGroupingExpr(expr *Grouping) (any, error) {
	return a.parenthesize("group", expr.Expression)
}

func (a *AstPrinter) VisitLiteralExpr(expr *Literal) (any, error) {
	if expr.Value == nil {
		return "nil", nil
	}
	return fmt.Sprintf("%v", expr.Value), nil
}

func (a *AstPrinter) VisitUnaryExpr(expr *Unary) (any, error) {
	return a.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (a *AstPrinter) VisitExpressionStmt(expr *Expression) (any, error) {
	return a.parenthesize("expr", expr.Expression)
}

func (a *AstPrinter) VisitPrintStmt(expr *Print) (any, error) {
	return a.parenthesize("print", expr.Expression)
}

func (a *AstPrinter) VisitVarStmt(stmt *Var) (any, error) {
	if stmt.Initializer != nil {
		return a.parenthesize("var "+stmt.Name.Lexeme, stmt.Initializer)
	}
	return a.parenthesize("var " + stmt.Name.Lexeme)
}

func (a *AstPrinter) VisitVariableExpr(expr *Variable) (any, error) {
	return expr.Name.Lexeme, nil
}

func (a *AstPrinter) parenthesize(name string, exprs ...Expr) (string, error) {
	result := "(" + name
	for _, expr := range exprs {
		subResult, _ := expr.Accept(a)
		result += " " + subResult.(string)
	}
	result += ")"
	return result, nil
}
