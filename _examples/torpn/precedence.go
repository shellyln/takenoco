package torpn

import (
	"errors"

	. "github.com/shellyln/takenoco/base"
	objparser "github.com/shellyln/takenoco/object"
)

// Production rule (Precedence = 3)
var expressionRule3 = Precedence{
	Rules: []ParserFn{
		Trans(
			FlatGroup(
				isOperator("UnaryOperator", []string{"-"}),
				anyOperand(),
			),
			transformUnaryOp,
		),
	},
	Rtol: true,
}

// Production rule (Precedence = 2)
var expressionRule2 = Precedence{
	Rules: []ParserFn{
		Trans(
			FlatGroup(
				anyOperand(),
				isOperator("BinaryOperator", []string{"*", "/"}),
				anyOperand(),
			),
			transformBinaryOp,
		),
	},
}

// Production rule (Precedence = 1)
var expressionRule1 = Precedence{
	Rules: []ParserFn{
		Trans(
			FlatGroup(
				anyOperand(),
				isOperator("BinaryOperator", []string{"+", "-"}),
				anyOperand(),
			),
			transformBinaryOp,
		),
	},
}

func transformUnaryOp(ctx ParserContext, asts AstSlice) (AstSlice, error) {
	switch asts[1].Type {
	case AstType_Int:
		opcode := asts[0].Value.(string)
		op1 := asts[1].Value.(int64)

		var v int64
		switch opcode {
		case "-":
			v = -op1
		}

		return AstSlice{{
			Type:      AstType_Int,
			ClassName: "Number",
			Value:     v,
		}}, nil
	case AstType_ListOfAst:
		ret := asts[1].Value.(AstSlice)
		ret = append(ret, asts[0])
		return AstSlice{{
			Type:      AstType_ListOfAst,
			ClassName: "List",
			Value:     ret,
		}}, nil
	default:
		return nil, errors.New("Cannot apply unary - operator")
	}
}

func transformBinaryOp(ctx ParserContext, asts AstSlice) (AstSlice, error) {
	opcode := asts[1].Value.(string)
	ret := make(AstSlice, 0, 3)

	switch asts[0].Type {
	case AstType_Int:
		ret = append(ret, asts[0])
	case AstType_ListOfAst:
		ret = append(ret, asts[0].Value.(AstSlice)...)
	default:
		return nil, errors.New("Cannot apply binary " + opcode + " operator")
	}

	switch asts[2].Type {
	case AstType_Int:
		ret = append(ret, asts[2])
	case AstType_ListOfAst:
		ret = append(ret, asts[2].Value.(AstSlice)...)
	default:
		return nil, errors.New("Cannot apply binary " + opcode + " operator")
	}

	ret = append(ret, asts[1])

	return AstSlice{{
		Type:      AstType_ListOfAst,
		ClassName: "List",
		Value:     ret,
	}}, nil
}

// Production rules
var precedences = []Precedence{
	expressionRule3,
	expressionRule2,
	expressionRule1,
}

// Production rules
func formulaProductionRules() TransformerFn {
	return ProductionRule(
		precedences,
		FlatGroup(Start(), objparser.Any(), objparser.End()),
	)
}
