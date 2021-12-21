package parser

import (
	"errors"

	clsz "github.com/shellyln/takenoco/base/classes"
)

// Avoid errors by not making recursive calls at parser construction time,
// but by delaying them at runtime.
func Indirect(fn func() ParserFn) ParserFn {
	const ClassName = clsz.Indirect
	var parser ParserFn
	return LightBaseParser(ClassName, func(ctx ParserContext) (ParserContext, error) {
		if parser == nil {
			parser = fn()
		}
		return parser(ctx)
	})
}

// Conditional expression for parser construction time.
func If(b bool, fnT ParserFn, fnF ParserFn) ParserFn {
	if b {
		return fnT
	} else {
		return fnF
	}
}

// Zero-width assertion (always error)
func Error(msg string) ParserFn {
	const ClassName = clsz.Error
	return LightBaseParser(ClassName, func(ctx ParserContext) (ParserContext, error) {
		ctx.Length = 0
		ctx.MatchStatus = MatchStatus_Error
		return ctx, errors.New(msg)
	})
}

// Zero-width assertion (always unmatched)
func Unmatched() ParserFn {
	const ClassName = clsz.Unmatched
	return LightBaseParser(ClassName, func(ctx ParserContext) (ParserContext, error) {
		ctx.Length = 0
		ctx.MatchStatus = MatchStatus_Unmatched
		return ctx, nil
	})
}

// Zero-width assertion (always matched)
func Zero(astsToInsert ...Ast) ParserFn {
	const ClassName = clsz.Zero
	return LightBaseParser(ClassName, func(ctx ParserContext) (ParserContext, error) {
		savedLen := len(ctx.AstStack)
		ctx.AstStack = append(ctx.AstStack, astsToInsert...)
		for i := 0; i < len(astsToInsert); i++ {
			ctx.AstStack[savedLen+i].SourcePosition = ctx.SourcePosition
		}
		ctx.Length = 0
		ctx.MatchStatus = MatchStatus_Matched
		return ctx, nil
	})
}

// Zero-width assertion of the source start position.
func Start() ParserFn {
	const ClassName = clsz.Start
	return LightBaseParser(ClassName, func(ctx ParserContext) (ParserContext, error) {
		ctx.MatchStatus = MatchStatus_Unmatched

		if ctx.Position == 0 {
			ctx.Length = 0
			ctx.MatchStatus = MatchStatus_Matched
		}

		return ctx, nil
	})
}

// Grouping assertion. The resulting AST will NOT be grouped by AstSlice.
// Flat ASTs will be returned.
func FlatGroup(children ...ParserFn) ParserFn {
	const ClassName = clsz.FlatGroup
	return BaseParser(ClassName, nil, nil, children, nil)
}

// Grouping assertion. The resulting AST will be grouped by AstSlice.
func Group(children ...ParserFn) ParserFn {
	const ClassName = clsz.Group
	return BaseParser(ClassName, nil, nil, children, []TransformerFn{GroupingTransform})
}

// First is the synonym for Or.
func First(children ...ParserFn) ParserFn {
	const ClassName = clsz.First
	return BaseParser(ClassName, nil, []interface{}{ThereExists{}}, children, nil)
}

//
// func Or(children ...ParserFn) ParserFn {
// 	const ClassName = "Or"
// 	return BaseParser(ClassName, nil, []interface{}{ThereExists{}}, children, nil)
// }

// TODO: ChoiceShortest
// TODO: ChoiceLongest

// Look-ahead assertion.
func LookAhead(children ...ParserFn) ParserFn {
	const ClassName = clsz.LookAhead
	return BaseParser(ClassName, nil, []interface{}{Rewind{}}, children, nil)
}

// Negation look-ahead assertion.
func LookAheadN(children ...ParserFn) ParserFn {
	const ClassName = clsz.LookAheadN
	return BaseParser(ClassName, nil, []interface{}{Rewind{}, Negative{}}, children, nil)
}

//
func lookBehindBase(negative bool, minN, maxN int, children ...ParserFn) ParserFn {
	var parser ParserFn
	if negative {
		parser = BaseParser(clsz.LookBehindN, nil, []interface{}{Rewind{}}, children, nil)
	} else {
		parser = BaseParser(clsz.LookBehind, nil, []interface{}{Rewind{}}, children, nil)
	}

	return func(ctx ParserContext) (ParserContext, error) {
		for i := minN; i <= maxN; i++ {
			out := ctx
			out.Position -= i
			if out.Position < 0 {
				continue
			}
			out, err := parser(out)
			if err != nil {
				return out, err
			}
			if out.MatchStatus == MatchStatus_Matched {
				if negative {
					out.MatchStatus = MatchStatus_Unmatched
					return out, nil
				} else {
					out = ctx
					out.MatchStatus = MatchStatus_Matched
					return ctx, nil
				}
			}
		}
		out := ctx
		if negative {
			out.MatchStatus = MatchStatus_Matched
		} else {
			out.MatchStatus = MatchStatus_Unmatched
		}
		return out, nil
	}
}

// Look-behind assertion.
func LookBehind(minN, maxN int, children ...ParserFn) ParserFn {
	return lookBehindBase(false, minN, maxN, children...)
}

// Negation look-behind assertion.
func LookBehindN(minN, maxN int, children ...ParserFn) ParserFn {
	return lookBehindBase(true, minN, maxN, children...)
}

// TODO: LookBehindFn

// Repetitive assertion. {n,m}
func Repeat(times Times, children ...ParserFn) ParserFn {
	const ClassName = clsz.Repeat
	return BaseParser(ClassName, nil, []interface{}{times}, children, nil)
}

// TODO: QtyShortest(qty, child ParserFn, subsequent ...ParserFn)

// Repetitive assertion. {1,1}
func Once(children ...ParserFn) ParserFn {
	return Repeat(qtyOnce, children...)
}

// Repetitive assertion. {0,1}
func ZeroOrOnce(children ...ParserFn) ParserFn {
	return Repeat(qtyZeroOrOnce, children...)
}

// Repetitive assertion. {0,}
func ZeroOrMoreTimes(children ...ParserFn) ParserFn {
	return Repeat(qtyZeroOrMoreTimes, children...)
}

// Repetitive assertion. {1,1}
func OneOrMoreTimes(children ...ParserFn) ParserFn {
	return Repeat(qtyOneOrMoreTimes, children...)
}

// Transform the result AST array.
func Trans(child ParserFn, tr ...TransformerFn) ParserFn {
	const ClassName = clsz.Trans
	return BaseParser(ClassName, nil, nil, []ParserFn{child}, tr)
}
