package golox

import (
	"mejroslav/golox/v2/internal/pkg/utils"
)

// Resolver performs static analysis to resolve variable bindings.
// It determines the scope depth of each variable and informs the interpreter.
type Resolver struct {
	interpreter      *Interpreter // The interpreter to resolve variables for
	scopeStack       *utils.Stack // Stack of scopes. Each scope is a map of variable names to their defined status
	currentFunction  FunctionType // The type of the current function being resolved
	currentClass     ClassType    // The type of the current class being resolved
	currentLoopDepth int          // The current depth of nested loops
}

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{
		interpreter:      interpreter,
		scopeStack:       utils.NewStack(),
		currentFunction:  FT_NONE,
		currentClass:     CT_NONE,
		currentLoopDepth: 0,
	}
}

func (r *Resolver) Resolve(statements []Stmt) ([]Stmt, error) {
	for _, statement := range statements {
		if err := r.resolveStmt(statement); err != nil {
			return nil, err
		}
	}
	return statements, nil
}

func (r *Resolver) BeginScope() {
	r.scopeStack.Push(make(map[string]bool))
}

func (r *Resolver) EndScope() {
	r.scopeStack.Pop()
}

func (r *Resolver) resolveStmt(statement Stmt) error {
	_, err := statement.Accept(r)
	return err
}

func (r *Resolver) resolveExpr(expression Expr) error {
	_, err := expression.Accept(r)
	return err
}

func (r *Resolver) VisitBlockStmt(s *Block) (any, error) {
	r.BeginScope()
	for _, stmt := range s.Statements {
		if err := r.resolveStmt(stmt); err != nil {
			return nil, err
		}
	}
	r.EndScope()
	return nil, nil
}

func (r *Resolver) VisitClassStmt(stmt *Class) (any, error) {
	enclosingClass := r.currentClass
	r.currentClass = CT_CLASS
	lastLoopDepth := r.currentLoopDepth
	r.currentLoopDepth = 0

	err := r.declare(stmt.Name)
	if err != nil {
		return nil, err
	}

	err = r.define(stmt.Name)
	if err != nil {
		return nil, err
	}

	if stmt.Superclass != nil && stmt.Name.Lexeme == stmt.Superclass.Name.Lexeme {
		return nil, NewRuntimeError(*stmt.Superclass.Name, "A class cannot inherit from itself.")
	} else if stmt.Superclass != nil {
		r.currentClass = CT_SUBCLASS
		r.resolveExpr(stmt.Superclass)
	}

	if stmt.Superclass != nil {
		// r.currentClass = CT_SUBCLASS
		r.BeginScope() // Scope for "super"
		r.scopeStack.Peek().(map[string]bool)["super"] = true
	}

	r.BeginScope() // Scope for "this"
	r.scopeStack.Peek().(map[string]bool)["this"] = true

	for _, method := range stmt.Methods {
		functionType := FT_METHOD
		if method.Name.Lexeme == "init" {
			functionType = FT_INITIALIZER
		}

		err = r.resolveFunction(&method, functionType)
		if err != nil {
			return nil, err
		}
	}

	r.EndScope() // End scope for "this"

	if stmt.Superclass != nil {
		r.EndScope() // End scope for "super"
	}
	r.currentClass = enclosingClass
	r.currentLoopDepth = lastLoopDepth
	return nil, nil
}

func (r *Resolver) VisitVarStmt(stmt *Var) (any, error) {
	err := r.declare(stmt.Name)
	if err != nil {
		return nil, err
	}

	if stmt.Initializer != nil {
		if err := r.resolveExpr(stmt.Initializer); err != nil {
			return nil, err
		}
	}

	err = r.define(stmt.Name)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *Resolver) VisitVariableExpr(expr *Variable) (any, error) {
	if !r.scopeStack.IsEmpty() {
		scope := r.scopeStack.Peek().(map[string]bool)
		if defined, ok := scope[expr.Name.Lexeme]; ok && !defined {
			return nil, NewRuntimeError(*expr.Name, "Cannot read local variable in its own initializer.")
		}
	}

	r.resolveLocal(expr, expr.Name)

	return nil, nil
}

func (r *Resolver) VisitAssignExpr(expr *Assign) (any, error) {
	if err := r.resolveExpr(expr.Value); err != nil {
		return nil, err
	}

	r.resolveLocal(expr, expr.Name)

	return nil, nil
}

func (r *Resolver) VisitFunctionStmt(stmt *Function) (any, error) {
	lastLoopDepth := r.currentLoopDepth
	r.currentLoopDepth = 0
	err := r.declare(stmt.Name)
	if err != nil {
		return nil, err
	}

	err = r.define(stmt.Name)
	if err != nil {
		return nil, err
	}

	err = r.resolveFunction(stmt, FT_FUNCTION)
	if err != nil {
		return nil, err
	}

	r.currentLoopDepth = lastLoopDepth
	return nil, nil
}

