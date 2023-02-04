package strparser

import (
	"strings"
	"unicode/utf8"

	. "github.com/shellyln/takenoco/base"
	clsz "github.com/shellyln/takenoco/string/classes"
)

// Assertion that always match.
func Any() ParserFn {
	const ClassName = clsz.Any
	return LightBaseParser(ClassName, func(ctx ParserContext) (ParserContext, error) {
		ctx.MatchStatus = MatchStatus_Unmatched

		ch, length := utf8.DecodeRuneInString(ctx.Str[ctx.Position:])
		if length == 0 {
			return ctx, nil
		}
		s := string(ch)

		ctx.AstStack = append(ctx.AstStack, Ast{
			ClassName:      ClassName,
			Type:           AstType_String,
			Value:          s,
			SourcePosition: ctx.SourcePosition,
		})
		ctx.Position += length
		ctx.Length = length
		ctx.MatchStatus = MatchStatus_Matched

		return ctx, nil
	})
}

// Zero-width assertion at the end of the source.
func End() ParserFn {
	const ClassName = clsz.End
	return LightBaseParser(ClassName, func(ctx ParserContext) (ParserContext, error) {
		ctx.MatchStatus = MatchStatus_Unmatched

		if ctx.Position == len(ctx.Str) {
			ctx.Length = 0
			ctx.MatchStatus = MatchStatus_Matched
		}

		return ctx, nil
	})
}

// Assertion that match a sequence of characters.
func Seq(s string) ParserFn {
	const ClassName = clsz.Seq
	return LightBaseParser(ClassName, func(ctx ParserContext) (ParserContext, error) {
		ctx.MatchStatus = MatchStatus_Unmatched

		length := len(s)
		if ctx.Position+length <= len(ctx.Str) {
			w := ctx.Str[ctx.Position : ctx.Position+length]
			if w == s {
				ctx.AstStack = append(ctx.AstStack, Ast{
					ClassName:      ClassName,
					Type:           AstType_String,
					Value:          w,
					SourcePosition: ctx.SourcePosition,
				})
				ctx.Position += length
				ctx.Length = length
				ctx.MatchStatus = MatchStatus_Matched
			}
		}
		return ctx, nil
	})
}

// Assertion that match a sequence of characters. (ignore case)
func SeqI(s string) ParserFn {
	const ClassName = clsz.SeqI
	return LightBaseParser(ClassName, func(ctx ParserContext) (ParserContext, error) {
		ctx.MatchStatus = MatchStatus_Unmatched

		length := len(s)
		if ctx.Position+length <= len(ctx.Str) {
			w := ctx.Str[ctx.Position : ctx.Position+length]
			if strings.EqualFold(w, s) {
				ctx.AstStack = append(ctx.AstStack, Ast{
					ClassName:      ClassName,
					Type:           AstType_String,
					Value:          w,
					SourcePosition: ctx.SourcePosition,
				})
				ctx.Position += length
				ctx.Length = length
				ctx.MatchStatus = MatchStatus_Matched
			}
		}
		return ctx, nil
	})
}

// Assertion that match a range of characters.
func CharRange(cr ...RuneRange) ParserFn {
	const ClassName = clsz.CharRange
	return LightBaseParser(ClassName, func(ctx ParserContext) (ParserContext, error) {
		ctx.MatchStatus = MatchStatus_Unmatched

		ch, length := utf8.DecodeRuneInString(ctx.Str[ctx.Position:])
		if length == 0 {
			return ctx, nil
		}
		s := string(ch)

		for _, rr := range cr {
			if rr.Start <= ch && ch <= rr.End {
				ctx.AstStack = append(ctx.AstStack, Ast{
					ClassName:      ClassName,
					Type:           AstType_String,
					Value:          s,
					SourcePosition: ctx.SourcePosition,
				})
				ctx.Position += length
				ctx.Length = length
				ctx.MatchStatus = MatchStatus_Matched
				break
			}
		}
		return ctx, nil
	})
}

