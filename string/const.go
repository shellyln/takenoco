package strparser

// Don't change these variables at runtime.
var (
	// https://en.wikipedia.org/wiki/Whitespace_character
	// HT, LF, VT, FF, CR, SP, NEL, NBSP
	WhitespaceCharacters = []string{"\u0009", "\u000a", "\u000b", "\u000c", "\u000d", "\u0020", "\u0085", "\u00a0"}
	// HT, SP, NBSP
	WhitespaceNoLineBreakCharacters = []string{"\u0009", "\u0020", "\u00a0"}
	// LF, VT, FF, CR, NEL
	LineBreakCharacters = []string{"\u000a", "\u000b", "\u000c", "\u000d", "\u0085"}
)
