package torpn_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/shellyln/takenoco/_examples/torpn"
)

type args struct {
	s string
}

type testMatrixItem struct {
	name    string
	args    args
	want    []interface{}
	wantErr bool
}

func runMatrixParse(t *testing.T, tests []testMatrixItem) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := torpn.Parse(tt.args.s)
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
		name:    "1",
		args:    args{s: "1+2*3"},
		want:    []interface{}{int64(1), int64(2), int64(3), "*", "+"},
		wantErr: false,
	}, {
		name:    "2",
		args:    args{s: "(1+2)*3"},
		want:    []interface{}{int64(1), int64(2), "+", int64(3), "*"},
		wantErr: false,
	}, {
		name:    "3",
		args:    args{s: "((((1+2)))*3)"},
		want:    []interface{}{int64(1), int64(2), "+", int64(3), "*"},
		wantErr: false,
	}, {
		name:    "4",
		args:    args{s: "-1-2*-3"},
		want:    []interface{}{int64(-1), int64(2), int64(-3), "*", "-"},
		wantErr: false,
	}, {
		name:    "5",
		args:    args{s: "-(1+2)*3"},
		want:    []interface{}{int64(1), int64(2), "+", "-", int64(3), "*"},
		wantErr: false,
	}, {
		name:    "6",
		args:    args{s: "--(1+2)*3"},
		want:    []interface{}{int64(1), int64(2), "+", "-", "-", int64(3), "*"},
		wantErr: false,
	}, {
		name:    "7",
		args:    args{s: "3+5*7+11"},
		want:    []interface{}{int64(3), int64(5), int64(7), "*", "+", int64(11), "+"},
		wantErr: false,
	}, {
		name:    "8",
		args:    args{s: "(3+5)*(7+11)"},
		want:    []interface{}{int64(3), int64(5), "+", int64(7), int64(11), "+", "*"},
		wantErr: false,
	}, {
		name:    "8",
		args:    args{s: "13"},
		want:    []interface{}{int64(13)},
		wantErr: false,
	}}

	runMatrixParse(t, tests)
}
