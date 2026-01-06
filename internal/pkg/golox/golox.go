package golox

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

// Main is the entry point for the golox command-line application.
func Main() {

	// Define flags
	help := flag.Bool("help", false, "Show help message and exit")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	// Add more flags here as needed

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "golox - Lox language interpreter in Go\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: golox [options] <file>\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	level := slog.LevelInfo
	if *verbose {
		fmt.Println("Verbose mode enabled")
		level = slog.LevelDebug
	} else {
		fmt.Println("Verbose mode disabled")
	}

	slog.SetDefault(slog.New(tint.NewHandler(os.Stderr, &tint.Options{
		Level:      level,
		TimeFormat: "2006-01-02 15:04:05.000",
	})))

	slog.Info("Starting golox interpreter")
	args := flag.Args()
	if len(args) < 1 {
		RunPrompt()
		os.Exit(0)
	}
	filePath := args[0]
	if err := RunFile(filePath, *verbose); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	slog.Info("Execution completed successfully")
}
