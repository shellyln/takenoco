package extra

import "unicode"

// Unicode word characters
func isUnicodeWord(r rune) bool {
	// ID_Continue + '$' + U+200C + U+200D
	// Alnum(), '_', '$', and ...
	return (unicode.Is(unicode.L, r) ||
		unicode.Is(unicode.Nl, r) ||
		unicode.Is(unicode.Other_ID_Start, r) ||
		unicode.Is(unicode.Mn, r) ||
		unicode.Is(unicode.Mc, r) ||
		unicode.Is(unicode.Nd, r) ||
		unicode.Is(unicode.Pc, r) ||
		unicode.Is(unicode.Other_ID_Continue, r) ||
		r == '$' || r == 0x0200c || r == 0x0200d) &&
		!unicode.Is(unicode.Pattern_Syntax, r) &&
		!unicode.Is(unicode.Pattern_White_Space, r)
}
