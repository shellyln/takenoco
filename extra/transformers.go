package extra

import (
	"time"

	. "github.com/shellyln/takenoco/base"
	clsz "github.com/shellyln/takenoco/extra/classes"
)

func ParseDate(ctx ParserContext, asts AstSlice) (AstSlice, error) {
	value := asts[len(asts)-1].Value.(string)
	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		return nil, err
	}
	return AstSlice{{
		ClassName: clsz.Date,
		Type:      AstType_Any,
		Value:     t.UTC(),
	}}, nil
}

func ParseDateTime(ctx ParserContext, asts AstSlice) (AstSlice, error) {
	// TODO: BUG: Cannot parse years with negative values or years greater than or equal to 10000.
	value := asts[len(asts)-1].Value.(string)
	t, err := time.Parse("2006-01-02T15:04:05.000000000-07:00", value)
	if err != nil {
		return nil, err
	}
	return AstSlice{{
		ClassName: clsz.DateTime,
		Type:      AstType_Any,
		Value:     t.UTC(),
	}}, nil
}

func ParseTime(ctx ParserContext, asts AstSlice) (AstSlice, error) {
	value := "1970-01-01T" + asts[len(asts)-1].Value.(string) + "+00:00"
	t, err := time.Parse("2006-01-02T15:04:05.000000000-07:00", value)
	if err != nil {
		return nil, err
	}
	return AstSlice{{
		ClassName: clsz.Time,
		Type:      AstType_Any,
		Value:     t.UTC(),
	}}, nil
}
