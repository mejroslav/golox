package golox

// LoxCallable represents any callable entity in the Lox language,
// such as functions and classes.
type LoxCallable interface {
	Arity() int                                                  // number of expected arguments
	Call(interpreter *Interpreter, arguments []any) (any, error) // execute the callable
}
