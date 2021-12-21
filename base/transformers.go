package parser

import (
	"errors"

	clsz "github.com/shellyln/takenoco/base/classes"
)

// Transform the result AST array.
// Erase the AST of the current token.
func Erase(_ ParserContext, asts AstSlice) (AstSlice, error) {
	return asts[0:0], nil
}

// Transform the result AST array.
// Raise an error.
func TransformError(s string) TransformerFn {
	return func(_ ParserContext, asts AstSlice) (AstSlice, error) {
		return asts, errors.New(s)
	}
}

// Transform the result AST array.
// Group the AST of the current token into AstSlice.
func GroupingTransform(_ ParserContext, asts AstSlice) (AstSlice, error) {
	w := make(AstSlice, 0, len(asts))
	w = append(w, asts...)
	return AstSlice{{
		ClassName: clsz.Group,
		Type:      AstType_ListOfAst,
		Value:     w,
	}}, nil
}

// Transform the result AST array.
// Convert to a slice of the source element type.
func ToSlice(ctx ParserContext, asts AstSlice) (AstSlice, error) {
	if len(asts) == 0 {
		return AstSlice{{
			Type:  AstType_ListOfAny,
			Value: ctx.Slice.Make(0, 0),
		}}, nil
	}
	s := ctx.Slice.Make(len(asts), len(asts))
	for i, w := range asts {
		s.Set(i, w.Value)
	}
	return AstSlice{{
		ClassName:      asts[0].ClassName,
		Type:           AstType_ListOfAny,
		Value:          s,
		SourcePosition: asts[0].SourcePosition,
	}}, nil
}

// Transform the result AST array.
// Set the opcode and class name to the ASTs.
func SetOpCodeAndClassName(opcode AstOpCodeType, name string) TransformerFn {
	return func(_ ParserContext, asts AstSlice) (AstSlice, error) {
		w := make(AstSlice, 0, len(asts))
		w = append(w, asts...)
		if 0 < len(w) {
			w[0].OpCode = opcode
			w[0].ClassName = name
		}
		return w, nil
	}
}

// Transform the result AST array.
// Set the opcode to the ASTs.
func SetOpCode(opcode AstOpCodeType) TransformerFn {
	return func(_ ParserContext, asts AstSlice) (AstSlice, error) {
		w := make(AstSlice, 0, len(asts))
		w = append(w, asts...)
		if 0 < len(w) {
			w[0].OpCode = opcode
		}
		return w, nil
	}
}

// Transform the result AST array.
// Set the class name to the ASTs.
func ChangeClassName(name string) TransformerFn {
	return func(_ ParserContext, asts AstSlice) (AstSlice, error) {
		w := make(AstSlice, 0, len(asts))
		w = append(w, asts...)
		if 0 < len(w) {
			w[0].ClassName = name
		}
		return w, nil
	}
}

// Transform the result AST array.
// Set the value to the ASTs.
func SetValue(typ AstType, v interface{}) TransformerFn {
	return func(_ ParserContext, asts AstSlice) (AstSlice, error) {
		w := make(AstSlice, 0, len(asts))
		w = append(w, asts...)
		if 0 < len(w) {
			w[0].Type = typ
			w[0].Value = v
		}
		return w, nil
	}
}

// Transform the result AST array.
// Prepend a AST to the AST slice.
func Prepend(ast Ast) TransformerFn {
	return func(_ ParserContext, asts AstSlice) (AstSlice, error) {
		if 0 < len(asts) {
			ast.SourcePosition = asts[0].SourcePosition
		}
		w := make(AstSlice, 1, len(asts)+1)
		w[0] = ast
		w = append(w, asts...)
		return w, nil
	}
}

// Transform the result AST array.
// Append a AST to the AST slice.
func Push(ast Ast) TransformerFn {
	return func(_ ParserContext, asts AstSlice) (AstSlice, error) {
		if 0 < len(asts) {
			ast.SourcePosition = asts[len(asts)-1].SourcePosition
		}
		w := make(AstSlice, 0, len(asts)+1)
		w = append(w, asts...)
		w = append(w, ast)
		return w, nil
	}
}

// Transform the result AST array.
// Pop a AST from the AST slice.
func Pop(_ ParserContext, asts AstSlice) (AstSlice, error) {
	if 0 < len(asts) {
		return asts[:len(asts)-1], nil
	} else {
		return asts, nil
	}
}

// Transform the result AST array.
// Exchange the top with one below it.
func Exchange(_ ParserContext, asts AstSlice) (AstSlice, error) {
	length := len(asts)
	if 2 <= length {
		w := make(AstSlice, 0, length)
		w = append(w, asts[:length-2]...)
		w[length-2] = asts[length-1]
		w[length-1] = asts[length-2]
		return w, nil
	} else {
		return asts, nil
	}
}

// Transform the result AST array.
// Roll the elements.
//
// e.g.) n==2
// asts bottom |0  |1  |2  |3  |...|n-2|n-1| top
//             <=  <=
//             |2  |3  |...|n-2|n-1|0  |1  |
//
// e.g.) n==-2
// asts bottom |0  |1  |...|n-4|n-3|n-2|n-1| top
//                                    =>  =>
//             |n-2|n-1|0  |1  |...|n-4|n-3|
func Roll(n int) TransformerFn {
	return func(_ ParserContext, asts AstSlice) (AstSlice, error) {
		length := len(asts)
		if 0 <= n {
			if n < length {
				return asts, nil
			} else {
				return asts, nil
			}
		} else {
			m := -n
			if m < length {
				return asts, nil
			} else {
				return asts, nil
			}
		}
	}
}
