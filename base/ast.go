package parser

import (
	"reflect"
	"unsafe"
)

// Ast.Value's type
type AstType uint

// Ast.OpCode's type
type AstOpCodeType uint64

const (
	// Ast.Value is Nil
	AstType_Nil AstType = iota
	// Ast.Value is rune (32bit unsigned)
	AstType_Rune
	// Ast.Value is int64
	AstType_Int
	// Ast.Value is uint64
	AstType_Uint
	// Ast.Value is float64
	AstType_Float
	// Ast.Value is bool
	AstType_Bool
	// Ast.Value is string
	AstType_String
	// Ast.Value is AstCons
	AstType_AstCons
	// Ast.Value is AstSlice
	AstType_ListOfAst
	// Ast.Value is []interface{}
	AstType_ListOfAny
	// Ast.Value is function
	AstType_Function
	// Ast.Value is interface{}
	AstType_Any
)

// Convert AstType to a string.
func (t AstType) String() string {
	switch t {
	case AstType_Nil:
		return "Nil"
	case AstType_Rune:
		return "Rune"
	case AstType_Int:
		return "Int"
	case AstType_Uint:
		return "Uint"
	case AstType_Float:
		return "Float"
	case AstType_Bool:
		return "Bool"
	case AstType_String:
		return "String"
	case AstType_AstCons:
		return "AstCons"
	case AstType_ListOfAst:
		return "ListOfAst"
	case AstType_ListOfAny:
		return "ListOfAny"
	case AstType_Function:
		return "Function"
	case AstType_Any:
		return "Any"
	default:
		return "Unknown"
	}
}

// Ast is an AST (Abstract Syntax Tree) node object.
// Implements the SliceLike interface.
type Ast struct {
	// User defined token's opcode
	OpCode AstOpCodeType `json:"op,omitempty"`
	// User defined token's class name, stereotype, etc.
	ClassName string `json:"cn,omitempty"`
	// Value's type
	Type AstType `json:"ty,omitempty"`
	// Token's value
	Value interface{} `json:"v,omitempty"`
	// Address of value
	Address Box `json:"-"`
	// Token's source position
	SourcePosition `json:",omitempty"`
}

// TODO: Ast.UnmarshalJSON()

// Mode of the outbound trip of the traverse.
type WayThereMode int

const (
	// If the traverse target is a list, returns the list.
	WayThereMode_None WayThereMode = iota
	// Returns the traverse target as is.
	// Calls to fnWayThere() and fnWayBack() are made.
	WayThereMode_Lazy
	// If the traverse target is a list, returns the last item.
	WayThereMode_Last // Scope / Breakable scope / Function
)

// The opcode for controlling the next execution position in the traverse function.
type TraverseOpcode uint16

const (
	// NOP
	TraverseOpcode_Nop TraverseOpcode = iota
	// Move the execution position to the end.
	TraverseOpcode_Break
	// Move the execution position to the first.
	TraverseOpcode_Continue
)

// Callback handler for Traverse.
// It is called on the outbound trip of the traverse.
type FnWayThere func(ctx interface{}, ast Ast) (
	out Ast, lazy WayThereMode, err error)

// Callback handler for Traverse.
// It is called on the return trip of the traverse.
type FnWayBack func(ctx interface{}, ast Ast) (
	out Ast, pcIncr int16, uOpcode TraverseOpcode, thrown interface{}, err error)

// Callback handler for Traverse.
// It is called when an error or throw occurs.
// You have to do the cleanup that you were going to do on FnWayBack.
type FnChildrenErr func(ctx interface{}, ast Ast, child Ast, thrown interface{}, err error)

