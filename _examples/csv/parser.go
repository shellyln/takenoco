package csv

import (
	"errors"
	"strconv"

	. "github.com/shellyln/takenoco/base"
	. "github.com/shellyln/takenoco/string"
)

var (
	// Comma and line break characters
	cellBreakCharacters []string
	documentParser      ParserFn
)

func init() {
	cellBreakCharacters = make([]string, 0, len(LineBreakCharacters)+1)
	cellBreakCharacters = append(cellBreakCharacters, ",")
	cellBreakCharacters = append(cellBreakCharacters, LineBreakCharacters...)
	documentParser = document()
}

// Remove the resulting AST.
func erase(fn ParserFn) ParserFn {
	return Trans(fn, Erase)
}

// Whitespaces
func sp() ParserFn {
	return erase(ZeroOrMoreTimes(WhitespaceNoLineBreak()))
}

func quotedCell() ParserFn {
	return Trans(
		OneOrMoreTimes(
			FlatGroup(
				sp(),
				erase(Seq("\"")),
				ZeroOrMoreTimes(
					First(
						erase(Seq("\"\"")),
						CharClassN("\""),
					),
				),
				First(
					erase(Seq("\"")),
					FlatGroup(End(), Error("Unexpected EOF")),
				),
				sp(),
			),
		),
		Concat,
	)
}

func cell() ParserFn {
	return Trans(
		ZeroOrMoreTimes(CharClassN(cellBreakCharacters...)),
		Trim,
	)
}

// Convert AST to array data. (line)
func lineTransform(_ ParserContext, asts AstSlice) (AstSlice, error) {
	w := make([]string, len(asts))
	length := len(asts)

	for i := 0; i < length; i++ {
		w[i] = asts[i].Value.(string)
	}

	return AstSlice{{
		ClassName: "*Line",
		Type:      AstType_Any,
		Value:     w,
	}}, nil
}

func line() ParserFn {
	return Trans(
		FlatGroup(
			ZeroOrMoreTimes(
				First(quotedCell(), cell()),
				erase(Seq(",")),
			),
			First(quotedCell(), cell()),
		),
		lineTransform,
	)
}

// Convert AST to array data. (Entire document)
func documentTransform(_ ParserContext, asts AstSlice) (AstSlice, error) {
	length := len(asts)
	w := make([][]string, length)

	for i := 0; i < length; i++ {
		w[i] = asts[i].Value.([]string)
	}
	for i := length - 1; i >= 0; i-- {
		if len(w[i]) == 0 || len(w[i]) == 1 && w[i][0] == "" {
			w = w[:i]
		} else {
			break
		}
	}

	return AstSlice{{
		ClassName: "*Document",
		Type:      AstType_Any,
		Value:     w,
	}}, nil
}

func document() ParserFn {
	return Trans(
		FlatGroup(
			ZeroOrMoreTimes(
				line(),
				erase(OneOrMoreTimes(LineBreak())),
			),
			line(),
			End(),
		),
		documentTransform,
	)
}

func Parse(s string) ([][]string, error) {
	out, err := documentParser(*NewStringParserContext(s))
	if err != nil {
		return nil, err
	} else {
		if out.MatchStatus == MatchStatus_Matched {
			return out.AstStack[0].Value.([][]string), nil
		} else {
			return nil, errors.New("Parse failed at " + strconv.Itoa(out.SourcePosition.Position))
		}
	}
}
