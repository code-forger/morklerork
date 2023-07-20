package executor

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"morklerork/ast"
	"morklerork/symbols"
	"os"
	"strconv"
)

type ResultType int

const (
	String ResultType = iota
	Int
	Bool
)

type ExpressionResult struct {
	String string
	Int    int
	Bool   bool
	Type   ResultType
}

type scope []map[string]ExpressionResult

type programs map[string]ast.Program

func assignInScope(name string, val ExpressionResult, scope scope) {
	for i := len(scope) - 1; i >= 0; i-- {
		_, ok := scope[i][name]
		if ok {
			scope[i][name] = val
			return
		}
	}
	log.Fatal("Could not find variable " + name + " in scope. Did you declare it first?")
}

func defineInScope(name string, val ExpressionResult, scope scope) {
	for i := len(scope) - 1; i >= 0; i-- {
		_, ok := scope[i][name]
		if ok {
			log.Fatal("VariableName " + name + " is already defined in this scope")
			return
		}
	}
	scope[len(scope)-1][name] = val
}

func getInScope(name string, scope scope) ExpressionResult {
	for i := len(scope) - 1; i >= 0; i-- {
		val, ok := scope[i][name]
		if ok {
			return val
		}
	}
	log.Fatal("Could not find variable " + name + " in scope. Did you declare it first?")
	return ExpressionResult{}
}

func isInScope(name string, scope scope) bool {
	for i := len(scope) - 1; i >= 0; i-- {
		_, ok := scope[i][name]
		if ok {
			return true
		}
	}
	return false
}

func addScope(scope scope) scope {
	return append(scope, make(map[string]ExpressionResult))
}

func popScope(scope scope) scope {
	return scope[:len(scope)-1]
}

var heap [10000]ExpressionResult

func executeBinaryOperatorOnString(lhs ExpressionResult, rhs ExpressionResult, operatorType symbols.BinaryOperatorType) (ExpressionResult, error) {
	switch rhs.Type {
	case String:
		switch operatorType {
		case symbols.PlusOperator:
			return ExpressionResult{String: lhs.String + rhs.String, Type: String}, nil
		case symbols.EqualOperator:
			return ExpressionResult{Bool: lhs.String == rhs.String, Type: Bool}, nil
		case symbols.NotEqualOperator:
			return ExpressionResult{Bool: lhs.String != rhs.String, Type: Bool}, nil
		default:
			return ExpressionResult{}, errors.New("Cannot use " + symbols.BinaryOperatorTypeNames[operatorType] + " on two strings")
		}
	case Int:
		switch operatorType {
		case symbols.PlusOperator:
			return ExpressionResult{String: lhs.String + strconv.Itoa(rhs.Int), Type: String}, nil
		case symbols.ModuloOperator:
			return ExpressionResult{String: string([]rune(lhs.String)[rhs.Int]), Type: String}, nil
		case symbols.LTOperator:
			return ExpressionResult{Bool: len(lhs.String) < rhs.Int, Type: Bool}, nil
		default:
			return ExpressionResult{}, errors.New("Cannot use " + symbols.BinaryOperatorTypeNames[operatorType] + " on string and int")
		}
	case Bool:
		switch operatorType {
		case symbols.PlusOperator:
			return ExpressionResult{String: lhs.String + strconv.FormatBool(rhs.Bool), Type: String}, nil
		default:
			return ExpressionResult{}, errors.New("Cannot use " + symbols.BinaryOperatorTypeNames[operatorType] + " on string and bool")
		}
	}

	return ExpressionResult{}, errors.New("RHS expression is of unrecognised type")
}

func executeBinaryOperatorOnInt(lhs ExpressionResult, rhs ExpressionResult, operatorType symbols.BinaryOperatorType) (ExpressionResult, error) {
	switch rhs.Type {
	case String:
		switch operatorType {
		case symbols.LTOperator:
			return ExpressionResult{Bool: lhs.Int < len(rhs.String), Type: Bool}, nil
		default:
			return ExpressionResult{}, errors.New("Cannot use " + symbols.BinaryOperatorTypeNames[operatorType] + " on int and string")
		}
	case Int:
		switch operatorType {
		case symbols.EqualOperator:
			return ExpressionResult{Bool: lhs.Int == rhs.Int, Type: Bool}, nil
		case symbols.NotEqualOperator:
			return ExpressionResult{Bool: lhs.Int != rhs.Int, Type: Bool}, nil
		case symbols.LTOperator:
			return ExpressionResult{Bool: lhs.Int < rhs.Int, Type: Bool}, nil
		case symbols.PlusOperator:
			return ExpressionResult{Int: lhs.Int + rhs.Int, Type: Int}, nil
		case symbols.MinusOperator:
			return ExpressionResult{Int: lhs.Int - rhs.Int, Type: Int}, nil
		case symbols.TimesOperator:
			return ExpressionResult{Int: lhs.Int * rhs.Int, Type: Int}, nil
		case symbols.DivideOperator:
			return ExpressionResult{Int: lhs.Int / rhs.Int, Type: Int}, nil
		case symbols.ModuloOperator:
			return ExpressionResult{Int: lhs.Int % rhs.Int, Type: Int}, nil
		default:
			return ExpressionResult{}, errors.New("Cannot use '" + symbols.BinaryOperatorTypeNames[operatorType] + "' on two ints")
		}
	case Bool:
		switch operatorType {
		default:
			return ExpressionResult{}, errors.New("Cannot use " + symbols.BinaryOperatorTypeNames[operatorType] + " on int and bool")
		}
	}

	return ExpressionResult{}, errors.New("RHS expression is of unrecognised type")
}