// Perform depth-first traversing and tree transformation.
//
// Parameters:
//   fnWayThere    FnWayThere    - Callback handler for the outbound trip.
//   fnWayBack     FnWayBack     - Callback handler for the return trip.
//   fnChildrenErr FnChildrenErr - Callback handler for errors that occur in the child tree.
//   ctx           interface{}   - Context
//
// Returns:
//   ast     Ast           - AST to be traversed.
//   pcIncr  int16         -
//   uOpcode TraverseFlags -
//   thrown  interface{}   -
//   err     error         -
//
func (s Ast) Traverse(
	fnWayThere FnWayThere,
	fnWayBack FnWayBack,
	fnChildrenErr FnChildrenErr,
	ctx interface{}) (Ast, int16, TraverseOpcode, interface{}, error) {

	var mode WayThereMode
	var pcIncr int16
	var uOpcode TraverseOpcode
	var thrown interface{}
	var err error

	s, mode, err = fnWayThere(ctx, s)
	if err != nil {
		return s, 0, 0, nil, err
	}

	if mode != WayThereMode_Lazy {

		switch s.Type {
		case AstType_AstCons:
			{
				// cons := s.Value.(AstCons)
				cons := *(*AstCons)((*rawInterface2)(unsafe.Pointer(&s.Value)).Ptr)

				cons.Car, _, _, thrown, err =
					cons.Car.Traverse(fnWayThere, fnWayBack, fnChildrenErr, ctx)

				if thrown != nil || err != nil {
					fnChildrenErr(ctx, s, cons.Car, thrown, err)
					return s, 0, 0, thrown, err
				}

				cons.Cdr, _, _, thrown, err =
					cons.Cdr.Traverse(fnWayThere, fnWayBack, fnChildrenErr, ctx)

				if thrown != nil || err != nil {
					fnChildrenErr(ctx, s, cons.Cdr, thrown, err)
					return s, 0, 0, thrown, err
				}
				s.Value = cons
			}

		case AstType_ListOfAst:
			if mode == WayThereMode_Last {
				// orig := s.Value.(AstSlice)
				orig := *(*AstSlice)((*rawInterface2)(unsafe.Pointer(&s.Value)).Ptr)
				origLen := len(orig)
				sliceLen := origLen
				if sliceLen != 0 {
					sliceLen = 1
				}

				slice := make(AstSlice, sliceLen, sliceLen)
				var buf Ast

				for i := 0; i < origLen; i++ {
					buf, pcIncr, uOpcode, thrown, err =
						orig[i].Traverse(fnWayThere, fnWayBack, fnChildrenErr, ctx)

					if thrown != nil || err != nil {
						fnChildrenErr(ctx, s, orig[i], thrown, err)
						return s, 0, 0, thrown, err
					}

					switch uOpcode {
					case TraverseOpcode_Nop:
						slice[0] = buf
					case TraverseOpcode_Break:
						i = origLen - 1
					case TraverseOpcode_Continue:
						i = -1
					}

					i += int(pcIncr)
				}
				s.Value = slice
			} else {
				// orig := s.Value.(AstSlice)
				orig := *(*AstSlice)((*rawInterface2)(unsafe.Pointer(&s.Value)).Ptr)
				origLen := len(orig)
				slice := make(AstSlice, origLen, origLen)

				for i := 0; i < origLen; i++ {
					slice[i], _, _, thrown, err =
						orig[i].Traverse(fnWayThere, fnWayBack, fnChildrenErr, ctx)

					if thrown != nil || err != nil {
						fnChildrenErr(ctx, s, orig[i], thrown, err)
						return s, 0, 0, thrown, err
					}
				}
				s.Value = slice
			}
		}
	}

	return fnWayBack(ctx, s)
}

// AstCons is a Cons cell for Ast objects.
type AstCons struct {
	Car Ast `json:"car"`
	Cdr Ast `json:"cdr,omitempty"`
}

// AstSlice is a slice of Ast object.
// Imprements the interface SliceLike.
type AstSlice []Ast

// Imprements SliceLike.Len().
func (s AstSlice) Len() int {
	return len(s)
}

// Imprements SliceLike.Get().
func (s AstSlice) Get(i int) interface{} {
	return s[i]
}

// Imprements SliceLike.Set().
func (s AstSlice) Set(i int, v interface{}) {
	s[i] = v.(Ast)
}

// Imprements SliceLike.Reslice().
func (s AstSlice) Reslice(start, end int) SliceLike {
	return s[start:end]
}

// Imprements SliceLike.Copy().
func (s AstSlice) Copy(start, end int) SliceLike {
	w := make(AstSlice, 0, end-start)
	return append(w, s[start:end]...)
}

// Imprements SliceLike.Make().
func (s AstSlice) Make(len, cap int) SliceLike {
	return make(AstSlice, len, cap)
}

// Imprements SliceLike.ItemEquals().
func (s AstSlice) ItemEquals(a, b interface{}) bool {
	wa := a.(Ast)
	wb := b.(Ast)
	if wa.ClassName == wb.ClassName && wa.Type == wb.Type {
		switch wa.Type {
		case AstType_AstCons:
			ca := wa.Value.(AstCons)
			cb := wb.Value.(AstCons)
			if !tempEmptyAstSlice.ItemEquals(ca.Car, cb.Car) {
				return false
			}
			if !tempEmptyAstSlice.ItemEquals(ca.Cdr, cb.Cdr) {
				return false
			}
			return true
		case AstType_ListOfAst:
			la := wa.Value.(AstSlice)
			lb := wb.Value.(AstSlice)
			if len(la) == len(lb) {
				for i := 0; i < len(la); i++ {
					if !la.ItemEquals(la[i], lb[i]) {
						return false
					}
				}
				return true
			} else {
				return false
			}
		default:
			return reflect.DeepEqual(wa.Value, wb.Value)
		}
	} else {
		return false
	}
}
