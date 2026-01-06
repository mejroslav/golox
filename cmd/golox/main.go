package main

import (
	"flag"
	"fmt"
	"mejroslav/golox/v2/internal/pkg/golox"
	"os"
)

func main() {
	// Define flags
	help := flag.Bool("help", false, "Show help message and exit")
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

	args := flag.Args()
	if len(args) < 1 {
		golox.RunPrompt()
		os.Exit(0)
	}
	filePath := args[0]
	if err := golox.RunFile(filePath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
