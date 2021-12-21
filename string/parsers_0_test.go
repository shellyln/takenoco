package strparser

import (
	. "github.com/shellyln/takenoco/base"
)

type parserWant struct {
	hasErr      bool
	matchStatus MatchStatusType
	astStack    AstSlice
}

func astSliceEquals(a AstSlice, b AstSlice) bool {
	if a == nil && b == nil {
		return true
	}
	if a != nil && b == nil || a == nil && b != nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if !a.ItemEquals(a[i], b[i]) {
			return false
		}
	}
	return true
}