func executeBinaryOperatorOnBool(lhs ExpressionResult, rhs ExpressionResult, operatorType symbols.BinaryOperatorType) (ExpressionResult, error) {
	switch rhs.Type {
	case String:
		switch operatorType {
		default:
			return ExpressionResult{}, errors.New("Cannot use " + symbols.BinaryOperatorTypeNames[operatorType] + " on bool and string")
		}
	case Int:
		switch operatorType {
		default:
			return ExpressionResult{}, errors.New("Cannot use '" + symbols.BinaryOperatorTypeNames[operatorType] + "' on bool and int")
		}
	case Bool:
		switch operatorType {
		case symbols.LogicalAndOperator:
			return ExpressionResult{Bool: lhs.Bool && rhs.Bool, Type: Bool}, nil
		case symbols.LogicalOrOperator:
			return ExpressionResult{Bool: lhs.Bool || rhs.Bool, Type: Bool}, nil
		case symbols.EqualOperator:
			return ExpressionResult{Bool: lhs.Bool == rhs.Bool, Type: Bool}, nil
		case symbols.NotEqualOperator:
			return ExpressionResult{Bool: lhs.Bool != rhs.Bool, Type: Bool}, nil
		default:
			return ExpressionResult{}, errors.New("Cannot use " + symbols.BinaryOperatorTypeNames[operatorType] + " on two bools")
		}
	}

	return ExpressionResult{}, errors.New("RHS expression is of unrecognised type")
}

func evaluateBinaryOperator(expression ast.BinaryOperator, scope scope) (ExpressionResult, error) {
	lhs, err := evaluateExpression(expression.Lhs, scope)

	if err != nil {
		log.Fatal(err)
	}

	rhs, err := evaluateExpression(expression.Rhs, scope)

	if err != nil {
		log.Fatal(err)
	}

	switch lhs.Type {
	case String:
		return executeBinaryOperatorOnString(lhs, rhs, expression.BinaryOperatorType)
	case Int:
		return executeBinaryOperatorOnInt(lhs, rhs, expression.BinaryOperatorType)
	case Bool:
		return executeBinaryOperatorOnBool(lhs, rhs, expression.BinaryOperatorType)
	}

	return ExpressionResult{}, errors.New("LHS expression is of unrecognised type")
}

func evaluateExpression(expression ast.Expression, scope scope) (ExpressionResult, error) {
	switch expression := expression.(type) {
	case ast.StringLiteral:
		return ExpressionResult{
			Type:   String,
			String: expression.Value,
		}, nil
	case ast.IntLiteral:
		return ExpressionResult{
			Type: Int,
			Int:  expression.Value,
		}, nil
	case ast.BooleanLiteral:
		return ExpressionResult{
			Type: Bool,
			Bool: expression.Value,
		}, nil
	case ast.VariableName:
		return getInScope(expression.Name, scope), nil
	case ast.HeapAccess:
		targetValue, err := evaluateExpression(expression.IndexExpression, scope)
		if err != nil {
			log.Fatal(err)
		}

		if targetValue.Type != Int {
			log.Fatal("tried to access the heap with value that is not an int")
		}
		return heap[targetValue.Int], nil
	case ast.BinaryOperator:
		return evaluateBinaryOperator(expression, scope)
	}
	return ExpressionResult{}, errors.New("tried to evaluate an unrecognised AST node")
}

func runPrint(print ast.Log, scope scope) {
	result, err := evaluateExpression(print.Expr, scope)

	if err != nil {
		log.Fatal(err)
	}

	switch result.Type {
	case String:
		fmt.Print(result.String)
		break
	case Int:
		fmt.Print(result.Int)
		break
	}
}

func readRune() (rune, error) {
	state, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return 0, err
	}

	r := bufio.NewReader(os.Stdin)
	ru, _, err := r.ReadRune()
	terminal.Restore(0, state)
	return ru, err
}

func runRead(read ast.Read, scope scope) {
	rawRune, err := readRune()
	if err != nil {
		log.Fatal(err)
	}

	str := ExpressionResult{
		Type:   String,
		String: string(rawRune),
	}

	switch assignTarget := read.Target.(type) {
	case ast.VariableName:
		assignInScope(assignTarget.Name, str, scope)
		break
	case ast.HeapAccess:
		targetValue, err := evaluateExpression(assignTarget.IndexExpression, scope)
		if err != nil {
			log.Fatal(err)
		}

		if targetValue.Type != Int {
			log.Fatal("tried to access the heap with value that is not an int")
		}
		heap[targetValue.Int] = str
	}
}

