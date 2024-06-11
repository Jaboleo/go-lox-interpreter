package tool

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func defineType(f *os.File, baseName string, className string, fieldList string) {
	typeName := baseName + "_" + className
	f.WriteString("type " + typeName + " struct{\n")
	fields := strings.Split(fieldList, ", ")
	for _, v := range fields {
		// name := strings.Split(v, " ")[1]
		f.WriteString(v + "\n")
	}
	f.WriteString("}\n")
	f.WriteString("\n")
	f.WriteString("func (e *" + typeName + ") accept(Visitor Visitor) any{\nreturn Visitor.Visit" + typeName + "(e)\n}\n")
}

func defineVisitor(f *os.File, baseName string, types []string) {
	f.WriteString("type Visitor interface{\n")

	for _, v := range types {
		typeName := strings.Trim(strings.Split(v, ":")[0], " ")
		typeName = baseName + "_" + typeName
		f.WriteString("Visit" + typeName + "(e *" + typeName + ")" + " any\n")
	}
	f.WriteString("}\n")
	f.WriteString("\n")
}

func defineAst(outputDir string, baseName string, types []string) {
	path := filepath.Join("..", outputDir, baseName+".go")
	fmt.Println(path)
	f, err := os.Create(path)

	if err != nil {
		panic(err)
	}

	defer f.Close()

	f.WriteString("package main\n")
	f.WriteString("\n")
	f.WriteString("type Expr interface{\naccept(Visitor Visitor) any\n}\n\n")

	defineVisitor(f, baseName, types)

	for _, v := range types {
		className := strings.Trim(strings.Split(v, ":")[0], " ")
		fields := strings.Trim(strings.Split(v, ":")[1], " ")
		defineType(f, baseName, className, fields)
	}
}

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Println("Usage: generate_ast [outputDir]")
		os.Exit(64)
	}
	outputDir := args[1]

	defineAst(outputDir, "Expr", []string{"Binary   : left Expr, operator Token, right Expr",
		"Grouping : expression Expr",
		"Literal  : value any",
		"Unary    : operator Token, right Expr"})

	defineAst(outputDir, "Stmt", []string{
		"Expression : Expr expression",
		"Print      : Expr expression"})
}

// TODO poprawić obsługę błędów
