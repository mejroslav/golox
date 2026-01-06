package golox

import "fmt"

func Error(line int, context, message string) bool {
	report(line, context, message)
	return true
}

func report(line int, context string, message string) {
	fmt.Printf("[line %d] Error: %s\n%s\n\n", line, message, context)
}