func (r *Resolver) VisitExpressionStmt(stmt *Expression) (any, error) {
	if err := r.resolveExpr(stmt.Expression); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitPrintStmt(stmt *Print) (any, error) {
	if err := r.resolveExpr(stmt.Expression); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitIfStmt(stmt *If) (any, error) {
	if err := r.resolveExpr(stmt.Condition); err != nil {
		return nil, err
	}
	if err := r.resolveStmt(stmt.ThenBranch); err != nil {
		return nil, err
	}
	if stmt.ElseBranch != nil {
		if err := r.resolveStmt(stmt.ElseBranch); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (r *Resolver) VisitWhileStmt(stmt *While) (any, error) {
	r.currentLoopDepth++
	defer func() {
		r.currentLoopDepth--
	}()

	if err := r.resolveExpr(stmt.Condition); err != nil {
		return nil, err
	}
	if err := r.resolveStmt(stmt.Body); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitBreakStmt(stmt *Break) (any, error) {
	if r.currentLoopDepth == 0 {
		return nil, NewRuntimeError(*stmt.Keyword, "Cannot use 'break' outside of a loop.")
	}
	return nil, nil
}

func (r *Resolver) VisitReturnStmt(stmt *Return) (any, error) {
	if r.currentFunction == FT_NONE {
		return nil, NewRuntimeError(*stmt.Keyword, "Cannot return from top-level code.")
	}
	if stmt.Value != nil {
		if r.currentFunction == FT_INITIALIZER {
			return nil, NewRuntimeError(*stmt.Keyword, "Cannot return a value from an initializer.")
		}
		if err := r.resolveExpr(stmt.Value); err != nil {
			return nil, err
		}
	}
	return nil, nil
}
func (r *Resolver) VisitBinaryExpr(expr *Binary) (any, error) {
	if err := r.resolveExpr(expr.Left); err != nil {
		return nil, err
	}
	if err := r.resolveExpr(expr.Right); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitCallExpr(expr *Call) (any, error) {
	if err := r.resolveExpr(expr.Callee); err != nil {
		return nil, err
	}
	for _, argument := range expr.Arguments {
		if err := r.resolveExpr(argument); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (r *Resolver) VisitGetExpr(expr *Get) (any, error) {
	if err := r.resolveExpr(expr.Object); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitSetExpr(expr *Set) (any, error) {
	if err := r.resolveExpr(expr.Value); err != nil {
		return nil, err
	}
	if err := r.resolveExpr(expr.Object); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitSuperExpr(expr *Super) (any, error) {
	if r.currentClass == CT_NONE {
		return nil, NewRuntimeError(*expr.Keyword, "Cannot use 'super' outside of a class.")
	} else if r.currentClass != CT_SUBCLASS {
		return nil, NewRuntimeError(*expr.Keyword, "Cannot use 'super' in a class with no superclass.")
	}
	r.resolveLocal(expr, expr.Keyword)
	return nil, nil
}

func (r *Resolver) VisitThisExpr(expr *This) (any, error) {
	if r.currentClass == CT_NONE {
		return nil, NewRuntimeError(*expr.Keyword, "Cannot use 'this' outside of a class.")
	}
	r.resolveLocal(expr, expr.Keyword)
	return nil, nil
}

func (r *Resolver) VisitGroupingExpr(expr *Grouping) (any, error) {
	if err := r.resolveExpr(expr.Expression); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitLiteralExpr(expr *Literal) (any, error) {
	return nil, nil
}

func (r *Resolver) VisitUnaryExpr(expr *Unary) (any, error) {
	if err := r.resolveExpr(expr.Right); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitLogicalExpr(expr *Logical) (any, error) {
	if err := r.resolveExpr(expr.Left); err != nil {
		return nil, err
	}
	if err := r.resolveExpr(expr.Right); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) declare(name *Token) error {
	if r.scopeStack.IsEmpty() {
		return nil
	}
	scope := r.scopeStack.Peek().(map[string]bool)
	if _, ok := scope[name.Lexeme]; ok {
		return NewRuntimeError(*name, "Variable with name '"+name.Lexeme+"' already declared in this scope.")
	}
	scope[name.Lexeme] = false
	return nil
}

func (r *Resolver) define(name *Token) error {
	if r.scopeStack.IsEmpty() {
		return nil
	}
	scope := r.scopeStack.Peek().(map[string]bool)
	scope[name.Lexeme] = true
	return nil
}

// resolveLocal resolves a local variable by determining its scope depth.
func (r *Resolver) resolveLocal(expr Expr, name *Token) error {
	for i := r.scopeStack.Size() - 1; i >= 0; i-- {
		scope, ok := r.scopeStack.Get(i)
		if !ok {
			continue
		}
		if _, ok := scope.(map[string]bool)[name.Lexeme]; ok {
			r.interpreter.resolve(expr, r.scopeStack.Size()-1-i)
			return nil
		}
	}

	return nil
}

func (r *Resolver) resolveFunction(function *Function, functionType FunctionType) error {
	enclosingFunction := r.currentFunction
	r.currentFunction = functionType

	r.BeginScope()
	for _, param := range function.Params {
		err := r.declare(param)
		if err != nil {
			return err
		}
		err = r.define(param)
		if err != nil {
			return err
		}
	}
	for _, bodyStmt := range function.Body {
		if err := r.resolveStmt(bodyStmt); err != nil {
			return err
		}
	}
	r.EndScope()

	r.currentFunction = enclosingFunction
	return nil
}
