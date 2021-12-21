package strparser

import (
	"errors"
	"strconv"
	"strings"

	. "github.com/shellyln/takenoco/base"
)

const (
	// https://en.wikipedia.org/wiki/Whitespace_character
	// HT, LF, VT, FF, CR, SP, NEL, NBSP
	cutset string = "\u0009\u000a\u000b\u000c\u000d\u0020\u0085\u00a0"
)

// Transform the result AST array.
// Concatenate strings.
func Concat(_ ParserContext, asts AstSlice) (AstSlice, error) {
	if len(asts) == 0 {
		return AstSlice{{
			Type:  AstType_String,
			Value: "",
		}}, nil
	}
	var sb strings.Builder
	for _, w := range asts {
		if w.Type != AstType_String {
			return nil, errors.New("Transformer:Concat: Bad source type:" + w.Type.String())
		}
		sb.WriteString(w.Value.(string))
	}
	return AstSlice{{
		ClassName:      asts[0].ClassName,
		Type:           AstType_String,
		Value:          sb.String(),
		SourcePosition: asts[0].SourcePosition,
	}}, nil
}

// Transform the result AST array.
// Trim a string.
func Trim(ctx ParserContext, asts AstSlice) (AstSlice, error) {
	asts, err := Concat(ctx, asts)
	if err != nil {
		return nil, err
	}
	asts[0].Value = strings.Trim(asts[0].Value.(string), cutset)
	return asts, nil
}

// Transform the result AST array.
// Trims the beginning of a string.
func TrimStart(ctx ParserContext, asts AstSlice) (AstSlice, error) {
	asts, err := Concat(ctx, asts)
	if err != nil {
		return nil, err
	}
	asts[0].Value = strings.TrimLeft(asts[0].Value.(string), cutset)
	return asts, nil
}

// Transform the result AST array.
// Trim the end of a string.
func TrimEnd(ctx ParserContext, asts AstSlice) (AstSlice, error) {
	asts, err := Concat(ctx, asts)
	if err != nil {
		return nil, err
	}
	asts[0].Value = strings.TrimRight(asts[0].Value.(string), cutset)
	return asts, nil
}

// Transform the result AST array.
func parseIntImpl(base int, ctx ParserContext, asts AstSlice) (AstSlice, error) {
	asts, err := Concat(ctx, asts)
	if err != nil {
		return nil, err
	}
	str := asts[0].Value.(string)
	num, err := strconv.ParseInt(str, base, 64)
	if err != nil {
		return nil, errors.New("Transformer:parseIntImpl:Bad number format:" + str + ":" + err.Error())
	}
	asts[0].Type = AstType_Int
	asts[0].Value = num
	return asts, nil
}

// Transform the result AST array.
// Parses a signed integer by specifying a radix.
func ParseIntRadix(base int) TransformerFn {
	return func(ctx ParserContext, asts AstSlice) (AstSlice, error) {
		return parseIntImpl(base, ctx, asts)
	}
}

// Transform the result AST array.
// Parses a signed integer.
func ParseInt(ctx ParserContext, asts AstSlice) (AstSlice, error) {
	return parseIntImpl(10, ctx, asts)
}

// Transform the result AST array.
func parseUintImpl(base int, ctx ParserContext, asts AstSlice) (AstSlice, error) {
	asts, err := Concat(ctx, asts)
	if err != nil {
		return nil, err
	}
	str := asts[0].Value.(string)
	str2 := str
	if str2[0] == '+' {
		str2 = str2[1:]
	}
	num, err := strconv.ParseUint(str2, base, 64)
	if err != nil {
		return nil, errors.New("Transformer:parseUintImpl:Bad number format:" + str + ":" + err.Error())
	}
	asts[0].Type = AstType_Uint
	asts[0].Value = num
	return asts, nil
}

// Transform the result AST array.
// Parses an unsigned integer by specifying a radix.
func ParseUintRadix(base int) TransformerFn {
	return func(ctx ParserContext, asts AstSlice) (AstSlice, error) {
		return parseUintImpl(base, ctx, asts)
	}
}

// Transform the result AST array.
// Parses an unsigned integer.
func ParseUint(ctx ParserContext, asts AstSlice) (AstSlice, error) {
	return parseUintImpl(10, ctx, asts)
}

// Transform the result AST array.
// Parses a floating point number.
func ParseFloat(ctx ParserContext, asts AstSlice) (AstSlice, error) {
	asts, err := Concat(ctx, asts)
	if err != nil {
		return nil, err
	}
	str := asts[0].Value.(string)
	num, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return nil, errors.New("Transformer:ParseFloat:Bad number format:" + str + ":" + err.Error())
	}
	asts[0].Type = AstType_Float
	asts[0].Value = num
	return asts, nil
}

// Transform the result AST array.
// Convert integer to character.
func RuneFromInt(ctx ParserContext, asts AstSlice) (AstSlice, error) {
	return AstSlice{{
		ClassName:      asts[0].ClassName,
		Type:           AstType_Any,
		Value:          rune(asts[0].Value.(int64)),
		SourcePosition: asts[0].SourcePosition,
	}}, nil
}

// Transform the result AST array.
// Convert character to integer.
func IntFromRune(ctx ParserContext, asts AstSlice) (AstSlice, error) {
	return AstSlice{{
		ClassName:      asts[0].ClassName,
		Type:           AstType_String,
		Value:          int64(asts[0].Value.(rune)),
		SourcePosition: asts[0].SourcePosition,
	}}, nil
}

// Transform the result AST array.
// Convert integer to string.
func StringFromInt(ctx ParserContext, asts AstSlice) (AstSlice, error) {
	return AstSlice{{
		ClassName:      asts[0].ClassName,
		Type:           AstType_String,
		Value:          string(rune(asts[0].Value.(int64))),
		SourcePosition: asts[0].SourcePosition,
	}}, nil
}

// Transform the result AST array.
// Convert character to string.
func StringFromRune(ctx ParserContext, asts AstSlice) (AstSlice, error) {
	return AstSlice{{
		ClassName:      asts[0].ClassName,
		Type:           AstType_String,
		Value:          string(asts[0].Value.(rune)),
		SourcePosition: asts[0].SourcePosition,
	}}, nil
}
