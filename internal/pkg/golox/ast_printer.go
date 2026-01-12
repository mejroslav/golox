package golox

import "fmt"

type AstPrinter struct{}

func NewASTPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (a *AstPrinter) Print(expr Expr) string {
	expr.Accept(a)
	result, _ := expr.Accept(a)
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

func (a *AstPrinter) parenthesize(name string, exprs ...Expr) (string, error) {
	result := "(" + name
	for _, expr := range exprs {
		subResult, _ := expr.Accept(a)
		result += " " + subResult.(string)
	}
	result += ")"
	return result, nil
}
