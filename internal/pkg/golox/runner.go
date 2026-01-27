package golox

import (
	"bufio"
	"fmt"
	"log/slog"
	"mejroslav/golox/v2/internal/pkg/golox/interpreter"
	"mejroslav/golox/v2/internal/pkg/golox/resolver"
	lox_scanner "mejroslav/golox/v2/internal/pkg/golox/scanner"
	"os"
)

// runFile reads a file line by line and prints each line to stdout.
func RunFile(path string, showTokens bool, showAST bool) error {

	slog.Debug("Running file", "path", path)
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
	codeScanner := lox_scanner.NewCodeScanner(1, path)
	tokens, scanErr := codeScanner.Run(source)
	if scanErr {
		return fmt.Errorf("scanning errors")
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
	if parseErr {
		return fmt.Errorf("parsing errors")
	}

	if showAST {
		astPrinter := NewASTPrinter()
		astPrinterResult := astPrinter.Print(statements)
		fmt.Println(astPrinterResult)
		fmt.Println()
	}

	// Resolve the statements
	interpreter := interpreter.NewInterpreter()
	resolver := resolver.NewResolver(interpreter)
	statements, err = resolver.Resolve(statements)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	// Interpret the statements
	_, runtimeErr := interpreter.Interpret(statements)
	if runtimeErr != nil {
		return fmt.Errorf("%w", runtimeErr)
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
		codeScanner := lox_scanner.NewCodeScanner(lineNumber, "<stdin>")
		codeScanner.Run(line)
		lineNumber++
	}
}
