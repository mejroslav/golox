package golox

import "fmt"

type Interpreter struct {
}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

func (i *Interpreter) Interpret(expr Expr) (any, error) {
	defer func() {
		if r := recover(); r != nil {
			if runtimeErr, ok := r.(RuntimeError); ok {
				// Handle the runtime error as needed, e.g., log it or return it
				_ = runtimeErr
			} else {
				panic(r) // re-panic if it's not a RuntimeError
			}
		}
	}()

	result := i.evaluate(expr)
	return result, nil
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
		checkNumberOperand(e.Operator, right)
		return -right.(float64)
	case BANG:
		return !isTruthy(right)
	}

	return nil
}

func checkNumberOperand(operator *Token, operand any) {
	if _, ok := operand.(float64); !ok {
		err := NewRuntimeError(*operator, "Operand must be a number.")
		fmt.Printf("%s", err.Error())
		panic(err)
	}
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
			} else {
				err := NewRuntimeError(*e.Operator, "Cannot concatenate string "+left.(string)+" with "+fmt.Sprintf("%T", right))
				fmt.Printf("%s", err.Error())
				panic(err)
			}
		} else if _, ok := left.(float64); ok {
			if _, ok := right.(float64); ok {
				return left.(float64) + right.(float64)
			} else {
				err := NewRuntimeError(*e.Operator, "Cannot add number "+fmt.Sprintf("%T", left)+" with "+fmt.Sprintf("%T", right))
				fmt.Printf("%s", err.Error())
				panic(err)
			}
		} else {
			err := NewRuntimeError(*e.Operator, "Cannot add "+fmt.Sprintf("%T", left)+" with "+fmt.Sprintf("%T", right))
			fmt.Printf("%s", err.Error())
			panic(err)
		}
	case MINUS:
		checkNumberOperands(e.Operator, left, right)
		return left.(float64) - right.(float64)
	case STAR:
		checkNumberOperands(e.Operator, left, right)
		return left.(float64) * right.(float64)
	case SLASH:
		checkNumberOperands(e.Operator, left, right)
		return left.(float64) / right.(float64)
	case GREATER:
		checkNumberOperands(e.Operator, left, right)
		return left.(float64) > right.(float64)
	case GREATER_EQUAL:
		checkNumberOperands(e.Operator, left, right)
		return left.(float64) >= right.(float64)
	case LESS:
		checkNumberOperands(e.Operator, left, right)
		return left.(float64) < right.(float64)
	case LESS_EQUAL:
		checkNumberOperands(e.Operator, left, right)
		return left.(float64) <= right.(float64)
	case BANG_EQUAL:
		return !isEqual(left, right)
	case EQUAL_EQUAL:
		return isEqual(left, right)
	}

	return nil
}

func checkNumberOperands(operator *Token, left, right any) {
	if _, ok := left.(float64); !ok {
		err := NewRuntimeError(*operator, "Operand "+left.(string)+" must be a number.")
		fmt.Printf("%s", err.Error())
		panic(err)
	}
	if _, ok := right.(float64); !ok {
		err := NewRuntimeError(*operator, "Operand "+right.(string)+" must be a number.")
		fmt.Printf("%s", err.Error())
		panic(err)
	}
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
