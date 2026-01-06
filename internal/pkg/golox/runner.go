package golox

import (
	"bufio"
	"fmt"
	"os"
)

// runFile reads a file line by line and prints each line to stdout.
func RunFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 1
	for scanner.Scan() {
		fmt.Printf("%d: %s\n", lineNumber, scanner.Text())
		lineNumber++
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}
	return nil
}

// RunPrompt starts a REPL that reads lines from stdin and echoes them back.
func RunPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	lineNumber := 1
	fmt.Println("Lox REPL. Type 'exit' to quit.")
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		if line == "exit" {
			break
		}
		fmt.Printf("%d: %s\n", lineNumber, line)
		lineNumber++
	}
}

func Error(line int, message string) bool {
	report(line, "", message)
	return true
}

func report(line int, where string, message string) {
	fmt.Printf("[line %d] Error%s: %s\n", line, where, message)
}

type CodeScanner struct{}

func NewCodeScanner() *CodeScanner {
	return &CodeScanner{}
}

func (s *CodeScanner) ScanTokens(source string) []string {
	// Dummy implementation for illustration
	return []string{"token1", "token2"}
}

func (s *CodeScanner) Run() bool {
	var codeScanner CodeScanner = *NewCodeScanner()
	codeScanner.ScanTokens("dummy source")
	return false
}
