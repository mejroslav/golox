package golox

type Interpreter struct {
}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

func (i *Interpreter) VisitLiteralExpr(e *Literal) any {
	return e.Value
}

func (i *Interpreter) VisitGroupingExpr(e *Grouping) any {
	return i.evaluate(e.Expression)
}

func (i *Interpreter) evaluate(expr Expr) any {
	return expr.Accept(i)
}

func (i *Interpreter) VisitUnaryExpr(e *Unary) any {
	right := i.evaluate(e.Right)

	switch e.Operator.Type {
	case MINUS:
		return -right.(float64)
	case BANG:
		return !isTruthy(right)
	}

	return nil
}

func isTruthy(object any) bool {
	if object == nil {
		return false
	}
	if b, ok := object.(bool); ok {
		return b
	}
	return true
}

func (i *Interpreter) VisitBinaryExpr(e *Binary) any {
	left := i.evaluate(e.Left)
	right := i.evaluate(e.Right)

	switch e.Operator.Type {
	case PLUS:
		if l, ok := left.(string); ok {
			if r, ok := right.(string); ok {
				return l + r
			}
		} else {
			return left.(float64) + right.(float64)
		}
	case MINUS:
		return left.(float64) - right.(float64)
	case STAR:
		return left.(float64) * right.(float64)
	case SLASH:
		return left.(float64) / right.(float64)
	case GREATER:
		return left.(float64) > right.(float64)
	case GREATER_EQUAL:
		return left.(float64) >= right.(float64)
	case LESS:
		return left.(float64) < right.(float64)
	case LESS_EQUAL:
		return left.(float64) <= right.(float64)
	case BANG_EQUAL:
		return !isEqual(left, right)
	case EQUAL_EQUAL:
		return isEqual(left, right)
	}

	return nil
}

func isEqual(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	return a == b
}