// Assertion that does not match a range of characters.
func CharRangeN(cr ...RuneRange) ParserFn {
	const ClassName = clsz.CharRangeN
	return LightBaseParser(ClassName, func(ctx ParserContext) (ParserContext, error) {
		ctx.MatchStatus = MatchStatus_Unmatched

		ch, length := utf8.DecodeRuneInString(ctx.Str[ctx.Position:])
		if length == 0 {
			return ctx, nil
		}
		s := string(ch)

		for _, rr := range cr {
			if rr.Start <= ch && ch <= rr.End {
				return ctx, nil
			}
		}
		ctx.AstStack = append(ctx.AstStack, Ast{
			ClassName:      ClassName,
			Type:           AstType_String,
			Value:          s,
			SourcePosition: ctx.SourcePosition,
		})
		ctx.Position += length
		ctx.Length = length
		ctx.MatchStatus = MatchStatus_Matched
		return ctx, nil
	})
}

// Assertion that match if a value belongs to a set of characters.
func CharClass(cc ...string) ParserFn {
	const ClassName = clsz.CharClass
	return LightBaseParser(ClassName, func(ctx ParserContext) (ParserContext, error) {
		ctx.MatchStatus = MatchStatus_Unmatched

		for _, s := range cc {
			length := len(s)
			if ctx.Position+length <= len(ctx.Str) {
				w := ctx.Str[ctx.Position : ctx.Position+length]
				if w == s {
					ctx.AstStack = append(ctx.AstStack, Ast{
						ClassName:      ClassName,
						Type:           AstType_String,
						Value:          w,
						SourcePosition: ctx.SourcePosition,
					})
					ctx.Position += length
					ctx.Length = length
					ctx.MatchStatus = MatchStatus_Matched
					break
				}
			}
		}
		return ctx, nil
	})
}

// Assertion that match if a value does not belong to a set of characters.
func CharClassN(cc ...string) ParserFn {
	const ClassName = clsz.CharClassN
	return LightBaseParser(ClassName, func(ctx ParserContext) (ParserContext, error) {
		ctx.MatchStatus = MatchStatus_Unmatched

		for _, s := range cc {
			length := len(s)
			if ctx.Position+length <= len(ctx.Str) {
				w := ctx.Str[ctx.Position : ctx.Position+length]
				if w == s {
					return ctx, nil
				}
			}
		}

		ch, length := utf8.DecodeRuneInString(ctx.Str[ctx.Position:])
		if length == 0 {
			return ctx, nil
		}
		s := string(ch)

		ctx.AstStack = append(ctx.AstStack, Ast{
			ClassName:      ClassName,
			Type:           AstType_String,
			Value:          s,
			SourcePosition: ctx.SourcePosition,
		})
		ctx.Position += length
		ctx.Length = length
		ctx.MatchStatus = MatchStatus_Matched

		return ctx, nil
	})
}

// Assertion that match if a value belongs to the set defined by the function.
func CharClassFn(fn func(c rune) bool) ParserFn {
	const ClassName = clsz.CharClassFn
	return LightBaseParser(ClassName, func(ctx ParserContext) (ParserContext, error) {
		ctx.MatchStatus = MatchStatus_Unmatched

		ch, length := utf8.DecodeRuneInString(ctx.Str[ctx.Position:])
		if length == 0 {
			return ctx, nil
		}
		s := string(ch)

		if fn(ch) {
			ctx.AstStack = append(ctx.AstStack, Ast{
				ClassName:      ClassName,
				Type:           AstType_String,
				Value:          s,
				SourcePosition: ctx.SourcePosition,
			})
			ctx.Position += length
			ctx.Length = length
			ctx.MatchStatus = MatchStatus_Matched
		}
		return ctx, nil
	})
}

// Character class of ASCII and Latin-1 whitespace characters.
func Whitespace() ParserFn {
	const ClassName = clsz.Whitespace
	return LightBaseParser(ClassName, func(ctx ParserContext) (ParserContext, error) {
		ctx.MatchStatus = MatchStatus_Unmatched

		ch, length := utf8.DecodeRuneInString(ctx.Str[ctx.Position:])
		if length == 0 {
			return ctx, nil
		}
		s := string(ch)

		if isWhitespace(ch) {
			ctx.AstStack = append(ctx.AstStack, Ast{
				ClassName:      ClassName,
				Type:           AstType_String,
				Value:          s,
				SourcePosition: ctx.SourcePosition,
			})
			ctx.Position += length
			ctx.Length = length
			ctx.MatchStatus = MatchStatus_Matched
		}
		return ctx, nil
	})
}

