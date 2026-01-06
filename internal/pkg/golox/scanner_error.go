package golox

import "fmt"

func Error(file string, line int, context, message string) bool {
	report(file, line, context, message)
	return true
}

func report(file string, line int, context string, message string) {
	fmt.Printf("LEXER ERROR [%s:%d] %s\n%s\n\n", file, line, message, context)
}
