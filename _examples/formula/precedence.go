package formula

import (
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
			func(ctx ParserContext, asts AstSlice) (AstSlice, error) {
				opcode := asts[0].Value.(string)
				op1 := asts[1].Value.(int64)

				var v int64
				switch opcode {
				case "-":
					v = -op1
				}

				return AstSlice{{
					ClassName: "Number",
					Value:     v,
				}}, nil
			},
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
			func(ctx ParserContext, asts AstSlice) (AstSlice, error) {
				opcode := asts[1].Value.(string)
				op1 := asts[0].Value.(int64)
				op2 := asts[2].Value.(int64)

				var v int64
				switch opcode {
				case "*":
					v = op1 * op2
				case "/":
					v = op1 / op2
				}

				return AstSlice{{
					ClassName: "Number",
					Value:     v,
				}}, nil
			},
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
			func(ctx ParserContext, asts AstSlice) (AstSlice, error) {
				opcode := asts[1].Value.(string)
				op1 := asts[0].Value.(int64)
				op2 := asts[2].Value.(int64)

				var v int64
				switch opcode {
				case "+":
					v = op1 + op2
				case "-":
					v = op1 - op2
				}

				return AstSlice{{
					ClassName: "Number",
					Value:     v,
				}}, nil
			},
		),
	},
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
