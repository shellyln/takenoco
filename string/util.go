package strparser

// ASCII and Latin-1 whitespace characters
func isWhitespace(r rune) bool {
	// https://en.wikipedia.org/wiki/Whitespace_character
	// HT, LF, VT, FF, CR, SP, NEL, NBSP
	if r == 0x0009 || r == 0x000a || r == 0x000b || r == 0x000c || r == 0x000d || r == 0x0020 || r == 0x0085 || r == 0x00a0 {
		return true
	} else {
		return false
	}
}

// ASCII and Latin-1 whitespace characters
func isWhitespaceNoLineBreak(r rune) bool {
	// https://en.wikipedia.org/wiki/Whitespace_character
	// HT, SP, NBSP
	if r == 0x0009 || r == 0x0020 || r == 0x00a0 {
		return true
	} else {
		return false
	}
}

// ASCII and Latin-1 whitespace characters
func isLineBreak(r rune) bool {
	// https://en.wikipedia.org/wiki/Whitespace_character
	// LF, VT, FF, CR, NEL
	if r == 0x000a || r == 0x000b || r == 0x000c || r == 0x000d || r == 0x0085 {
		return true
	} else {
		return false
	}
}

// ASCII word characters
func isWord(r rune) bool {
	if 'A' <= r && r <= 'Z' || 'a' <= r && r <= 'z' || '0' <= r && r <= '9' || r == '_' {
		return true
	} else {
		return false
	}
}
