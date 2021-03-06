package extra

import (
	. "github.com/shellyln/takenoco/base"
	clsz "github.com/shellyln/takenoco/extra/classes"
	. "github.com/shellyln/takenoco/string"
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
