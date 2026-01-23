# GoLox = A Lox Interpreter in Go

GoLox is an interpreter for the Lox programming language, implemented in Go. It is based on the book "Crafting Interpreters" by Robert Nystrom.

## Differences from the original language

- Added native function `input()` to read user input from the console.
- Added keyword `function` as an alias for `fun` when declaring functions.

## Differences from the original implementation

- The scanner keeps track of line numbers and columns for better error reporting.
- The interpreter uses Go's error handling instead of exceptions.
- The interpreter is structured to leverage Go's type system and interfaces.
