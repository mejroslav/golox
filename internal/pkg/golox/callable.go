package golox

type LoxCallable interface {
	Arity() int
	Call(interpreter *Interpreter, arguments []any) (any, error)
}

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

func (lf *LoxFunction) Arity() int {
	return len(lf.Declaration.Params)
}

func (lf *LoxFunction) String() string {
	return "<fn " + lf.Declaration.Name.Lexeme + ">"
}

func (lf *LoxFunction) Call(interpreter *Interpreter, arguments []any) (any, error) {
	environment := NewEnvironment(lf.Closure)
	for i, param := range lf.Declaration.Params {
		environment.Define(param.Lexeme, arguments[i])
	}

	_, err := interpreter.executeBlock(lf.Declaration.Body, environment)
	if err != nil {
		if returnErr, ok := err.(*ReturnValue); ok {
			return returnErr.Value, nil
		}
		return nil, err
	}
	return nil, nil
}
