package extra

import (
	"unicode"
	"unicode/utf8"

	. "github.com/shellyln/takenoco/base"
	clsz "github.com/shellyln/takenoco/extra/classes"
	. "github.com/shellyln/takenoco/string"
	strclsz "github.com/shellyln/takenoco/string/classes"
)

// Exclude parsed ASTs from the results.
func erase(fn ParserFn) ParserFn {
	return Trans(fn, Erase)
}

// Parse the binary number.
func BinaryNumberStr() ParserFn {
	return Trans(
		FlatGroup(
			CharRange(RuneRange{Start: '0', End: '1'}),
			ZeroOrMoreTimes(
				First(
					CharRange(RuneRange{Start: '0', End: '1'}),
					erase(CharClass("_")),
				),
			),
		),
		Concat,
		ChangeClassName(clsz.BinaryNumberStr),
	)
}

// Parse the octal number.
func OctalNumberStr() ParserFn {
	return Trans(
		FlatGroup(
			CharRange(RuneRange{Start: '0', End: '7'}),
			ZeroOrMoreTimes(
				First(
					CharRange(RuneRange{Start: '0', End: '7'}),
					erase(CharClass("_")),
				),
			),
		),
		Concat,
		ChangeClassName(clsz.OctalNumberStr),
	)
}

// Parse the hexadecimal number.
func HexNumberStr() ParserFn {
	return Trans(
		FlatGroup(
			First(
				Number(),
				CharRange(RuneRange{Start: 'A', End: 'F'}),
				CharRange(RuneRange{Start: 'a', End: 'f'}),
			),
			ZeroOrMoreTimes(
				First(
					Number(),
					CharRange(RuneRange{Start: 'A', End: 'F'}),
					CharRange(RuneRange{Start: 'a', End: 'f'}),
					erase(CharClass("_")),
				),
			),
		),
		Concat,
		ChangeClassName(clsz.HexNumberStr),
	)
}

// Parse the integer number.
func IntegerNumberStr() ParserFn {
	return Trans(
		FlatGroup(
			ZeroOrOnce(CharClass("+", "-")),
			First(
				FlatGroup(
					CharClass("0"),
					ZeroOrMoreTimes(erase(CharClass("_"))),
				),
				FlatGroup(
					Number(),
					ZeroOrMoreTimes(First(Number(), erase(CharClass("_")))),
				),
			),
		),
		Concat,
		ChangeClassName(clsz.IntegerNumberStr),
	)
}

// Parse the exponent of floating point number.
func exponentStr() ParserFn {
	return FlatGroup(
		CharClass("E", "e"),
		ZeroOrOnce(CharClass("+", "-")),
		FlatGroup(
			Number(),
			ZeroOrMoreTimes(First(Number(), erase(CharClass("_")))),
		),
	)
}

// Parse the floating point number.
func FloatNumberStr() ParserFn {
	return Trans(
		FlatGroup(
			ZeroOrOnce(CharClass("+", "-")),
			First(
				FlatGroup(
					First(
						FlatGroup(
							CharClass("0"),
							ZeroOrMoreTimes(erase(CharClass("_"))),
						),
						FlatGroup(
							Number(),
							ZeroOrMoreTimes(First(Number(), erase(CharClass("_")))),
						),
					),
					FlatGroup(CharClass("."), ZeroOrMoreTimes(First(Number(), erase(CharClass("_"))))),
					ZeroOrOnce(exponentStr()),
				),
				FlatGroup(
					CharClass("."), OneOrMoreTimes(First(Number(), erase(CharClass("_")))),
					ZeroOrOnce(exponentStr()),
				),
				FlatGroup(
					First(
						FlatGroup(
							CharClass("0"),
							ZeroOrMoreTimes(erase(CharClass("_"))),
						),
						FlatGroup(
							Number(),
							ZeroOrMoreTimes(First(Number(), erase(CharClass("_")))),
						),
					),
					exponentStr(),
				),
			),
		),
		Concat,
		ChangeClassName(clsz.FloatNumberStr),
	)
}

