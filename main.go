package main

import (
	"morklerork/executor"
	"morklerork/lexer"
	"morklerork/loader"
	"morklerork/parser"
)

func main() {
	programString := loader.Load()
	programSymbols := lexer.Lex(programString)
	programAst, _ := parser.ParseBlock(programSymbols, 0)
	executor.ExecuteProgram(programAst)
}
