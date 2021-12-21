package parser

// Common implementation for parsers that do not have child sub-parsers.
func LightBaseParser(className string, fn LightParserImplFn) ParserFn {
	parser := func(ctx ParserContext) (ParserContext, error) {
		out, err := fn(ctx)
		if err != nil && out.MatchStatus < MatchStatus_Error {
			out.MatchStatus = MatchStatus_Error
		}
		if out.MatchStatus == MatchStatus_Matched {
			out.Quantity = 1
		}
		out.ClassName = className
		return out, err
	}

	if traceEnabled {
		traceParserTrackingNo++
		return tracer(traceScope, parserTracer, traceParserTrackingNo, className, parser)
	} else {
		return parser
	}
}

// Common implementation for parsers.
func BaseParser(
	className string,
	fn ParserImplFn, props []interface{},
	children []ParserFn, tr []TransformerFn) ParserFn {

	negative := false
	thereExists := false
	rewind := false
	qty := qtyOnce
	if props != nil {
		for _, p := range props {
			switch w := p.(type) {
			case Times:
				qty = w
			case Negative:
				negative = true
			case ThereExists:
				thereExists = true
			case Rewind:
				rewind = true
			}
		}
	}

	parser := func(ctx ParserContext) (ParserContext, error) {
		ctx.ClassName = className

		out := ctx
		bottomOfAst := len(ctx.AstStack)
		count := 0
		var err error = nil

		out.MatchStatus = MatchStatus_Matched

	PARENT:
		for ; qty.Max < 0 || count < qty.Max; count++ {
			saved := out
			numChildrenMatched := 0

		CHILDREN:
			for _, child := range children {
				prev := out

				out, err = child(out)

				if err != nil && out.MatchStatus < MatchStatus_Error {
					out.MatchStatus = MatchStatus_Error
				}

				switch out.MatchStatus {
				case MatchStatus_Error:
					return out, err
				case MatchStatus_Unmatched:
					if thereExists {
						// rewind current child
						out = prev
						continue CHILDREN
					} else {
						// rewind all children
						out = saved
						break PARENT
					}
				}
				numChildrenMatched++
				if thereExists {
					break CHILDREN
				}
			}

			if thereExists {
				if 0 < len(children) && numChildrenMatched == 0 {
					// rewind all children
					out = saved
					break PARENT
				}
			}

			out.Quantity = count
			if fn != nil {
				out, err = fn(numChildrenMatched, bottomOfAst, out)
				if err != nil && out.MatchStatus < MatchStatus_Error {
					out.MatchStatus = MatchStatus_Error
				}
			}

			switch out.MatchStatus {
			case MatchStatus_Error:
				return out, err
			case MatchStatus_Unmatched:
				// rewind all children
				out = saved
				break PARENT
			}
		}

		if 0 <= qty.Min && count < qty.Min {
			out.MatchStatus = MatchStatus_Unmatched
		}

		if negative {
			switch out.MatchStatus {
			case MatchStatus_Unmatched:
				out.MatchStatus = MatchStatus_Matched
			case MatchStatus_Matched:
				out.MatchStatus = MatchStatus_Unmatched
			}
		}

		if MatchStatus_Unmatched <= out.MatchStatus {
			return out, nil
		}

		out.Quantity = count
		if fn != nil {
			out, err = fn(-1, bottomOfAst, out)
			if err != nil && out.MatchStatus < MatchStatus_Error {
				out.MatchStatus = MatchStatus_Error
			}
		}

		if MatchStatus_Unmatched <= out.MatchStatus {
			return out, nil
		}

		if rewind {
			out = ctx
		}

		if tr != nil {
			var asts []Ast = out.AstStack[len(ctx.AstStack):]
			for _, transform := range tr {
				asts, err = transform(ctx, asts)
				if err != nil {
					out.MatchStatus = MatchStatus_Error
					return out, err
				}
			}
			out.AstStack = append(out.AstStack[:len(ctx.AstStack)], asts...)
		}

		return out, nil
	}

	if traceEnabled {
		traceParserTrackingNo++
		return tracer(traceScope, parserTracer, traceParserTrackingNo, className, parser)
	} else {
		return parser
	}
}
