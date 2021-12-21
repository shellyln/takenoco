package parser

import "reflect"

// Case matcher for Match().
type MatchCaseIfFn func(asts AstSlice) bool

// Case matcher and AST transformer for Match().
type Case struct {
	If  MatchCaseIfFn
	Let []TransformerFn
}

// Test which case matches for up to n previous (parsed) AST elements.
func Match(n int) func(cases ...Case) ParserFn {
	const ClassName = "Match"
	return func(cases ...Case) ParserFn {
		return LightBaseParser(ClassName, func(ctx ParserContext) (ParserContext, error) {
			sp := len(ctx.AstStack) - n
			ctx.Length = 0
			ctx.MatchStatus = MatchStatus_Unmatched

			for _, c := range cases {
				if c.If(ctx.AstStack) {
					var asts []Ast = ctx.AstStack[sp:]
					var err error

					for _, transform := range c.Let {
						asts, err = transform(ctx, asts)
						if err != nil {
							ctx.MatchStatus = MatchStatus_Error
							return ctx, err
						}
					}
					ctx.AstStack = append(ctx.AstStack[:sp], asts...)

					ctx.MatchStatus = MatchStatus_Matched
					return ctx, nil
				}
			}
			return ctx, nil
		})
	}
}

// Case matcher for Match(). The Ast is rune.
func TopIsRune(v ...rune) MatchCaseIfFn {
	return func(asts AstSlice) bool {
		if 0 < len(asts) {
			s, ok := asts[len(asts)-1].Value.(rune)
			if !ok {
				return false
			}
			for _, z := range v {
				if s == z {
					return true
				}
			}
			return false
		} else {
			return false
		}
	}
}

// Case matcher for Match(). The Ast is signed integer.
func TopIsInt(v ...int64) MatchCaseIfFn {
	return func(asts AstSlice) bool {
		if 0 < len(asts) {
			s, ok := asts[len(asts)-1].Value.(int64)
			if !ok {
				return false
			}
			for _, z := range v {
				if s == z {
					return true
				}
			}
			return false
		} else {
			return false
		}
	}
}

// Case matcher for Match(). The Ast is unsigned integer.
func TopIsUint(v ...uint64) MatchCaseIfFn {
	return func(asts AstSlice) bool {
		if 0 < len(asts) {
			s, ok := asts[len(asts)-1].Value.(uint64)
			if !ok {
				return false
			}
			for _, z := range v {
				if s == z {
					return true
				}
			}
			return false
		} else {
			return false
		}
	}
}

// Case matcher for Match(). The Ast is float.
func TopIsFloat(v ...float64) MatchCaseIfFn {
	return func(asts AstSlice) bool {
		if 0 < len(asts) {
			s, ok := asts[len(asts)-1].Value.(float64)
			if !ok {
				return false
			}
			for _, z := range v {
				if s == z {
					return true
				}
			}
			return false
		} else {
			return false
		}
	}
}

// Case matcher for Match(). The Ast is boolean.
func TopIsBool(v ...bool) MatchCaseIfFn {
	return func(asts AstSlice) bool {
		if 0 < len(asts) {
			s, ok := asts[len(asts)-1].Value.(bool)
			if !ok {
				return false
			}
			for _, z := range v {
				if s == z {
					return true
				}
			}
			return false
		} else {
			return false
		}
	}
}

// Case matcher for Match(). The Ast is string.
func TopIsStr(v ...string) MatchCaseIfFn {
	return func(asts AstSlice) bool {
		if 0 < len(asts) {
			s, ok := asts[len(asts)-1].Value.(string)
			if !ok {
				return false
			}
			for _, z := range v {
				if s == z {
					return true
				}
			}
			return false
		} else {
			return false
		}
	}
}

// Case matcher for Match().
func TopIs(v ...interface{}) MatchCaseIfFn {
	return func(asts AstSlice) bool {
		if 0 < len(asts) {
			s := asts[len(asts)-1].Value
			for _, z := range v {
				if reflect.DeepEqual(s, z) {
					return true
				}
			}
			return false
		} else {
			return false
		}
	}
}

// Case matcher for Match(). The Ast is any. It always matches.
func TopIsAny(asts AstSlice) bool {
	return true
}
