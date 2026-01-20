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
	statements, parseErr := parser.Parse()
	if parseErr != nil {
		return fmt.Errorf("parsing errors occurred: %w", parseErr)
	}

	if showAST {
		astPrinter := NewASTPrinter()
		fmt.Println("AST:")
		for _, stmt := range statements {
			astStr := astPrinter.Print(stmt)
			fmt.Println(astStr)
		}
		fmt.Println()
	}

	interpreter := NewInterpreter()
	// Resolve the statements
	resolver := NewResolver(interpreter)
	statements = resolver.Resolve(statements)

	// Interpret the statements
	_, runtimeErr := interpreter.Interpret(statements)
	if runtimeErr != nil {
		return fmt.Errorf("runtime error: %w", runtimeErr)
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
