package parser

// Wrapper of  []T
type SliceLike interface {
	// Get length of the slice.
	Len() int
	// Get the element of the slice.
	Get(i int) interface{}
	// Set the element in the slice.
	Set(i int, v interface{})
	// Resize the slice.
	Reslice(start, end int) SliceLike
	// Copy the slice.
	Copy(start, end int) SliceLike
	// Make a new slice.
	Make(len, cap int) SliceLike
	// Compare the two elements.
	ItemEquals(a, b interface{}) bool
}

// Source position
type SourcePosition struct {
	// Index from the start of the source string or slice.
	Position int `json:"pos,omitempty"`
	// Length of the token in the source string or slice.
	Length int `json:"len,omitempty"`
}

// Source position (for error reporting)
type LineAndColPosition struct {
	LineIndex int
	Line      int
	Col       int
	Position  int
	ErrSource string
}

// Match result status
type MatchStatusType int

const (
	// Match result status: Matched
	MatchStatus_Matched MatchStatusType = iota
	// Match result status: Unmatched
	MatchStatus_Unmatched
	// Match result status: Error
	MatchStatus_Error
)

//
func (t MatchStatusType) String() string {
	switch t {
	case MatchStatus_Matched:
		return "Matched"
	case MatchStatus_Unmatched:
		return "Unmatched"
	case MatchStatus_Error:
		return "Error"
	default:
		return "Unknown"
	}
}

// Parser context
type ParserContext struct {
	// Source for string parser
	Str string
	// Source for object parser ([]T)
	Slice SliceLike
	// Position is next source position. Length is matched source length.
	SourcePosition
	// Number of times the assertion is matched.
	Quantity int
	// AST stacks. Push when grown, pop when rolled back.
	AstStack AstSlice
	// Match result status
	MatchStatus MatchStatusType
	// Class or stereotype of the matched token
	ClassName string
	// User-defined tag
	Tag interface{}
}

// Quantifier property of BaseParser().
type Times struct {
	// Minimum number of times
	Min int
	// Maximum number of times
	Max int
}

// Negation match property of BaseParser(). This is a marker object.
type Negative struct {
}

// Alternation (OR) match property of BaseParser(). This is a marker object.
type ThereExists struct {
}

// Zero length look-ahead match property of BaseParser(). This is a marker object.
type Rewind struct {
}

// Character code range property of BaseParser().
type RuneRange struct {
	// The smallest character code in the range.
	Start rune
	// The largest character code in the range.
	End rune
}

// Type of the parser's implementation function.
type LightParserImplFn func(ctx ParserContext) (ParserContext, error)

// Type of the parser's implementation function.
type ParserImplFn func(numChildrenMatched, bottomOfAst int, ctx ParserContext) (ParserContext, error)

// Type of the parser function.
type ParserFn func(ctx ParserContext) (ParserContext, error)

// Type of the AST transformer function.
type TransformerFn func(ctx ParserContext, asts AstSlice) (AstSlice, error)
