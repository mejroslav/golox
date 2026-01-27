package golox

// LoxClass represents a class in the Lox language.
type LoxClass struct {
	Name       string                  // The name of the class
	Superclass *LoxClass               // The superclass of the class, if any
	Methods    map[string]*LoxFunction // The methods defined in the class
}

func NewLoxClass(name string, superclass *LoxClass, methods map[string]*LoxFunction) *LoxClass {
	return &LoxClass{
		Name:       name,
		Superclass: superclass,
		Methods:    methods,
	}
}

// String returns a string representation of the Lox class.
func (lc *LoxClass) String() string {
	return "<class " + lc.Name + ">"
}

// Arity returns the number of parameters required to call the class (i.e., its initializer's arity).
func (lc *LoxClass) Arity() int {
	if initializer, ok := lc.getInitializer(); ok {
		return initializer.Arity()
	}
	return 0
}

// Call creates a new instance of the class and initializes it if there is an initializer.
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

// GetMethod looks up a method by name, checking superclasses if necessary.
func (lc *LoxClass) GetMethod(name string) (*LoxFunction, bool) {
	method, ok := lc.Methods[name]
	if !ok && lc.Superclass != nil {
		return lc.Superclass.GetMethod(name)
	}
	return method, ok
}
