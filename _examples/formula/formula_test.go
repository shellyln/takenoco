package formula_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/shellyln/takenoco/_examples/formula"
)

type args struct {
	s string
}

type testMatrixItem struct {
	name    string
	args    args
	want    int64
	wantErr bool
}

func runMatrixParse(t *testing.T, tests []testMatrixItem) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := formula.Parse(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("%v: Parse() error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%v: Parse() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []testMatrixItem{{
		name:    "1a",
		args:    args{s: "17"},
		want:    17,
		wantErr: false,
	}, {
		name:    "1b",
		args:    args{s: " 17 "},
		want:    17,
		wantErr: false,
	}, {
		name:    "1c",
		args:    args{s: "----17"},
		want:    17,
		wantErr: false,
	}, {
		name:    "2a",
		args:    args{s: "17+19"},
		want:    36,
		wantErr: false,
	}, {
		name:    "2b",
		args:    args{s: " 17 + 19 "},
		want:    36,
		wantErr: false,
	}, {
		name:    "3a",
		args:    args{s: "(1*2+3)*(4-5*6)+7"},
		want:    -123,
		wantErr: false,
	}, {
		name:    "3b",
		args:    args{s: " ( 1 * 2 + 3) * ( 4 - 5 * 6 ) + 7 "},
		want:    -123,
		wantErr: false,
	}, {
		name:    "4a",
		args:    args{s: "7+(1*2+3)*(4-5*6)"},
		want:    -123,
		wantErr: false,
	}, {
		name:    "4b",
		args:    args{s: " 7 + ( 1 * 2 + 3 ) * ( 4 - 5 * 6 ) "},
		want:    -123,
		wantErr: false,
	}, {
		name:    "5a",
		args:    args{s: "-7+(1*2+3)*(4-5*6)"},
		want:    -137,
		wantErr: false,
	}, {
		name:    "6a",
		args:    args{s: "-7+-(1*2+3)*(4-5*6)"},
		want:    123,
		wantErr: false,
	}, {
		name:    "7a",
		args:    args{s: "-7+-(1*2+3)*-(4-5*6)"},
		want:    -137,
		wantErr: false,
	}, {
		name:    "8a",
		args:    args{s: "-7+-(1*2+a3)*-(4-5*6)"},
		wantErr: true,
	}}

	runMatrixParse(t, tests)
}
