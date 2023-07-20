package parser

import (
	"errors"
	"fmt"
	"log"
	"morklerork/ast"
	"morklerork/symbols"
)

func splitByLowestPrecedence(expressionSymbols []symbols.Symbol) ([]symbols.Symbol, symbols.BinaryOperator, []symbols.Symbol, error) {
	if len(expressionSymbols) < 3 {
		return nil, symbols.BinaryOperator{}, nil, errors.New("not enough symbols in expression slice to split: " + fmt.Sprint(expressionSymbols))
	}

	index := -1
	for i, symbol := range expressionSymbols {
		switch symbol.(type) {
		case symbols.BinaryOperator:
			if index == -1 {
				index = i
			} else {
				// we can directly compare the enum values since they are ints under the hood, and are in the right order
				// also, if the precedence is equal, take the _later_ operator so that _earlier_ operators execute first
				if symbol.(symbols.BinaryOperator).BinaryOperatorType <= expressionSymbols[index].(symbols.BinaryOperator).BinaryOperatorType {
					index = i
				}
			}
		}
	}

	return expressionSymbols[:index], expressionSymbols[index].(symbols.BinaryOperator), expressionSymbols[index+1:], nil
}

func parseSingleSymbolExpression(expressionSymbol symbols.Symbol) (ast.Expression, error) {
	switch expressionSymbol := expressionSymbol.(type) {
	case symbols.StringLiteral:
		return ast.StringLiteral{Value: expressionSymbol.Value}, nil
	case symbols.IntLiteral:
		return ast.IntLiteral{Value: expressionSymbol.Value}, nil
	case symbols.BooleanLiteral:
		return ast.BooleanLiteral{Value: expressionSymbol.Value}, nil
	case symbols.VariableName:
		return ast.VariableName{Name: expressionSymbol.Name}, nil
	case symbols.HeapAccess:
		expr, err := parseSingleSymbolExpression(expressionSymbol.IndexExpressionSymbol)
		if err != nil {
			return expr, err
		}
		return ast.HeapAccess{IndexExpression: expr}, nil
	}
	fmt.Println(expressionSymbol)
	return nil, errors.New("the value symbol was not a String or an Int")
}

func parseExpression(expressionSymbols []symbols.Symbol) (ast.Expression, error) {

	if len(expressionSymbols) == 0 {
		return nil, errors.New("tried to parse and empty expression")
	}

	if len(expressionSymbols) == 1 {
		return parseSingleSymbolExpression(expressionSymbols[0])
	}

	// counter intuitively, the first thing we split on will be executed last
	// So split on the _lowest_ precedence
	lhs, operator, rhs, err := splitByLowestPrecedence(expressionSymbols)

	if err != nil {
		return nil, err
	}

	parsedLhs, err := parseExpression(lhs)

	if err != nil {
		return nil, err
	}

	parsedRhs, err := parseExpression(rhs)

	if err != nil {
		return nil, err
	}

	return ast.BinaryOperator{
		Lhs:                parsedLhs,
		Rhs:                parsedRhs,
		BinaryOperatorType: operator.BinaryOperatorType,
	}, nil
}

func parseLog(logSymbols []symbols.Symbol, indent int) ast.Log {
	expr, err := parseExpression(logSymbols)
	if err != nil {
		log.Fatal(err)
	}
	return ast.Log{Expr: expr, Indent: indent}
}

func parseRead(readSymbols []symbols.Symbol, indent int) ast.Read {
	var target ast.Expression

	switch targetSymbol := readSymbols[0].(type) {
	case symbols.VariableName:
		target = ast.VariableName{Name: targetSymbol.Name}
		break
	case symbols.HeapAccess:
		// Since HeapAccess can have more HeapAccesses inside it needs to be fully parsed
		parsedHeapAccess, err := parseSingleSymbolExpression(targetSymbol)
		if err != nil {
			log.Fatal(err)
		}
		target = parsedHeapAccess
		break
	default:
		log.Fatal("The first symbol in a read must be a variable or heap access")
	}

	if len(readSymbols) != 1 {
		log.Fatal("Read should only be given 2 symbol")
	}

	return ast.Read{Target: target, Indent: indent}
}

func parseAssign(assignSymbols []symbols.Symbol, indent int) ast.Assign {
	var target ast.Expression

	switch targetSymbol := assignSymbols[0].(type) {
	case symbols.VariableName:
		target = ast.VariableName{Name: targetSymbol.Name}
		break
	case symbols.HeapAccess:
		// Since HeapAccess can have more HeapAccesses inside it needs to be fully parsed
		parsedHeapAccess, err := parseSingleSymbolExpression(targetSymbol)
		if err != nil {
			log.Fatal(err)
		}
		target = parsedHeapAccess
		break
	default:
		log.Fatal("The first symbol in an assignment must be a variable or heap access")
	}

	expr, err := parseExpression(assignSymbols[1:])
	if err != nil {
		log.Fatal(err)
	}

	return ast.Assign{Target: target, Expr: expr, Indent: indent}
}

func parseNew(newSymbols []symbols.Symbol, indent int) ast.New {
	variableName := newSymbols[0]
	expr, err := parseExpression(newSymbols[1:])
	if err != nil {
		log.Fatal(err)
	}
	return ast.New{VariableName: variableName.(symbols.VariableName).Name, Expr: expr, Indent: indent}
}

func parseIf(IfSymbols []symbols.Symbol, indent int) ast.If {
	expr, err := parseExpression(IfSymbols)
	if err != nil {
		log.Fatal(err)
	}
	return ast.If{Cond: expr, Indent: indent}
}

