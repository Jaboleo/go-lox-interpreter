package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

var hadError bool = false
var hadRuntimeError bool = false

func run(source string) {
	scanner := NewScanner(source)
	tokens := scanner.ScanTokens()
	parser := NewParser(tokens)
	statements := parser.Parse()
	interpreter := newInterpreter()

	if hadError {
		os.Exit(65)
	}
	if hadRuntimeError {
		os.Exit(70)
	}

	// printer := AstPrinter{}
	// fmt.Println(printer.Print(expression))
	interpreter.interpret(statements)
}

func ReportError(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	log.Fatalf("[line %d] Error%s: %s", line, where, message)
	hadError = true
}

func loxError(token Token, message string) {
	if token.ttype == EOF {
		report(token.line, " at end", message)
	} else {
		report(token.line, " at '"+token.lexeme+"'", message)
	}
}

func runtimeError(err RuntimeError) {
	fmt.Printf(err.message + "\n[line" + fmt.Sprintf("%f", err.token.line) + "]")
	hadRuntimeError = true
}

func runFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	run(string(data))
	if hadError {
		os.Exit(65)
	}
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, _ := reader.ReadString('\n')
		if line == "" {
			break
		}
		run(line)
		hadError = false
	}
}

func main() {
	args := os.Args
	if len(args) > 2 {
		fmt.Println("Usage: jlox [script]")
		os.Exit(64)
	} else if len(args) == 2 {
		fmt.Printf("runFile(%v)\n", args[1])
		runFile(args[1])
	} else {
		runPrompt()
	}
}
