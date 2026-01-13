package golox

type Environment struct {
	values map[string]any
}

func NewEnvironment() *Environment {
	return &Environment{
		values: make(map[string]any),
	}
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) Get(name *Token) (any, error) {
	if value, ok := e.values[name.Lexeme]; ok {
		return value, nil
	}
	return nil, RuntimeError{
		Token:   *name,
		Message: "Undefined variable '" + name.Lexeme + "'.",
	}
}

func (e *Environment) Assign(name *Token, value any) error {
	if _, ok := e.values[name.Lexeme]; ok {
		e.values[name.Lexeme] = value
		return nil
	}
	return RuntimeError{
		Token:   *name,
		Message: "Undefined variable '" + name.Lexeme + "'.",
	}
}
