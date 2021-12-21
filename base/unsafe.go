package parser

import "unsafe"

// Internal memory layout of `interface{}`
type rawInterface2 struct {
	Typ uintptr
	Ptr unsafe.Pointer
}
