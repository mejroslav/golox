package golox

import (
	"fmt"
	"mejroslav/golox/v2/internal/pkg/golox/ast"
)

type AstPrinter struct{}

func NewASTPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (a *AstPrinter) Print(statements []ast.Stmt) string {
	result := ""
	for _, statement := range statements {
		subResult, _ := statement.Accept(a)
		result += subResult.(string) + "\n"
	}
	return result
}

func (a *AstPrinter) VisitBinaryExpr(expr *ast.Binary) (any, error) {
	return a.parenthesizeExprs(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (a *AstPrinter) VisitGroupingExpr(expr *ast.Grouping) (any, error) {
	return a.parenthesizeExprs("group", expr.Expression)
}

func (a *AstPrinter) VisitLiteralExpr(expr *ast.Literal) (any, error) {
	if expr.Value == nil {
		return "nil", nil
	}
	return fmt.Sprintf("%v", expr.Value), nil
}

func (a *AstPrinter) VisitUnaryExpr(expr *ast.Unary) (any, error) {
	return a.parenthesizeExprs(expr.Operator.Lexeme, expr.Right)
}

func (a *AstPrinter) VisitExpressionStmt(expr *ast.Expression) (any, error) {
	return a.parenthesizeExprs("expr", expr.Expression)
}

func (a *AstPrinter) VisitPrintStmt(expr *ast.Print) (any, error) {
	return a.parenthesizeExprs("print", expr.Expression)
}

func (a *AstPrinter) VisitVarStmt(stmt *ast.Var) (any, error) {
	if stmt.Initializer != nil {
		return a.parenthesizeExprs("var "+stmt.Name.Lexeme, stmt.Initializer)
	}
	return "(var " + stmt.Name.Lexeme + ")", nil
}

func (a *AstPrinter) VisitVariableExpr(expr *ast.Variable) (any, error) {
	return expr.Name.Lexeme, nil
}

func (a *AstPrinter) VisitAssignExpr(expr *ast.Assign) (any, error) {
	return a.parenthesizeExprs("assign "+expr.Name.Lexeme, expr.Value)
}

func (a *AstPrinter) VisitBlockStmt(stmt *ast.Block) (any, error) {
	return a.parenthesizeStmts("block", stmt.Statements...)
}

func (a *AstPrinter) VisitIfStmt(stmt *ast.If) (any, error) {
	if stmt.ElseBranch != nil {
		return a.parenthesize("if-else", stmt.Condition, stmt.ThenBranch, stmt.ElseBranch)
	}
	return a.parenthesize("if", stmt.Condition, stmt.ThenBranch)
}

func (a *AstPrinter) VisitLogicalExpr(expr *ast.Logical) (any, error) {
	return a.parenthesizeExprs(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (a *AstPrinter) VisitWhileStmt(stmt *ast.While) (any, error) {
	return a.parenthesize("while", stmt.Condition, stmt.Body)
}

func (a *AstPrinter) VisitCallExpr(expr *ast.Call) (any, error) {
	return a.parenthesizeExprs("call", expr.Callee)
}

func (a *AstPrinter) VisitFunctionStmt(stmt *ast.Function) (any, error) {
	parts := []any{stmt.Name}
	for _, param := range stmt.Params {
		parts = append(parts, param)
	}
	for _, bodyStmt := range stmt.Body {
		parts = append(parts, bodyStmt)
	}
	return a.parenthesize("fun", parts...)
}

func (a *AstPrinter) VisitReturnStmt(stmt *ast.Return) (any, error) {
	if stmt.Value != nil {
		return a.parenthesizeExprs("return", stmt.Value)
	}
	return "(return)", nil
}

func (a *AstPrinter) VisitClassStmt(stmt *ast.Class) (any, error) {
	parts := []any{stmt.Name}
	for _, method := range stmt.Methods {
		parts = append(parts, method)
	}
	return a.parenthesize("class", parts...)
}

func (a *AstPrinter) VisitGetExpr(expr *ast.Get) (any, error) {
	return a.parenthesizeExprs("get "+expr.Name.Lexeme, expr.Object)
}

func (a *AstPrinter) VisitSetExpr(expr *ast.Set) (any, error) {
	return a.parenthesizeExprs("set "+expr.Name.Lexeme, expr.Object, expr.Value)
}

func (a *AstPrinter) VisitThisExpr(expr *ast.This) (any, error) {
	return "this", nil
}

func (a *AstPrinter) VisitSuperExpr(expr *ast.Super) (any, error) {
	return a.parenthesizeExprs("super " + expr.Method.Lexeme)
}

func (a *AstPrinter) VisitBreakStmt(stmt *ast.Break) (any, error) {
	return "(break)", nil
}

// Helper methods

func (a *AstPrinter) parenthesizeExprs(name string, exprs ...ast.Expr) (string, error) {
	result := "(" + name
	for _, expr := range exprs {
		subResult, _ := expr.Accept(a)
		result += " " + subResult.(string)
	}
	result += ")"
	return result, nil
}

func (a *AstPrinter) parenthesizeStmts(name string, stmts ...ast.Stmt) (string, error) {
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
		case ast.Expr:
			subResult, _ = v.Accept(a)
		case ast.Stmt:
			subResult, _ = v.Accept(a)
		case nil:
			subResult = "nil"
		default:
			subResult = fmt.Sprintf("%v", v)
		}
		result += " " + subResult.(string)
	}
	result += ")"
	return result, nil
}
