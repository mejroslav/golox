package golox

import "fmt"

// LoxInstance represents an instance of a Lox class.
type LoxInstance struct {
	class  *LoxClass
	fields map[string]any
}

func NewLoxInstance(class *LoxClass) *LoxInstance {
	return &LoxInstance{
		class:  class,
		fields: make(map[string]any),
	}
}

func (li *LoxInstance) String() string {
	return "<instance of " + li.class.Name + ">"
}

func (li *LoxInstance) Get(name Token) (any, error) {
	if value, ok := li.fields[name.Lexeme]; ok {
		return value, nil
	}

	if method, ok := li.class.Methods[name.Lexeme]; ok {
		return method.Bind(li), nil
	}

	err := NewRuntimeError(name, fmt.Sprintf("Class '%s' has not defined property '%s'.", li.class.Name, name.Lexeme))
	return nil, err
}

func (li *LoxInstance) Set(name Token, value any) {
	li.fields[name.Lexeme] = value
}
