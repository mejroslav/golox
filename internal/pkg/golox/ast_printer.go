package golox

import "fmt"

type AstPrinter struct{}

func NewASTPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (a *AstPrinter) Print(expr Expr) string {
	return expr.Accept(a).(string)
}

func (a *AstPrinter) VisitBinaryExpr(expr *Binary) any {
	return a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (a *AstPrinter) VisitGroupingExpr(expr *Grouping) any {
	return a.parenthesize("group", expr.Expression)
}

func (a *AstPrinter) VisitLiteralExpr(expr *Literal) any {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", expr.Value)
}

func (a *AstPrinter) VisitUnaryExpr(expr *Unary) any {
	return a.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (a *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	result := "(" + name
	for _, expr := range exprs {
		result += " " + expr.Accept(a).(string)
	}
	result += ")"
	return result
}
