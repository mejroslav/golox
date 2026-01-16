package golox

import (
	"mejroslav/golox/v2/internal/pkg/utils"
)

type Resolver struct {
	interpreter *Interpreter
	scopeStack  *utils.Stack
}

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{
		interpreter: interpreter,
		scopeStack:  utils.NewStack(),
	}
}

func (r *Resolver) VisitBlockStmt(s *Block) (any, error) {
	r.BeginScope()
	for _, stmt := range s.Statements {
		if err := r.resolveStmt(stmt); err != nil {
			return nil, err
		}
	}
	r.EndScope()
	return nil, nil
}

func (r *Resolver) resolveStmt(statement Stmt) error {
	_, err := statement.Accept(r)
	return err
}

func (r *Resolver) resolveExpr(expression Expr) error {
	_, err := expression.Accept(r)
	return err
}

func (r *Resolver) BeginScope() {
	r.scopeStack.Push(make(map[string]bool))
}

func (r *Resolver) EndScope() {
	r.scopeStack.Pop()
}

func (r *Resolver) VisitVarStmt(s *Var) (any, error) {
	if !r.scopeStack.IsEmpty() {
		scope := r.scopeStack.Peek().(map[string]bool)
		scope[s.Name.Lexeme] = false
	}
	if s.Initializer != nil {
		if err := r.resolveExpr(s.Initializer); err != nil {
			return nil, err
		}
	}
	if !r.scopeStack.IsEmpty() {
		scope := r.scopeStack.Peek().(map[string]bool)
		scope[s.Name.Lexeme] = true
	}
	return nil, nil
}

func (r *Resolver) declare(name *Token) {
	if r.scopeStack.IsEmpty() {
		return
	}
	scope := r.scopeStack.Peek().(map[string]bool)
	scope[name.Lexeme] = false
}

func (r *Resolver) define(name *Token) {
	if r.scopeStack.IsEmpty() {
		return
	}
	scope := r.scopeStack.Peek().(map[string]bool)
	scope[name.Lexeme] = true
}
