package golox

import "fmt"

type Interpreter struct {
	environment *Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		environment: NewEnvironment(),
	}
}

func (i *Interpreter) Interpret(statements []Stmt) (any, error) {
	for _, stmt := range statements {
		_, err := i.execute(stmt)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (i *Interpreter) execute(stmt Stmt) (any, error) {
	return stmt.Accept(i)
}

func (i *Interpreter) evaluate(expr Expr) (any, error) {
	return expr.Accept(i)
}

func (i *Interpreter) VisitLiteralExpr(e *Literal) (any, error) {
	return e.Value, nil
}

func (i *Interpreter) VisitGroupingExpr(e *Grouping) (any, error) {
	return i.evaluate(e.Expression)
}

func (i *Interpreter) VisitUnaryExpr(e *Unary) (any, error) {
	right, err := i.evaluate(e.Right)
	if err != nil {
		return nil, err
	}

	switch e.Operator.Type {
	case MINUS:
		i.checkFloat64Operand(e.Operator, right)
		return -right.(float64), nil
	case BANG:
		return !isTruthy(right), nil
	}

	return nil, nil
}

func (i *Interpreter) VisitBinaryExpr(e *Binary) (any, error) {
	left, err := i.evaluate(e.Left)
	if err != nil {
		return nil, err
	}
	right, err := i.evaluate(e.Right)
	if err != nil {
		return nil, err
	}

	switch e.Operator.Type {
	case PLUS:
		if l, ok := left.(string); ok {
			if r, ok := right.(string); ok {
				return l + r, nil
			} else {
				err := NewRuntimeError(*e.Operator, "Cannot concatenate string "+left.(string)+" with "+fmt.Sprintf("%T", right))
				return nil, err
			}
		} else if _, ok := left.(float64); ok {
			if _, ok := right.(float64); ok {
				return left.(float64) + right.(float64), nil
			} else {
				err := NewRuntimeError(*e.Operator, "Cannot add number "+fmt.Sprintf("%T", left)+" with "+fmt.Sprintf("%T", right))
				return nil, err
			}
		} else {
			err := NewRuntimeError(*e.Operator, "Cannot add "+fmt.Sprintf("%T", left)+" with "+fmt.Sprintf("%T", right))
			return nil, err
		}
	case MINUS:
		i.checkFloat64Operands(e.Operator, left, right)
		return left.(float64) - right.(float64), nil
	case STAR:
		i.checkFloat64Operands(e.Operator, left, right)
		return left.(float64) * right.(float64), nil
	case SLASH:
		i.checkFloat64Operands(e.Operator, left, right)
		return left.(float64) / right.(float64), nil
	case GREATER:
		i.checkFloat64Operands(e.Operator, left, right)
		return left.(float64) > right.(float64), nil
	case GREATER_EQUAL:
		i.checkFloat64Operands(e.Operator, left, right)
		return left.(float64) >= right.(float64), nil
	case LESS:
		i.checkFloat64Operands(e.Operator, left, right)
		return left.(float64) < right.(float64), nil
	case LESS_EQUAL:
		i.checkFloat64Operands(e.Operator, left, right)
		return left.(float64) <= right.(float64), nil
	case BANG_EQUAL:
		return !isEqual(left, right), nil
	case EQUAL_EQUAL:
		return isEqual(left, right), nil
	}

	return nil, nil
}

func (i *Interpreter) VisitExpressionStmt(e *Expression) (any, error) {
	return i.evaluate(e.Expression)
}

func (i *Interpreter) VisitPrintStmt(e *Print) (any, error) {
	value, err := i.evaluate(e.Expression)
	if err != nil {
		return nil, err
	}
	fmt.Println(stringify(value))
	return nil, nil
}

func (i *Interpreter) VisitVarStmt(e *Var) (any, error) {
	var value any // The default is nil. Another choice would be to raise an error, but Lox allows uninitialized variables.
	var err error
	if e.Initializer != nil {
		value, err = i.evaluate(e.Initializer)
		if err != nil {
			return nil, err
		}
	}
	i.environment.Define(e.Name.Lexeme, value)
	return nil, nil
}

func (i *Interpreter) VisitVariableExpr(e *Variable) (any, error) {
	return i.environment.Get(&e.Name)
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

func isTruthy(object any) bool {
	if object == nil {
		return false
	}
	if b, ok := object.(bool); ok {
		return b
	}
	return true
}

func stringify(object any) string {
	if object == nil {
		return "nil"
	}
	switch v := object.(type) {
	case float64:
		s := fmt.Sprintf("%v", v)
		if s[max(0, len(s)-2):] == ".0" {
			s = s[:len(s)-2]
		}
		return s
	case bool:
		if v {
			return "true"
		}
		return "false"
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

// checkFloat64Operand checks if the operand is float64, returns an error if not.
func (i *Interpreter) checkFloat64Operand(operator *Token, operand any) error {
	if _, ok := operand.(float64); !ok {
		return NewRuntimeError(*operator, "Operand must be a number.")
	}
	return nil
}

// checkFloat64Operands checks if both operands are float64, returns an error if not.
func (i *Interpreter) checkFloat64Operands(operator *Token, left, right any) error {
	if _, ok := left.(float64); !ok {
		return NewRuntimeError(*operator, "Operand "+left.(string)+" must be a number.")
	}
	if _, ok := right.(float64); !ok {
		return NewRuntimeError(*operator, "Operand "+right.(string)+" must be a number.")
	}
	return nil
}
