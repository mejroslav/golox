package golox

import (
	"fmt"
)

// Interpreter interprets and executes Lox code.
type Interpreter struct {
	globals     *Environment // The global environment
	environment *Environment // The current environment
	locals      map[Expr]int // Maps expressions to their scope depth
}

func NewInterpreter() *Interpreter {
	globals := NewEnvironment(nil)

	// Add built-in functions to the global environment
	clockCallable := Clock{}
	globals.Define("clock", clockCallable)

	environment := globals
	return &Interpreter{
		globals:     globals,
		environment: environment,
		locals:      make(map[Expr]int),
	}
}

// Interpret interprets and executes a list of statements.
func (i *Interpreter) Interpret(statements []Stmt) (any, error) {
	for _, stmt := range statements {
		_, err := i.execute(stmt)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
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
	return i.lookupVariable(*e.Name, e)
}

func (i *Interpreter) lookupVariable(name Token, expr Expr) (any, error) {
	distance, ok := i.locals[expr]
	if ok {
		return i.environment.GetAt(distance, name.Lexeme)
	} else {
		return i.globals.Get(&name)
	}
}

func (i *Interpreter) VisitAssignExpr(e *Assign) (any, error) {
	value, err := i.evaluate(e.Value)
	if err != nil {
		return nil, err
	}

	distance, ok := i.locals[e]
	if ok {
		err = i.environment.AssignAt(distance, e.Name, value)
	} else {
		err = i.globals.Assign(e.Name, value)
	}
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (i *Interpreter) VisitBlockStmt(stmt *Block) (any, error) {
	return i.executeBlock(stmt.Statements, NewEnvironment(i.environment))
}

func (i *Interpreter) VisitClassStmt(stmt *Class) (any, error) {
	i.environment.Define(stmt.Name.Lexeme, nil)
	var loxClass *LoxClass = NewLoxClass(stmt, i.environment)
	i.environment.Assign(stmt.Name, loxClass)
	return nil, nil
}

func (i *Interpreter) VisitIfStmt(stmt *If) (any, error) {
	condition, err := i.evaluate(stmt.Condition)
	if err != nil {
		return nil, err
	}
	if isTruthy(condition) {
		return i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		return i.execute(stmt.ElseBranch)
	}
	return nil, nil
}

func (i *Interpreter) VisitLogicalExpr(expr *Logical) (any, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	if expr.Operator.Type == OR {
		if isTruthy(left) {
			return left, nil
		}
	} else {
		if !isTruthy(left) {
			return left, nil
		}
	}

	return i.evaluate(expr.Right)
}

func (i *Interpreter) VisitWhileStmt(stmt *While) (any, error) {
	for {
		condition, err := i.evaluate(stmt.Condition)
		if err != nil {
			return nil, err
		}
		if !isTruthy(condition) {
			break
		}
		_, err = i.execute(stmt.Body)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (i *Interpreter) VisitCallExpr(expr *Call) (any, error) {
	callee, err := i.evaluate(expr.Callee)
	if err != nil {
		return nil, err
	}

	arguments := []any{}
	for _, argument := range expr.Arguments {
		arg, err := i.evaluate(argument)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, arg)
	}

	if _, ok := callee.(LoxCallable); !ok {
		return nil, NewRuntimeError(*expr.Paren, "Can only call functions and classes.")
	}

	function, ok := callee.(LoxCallable)
	if !ok {
		return nil, NewRuntimeError(*expr.Paren, "Can only call functions and classes.")
	}

	if len(arguments) != function.Arity() {
		return nil, NewRuntimeError(*expr.Paren, fmt.Sprintf("Expected %d arguments but got %d.", function.Arity(), len(arguments)))
	}

	return function.Call(i, arguments)
}

func (i *Interpreter) VisitGetExpr(expr *Get) (any, error) {
	object, err := i.evaluate(expr.Object)
	if err != nil {
		return nil, err
	}

	loxInstance, ok := object.(*LoxInstance)
	if !ok {
		return nil, NewRuntimeError(*expr.Name, "Only instances have properties.")
	}

	return loxInstance.Get(*expr.Name)
}

func (i *Interpreter) VisitFunctionStmt(stmt *Function) (any, error) {
	function := NewLoxFunction(stmt, i.environment)
	i.environment.Define(stmt.Name.Lexeme, function)
	return nil, nil
}

func (i *Interpreter) VisitReturnStmt(stmt *Return) (any, error) {
	var value any
	var err error
	if stmt.Value != nil {
		value, err = i.evaluate(stmt.Value)
		if err != nil {
			return nil, err
		}
	}
	return nil, &ReturnValue{Value: value}
}

// ---------------------------------------------------------------------
// Execute a statement

func (i *Interpreter) execute(stmt Stmt) (any, error) {
	return stmt.Accept(i)
}

func (i *Interpreter) resolve(expr Expr, depth int) error {
	i.locals[expr] = depth
	return nil
}

func (i *Interpreter) evaluate(expr Expr) (any, error) {
	return expr.Accept(i)
}

func (i *Interpreter) executeBlock(statements []Stmt, environment *Environment) (any, error) {
	previous := i.environment
	i.environment = environment
	defer func() {
		i.environment = previous
	}()

	for _, statement := range statements {
		_, err := i.execute(statement)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

// ---------------------------------------------------------------------
// Helpers

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

// stringify converts an object to its string representation.
//
// For nil, it returns "nil".
// For float64, it removes the decimal part if it's zero.
// For bool, it returns "true" or "false".
// For string, it returns the string itself.
// For other types, it uses fmt.Sprintf to convert to string.
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
