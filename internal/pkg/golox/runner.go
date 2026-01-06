package golox

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
)

// runFile reads a file line by line and prints each line to stdout.
func RunFile(path string, verbose bool) error {

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
	scannerHadErrors := codeScanner.Run(source, verbose)
	if scannerHadErrors {
		return fmt.Errorf("scanning errors occurred")
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
		codeScanner := NewCodeScanner(lineNumber, "<stdin>")
		codeScanner.Run(line, true)
		lineNumber++
	}
}
