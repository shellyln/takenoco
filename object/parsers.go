package objparser

import (
	. "github.com/shellyln/takenoco/base"
	clsz "github.com/shellyln/takenoco/object/classes"
)

// Assertion that always match.
func Any() ParserFn {
	const ClassName = clsz.Any
	return LightBaseParser(ClassName, func(ctx ParserContext) (ParserContext, error) {
		ctx.MatchStatus = MatchStatus_Unmatched

		if ctx.Position+1 <= ctx.Slice.Len() {
			w := ctx.Slice.Get(ctx.Position)
			ctx.AstStack = append(ctx.AstStack, Ast{
				ClassName:      ClassName,
				Type:           AstType_Any,
				Value:          w,
				SourcePosition: ctx.SourcePosition,
			})
			ctx.Position += 1
			ctx.Length = 1
			ctx.MatchStatus = MatchStatus_Matched
		}
		return ctx, nil
	})
}

// Zero-width assertion at the end of the source.
func End() ParserFn {
	const ClassName = clsz.End
	return LightBaseParser(ClassName, func(ctx ParserContext) (ParserContext, error) {
		ctx.MatchStatus = MatchStatus_Unmatched

		if ctx.Position == ctx.Slice.Len() {
			ctx.Length = 0
			ctx.MatchStatus = MatchStatus_Matched
		}

		return ctx, nil
	})
}

// Assertion that match a sequence of values.
func Seq(seq ...interface{}) ParserFn {
	const ClassName = clsz.Seq
	return LightBaseParser(ClassName, func(ctx ParserContext) (ParserContext, error) {
		ctx.MatchStatus = MatchStatus_Unmatched

		length := len(seq)
		if ctx.Position+length <= ctx.Slice.Len() {
			w := ctx.Slice.Reslice(ctx.Position, ctx.Position+length)
			for i := 0; i < length; i++ {
				if !w.ItemEquals(w.Get(i), seq[i]) {
					break
				}
			}
			ctx.AstStack = append(ctx.AstStack, Ast{
				ClassName:      ClassName,
				Type:           AstType_ListOfAny,
				Value:          w,
				SourcePosition: ctx.SourcePosition,
			})
			ctx.Position += length
			ctx.Length = length
			ctx.MatchStatus = MatchStatus_Matched
		}
		return ctx, nil
	})
}

// Assertion that match if a value belongs to a set of values.
func ObjClass(oc ...interface{}) ParserFn {
	const ClassName = clsz.ObjClass
	return LightBaseParser(ClassName, func(ctx ParserContext) (ParserContext, error) {
		ctx.MatchStatus = MatchStatus_Unmatched

		if ctx.Position+1 <= ctx.Slice.Len() {
			w := ctx.Slice.Get(ctx.Position)
			for _, s := range oc {
				if ctx.Slice.ItemEquals(w, s) {
					ctx.AstStack = append(ctx.AstStack, Ast{
						ClassName:      ClassName,
						Type:           AstType_Any,
						Value:          w,
						SourcePosition: ctx.SourcePosition,
					})
					ctx.Position += 1
					ctx.Length = 1
					ctx.MatchStatus = MatchStatus_Matched
					break
				}
			}
		}
		return ctx, nil
	})
}

// Assertion that match if a value does not belong to a set of values.
func ObjClassN(oc ...interface{}) ParserFn {
	const ClassName = clsz.ObjClassN
	return LightBaseParser(ClassName, func(ctx ParserContext) (ParserContext, error) {
		ctx.MatchStatus = MatchStatus_Unmatched

		if ctx.Position+1 <= ctx.Slice.Len() {
			w := ctx.Slice.Get(ctx.Position)
			for _, s := range oc {
				if ctx.Slice.ItemEquals(w, s) {
					return ctx, nil
				}
			}

			ctx.AstStack = append(ctx.AstStack, Ast{
				ClassName:      ClassName,
				Type:           AstType_String,
				Value:          w,
				SourcePosition: ctx.SourcePosition,
			})
			ctx.Position += 1
			ctx.Length = 1
			ctx.MatchStatus = MatchStatus_Matched
		}

		return ctx, nil
	})
}

// Assertion that match if a value belongs to the set defined by the function.
func ObjClassFn(fn func(c interface{}) bool) ParserFn {
	const ClassName = clsz.ObjClassFn
	return LightBaseParser(ClassName, func(ctx ParserContext) (ParserContext, error) {
		ctx.MatchStatus = MatchStatus_Unmatched

		if ctx.Position+1 <= ctx.Slice.Len() {
			w := ctx.Slice.Get(ctx.Position)
			if fn(w) {
				ctx.AstStack = append(ctx.AstStack, Ast{
					ClassName:      ClassName,
					Type:           AstType_Any,
					Value:          w,
					SourcePosition: ctx.SourcePosition,
				})
				ctx.Position += 1
				ctx.Length = 1
				ctx.MatchStatus = MatchStatus_Matched
			}
		}
		return ctx, nil
	})
}
