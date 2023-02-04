package csv_test

import (
	"fmt"
	"reflect"
	"testing"

	. "github.com/shellyln/takenoco/_examples/csv"
)

type args struct {
	s string
}

type testMatrixItem struct {
	name    string
	args    args
	want    interface{}
	wantErr bool
}

func runMatrixParse(t *testing.T, tests []testMatrixItem) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
