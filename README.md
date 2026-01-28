# GoLox -  A Lox Interpreter in Go

GoLox is an interpreter for the Lox programming language, implemented in Go. It is based on the book ["Crafting Interpreters" by Robert Nystrom](https://craftinginterpreters.com/).

## Usage

To run a Lox script, use the following command:

```shell
./golox path/to/script.lox
```

## Differences from the original language

- Added native function `input()` to read user input from the console.
- Added keyword `function` as an alias for `fun` when declaring functions.
- Added `break` statement to exit loops early.

## Differences from the original implementation

- The scanner keeps track of line numbers and columns of each token.

    Scanner error messages include line and column numbers, as shown below:

    ```
    SCANNER ERROR [examples/04-syntax-err.lox:3:12] Unexpected character '?'.

    var xy = 21?;
    -----------^-
    ```

    Parser error messages also include line and column numbers:

    ```
    PARSER ERROR [examples/04-syntax-err.lox:3:9] at '==': Expect ';' after variable declaration.
    ```

- The interpreter uses Go's error handling instead of exceptions.
- The interpreter is structured to leverage Go's type system and interfaces.

## In progress

- REPL mode is broken.
- Scanner errors are nicely printed, but parser and runtime errors are not yet formatted similarly.
- Implementation of `continue` keyword is pending, as it cannot be simply added to the for loop without restructuring.
