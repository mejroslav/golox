package golox

import "fmt"

type ReturnValue struct {
	Keyword *Token
	Value   any
}

func (r *ReturnValue) Error() string {
	return fmt.Sprintf("Return value: %v", r.Value)
}
