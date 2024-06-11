package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/kljablon/golox/ast"
	"github.com/kljablon/golox/interpret"
	"github.com/kljablon/golox/parse"
	"github.com/kljablon/golox/resolve"
	"github.com/kljablon/golox/utils"
)

var hadError bool = false
var hadRuntimeError bool = false

func run(source string) {
	scanner := parse.NewScanner(source)
	tokens := scanner.ScanTokens()
	parser := parse.NewParser(tokens)
	statements := parser.Parse()
	interpreter := interpret.NewInterpreter()

	if hadError {
		os.Exit(65)
	}
	if hadRuntimeError {
		os.Exit(70)
	}

	resolver := resolve.NewResover()
	resolver.ResolveStmts(statements)

	// printer := AstPrinter{}
	// fmt.Println(printer.Print(expression))
	interpreter.Interpret(statements)
}

func ReportError(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	log.Fatalf("[line %d] Error%s: %s", line, where, message)
	hadError = true
}

func loxError(token ast.Token, message string) {
	if token.TokenType == ast.EOF {
		report(token.Line, " at end", message)
	} else {
		report(token.Line, " at '"+token.Lexeme+"'", message)
	}
}

func runtimeError(err utils.RuntimeError) {
	fmt.Printf(err.Message + "\n[line" + fmt.Sprintf("%d", err.Token.Line) + "]")
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
