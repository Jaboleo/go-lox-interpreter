package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

var hadError bool = false

func run(source string) {
	scanner := NewScanner(source)
	tokens := scanner.ScanTokens()

	for _, v := range tokens {
		fmt.Println(v.ToString())
	}
}

func ReportError(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	log.Fatalf("[line %d] Error%s: %s", line, where, message)
	hadError = true
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