// Character class of ASCII and Latin-1 whitespace characters.
func WhitespaceNoLineBreak() ParserFn {
	const ClassName = clsz.WhitespaceNoLineBreak
	return LightBaseParser(ClassName, func(ctx ParserContext) (ParserContext, error) {
		ctx.MatchStatus = MatchStatus_Unmatched

		ch, length := utf8.DecodeRuneInString(ctx.Str[ctx.Position:])
		if length == 0 {
			return ctx, nil
		}
		s := string(ch)

		if isWhitespaceNoLineBreak(ch) {
			ctx.AstStack = append(ctx.AstStack, Ast{
				ClassName:      ClassName,
				Type:           AstType_String,
				Value:          s,
				SourcePosition: ctx.SourcePosition,
			})
			ctx.Position += length
			ctx.Length = length
			ctx.MatchStatus = MatchStatus_Matched
		}
		return ctx, nil
	})
}

// Character class of ASCII and Latin-1 whitespace characters.
func LineBreak() ParserFn {
	const ClassName = clsz.LineBreak
	return LightBaseParser(ClassName, func(ctx ParserContext) (ParserContext, error) {
		ctx.MatchStatus = MatchStatus_Unmatched

		ch, length := utf8.DecodeRuneInString(ctx.Str[ctx.Position:])
		if length == 0 {
			return ctx, nil
		}
		s := string(ch)

		if isLineBreak(ch) {
			ctx.AstStack = append(ctx.AstStack, Ast{
				ClassName:      ClassName,
				Type:           AstType_String,
				Value:          s,
				SourcePosition: ctx.SourcePosition,
			})
			ctx.Position += length
			ctx.Length = length
			ctx.MatchStatus = MatchStatus_Matched
		}
		return ctx, nil
	})
}

// Character class of ASCII alphabet characters.
func Alpha() ParserFn {
	const ClassName = clsz.Alpha
	return LightBaseParser(ClassName, func(ctx ParserContext) (ParserContext, error) {
		ctx.MatchStatus = MatchStatus_Unmatched

		ch, length := utf8.DecodeRuneInString(ctx.Str[ctx.Position:])
		if length == 0 {
			return ctx, nil
		}
		s := string(ch)

		if 'A' <= ch && ch <= 'Z' || 'a' <= ch && ch <= 'z' {
			ctx.AstStack = append(ctx.AstStack, Ast{
				ClassName:      ClassName,
				Type:           AstType_String,
				Value:          s,
				SourcePosition: ctx.SourcePosition,
			})
			ctx.Position += length
			ctx.Length = length
			ctx.MatchStatus = MatchStatus_Matched
		}
		return ctx, nil
	})
}

// Character class of ASCII number characters.
func Number() ParserFn {
	const ClassName = clsz.Number
	return LightBaseParser(ClassName, func(ctx ParserContext) (ParserContext, error) {
		ctx.MatchStatus = MatchStatus_Unmatched

		ch, length := utf8.DecodeRuneInString(ctx.Str[ctx.Position:])
		if length == 0 {
			return ctx, nil
		}
		s := string(ch)

		if '0' <= ch && ch <= '9' {
			ctx.AstStack = append(ctx.AstStack, Ast{
				ClassName:      ClassName,
				Type:           AstType_String,
				Value:          s,
				SourcePosition: ctx.SourcePosition,
			})
			ctx.Position += length
			ctx.Length = length
			ctx.MatchStatus = MatchStatus_Matched
		}
		return ctx, nil
	})
}

// Character class of ASCII alphabet and number characters.
func Alnum() ParserFn {
	const ClassName = clsz.Alnum
	return LightBaseParser(ClassName, func(ctx ParserContext) (ParserContext, error) {
		ctx.MatchStatus = MatchStatus_Unmatched

		ch, length := utf8.DecodeRuneInString(ctx.Str[ctx.Position:])
		if length == 0 {
			return ctx, nil
		}
		s := string(ch)

		if 'A' <= ch && ch <= 'Z' || 'a' <= ch && ch <= 'z' || '0' <= ch && ch <= '9' {
			ctx.AstStack = append(ctx.AstStack, Ast{
				ClassName:      ClassName,
				Type:           AstType_String,
				Value:          s,
				SourcePosition: ctx.SourcePosition,
			})
			ctx.Position += length
			ctx.Length = length
			ctx.MatchStatus = MatchStatus_Matched
		}
		return ctx, nil
	})
}

