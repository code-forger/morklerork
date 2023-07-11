package symbols

// Since go lacks union types, a 'Symbol' can technically be anything
// Its up to functions that receive Symbols to check they are one of the
// correct types
type Symbol interface{}

type Indent struct {
	Level int
}

type Print struct{}
type Assign struct{}
type Define struct{}
type If struct{}
type While struct{}
type Program struct{}
type Call struct{}
type Return struct{}

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
	IndexExpressionSymbol Symbol
}

// By having a single 'BinaryOperator' Symbol with an enum type inside,
// We can get precedence for free by using the numerical underpinning of
// the enum as the precedence
type BinaryOperatorType int

const (
	LogicalAndOperator BinaryOperatorType = iota
	LogicalOrOperator
	EqualOperator
	NotEqualOperator
	LTOperator
	PlusOperator
	MinusOperator
	DivideOperator
	TimesOperator
)

var BinaryOperatorTypeNames = map[BinaryOperatorType]string{
	0: "&",
	1: "|",
	2: "==",
	3: "!=",
	4: "<",
	5: "+",
	6: "-",
	7: "*",
	8: "/",
}

type BinaryOperator struct {
	BinaryOperatorType BinaryOperatorType
}