// Parse the number.
func NumericStr() ParserFn {
	return Trans(
		FlatGroup(
			ZeroOrOnce(CharClass("+", "-")),
			First(
				FlatGroup(
					First(
						FlatGroup(
							CharClass("0"),
							ZeroOrMoreTimes(erase(CharClass("_"))),
						),
						FlatGroup(
							Number(),
							ZeroOrMoreTimes(First(Number(), erase(CharClass("_")))),
						),
					),
					ZeroOrOnce(CharClass("."), ZeroOrMoreTimes(First(Number(), erase(CharClass("_"))))),
				),
				FlatGroup(
					CharClass("."), OneOrMoreTimes(First(Number(), erase(CharClass("_")))),
				),
			),
			ZeroOrOnce(exponentStr()),
		),
		Concat,
		ChangeClassName(clsz.NumericStr),
	)
}

// Parse the ASCII identifier
func AsciiIdentifierStr() ParserFn {
	return Trans(
		FlatGroup(
			Once(First(
				Alpha(),
				CharClass("_", "$"),
			)),
			ZeroOrMoreTimes(First(
				Alnum(),
				CharClass("_", "$"),
			)),
		),
		Concat,
		ChangeClassName(clsz.IdentifierStr),
	)
}

// Parse the Unicode identifier
func UnicodeIdentifierStr() ParserFn {
	return Trans(
		FlatGroup(
			Once(First(
				// ID_Start + '_' + '$'
				CharClass("_", "$"),
				CharClassFn(func(c rune) bool {
					// ID_Start: Alpha(), and ...
					return (unicode.Is(unicode.L, c) ||
						unicode.Is(unicode.Nl, c) ||
						unicode.Is(unicode.Other_ID_Start, c)) &&
						!unicode.Is(unicode.Pattern_Syntax, c) &&
						!unicode.Is(unicode.Pattern_White_Space, c)
				}),
			)),
			ZeroOrMoreTimes(First(
				// ID_Continue + '$' + U+200C + U+200D
				CharClass("$"),
				CharClassFn(func(c rune) bool {
					// Alnum(), '_', and ...
					return (unicode.Is(unicode.L, c) ||
						unicode.Is(unicode.Nl, c) ||
						unicode.Is(unicode.Other_ID_Start, c) ||
						unicode.Is(unicode.Mn, c) ||
						unicode.Is(unicode.Mc, c) ||
						unicode.Is(unicode.Nd, c) ||
						unicode.Is(unicode.Pc, c) ||
						unicode.Is(unicode.Other_ID_Continue, c) ||
						c == 0x0200c || c == 0x0200d) &&
						!unicode.Is(unicode.Pattern_Syntax, c) &&
						!unicode.Is(unicode.Pattern_White_Space, c)
				}),
			)),
		),
		Concat,
		ChangeClassName(clsz.IdentifierStr),
	)
}

