package lexer

import (
	"log"
	"morklerork/symbols"
	"strconv"
	"strings"
)

// If not for strings with spaces in them, we could just split on ' '
// Instead go rune by rune, keeping track of if we are
// in a string or not
func splitIntoSymbols(line []rune) []string {
	// rune queue to process one by one
	runeQueue := line[:]

	// array of string builders to become each Symbol
	santizedSymbols := []strings.Builder{
		{},
	}

	lastSeenRune := '\000'
	isInString := false
	for len(runeQueue) > 0 {
		if runeQueue[0] == '\'' { // we see a quote, starting or ending a string literal
			if !isInString { // we are not in a string, so start one
				isInString = true
				santizedSymbols[len(santizedSymbols)-1].WriteRune(runeQueue[0])
				lastSeenRune = runeQueue[0]
			} else { // we are in a string, consider ending it
				if lastSeenRune == '\\' { // escaped quote, dont end the string
					santizedSymbols[len(santizedSymbols)-1].WriteRune(runeQueue[0])
					lastSeenRune = runeQueue[0]
				} else { // end the string
					isInString = false
					santizedSymbols[len(santizedSymbols)-1].WriteRune(runeQueue[0])
					lastSeenRune = runeQueue[0]
				}
			}
		} else if runeQueue[0] == ' ' { // we have seen a space
			if isInString { // in a string, preserve the space for the string
				santizedSymbols[len(santizedSymbols)-1].WriteRune(runeQueue[0])
				lastSeenRune = runeQueue[0]
			} else { // we are not in a string, start a new symbol
				santizedSymbols = append(santizedSymbols, strings.Builder{})
			}
		} else { // any other caracter, write the rune into the symbol
			santizedSymbols[len(santizedSymbols)-1].WriteRune(runeQueue[0])
			lastSeenRune = runeQueue[0]
		}

		runeQueue = runeQueue[1:]
	}

	// convert the builders into strings
	programSymbols := make([]string, 0)
	for _, symbol := range santizedSymbols {
		programSymbols = append(programSymbols, symbol.String())
	}
	return programSymbols
}

func lexIndent(line []rune) (int, []rune) {
	indent := 0
	for len(line) > 0 && line[0] == ' ' {
		line = line[1:]
		indent++
	}
	return indent, line
}

func lexLiteralsAndUserDefinedSymbols(symbol string) symbols.Symbol {
	if symbol[0] == ':' {
		return symbols.VariableName{Name: symbol}
	} else if symbol[0] == '$' {
		return symbols.ProgramName{Name: symbol}
	} else if symbol[0] == '?' {
		if symbol[1:] == "true" {
			return symbols.BooleanLiteral{Value: true}
		} else if symbol[1:] == "false" {
			return symbols.BooleanLiteral{Value: false}
		} else {
			log.Fatal("Bool literal should either be `?true` or ?false`")
		}
	} else if symbol[0] == '\'' && symbol[len(symbol)-1] == '\'' {
		unescapedString := strings.Replace(symbol[1:len(symbol)-1], "\\'", "'", -1)
		stringVal, err := strconv.Unquote(`"` + unescapedString + `"`)
		if err != nil {
			log.Println(err.Error())
			log.Fatal(err)
		}
		return symbols.StringLiteral{Value: stringVal}
	} else if symbol[0] == '[' && symbol[len(symbol)-1] == ']' {
		innerSymbol := lexSymbol(symbol[1 : len(symbol)-1])
		return symbols.HeapAccess{IndexExpressionSymbol: innerSymbol}
	} else if num, err := strconv.Atoi(symbol); err == nil {
		return symbols.IntLiteral{Value: num}
	} else {
		log.Fatal("Unrecognised symbol: " + symbol + ", did you mean to use a variable? try `:" + symbol + "`, or a string? try `'" + symbol + "'`")
	}

	// this is unreachable, why is it needed?
	return nil
}

func lexSymbol(symbol string) symbols.Symbol {
	switch symbol {
	// CommandSymbols
	case "log":
		return symbols.Print{}
	case "=":
		return symbols.Assign{}
	case "new":
		return symbols.Define{}
	case "if":
		return symbols.If{}
	case "while":
		return symbols.While{}
	case "program":
		return symbols.Program{}
	case "call":
		return symbols.Call{}
	case "return":
		return symbols.Return{}
	// OperatorSymbols
	case "&":
		return symbols.BinaryOperator{BinaryOperatorType: symbols.LogicalAndOperator}
	case "|":
		return symbols.BinaryOperator{BinaryOperatorType: symbols.LogicalOrOperator}
	case "==":
		return symbols.BinaryOperator{BinaryOperatorType: symbols.EqualOperator}
	case "!=":
		return symbols.BinaryOperator{BinaryOperatorType: symbols.NotEqualOperator}
	case "<":
		return symbols.BinaryOperator{BinaryOperatorType: symbols.LTOperator}
	case "+":
		return symbols.BinaryOperator{BinaryOperatorType: symbols.PlusOperator}
	case "-":
		return symbols.BinaryOperator{BinaryOperatorType: symbols.MinusOperator}
	case "*":
		return symbols.BinaryOperator{BinaryOperatorType: symbols.TimesOperator}
	case "/":
		return symbols.BinaryOperator{BinaryOperatorType: symbols.DivideOperator}
	case "%":
		return symbols.BinaryOperator{BinaryOperatorType: symbols.ModuloOperator}
	// literals and user defined symbols
	default:
		return lexLiteralsAndUserDefinedSymbols(symbol)
	}
	// This is unreachable, why is it needed?
	return nil
}

func Lex(programString string) [][]symbols.Symbol {
	program := make([][]symbols.Symbol, 0)
	programLines := strings.Split(programString, "\n")
	for _, line := range programLines {
		if line == "" { // ignore blank lines
			continue
		}
		if line[0] == '#' { // ignore comments
			continue
		}
		indent, unindentedLine := lexIndent([]rune(line))
		if len(unindentedLine) == 0 { // only had whitespace
			continue
		}

		programCommand := make([]symbols.Symbol, 0)

		programCommand = append(programCommand, symbols.Indent{Level: indent})

		for _, symbol := range splitIntoSymbols(unindentedLine) {
			programCommand = append(programCommand, lexSymbol(symbol))
		}
		program = append(program, programCommand)
	}
	return program
}
