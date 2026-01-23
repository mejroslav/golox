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
	if initializer, ok := lc.getInitializer(); ok {
		return initializer.Arity()
	}
	return 0
}

func (lc *LoxClass) Call(interpreter *Interpreter, arguments []any) (any, error) {
	instance := NewLoxInstance(lc)
	if initializer, ok := lc.getInitializer(); ok {
		_, err := initializer.Bind(instance).Call(interpreter, arguments)
		if err != nil {
			return nil, err
		}
	}
	return instance, nil
}

func (lc *LoxClass) getInitializer() (*LoxFunction, bool) {
	initializer, ok := lc.Methods["init"]
	if !ok {
		return nil, false
	}
	return initializer, true
}
