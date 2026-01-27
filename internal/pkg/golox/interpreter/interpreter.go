package interpreter

import (
	"fmt"
	"mejroslav/golox/v2/internal/pkg/golox/ast"
	"mejroslav/golox/v2/internal/pkg/golox/lox_error"
	"mejroslav/golox/v2/internal/pkg/golox/token"
	"mejroslav/golox/v2/internal/pkg/golox/types"
)

// Interpreter interprets and executes Lox code.
type Interpreter struct {
	globals     *Environment     // The global environment
	environment *Environment     // The current environment
	locals      map[ast.Expr]int // Maps ast.Expressions to their scope depth
}

func NewInterpreter() *Interpreter {
	globals := NewEnvironment(nil)

	// Add built-in functions to the global environment
	clockCallable := &Clock{}
	globals.Define("clock", clockCallable)
	inputCallable := &Input{}
	globals.Define("input", inputCallable)

	environment := globals
	return &Interpreter{
		globals:     globals,
		environment: environment,
		locals:      make(map[ast.Expr]int),
	}
}

// Interpret interprets and executes a list of statements.
func (i *Interpreter) Interpret(statements []ast.Stmt) (any, error) {
	for _, stmt := range statements {
		_, err := i.execute(stmt)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (i *Interpreter) VisitLiteralExpr(e *ast.Literal) (any, error) {
	return e.Value, nil
}

func (i *Interpreter) VisitGroupingExpr(e *ast.Grouping) (any, error) {
	return i.evaluate(e.Expression)
}

func (i *Interpreter) VisitUnaryExpr(e *ast.Unary) (any, error) {
	right, err := i.evaluate(e.Right)
	if err != nil {
		return nil, err
	}

	switch e.Operator.Type {
	case token.MINUS:
		i.checkFloat64Operand(e.Operator, right)
		return -right.(float64), nil
	case token.BANG:
		return !isTruthy(right), nil
	}

	return nil, nil
}

func (i *Interpreter) VisitBinaryExpr(e *ast.Binary) (any, error) {
	left, err := i.evaluate(e.Left)
	if err != nil {
		return nil, err
	}
	right, err := i.evaluate(e.Right)
	if err != nil {
		return nil, err
	}

	switch e.Operator.Type {
	case token.PLUS:
		if l, ok := left.(string); ok {
			if r, ok := right.(string); ok {
				return l + r, nil
			} else {
				err := lox_error.NewRuntimeError(*e.Operator, "Cannot concatenate string "+left.(string)+" with "+fmt.Sprintf("%T", right))
				return nil, err
			}
		} else if _, ok := left.(float64); ok {
			if _, ok := right.(float64); ok {
				return left.(float64) + right.(float64), nil
			} else {
				err := lox_error.NewRuntimeError(*e.Operator, "Cannot add number "+fmt.Sprintf("%T", left)+" with "+fmt.Sprintf("%T", right))
				return nil, err
			}
		} else {
			err := lox_error.NewRuntimeError(*e.Operator, "Cannot add "+fmt.Sprintf("%T", left)+" with "+fmt.Sprintf("%T", right))
			return nil, err
		}
	case token.MINUS:
		i.checkFloat64Operands(e.Operator, left, right)
		return left.(float64) - right.(float64), nil
	case token.STAR:
		i.checkFloat64Operands(e.Operator, left, right)
		return left.(float64) * right.(float64), nil
	case token.SLASH:
		i.checkFloat64Operands(e.Operator, left, right)
		return left.(float64) / right.(float64), nil
	case token.GREATER:
		i.checkFloat64Operands(e.Operator, left, right)
		return left.(float64) > right.(float64), nil
	case token.GREATER_EQUAL:
		i.checkFloat64Operands(e.Operator, left, right)
		return left.(float64) >= right.(float64), nil
	case token.LESS:
		i.checkFloat64Operands(e.Operator, left, right)
		return left.(float64) < right.(float64), nil
	case token.LESS_EQUAL:
		i.checkFloat64Operands(e.Operator, left, right)
		return left.(float64) <= right.(float64), nil
	case token.BANG_EQUAL:
		return !isEqual(left, right), nil
	case token.EQUAL_EQUAL:
		return isEqual(left, right), nil
	}

	return nil, nil
}

func (i *Interpreter) VisitExpressionStmt(e *ast.Expression) (any, error) {
	return i.evaluate(e.Expression)
}

func (i *Interpreter) VisitPrintStmt(e *ast.Print) (any, error) {
	value, err := i.evaluate(e.Expression)
	if err != nil {
		return nil, err
	}
	fmt.Println(stringify(value))
	return nil, nil
}

func (i *Interpreter) VisitVarStmt(e *ast.Var) (any, error) {
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

func (i *Interpreter) VisitVariableExpr(e *ast.Variable) (any, error) {
	return i.lookupVariable(*e.Name, e)
}

func (i *Interpreter) lookupVariable(name token.Token, expr ast.Expr) (any, error) {
	distance, ok := i.locals[expr]
	if ok {
		return i.environment.GetAt(distance, name.Lexeme)
	} else {
		return i.globals.Get(&name)
	}
}

func (i *Interpreter) VisitAssignExpr(e *ast.Assign) (any, error) {
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

func (i *Interpreter) VisitBlockStmt(stmt *ast.Block) (any, error) {
	return i.executeBlock(stmt.Statements, NewEnvironment(i.environment))
}

func (i *Interpreter) VisitClassStmt(stmt *ast.Class) (any, error) {
	var superclass *LoxClass
	if stmt.Superclass != nil {
		superclassValue, err := i.evaluate(stmt.Superclass)
		if err != nil {
			return nil, err
		}
		var ok bool
		superclass, ok = superclassValue.(*LoxClass)
		if !ok {
			return nil, lox_error.NewRuntimeError(*stmt.Superclass.Name, "Superclass '"+stmt.Superclass.Name.Lexeme+"' must be a class.")
		}
	}

	i.environment.Define(stmt.Name.Lexeme, nil)

	if stmt.Superclass != nil {
		// Create a new environment for "super"
		i.environment = NewEnvironment(i.environment)
		i.environment.Define("super", superclass)
	}

	methods := make(map[string]*LoxFunction)
	for _, method := range stmt.Methods {
		if method.Name.Lexeme == "init" {
			function := NewInitializerFunction(&method, i.environment)
			methods[method.Name.Lexeme] = function
		} else {
			function := NewLoxFunction(&method, i.environment)
			methods[method.Name.Lexeme] = function
		}
	}

	var loxClass *LoxClass = NewLoxClass(stmt.Name.Lexeme, superclass, methods)
	i.environment.Assign(stmt.Name, loxClass)

	if stmt.Superclass != nil {
		i.environment = i.environment.GetEnclosing()
	}

	return nil, nil
}

func (i *Interpreter) VisitIfStmt(stmt *ast.If) (any, error) {
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

func (i *Interpreter) VisitLogicalExpr(e *ast.Logical) (any, error) {
	left, err := i.evaluate(e.Left)
	if err != nil {
		return nil, err
	}

	if e.Operator.Type == token.OR {
		if isTruthy(left) {
			return left, nil
		}
	} else {
		if !isTruthy(left) {
			return left, nil
		}
	}

	return i.evaluate(e.Right)
}

func (i *Interpreter) VisitWhileStmt(stmt *ast.While) (any, error) {
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
			if _, ok := err.(*types.BreakValue); ok {
				// Break out of the loop
				break
			}
			return nil, err
		}
	}
	return nil, nil
}

func (i *Interpreter) VisitCallExpr(e *ast.Call) (any, error) {
	callee, err := i.evaluate(e.Callee)
	if err != nil {
		return nil, err
	}

	arguments := []any{}
	for _, argument := range e.Arguments {
		arg, err := i.evaluate(argument)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, arg)
	}

	function, ok := callee.(LoxCallable)
	if !ok {
		return nil, lox_error.NewRuntimeError(*e.Paren, "Can only call functions and classes.")
	}

	if len(arguments) != function.Arity() {
		return nil, lox_error.NewRuntimeError(*e.Paren, fmt.Sprintf("Expected %d arguments but got %d.", function.Arity(), len(arguments)))
	}

	return function.Call(i, arguments)
}

func (i *Interpreter) VisitGetExpr(e *ast.Get) (any, error) {
	object, err := i.evaluate(e.Object)
	if err != nil {
		return nil, err
	}

	loxInstance, ok := object.(*LoxInstance)
	if !ok {
		return nil, lox_error.NewRuntimeError(*e.Name, "Only instances have properties.")
	}

	return loxInstance.Get(*e.Name)
}

func (i *Interpreter) VisitSetExpr(e *ast.Set) (any, error) {
	object, err := i.evaluate(e.Object)
	if err != nil {
		return nil, err
	}

	loxInstance, ok := object.(*LoxInstance)
	if !ok {
		return nil, lox_error.NewRuntimeError(*e.Name, "Only instances have fields.")
	}

	value, err := i.evaluate(e.Value)
	if err != nil {
		return nil, err
	}

	loxInstance.Set(*e.Name, value)
	return value, nil
}

func (i *Interpreter) VisitSuperExpr(e *ast.Super) (any, error) {
	distance, ok := i.locals[e]
	if !ok {
		return nil, lox_error.NewRuntimeError(*e.Method, "Undefined 'super' reference.")
	}

	superclassValue, err := i.environment.GetAt(distance, "super")
	if err != nil {
		return nil, err
	}
	superclass, ok := superclassValue.(*LoxClass)
	if !ok {
		return nil, lox_error.NewRuntimeError(*e.Method, "'super' is not a class.")
	}

	// We can't access 'this' directly from the environment because 'this' is stored
	// in the enclosing environment (one level up).
	objectValue, err := i.environment.GetAt(distance-1, "this")
	if err != nil {
		return nil, err
	}
	loxInstance, ok := objectValue.(*LoxInstance)
	if !ok {
		return nil, lox_error.NewRuntimeError(*e.Method, "'this' is not an instance.")
	}

	method, ok := superclass.GetMethod(e.Method.Lexeme)
	if !ok {
		return nil, lox_error.NewRuntimeError(*e.Method, fmt.Sprintf("Undefined property '%s'.", e.Method.Lexeme))
	}

	return method.Bind(loxInstance), nil
}

func (i *Interpreter) VisitThisExpr(e *ast.This) (any, error) {
	return i.lookupVariable(*e.Keyword, e)
}

func (i *Interpreter) VisitFunctionStmt(stmt *ast.Function) (any, error) {
	function := NewLoxFunction(stmt, i.environment)
	i.environment.Define(stmt.Name.Lexeme, function)
	return nil, nil
}

func (i *Interpreter) VisitReturnStmt(stmt *ast.Return) (any, error) {
	var value any
	var err error
	if stmt.Value != nil {
		value, err = i.evaluate(stmt.Value)
		if err != nil {
			return nil, err
		}
	}
	return nil, &types.ReturnValue{Value: value}
}

func (i *Interpreter) VisitBreakStmt(stmt *ast.Break) (any, error) {
	return nil, &types.BreakValue{Keyword: stmt.Keyword}
}

// ---------------------------------------------------------------------
// Execute a statement

func (i *Interpreter) execute(stmt ast.Stmt) (any, error) {
	return stmt.Accept(i)
}

func (i *Interpreter) Resolve(e ast.Expr, depth int) error {
	i.locals[e] = depth
	return nil
}

func (i *Interpreter) evaluate(e ast.Expr) (any, error) {
	return e.Accept(i)
}

func (i *Interpreter) executeBlock(statements []ast.Stmt, environment *Environment) (any, error) {
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
func (i *Interpreter) checkFloat64Operand(operator *token.Token, operand any) error {
	if _, ok := operand.(float64); !ok {
		return lox_error.NewRuntimeError(*operator, "Operand must be a number.")
	}
	return nil
}

// checkFloat64Operands checks if both operands are float64, returns an error if not.
func (i *Interpreter) checkFloat64Operands(operator *token.Token, left, right any) error {
	if _, ok := left.(float64); !ok {
		return lox_error.NewRuntimeError(*operator, "Operand "+left.(string)+" must be a number.")
	}
	if _, ok := right.(float64); !ok {
		return lox_error.NewRuntimeError(*operator, "Operand "+right.(string)+" must be a number.")
	}
	return nil
}
