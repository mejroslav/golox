package golox

// LoxClass represents a class in the Lox language.
type LoxClass struct {
	Name    string
	Methods map[string]*LoxFunction
}

func NewLoxClass(name string, methods map[string]*LoxFunction) *LoxClass {
	return &LoxClass{
		Name:    name,
		Methods: methods,
	}
}

func (lc *LoxClass) String() string {
	return "<class " + lc.Name + ">"
}

func (lc *LoxClass) Arity() int {
	// TODO: implement when we have initializers
	// if initializer, ok := lc.Methods["init"]; ok {
	// 	return initializer.Arity()
	// }
	return 0
}

func (lc *LoxClass) Call(interpreter *Interpreter, arguments []any) (any, error) {
	instance := NewLoxInstance(lc)
	return instance, nil
}
