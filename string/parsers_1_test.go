package strparser

import (
	"testing"

	. "github.com/shellyln/takenoco/base"
)

func TestAny(t *testing.T) {
}

func TestSeq(t *testing.T) {
	type args struct {
		s    string
		text string
	}

	tests := []struct {
		name string
		args args
		want parserWant
	}{{
		name: "Case 1",
		args: args{
			s:    "abc",
			text: "abc",
		},
		want: parserWant{
			hasErr:      false,
			matchStatus: MatchStatus_Matched,
			astStack: AstSlice{{
				ClassName: ":string:Seq",
				Type:      AstType_String,
				Value:     "abc",
			}},
		},
	}, {
		name: "Case 2",
		args: args{
			s:    "abc",
			text: "abd",
		},
		want: parserWant{
			hasErr:      false,
			matchStatus: MatchStatus_Unmatched,
			astStack:    nil,
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Seq(tt.args.s)(*NewStringParserContext(tt.args.text))
			if tt.want.hasErr != (err != nil) {
				t.Errorf("Seq().err is %v, want %v", err, tt.want.hasErr)
				return
			}
			if !tt.want.hasErr {
				if tt.want.matchStatus != got.MatchStatus {
					t.Errorf("Seq().got.MatchStatus is %v, want %v", got.MatchStatus, tt.want.matchStatus)
					return
				}
			}
			if tt.want.matchStatus == MatchStatus_Matched {
				if !astSliceEquals(got.AstStack, tt.want.astStack) {
					t.Errorf("Seq().got.AstStack = %v, want %v", got, tt.want.astStack)
					return
				}
			}
		})
	}
}

func TestSeqI(t *testing.T) {
}

func TestCharRange(t *testing.T) {
}

func TestCharRangeN(t *testing.T) {
}

func TestCharClass(t *testing.T) {
}

func TestCharClassN(t *testing.T) {
}

func TestCharClassFn(t *testing.T) {
}

func TestWhitespace(t *testing.T) {
}

func TestWhitespaceNoLineBreak(t *testing.T) {
}

func TestLineBreak(t *testing.T) {
}

func TestAlpha(t *testing.T) {
}

func TestNumber(t *testing.T) {
}

func TestAlnum(t *testing.T) {
}

func TestWordBoundary(t *testing.T) {
}