// Character class of ASCII binary number characters.
func BinNumber() ParserFn {
	const ClassName = clsz.BinNumber
	return LightBaseParser(ClassName, func(ctx ParserContext) (ParserContext, error) {
		ctx.MatchStatus = MatchStatus_Unmatched

		ch, length := utf8.DecodeRuneInString(ctx.Str[ctx.Position:])
		if length == 0 {
			return ctx, nil
		}
		s := string(ch)

		if '0' <= ch && ch <= '1' {
			ctx.AstStack = append(ctx.AstStack, Ast{
				ClassName:      ClassName,
				Type:           AstType_String,
				Value:          s,
				SourcePosition: ctx.SourcePosition,
			})
			ctx.Position += length
			ctx.Length = length
			ctx.MatchStatus = MatchStatus_Matched
		}
		return ctx, nil
	})
}

// Character class of ASCII octal number characters.
func OctNumber() ParserFn {
	const ClassName = clsz.OctNumber
	return LightBaseParser(ClassName, func(ctx ParserContext) (ParserContext, error) {
		ctx.MatchStatus = MatchStatus_Unmatched

		ch, length := utf8.DecodeRuneInString(ctx.Str[ctx.Position:])
		if length == 0 {
			return ctx, nil
		}
		s := string(ch)

		if '0' <= ch && ch <= '7' {
			ctx.AstStack = append(ctx.AstStack, Ast{
				ClassName:      ClassName,
				Type:           AstType_String,
				Value:          s,
				SourcePosition: ctx.SourcePosition,
			})
			ctx.Position += length
			ctx.Length = length
			ctx.MatchStatus = MatchStatus_Matched
		}
		return ctx, nil
	})
}

// Character class of ASCII hex number characters.
func HexNumber() ParserFn {
	const ClassName = clsz.HexNumber
	return LightBaseParser(ClassName, func(ctx ParserContext) (ParserContext, error) {
		ctx.MatchStatus = MatchStatus_Unmatched

		ch, length := utf8.DecodeRuneInString(ctx.Str[ctx.Position:])
		if length == 0 {
			return ctx, nil
		}
		s := string(ch)

		if '0' <= ch && ch <= '7' || 'A' <= ch && ch <= 'F' || 'a' <= ch && ch <= 'f' {
			ctx.AstStack = append(ctx.AstStack, Ast{
				ClassName:      ClassName,
				Type:           AstType_String,
				Value:          s,
				SourcePosition: ctx.SourcePosition,
			})
			ctx.Position += length
			ctx.Length = length
			ctx.MatchStatus = MatchStatus_Matched
		}
		return ctx, nil
	})
}

// Zero-width assertion on a word boundary.
func WordBoundary() ParserFn {
	const ClassName = clsz.WordBoundary
	return LightBaseParser(ClassName, func(ctx ParserContext) (ParserContext, error) {
		ctx.Length = 0
		ctx.MatchStatus = MatchStatus_Unmatched

		ch, length := utf8.DecodeRuneInString(ctx.Str[ctx.Position:])

		if ctx.Position == 0 {
			if 0 < length && isWord(ch) {
				ctx.MatchStatus = MatchStatus_Matched
			}
			return ctx, nil
		}

		var prevCh rune
		var prevChLength int
		for i := ctx.Position - 1; 0 <= i; i-- {
			b := ctx.Str[i]

			// TODO: BUG: Valid range is `b <= 0x7f || 0xc2 <= b && b <= 0xf4`
			// https://github.com/golang/go/blob/0a86cd6857b9fb12a798b3dbcfb6974384aa07d6/src/unicode/utf8/utf8.go#L65-L84

			if b <= 0x7f || 0xc2 <= b && b <= 0xf0 || b == 0xf3 {
				prevCh, prevChLength = utf8.DecodeRuneInString(ctx.Str[i:])
				break
			}
		}

		if ctx.Position == len(ctx.Str) {
			if 0 < prevChLength && isWord(prevCh) {
				ctx.MatchStatus = MatchStatus_Matched
			}
		} else {
			if length != 0 && prevChLength != 0 {
				if isWord(prevCh) && !isWord(ch) || !isWord(prevCh) && isWord(ch) {
					ctx.MatchStatus = MatchStatus_Matched
				}
			}
		}

		return ctx, nil
	})
}