func runAssign(assign ast.Assign, scope scope) {
	result, err := evaluateExpression(assign.Expr, scope)

	if err != nil {
		log.Fatal(err)
	}

	switch assignTarget := assign.Target.(type) {
	case ast.VariableName:
		assignInScope(assignTarget.Name, result, scope)
		break
	case ast.HeapAccess:
		targetValue, err := evaluateExpression(assignTarget.IndexExpression, scope)
		if err != nil {
			log.Fatal(err)
		}

		if targetValue.Type != Int {
			log.Fatal("tried to access the heap with value that is not an int")
		}
		heap[targetValue.Int] = result
	}
}

func runDefine(define ast.New, scope scope) {
	result, err := evaluateExpression(define.Expr, scope)

	if err != nil {
		log.Fatal(err)
	}

	defineInScope(define.VariableName, result, scope)
}

func runIf(ifCommand ast.If, scope scope, programs programs) (ExpressionResult, bool, bool) {
	result, err := evaluateExpression(ifCommand.Cond, scope)
	if err != nil {
		log.Fatal(err)
	}

	if result.Type != Bool {
		log.Fatal("If condition did not evaluate to a boolean")
	}

	if result.Bool {
		return ExecuteBlock(ifCommand.Commands, scope, programs)
	}
	return ExpressionResult{}, false, false
}

func runWhile(whileCommand ast.While, scope scope, programs programs) (ExpressionResult, bool, bool) {
	for {
		result, err := evaluateExpression(whileCommand.Cond, scope)
		if err != nil {
			log.Fatal(err)
		}

		if result.Type != Bool {
			log.Fatal("If condition did not evaluate to a boolean")
		}

		if result.Bool {
			val, didReturn, hasVal := ExecuteBlock(whileCommand.Commands, scope, programs)
			if didReturn {
				return val, didReturn, hasVal
			}
		} else {
			return ExpressionResult{}, false, false
		}
	}
}

func runProgram(programCommand ast.Program, programs programs) {
	programs[programCommand.Name.Name] = programCommand
}

func runCall(callCommand ast.Call, upperScope scope, programs programs) {
	program, ok := programs[callCommand.Name.Name]

	if !ok {
		log.Fatal("Tried to call " + callCommand.Name.Name + " but it has not been created")
	}

	if len(callCommand.Expressions) != len(program.Parameters) {
		log.Fatal("Tried to call " + callCommand.Name.Name + " with " + fmt.Sprint(len(callCommand.Expressions)) + " But it expects " + fmt.Sprint(len(program.Parameters)) + " parameters")
	}

	scope := scope{
		{},
	}

	for i, parameter := range program.Parameters {
		val, err := evaluateExpression(callCommand.Expressions[i], upperScope)
		if err != nil {
			log.Fatal(err)
		}
		defineInScope(parameter.Name, val, scope)
	}
	val, _, hasVal := ExecuteBlock(program.Commands, scope, programs)
	if hasVal {
		if callCommand.HasReturnTarget {
			assignInScope(callCommand.ReturnTarget.Name, val, upperScope)
		}
	}
}

func runReturn(returnCommand ast.Return, scope scope) (ExpressionResult, bool, bool) {

	if returnCommand.HasExpression {
		result, err := evaluateExpression(returnCommand.Expression, scope)

		if err != nil {
			log.Fatal(err)
		}

		return result, true, true
	}

	return ExpressionResult{}, true, false
}

func runCommand(command ast.Command, scope scope, programs programs) (ExpressionResult, bool, bool) {
	switch command := command.(type) {
	case ast.Log:
		runPrint(command, scope)
	case ast.Read:
		runRead(command, scope)
	case ast.Assign:
		runAssign(command, scope)
	case ast.New:
		runDefine(command, scope)
	case ast.If:
		val, didReturn, hasVal := runIf(command, scope, programs)
		if didReturn {
			return val, didReturn, hasVal
		}
	case ast.While:
		val, didReturn, hasVal := runWhile(command, scope, programs)
		if didReturn {
			return val, didReturn, hasVal
		}
	case ast.Program:
		runProgram(command, programs)
	case ast.Call:
		runCall(command, scope, programs)
	case ast.Return:
		return runReturn(command, scope)
	default:
		log.Fatal("Unrecognised command")
	}
	return ExpressionResult{}, false, false
}

func ExecuteBlock(program []ast.Command, scope scope, programs programs) (ExpressionResult, bool, bool) {
	scope = addScope(scope)
	for _, command := range program {
		val, didReturn, hasVal := runCommand(command, scope, programs)
		if didReturn {
			return val, didReturn, hasVal
		}
	}
	return ExpressionResult{}, false, false
}

func ExecuteProgram(program []ast.Command) {
	programs := make(programs)
	scope := scope{
		{},
	}
	ExecuteBlock(program, scope, programs)
}
