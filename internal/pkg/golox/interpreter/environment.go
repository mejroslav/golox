package interpreter

import (
	"mejroslav/golox/v2/internal/pkg/golox/lox_error"
	"mejroslav/golox/v2/internal/pkg/golox/token"
)

// Environment represents a variable scope in the Lox language.
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

// Define adds a new variable to the environment.
func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

// Get retrieves the value of a variable from the environment.
// It searches recursively in enclosing environments if the variable is not found.
func (e *Environment) Get(name *token.Token) (any, error) {
	if value, ok := e.values[name.Lexeme]; ok {
		return value, nil
	}

	if e.enclosing != nil {
		// Recursive lookup in the enclosing environment
		return e.enclosing.Get(name)
	}

	return nil, lox_error.RuntimeError{
		Token:   *name,
		Message: "Undefined variable '" + name.Lexeme + "'.",
	}
}

// GetAt retrieves the value of a variable at a specific distance
// from the current environment.
func (e *Environment) GetAt(distance int, name string) (any, error) {
	environment := e.ancestor(distance)
	if value, ok := environment.values[name]; ok {
		return value, nil
	}
	return nil, lox_error.RuntimeError{
		Message: "Undefined variable '" + name + "'.",
	}
}

func (e *Environment) ancestor(distance int) *Environment {
	for i := 0; i < distance; i++ {
		e = e.enclosing
	}
	return e
}

// Assign updates the value of an existing variable in the environment.
// It searches recursively in enclosing environments if the variable is not found.
func (e *Environment) Assign(name *token.Token, value any) error {
	if _, ok := e.values[name.Lexeme]; ok {
		e.values[name.Lexeme] = value
		return nil
	}

	if e.enclosing != nil {
		// Recursive assignment in the enclosing environment
		return e.enclosing.Assign(name, value)
	}

	return lox_error.RuntimeError{
		Token:   *name,
		Message: "Undefined variable '" + name.Lexeme + "'.",
	}
}

// AssignAt updates the value of a variable at a specific distance
// from the current environment.
func (e *Environment) AssignAt(distance int, name *token.Token, value any) error {
	environment := e.ancestor(distance)
	if _, ok := environment.values[name.Lexeme]; ok {
		environment.values[name.Lexeme] = value
		return nil
	}
	return lox_error.RuntimeError{
		Token:   *name,
		Message: "Undefined variable '" + name.Lexeme + "'.",
	}
}

// GetEnclosing returns the enclosing environment.
func (e *Environment) GetEnclosing() *Environment {
	return e.enclosing
}
