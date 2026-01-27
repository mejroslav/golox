package golox

import "fmt"

// LoxInstance represents an instance of a Lox class.
type LoxInstance struct {
	Class  *LoxClass
	Fields map[string]any
}

func NewLoxInstance(class *LoxClass) *LoxInstance {
	return &LoxInstance{
		Class:  class,
		Fields: make(map[string]any),
	}
}

// String returns a string representation of the Lox instance.
func (li *LoxInstance) String() string {
	return "<instance of " + li.Class.Name + ">"
}

// Get retrieves a property or method from the instance.
func (li *LoxInstance) Get(name Token) (any, error) {
	if value, ok := li.Fields[name.Lexeme]; ok {
		return value, nil
	}

	if method, ok := li.FindMethod(name.Lexeme); ok {
		return method.Bind(li), nil
	}

	err := NewRuntimeError(name, fmt.Sprintf("Class '%s' has not defined property '%s'.", li.Class.Name, name.Lexeme))
	return nil, err
}

// Set assigns a value to a property of the instance.
func (li *LoxInstance) Set(name Token, value any) {
	li.Fields[name.Lexeme] = value
}

// FindMethod looks up a method by name in the instance's class.
func (li *LoxInstance) FindMethod(name string) (*LoxFunction, bool) {
	return li.Class.GetMethod(name)
}
