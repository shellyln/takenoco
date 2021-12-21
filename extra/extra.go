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
