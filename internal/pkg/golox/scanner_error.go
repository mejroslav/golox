package golox

import "fmt"

func Error(line int, message string) bool {
	report(line, "", message)
	return true
}

func report(line int, where string, message string) {
	fmt.Printf("[line %d] Error%s: %s\n", line, where, message)
}
