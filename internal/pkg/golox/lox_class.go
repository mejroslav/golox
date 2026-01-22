package golox

// LoxClass represents a class in the Lox language.
type LoxClass struct {
	Name    string
	Methods map[string]*LoxFunction
	Env     *Environment
}

func NewLoxClass(classStmt *Class, env *Environment) *LoxClass {
	methods := make(map[string]*LoxFunction)
	for _, method := range classStmt.Methods {
		function := NewLoxFunction(&method, env)
		methods[method.Name.Lexeme] = function
	}
	return &LoxClass{
		Name:    classStmt.Name.Lexeme,
		Methods: methods,
		Env:     env,
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
