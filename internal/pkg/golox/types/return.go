package types

import (
	"fmt"
	"mejroslav/golox/v2/internal/pkg/golox/token"
)

// ReturnValue is a special error type used to handle return statements in functions.
// It carries the return value and allows it to be propagated up the call stack.
type ReturnValue struct {
	Keyword *token.Token
	Value   any
}

func (r *ReturnValue) Error() string {
	return fmt.Sprintf("Return value: %v", r.Value)
}
