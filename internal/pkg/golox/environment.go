package golox

type Environment struct {
	enclosing *Environment
	values    map[string]any
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		enclosing: enclosing,
		values:    make(map[string]any),
	}
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) Get(name *Token) (any, error) {
	if value, ok := e.values[name.Lexeme]; ok {
		return value, nil
	}

	if e.enclosing != nil {
		// Recursive lookup in the enclosing environment
		return e.enclosing.Get(name)
	}

	return nil, RuntimeError{
		Token:   *name,
		Message: "Undefined variable '" + name.Lexeme + "'.",
	}
}

func (e *Environment) GetAt(distance int, name string) (any, error) {
	environment := e.ancestor(distance)
	if value, ok := environment.values[name]; ok {
		return value, nil
	}
	return nil, RuntimeError{
		Message: "Undefined variable '" + name + "'.",
	}
}

func (e *Environment) ancestor(distance int) *Environment {
	for i := 0; i < distance; i++ {
		e = e.enclosing
	}
	return e
}

func (e *Environment) Assign(name *Token, value any) error {
	if _, ok := e.values[name.Lexeme]; ok {
		e.values[name.Lexeme] = value
		return nil
	}

	if e.enclosing != nil {
		// Recursive assignment in the enclosing environment
		return e.enclosing.Assign(name, value)
	}

	return RuntimeError{
		Token:   *name,
		Message: "Undefined variable '" + name.Lexeme + "'.",
	}
}

func (e *Environment) AssignAt(distance int, name *Token, value any) error {
	environment := e.ancestor(distance)
	if _, ok := environment.values[name.Lexeme]; ok {
		environment.values[name.Lexeme] = value
		return nil
	}
	return RuntimeError{
		Token:   *name,
		Message: "Undefined variable '" + name.Lexeme + "'.",
	}
}
