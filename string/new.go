package strparser

import (
	. "github.com/shellyln/takenoco/base"
)

// Constructor
func NewStringParserContext(s string) *ParserContext {
	return &ParserContext{
		Str:      s,
		AstStack: make(AstSlice, 0, 1024),
	}
}

// Constructor
func NewStringParserContextWithTag(s string, t interface{}) *ParserContext {
	return &ParserContext{
		Str:      s,
		AstStack: make(AstSlice, 0, 1024),
		Tag:      t,
	}
}
