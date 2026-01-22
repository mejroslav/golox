package golox

import (
	"fmt"
	"time"
)

// Clock is a native function that returns the current time in seconds since the Unix epoch.
type Clock struct{}

func (c *Clock) Arity() int {
	return 0
}

func (c *Clock) Call(interpreter *Interpreter, arguments []any) (any, error) {
	return float64(time.Now().UnixNano()) / 1e9, nil
}
func (c *Clock) String() string {
	return "<native fn clock>"
}

// Input is a native function that prompts the user for input and returns it as a string.
type Input struct{}

func (i *Input) Arity() int {
	return 1
}

func (i *Input) Call(interpreter *Interpreter, arguments []any) (any, error) {
	var prompt string
	if len(arguments) > 0 {
		if p, ok := arguments[0].(string); ok {
			prompt = p
		}
	}
	var input string
	if prompt != "" {
		print(prompt)
	}
	_, err := fmt.Scanln(&input)
	if err != nil {
		return "", err
	}
	return input, nil
}

func (i *Input) String() string {
	return "<native fn input>"
}
