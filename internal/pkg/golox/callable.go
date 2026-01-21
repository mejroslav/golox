package golox

// LoxCallable represents any callable entity in the Lox language,
// such as functions and classes.
type LoxCallable interface {
	Arity() int                                                  // number of expected arguments
	Call(interpreter *Interpreter, arguments []any) (any, error) // execute the callable
}

// LoxFunction represents a user-defined function in the Lox language.
type LoxFunction struct {
	Declaration *Function
	Closure     *Environment
}

func NewLoxFunction(declaration *Function, closure *Environment) *LoxFunction {
	return &LoxFunction{
		Declaration: declaration,
		Closure:     closure,
	}
}

// Arity returns the number of parameters the function expects.
func (lf *LoxFunction) Arity() int {
	return len(lf.Declaration.Params)
}

// String returns a string representation of the function.
func (lf *LoxFunction) String() string {
	return "<fn " + lf.Declaration.Name.Lexeme + ">"
}

// Call executes the function with the given arguments.
func (lf *LoxFunction) Call(interpreter *Interpreter, arguments []any) (any, error) {
	// Create a new environment for the function execution
	// with the function's closure as its parent.
	// This allows the function to access variables from its defining scope.
	environment := NewEnvironment(lf.Closure)
	for i, param := range lf.Declaration.Params {
		environment.Define(param.Lexeme, arguments[i])
	}

	_, err := interpreter.executeBlock(lf.Declaration.Body, environment)
	if err != nil {

		// 'return' statement can be anywhere in the function body.
		// ReturnValue is used to handle return statements in functions
		// and to propagate the return value up the call stack.
		// If we catch a ReturnValue error, we extract the value and return it.
		if returnErr, ok := err.(*ReturnValue); ok {
			return returnErr.Value, nil
		}
		return nil, err
	}
	return nil, nil
}
