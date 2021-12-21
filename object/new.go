package objparser

import (
	. "github.com/shellyln/takenoco/base"
)

// Constructor
func NewObjectParserContext(slice SliceLike) *ParserContext {
	return &ParserContext{
		Slice:    slice,
		AstStack: make(AstSlice, 0, 1024),
	}
}

// Constructor
func NewObjectParserContextWithTag(slice SliceLike, t interface{}) *ParserContext {
	return &ParserContext{
		Slice:    slice,
		AstStack: make(AstSlice, 0, 1024),
		Tag:      t,
	}
}