// Zero-width assertion on a word boundary.
func UnicodeWordBoundary() ParserFn {
	const ClassName = strclsz.WordBoundary
	return LightBaseParser(ClassName, func(ctx ParserContext) (ParserContext, error) {
		ctx.Length = 0
		ctx.MatchStatus = MatchStatus_Unmatched

		ch, length := utf8.DecodeRuneInString(ctx.Str[ctx.Position:])

		if ctx.Position == 0 {
			if 0 < length && isUnicodeWord(ch) {
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
			if 0 < prevChLength && isUnicodeWord(prevCh) {
				ctx.MatchStatus = MatchStatus_Matched
			}
		} else {
			if length != 0 && prevChLength != 0 {
				if isUnicodeWord(prevCh) && !isUnicodeWord(ch) || !isUnicodeWord(prevCh) && isUnicodeWord(ch) {
					ctx.MatchStatus = MatchStatus_Matched
				}
			}
		}

		return ctx, nil
	})
}

// Parse the ISO 8601 date string. (yyyy-MM-dd)
func DateStr() ParserFn {
	return Trans(
		FlatGroup(
			ZeroOrOnce(Seq("-")),
			Repeat(Times{Min: 4, Max: -1}, Number()),
			Seq("-"),
			CharRange(RuneRange{Start: '0', End: '1'}),
			CharRange(RuneRange{Start: '0', End: '9'}),
			Seq("-"),
			CharRange(RuneRange{Start: '0', End: '3'}),
			CharRange(RuneRange{Start: '0', End: '9'}),
		),
		Concat,
		ChangeClassName(clsz.DateStr),
	)
}

// Parse the ISO 8601 datetime string.
// (yyyy-MM-ddThh:mmZ , ... , yyyy-MM-ddThh:mm:ss.fffffffffZ)
// (yyyy-MM-ddThh:mm+00:00 , ... , yyyy-MM-ddThh:mm:ss.fffffffff+00:00)
func DateTimeStr() ParserFn {
	return Trans(
		FlatGroup(
			ZeroOrOnce(Seq("-")),
			Repeat(Times{Min: 4, Max: -1}, Number()),
			Seq("-"),
			CharRange(RuneRange{Start: '0', End: '1'}),
			CharRange(RuneRange{Start: '0', End: '9'}),
			Seq("-"),
			CharRange(RuneRange{Start: '0', End: '3'}),
			CharRange(RuneRange{Start: '0', End: '9'}),
			Seq("T"),
			CharRange(RuneRange{Start: '0', End: '2'}),
			CharRange(RuneRange{Start: '0', End: '9'}),
			Seq(":"),
			CharRange(RuneRange{Start: '0', End: '5'}),
			CharRange(RuneRange{Start: '0', End: '9'}),
			First(
				FlatGroup(
					Seq(":"),
					CharRange(RuneRange{Start: '0', End: '6'}),
					CharRange(RuneRange{Start: '0', End: '9'}),
					First(
						FlatGroup(
							Seq("."),
							Trans(
								Repeat(Times{Min: 1, Max: 9}, // 3: milli, 6: micro, 9: nano
									CharRange(RuneRange{Start: '0', End: '9'}),
								),
								Concat,
								func(ctx ParserContext, asts AstSlice) (AstSlice, error) {
									return AstSlice{{
										Type:  AstType_String,
										Value: (asts[len(asts)-1].Value.(string) + "000000000")[0:9],
									}}, nil
								},
							),
						),
						Zero(Ast{
							Type:  AstType_String,
							Value: ".000000000",
						}),
					),
				),
				Zero(Ast{
					Type:  AstType_String,
					Value: ":00.000000000",
				}),
			),
			First(
				FlatGroup(
					erase(Seq("Z")),
					Zero(Ast{
						Type:  AstType_String,
						Value: "+00:00",
					}),
				),
				FlatGroup(
					CharClass("+", "-"),
					Repeat(Times{Min: 2, Max: 2},
						CharRange(RuneRange{Start: '0', End: '9'}),
					),
					Seq(":"),
					CharRange(RuneRange{Start: '0', End: '5'}),
					CharRange(RuneRange{Start: '0', End: '9'}),
				),
			),
		),
		Concat,
		ChangeClassName(clsz.DateTimeStr),
	)
}

// Parse the ISO 8601 time string.
// (hh:mm , ... , hh:mm:ss.fffffffff)
func TimeStr() ParserFn {
	return Trans(
		FlatGroup(
			CharRange(RuneRange{Start: '0', End: '2'}),
			CharRange(RuneRange{Start: '0', End: '9'}),
			Seq(":"),
			CharRange(RuneRange{Start: '0', End: '5'}),
			CharRange(RuneRange{Start: '0', End: '9'}),
			First(
				FlatGroup(
					Seq(":"),
					CharRange(RuneRange{Start: '0', End: '6'}),
					CharRange(RuneRange{Start: '0', End: '9'}),
					First(
						FlatGroup(
							Seq("."),
							Trans(
								Repeat(Times{Min: 1, Max: 9}, // 3: milli, 6: micro, 9: nano
									CharRange(RuneRange{Start: '0', End: '9'}),
								),
								Concat,
								func(ctx ParserContext, asts AstSlice) (AstSlice, error) {
									return AstSlice{{
										Type:  AstType_String,
										Value: (asts[len(asts)-1].Value.(string) + "000000000")[0:9],
									}}, nil
								},
							),
						),
						Zero(Ast{
							Type:  AstType_String,
							Value: ".000000000",
						}),
					),
				),
				Zero(Ast{
					Type:  AstType_String,
					Value: ":00.000000000",
				}),
			),
		),
		Concat,
		ChangeClassName(clsz.TimeStr),
	)
}
