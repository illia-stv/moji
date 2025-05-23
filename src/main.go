package main

import (
	"fmt"
	"os"

	"moji/src/evaluator"
	"moji/src/parser"
	"moji/src/scanner"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh <command> <filename>")
		os.Exit(1)
	}

	command := os.Args[1]
	filename := os.Args[2]

	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	switch command {
	case "tokenize":
		scanner.Scan(fileContents)
	case "parse":
		s := scanner.NewScanner(string(fileContents))
		tokens := s.ScanTokens()
		if s.HasError() {
			os.Exit(65)
		}
		p := parser.NewParser(tokens)
		statements := p.ParseStatements()
		for _, stmt := range statements {
			fmt.Println(stmt)
		}
	case "evaluate":
		s := scanner.NewScanner(string(fileContents))
		tokens := s.ScanTokens()
		if s.HasError() {
			os.Exit(65)
		}
		p := parser.NewParser(tokens)
		e := evaluator.NewEvaluator(p)
		
		// The evaluator will handle runtime errors and exit with code 70 if needed
		result := e.Evaluate()
		fmt.Println(result)
	case "run":
		s := scanner.NewScanner(string(fileContents))
		tokens := s.ScanTokens()
		if s.HasError() {
			os.Exit(65)
		}
		p := parser.NewParser(tokens)
		e := evaluator.NewEvaluator(p)
		
		// Evaluate statements, including print statements
		e.EvaluateStatements()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}
