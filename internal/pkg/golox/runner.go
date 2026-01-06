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
