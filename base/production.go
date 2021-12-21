package parser

import "errors"

// Production rules precedence
type Precedence struct {
	// Slice of parsers that input:AST -> output:AST.
	// If a part of the input is matched by one of these parsers,
	// the production rule will be applied to the matched part.
	Rules []ParserFn
	// True if applying the rules from right to left.
	Rtol bool
}

// Constructor for production rule
func newContext(slice SliceLike, t interface{}) *ParserContext {
	return &ParserContext{
		Slice:    slice,
		AstStack: make(AstSlice, 0, 1024),
		SourcePosition: SourcePosition{
			Position: 0,
			Length:   0,
		},
		Tag: t,
	}
}

// Transform the slices of AST according to the production rules.
func ProductionRule(precedences []Precedence, check ParserFn) TransformerFn {
	return func(ctx ParserContext, asts AstSlice) (AstSlice, error) {
		for {
			matched := false
		PRECEDENCE:
			for _, precedence := range precedences {
				astCtx := *newContext(asts, ctx.Tag)

				for i := 0; i <= asts.Len(); i++ {
					for _, rule := range precedence.Rules {
						curCtx := astCtx
						if precedence.Rtol {
							curCtx.Position = asts.Len() - i
						} else {
							curCtx.Position = i
						}

						out, err := rule(curCtx)
						if err != nil {
							return nil, err
						}
						if out.MatchStatus >= MatchStatus_Unmatched {
							continue
						}

						z := make(AstSlice, 0, asts.Len())
						z = append(z, asts[0:curCtx.Position]...)
						z = append(z, out.AstStack...)
						z = append(z, asts[out.Position:]...)
						asts = z

						matched = true
						break PRECEDENCE
					}
				}
			}

			if out, err := check(*newContext(asts, ctx.Tag)); err == nil && out.MatchStatus == MatchStatus_Matched {
				break
			}

			if !matched {
				// TODO:
				return nil, errors.New("Production rules are not matched.")
			}
		}
		return asts, nil
	}
}
