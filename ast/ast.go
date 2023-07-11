package ast

import "morklerork/symbols"

type Expression interface{}

type StringLiteral struct {
	Value string
}

type IntLiteral struct {
	Value int
}

type BooleanLiteral struct {
	Value bool
}

type VariableName struct {
	Name string
}

type ProgramName struct {
	Name string
}

type HeapAccess struct {
	IndexExpression Expression
}

type BinaryOperator struct {
	// There is no need to re-state the BinaryOperatorType's just re-use the symbols
	BinaryOperatorType symbols.BinaryOperatorType
	Lhs                Expression
	Rhs                Expression
}

type Command interface{}

type Log struct {
	Indent int
	Expr   Expression
}

type New struct {
	Indent       int
	VariableName string
	Expr         Expression
}

type Assign struct {
	Indent int
	Target Expression
	Expr   Expression
}

type If struct {
	Indent   int
	Cond     Expression
	Commands []Command
}

type While struct {
	Indent   int
	Cond     Expression
	Commands []Command
}

type Program struct {
	Indent     int
	Name       ProgramName
	Parameters []VariableName
	Commands   []Command
}

type Call struct {
	Indent          int
	Name            ProgramName
	Expressions     []Expression
	ReturnTarget    VariableName
	HasReturnTarget bool
}

type Return struct {
	Indent        int
	Name          ProgramName
	Expression    Expression
	HasExpression bool
}
