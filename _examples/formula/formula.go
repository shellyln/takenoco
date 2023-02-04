package formula

import (
	"errors"
	"strconv"

	. "github.com/shellyln/takenoco/base"
	"github.com/shellyln/takenoco/extra"
	. "github.com/shellyln/takenoco/string"
)

var (
	rootParser ParserFn
)

func init() {
	rootParser = program()
}

// Remove the resulting AST.
func erase(fn ParserFn) ParserFn {
	return Trans(fn, Erase)
}

// Whitespaces
func sp0() ParserFn {
	return erase(ZeroOrMoreTimes(Whitespace()))
}

// Integer number operand
func number() ParserFn {
	return Trans(
		FlatGroup(
			extra.IntegerNumberStr(),
			WordBoundary(),
			erase(sp0()),
		),
		func(ctx ParserContext, asts AstSlice) (AstSlice, error) {
			v, err := strconv.ParseInt(asts[0].Value.(string), 10, 64)
			if err != nil {
				return nil, err
			}
			asts = AstSlice{{
				Type:      AstType_Int,
				ClassName: "Number",
				Value:     v,
			}}
			return asts, nil
		},
	)
}

// Unary operators
func unaryOperator() ParserFn {
	return Trans(
		FlatGroup(
			CharClass("-"),
			erase(sp0()),
		),
		ChangeClassName("UnaryOperator"),
	)
}

// Binary operators
func binaryOperator() ParserFn {
	return Trans(
		FlatGroup(
			CharClass("+", "-", "*", "/"),
			erase(sp0()),
		),
		ChangeClassName("BinaryOperator"),
	)
}

// Expression without parentheses
func simpleExpression() ParserFn {
	return FlatGroup(
		number(),
		ZeroOrMoreTimes(
			binaryOperator(),
			number(),
		),
	)
}

// Expression enclosed in parenthesis
func groupedExpresion() ParserFn {
	return FlatGroup(
		erase(CharClass("(")),
		First(
			FlatGroup(
				erase(sp0()),
				expression(),
				erase(CharClass(")")),
				erase(sp0()),
			),
			Error("Error in grouped expression"),
		),
	)
}

// Expression before applying production rules
func expressionInner() ParserFn {
	return FlatGroup(
		ZeroOrMoreTimes(unaryOperator()),
		First(
			simpleExpression(),
			Indirect(groupedExpresion),
			Error("Value required"),
		),
		ZeroOrMoreTimes(
			binaryOperator(),
			First(
				FlatGroup(
					ZeroOrMoreTimes(unaryOperator()),
					First(
						simpleExpression(),
						Indirect(groupedExpresion),
					),
				),
				Error("Error in the expression after the binary operator"),
			),
		),
	)
}

// Single expression
func expression() ParserFn {
	return Trans(
		expressionInner(),
		formulaProductionRules(),
	)
}

// Entire program
func program() ParserFn {
	return FlatGroup(
		Start(),
		erase(sp0()),
		expression(),
		End(),
	)
}

// Parser
func Parse(s string) (int64, error) {
	out, err := rootParser(*NewStringParserContext(s))
	if err != nil {
		pos := GetLineAndColPosition(s, out.SourcePosition, 4)
		return 0, errors.New(
			err.Error() +
				"\n --> Line " + strconv.Itoa(pos.Line) +
				", Col " + strconv.Itoa(pos.Col) + "\n" +
				pos.ErrSource)
	}

	if out.MatchStatus == MatchStatus_Matched {
		return out.AstStack[0].Value.(int64), nil
	} else {
		pos := GetLineAndColPosition(s, out.SourcePosition, 4)
		return 0, errors.New(
			"Parse failed" +
				"\n --> Line " + strconv.Itoa(pos.Line) +
				", Col " + strconv.Itoa(pos.Col) + "\n" +
				pos.ErrSource)
	}
}
