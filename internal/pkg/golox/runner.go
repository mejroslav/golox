package golox

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
)

// runFile reads a file line by line and prints each line to stdout.
func RunFile(path string, showTokens bool, showAST bool) error {

	slog.Info("Running file", "path", path)
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	var source string

	// Load the entire file in memory
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		source += scanner.Text() + "\n"
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// Run the code scanner on the loaded source
	codeScanner := NewCodeScanner(1, path)
	tokens, scanErr := codeScanner.Run(source)
	if scanErr {
		return fmt.Errorf("scanning errors occurred")
	}

	if showTokens {
		fmt.Println("Tokens:")
		for _, token := range tokens {
			fmt.Println(token)
		}
		fmt.Println()
	}

	// Run the parser on the tokens
	parser := NewParser(tokens)
	expression, parseErr := parser.Parse()
	if parseErr != nil {
		return fmt.Errorf("parsing errors occurred: %w", parseErr)
	}

	if showAST {
		astPrinter := NewASTPrinter()
		astStr := astPrinter.Print(expression)
		fmt.Println("AST:")
		fmt.Println(astStr)
		fmt.Println()
	}
	return nil
}

// TODO: Fix the error handling and reporting in the REPL
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
		codeScanner := NewCodeScanner(lineNumber, "<stdin>")
		codeScanner.Run(line)
		lineNumber++
	}
}