func parseWhile(WhileSymbols []symbols.Symbol, indent int) ast.While {
	expr, err := parseExpression(WhileSymbols)
	if err != nil {
		log.Fatal(err)
	}
	return ast.While{Cond: expr, Indent: indent}
}

func parseProgram(ProgramSymbols []symbols.Symbol, indent int) ast.Program {
	name := ast.ProgramName{Name: ProgramSymbols[0].(symbols.ProgramName).Name}
	parameterSymbols := ProgramSymbols[1:]
	variables := make([]ast.VariableName, 0)
	for i, _ := range parameterSymbols {
		expr, err := parseSingleSymbolExpression(parameterSymbols[i])
		if err != nil {
			log.Fatal(err)
		}
		variables = append(variables, expr.(ast.VariableName))
	}
	return ast.Program{Name: name, Parameters: variables, Indent: indent}
}

func parseCall(CallSymbols []symbols.Symbol, indent int) ast.Call {
	callSymbols := CallSymbols[:]
	returnTargetName, hasReturnTarget := callSymbols[0].(symbols.VariableName)
	if hasReturnTarget { // if a return target was specified, remove that symbol for the rest of the parsing
		callSymbols = CallSymbols[1:]
	}
	name := ast.ProgramName{Name: callSymbols[0].(symbols.ProgramName).Name}
	expressionsSymbols := callSymbols[1:]
	expressions := make([]ast.Expression, 0)
	if len(expressionsSymbols) > 0 {
		for _, symbol := range expressionsSymbols {
			expression, err := parseSingleSymbolExpression(symbol)
			if err != nil {
				log.Fatal(err)
			}
			expressions = append(expressions, expression)
		}
	}

	return ast.Call{Name: name, Expressions: expressions, ReturnTarget: ast.VariableName{Name: returnTargetName.Name}, HasReturnTarget: hasReturnTarget, Indent: indent}
}

func parseReturn(ReturnSymbols []symbols.Symbol, indent int) ast.Return {
	hasExpression := len(ReturnSymbols) != 0
	expr := ast.Expression(nil)
	if hasExpression {
		_expr, err := parseExpression(ReturnSymbols)
		if err != nil {
			log.Fatal(err)
		}
		expr = _expr
	}
	return ast.Return{Expression: expr, HasExpression: hasExpression, Indent: indent}
}

func parseCommand(commandSymbols []symbols.Symbol) (ast.Command, bool, error) {
	indent := commandSymbols[0].(symbols.Indent).Level
	switch commandSymbols[1].(type) {
	case symbols.Print:
		return parseLog(commandSymbols[2:], indent), false, nil
	case symbols.Read:
		return parseRead(commandSymbols[2:], indent), false, nil
	case symbols.Assign:
		return parseAssign(commandSymbols[2:], indent), false, nil
	case symbols.Define:
		return parseNew(commandSymbols[2:], indent), false, nil
	case symbols.If:
		return parseIf(commandSymbols[2:], indent), true, nil
	case symbols.While:
		return parseWhile(commandSymbols[2:], indent), true, nil
	case symbols.Program:
		return parseProgram(commandSymbols[2:], indent), true, nil
	case symbols.Call:
		return parseCall(commandSymbols[2:], indent), false, nil
	case symbols.Return:
		return parseReturn(commandSymbols[2:], indent), false, nil
	}
	return nil, false, errors.New("the first symbol in the command is not recognized")
}

func setCommandsInBlockCommand(command ast.Command, commands []ast.Command) (ast.Command, error) {
	switch blockCommand := command.(type) {
	case ast.If:
		blockCommand.Commands = commands
		return blockCommand, nil
	case ast.While:
		blockCommand.Commands = commands
		return blockCommand, nil
	case ast.Program:
		blockCommand.Commands = commands
		return blockCommand, nil
	}

	return nil, errors.New("tried to set commands on a non block command")
}

// ParseBlock
// Parse a single Block entirely.
// Sub blocks will be recursively parsed then assigned into the
// command they are related to within the outer block
func ParseBlock(program [][]symbols.Symbol, expectedIndentation int) ([]ast.Command, int) {
	commands := make([]ast.Command, 0)

	if program[0][0].(symbols.Indent).Level != expectedIndentation {
		log.Fatal("First command in block is not indented properly, expected: " + fmt.Sprint(expectedIndentation) + " got: " + fmt.Sprint(program[0][0].(symbols.Indent).Level))
	}

	// Track how many commands were parsed, so the parent block can skip over them
	parsed := 0

	// Allows for skipping commands when a recursive call has parsed some number of commands
	skip := 0
	for index, commandSymbols := range program {
		if skip > 0 { // skip some commands already parsed
			parsed++
			skip--
			continue
		}
		thisIndent := commandSymbols[0].(symbols.Indent).Level

		if thisIndent > expectedIndentation {
			log.Fatal("The indentation unexpectedly increased")
		} else if thisIndent < expectedIndentation {
			// We are done parsing this block early because it un-indented
			return commands, parsed
		}

		command, commandExpectsBlock, err := parseCommand(commandSymbols)
		if err != nil {
			log.Fatal(err)
		}
		if commandExpectsBlock {
			nextIndent := program[index+1][0].(symbols.Indent).Level
			if nextIndent > thisIndent {
				blockCommands, parsed := ParseBlock(program[index+1:], nextIndent)
				skip = parsed // skip the number of commands parsed by the recursive call
				command, err = setCommandsInBlockCommand(command, blockCommands)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				log.Fatal("The next command after a block command was not indented")
			}
		}
		commands = append(commands, command)
		parsed++
	}

	return commands, parsed
}
